# mautrix-twitter - A Matrix-Twitter DM puppeting bridge
# Copyright (C) 2020 Tulir Asokan
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.
from typing import (Dict, Tuple, Optional, List, Deque, Set, Any, Union, AsyncGenerator,
                    NamedTuple, TYPE_CHECKING, cast)
from collections import deque
from datetime import datetime
from os import path
import asyncio
import html

import magic

from mautwitdm.types import (ConversationType, Conversation, MessageData, Participant, ReactionKey,
                             MessageAttachmentMedia, TimelineStatus, ReactionCreateEntry,
                             ReactionDeleteEntry, MessageEntry)
from mautrix.appservice import AppService, IntentAPI
from mautrix.bridge import BasePortal, NotificationDisabler
from mautrix.types import (EventID, MessageEventContent, RoomID, EventType, MessageType,
                           TextMessageEventContent, MediaMessageEventContent, ImageInfo, VideoInfo,
                           ContentURI, EncryptedFile)
from mautrix.errors import MatrixError, MForbidden
from mautrix.util.simple_lock import SimpleLock

from .db import Portal as DBPortal, Message as DBMessage, Reaction as DBReaction
from .config import Config
from . import user as u, puppet as p, matrix as m

if TYPE_CHECKING:
    from .__main__ import TwitterBridge

try:
    from mautrix.crypto.attachments import encrypt_attachment, decrypt_attachment
except ImportError:
    encrypt_attachment = decrypt_attachment = None

StateBridge = EventType.find("m.bridge", EventType.Class.STATE)
StateHalfShotBridge = EventType.find("uk.half-shot.bridge", EventType.Class.STATE)
BackfillEntryTypes = Union[MessageEntry, ReactionCreateEntry, ReactionDeleteEntry]
ReuploadedMediaInfo = NamedTuple('ReuploadedMediaInfo', mxc=Optional[ContentURI],
                                 decryption_info=Optional[EncryptedFile],
                                 mime_type=str, file_name=str, size=int)


class Portal(DBPortal, BasePortal):
    by_mxid: Dict[RoomID, 'Portal'] = {}
    by_twid: Dict[Tuple[str, int], 'Portal'] = {}
    config: Config
    matrix: 'm.MatrixHandler'
    az: AppService

    _main_intent: Optional[IntentAPI]
    _create_room_lock: asyncio.Lock
    backfill_lock: SimpleLock
    _msgid_dedup: Deque[int]
    _reqid_dedup: Set[str]
    _reaction_dedup: Deque[Tuple[int, int, ReactionKey]]

    _main_intent: IntentAPI
    _last_participant_update: Set[int]
    _reaction_lock: asyncio.Lock

    def __init__(self, twid: str, receiver: int, conv_type: ConversationType,
                 other_user: Optional[int] = None, mxid: Optional[RoomID] = None,
                 name: Optional[str] = None, encrypted: bool = False) -> None:
        super().__init__(twid, receiver, conv_type, other_user, mxid, name, encrypted)
        self._create_room_lock = asyncio.Lock()
        self.log = self.log.getChild(twid)
        self._msgid_dedup = deque(maxlen=100)
        self._reaction_dedup = deque(maxlen=100)
        self._reqid_dedup = set()
        self._last_participant_update = set()

        self.backfill_lock = SimpleLock("Waiting for backfilling to finish before handling %s",
                                        log=self.log)
        self._main_intent = None
        self._reaction_lock = asyncio.Lock()

    @property
    def is_direct(self) -> bool:
        return self.conv_type == ConversationType.ONE_TO_ONE

    @property
    def main_intent(self) -> IntentAPI:
        if not self._main_intent:
            raise ValueError("Portal must be postinit()ed before main_intent can be used")
        return self._main_intent

    @classmethod
    def init_cls(cls, bridge: 'TwitterBridge') -> None:
        cls.config = bridge.config
        cls.matrix = bridge.matrix
        cls.az = bridge.az
        cls.loop = bridge.loop
        cls.bridge = bridge
        NotificationDisabler.puppet_cls = p.Puppet
        NotificationDisabler.config_enabled = cls.config["bridge.backfill.disable_notifications"]

    # region Misc

    async def _send_delivery_receipt(self, event_id: EventID) -> None:
        if event_id and self.config["bridge.delivery_receipts"]:
            try:
                await self.az.intent.mark_read(self.mxid, event_id)
            except Exception:
                self.log.exception("Failed to send delivery receipt for %s", event_id)

    async def _upsert_reaction(self, existing: DBReaction, intent: IntentAPI, mxid: EventID,
                               message: DBMessage, sender: Union['u.User', 'p.Puppet'],
                               reaction: ReactionKey) -> None:
        if existing:
            self.log.debug(f"_upsert_reaction redacting {existing.mxid} and inserting {mxid}"
                           f" (message: {message.mxid})")
            await intent.redact(existing.mx_room, existing.mxid)
            await existing.edit(reaction=reaction, mxid=mxid, mx_room=message.mx_room)
        else:
            self.log.debug(f"_upsert_reaction inserting {mxid} (message: {message.mxid})")
            await DBReaction(mxid=mxid, mx_room=message.mx_room, tw_msgid=message.twid,
                             tw_receiver=self.receiver, tw_sender=sender.twid,
                             reaction=reaction).insert()

    # endregion
    # region Matrix event handling

    async def handle_matrix_message(self, sender: 'u.User', message: MessageEventContent,
                                    event_id: EventID) -> None:
        if not sender.client:
            self.log.debug(f"Ignoring message {event_id} as user is not connected")
            return
        elif ((message.get(self.bridge.real_user_content_key,
                           False) and await p.Puppet.get_by_custom_mxid(sender.mxid))):
            self.log.debug(f"Ignoring puppet-sent message by confirmed puppet user {sender.mxid}")
            return
        request_id = str(sender.client.new_request_id())
        self._reqid_dedup.add(request_id)
        text = message.body
        media_id = None
        if message.msgtype == MessageType.EMOTE:
            text = f"/me {text}"
        elif message.msgtype.is_media:
            if message.file and decrypt_attachment:
                data = await self.main_intent.download_media(message.file.url)
                data = decrypt_attachment(data, message.file.key.key,
                                          message.file.hashes.get("sha256"), message.file.iv)
            else:
                data = await self.main_intent.download_media(message.url)
            mime_type = message.info.mimetype or magic.from_buffer(data, mime=True)
            # TODO this will throw errors if the mime type is not supported
            #      those errors should be sent back to the client
            upload_resp = await sender.client.upload(data, mime_type=mime_type)
            media_id = upload_resp.media_id
            text = ""
        resp = await sender.client.conversation(self.twid).send(text, media_id=media_id,
                                                                request_id=request_id)
        resp_msg_id = int(resp.entries[0].message.id)
        self._msgid_dedup.appendleft(resp_msg_id)
        msg = DBMessage(mxid=event_id, mx_room=self.mxid, twid=resp_msg_id, receiver=self.receiver)
        await msg.insert()
        self._reqid_dedup.remove(request_id)
        await self._send_delivery_receipt(event_id)
        self.log.debug(f"Handled Matrix message {event_id} -> {resp_msg_id}")

    async def handle_matrix_reaction(self, sender: 'u.User', event_id: EventID,
                                     reacting_to: EventID, reaction_val: str) -> None:
        try:
            reaction = ReactionKey.from_emoji(reaction_val)
        except ValueError:
            self.log.debug(f"Ignoring unsupported reaction {event_id} with value {reaction_val}")
            return

        message = await DBMessage.get_by_mxid(reacting_to, self.mxid)
        if not message:
            self.log.debug(f"Ignoring reaction to unknown event {reacting_to}")
            return

        existing = await DBReaction.get_by_twid(message.twid, message.receiver, sender.twid)
        if existing and existing.reaction == reaction:
            return

        dedup_id = (message.twid, sender.twid, reaction)
        self._reaction_dedup.appendleft(dedup_id)
        async with self._reaction_lock:
            await sender.client.conversation(self.twid).react(message.twid, reaction)
            await self._upsert_reaction(existing, self.main_intent, event_id, message, sender,
                                        reaction)
            self.log.trace(f"{sender.mxid} reacted to {message.twid} with {reaction}")
        await self._send_delivery_receipt(event_id)

    async def handle_matrix_redaction(self, sender: 'u.User', event_id: EventID,
                                      redaction_event_id: EventID) -> None:
        if not self.mxid:
            return

        reaction = await DBReaction.get_by_mxid(event_id, self.mxid)
        if reaction:
            try:
                await reaction.delete()
                await sender.client.conversation(self.twid).delete_reaction(reaction.tw_msgid,
                                                                            reaction.reaction)
                await self._send_delivery_receipt(redaction_event_id)
                self.log.trace(f"Removed {reaction} after Matrix redaction")
            except Exception:
                self.log.exception("Removing reaction failed")

    async def handle_matrix_leave(self, user: 'u.User') -> None:
        if self.is_direct:
            self.log.info(f"{user.mxid} left private chat portal with {self.twid}")
            if user.twid == self.receiver:
                self.log.info(f"{user.mxid} was the recipient of this portal. "
                              "Cleaning up and deleting...")
                await self.cleanup_and_delete()
        else:
            self.log.debug(f"{user.mxid} left portal to {self.twid}")
            # TODO cleanup if empty

    # endregion
    # region Twitter event handling

    async def handle_twitter_message(self, source: 'u.User', sender: 'p.Puppet',
                                     message: MessageData, request_id: str) -> None:
        msg_id = int(message.id)
        if request_id in self._reqid_dedup:
            self.log.debug(f"Ignoring message {message.id} by {message.sender_id}"
                           " as it was sent by us (request_id in dedup queue)")
        elif msg_id in self._msgid_dedup:
            self.log.debug(f"Ignoring message {message.id} by {message.sender_id}"
                           " as it was already handled (message.id in dedup queue)")
        elif await DBMessage.get_by_twid(msg_id, self.receiver) is not None:
            self.log.debug(f"Ignoring message {message.id} by {message.sender_id}"
                           " as it was already handled (message.id found in database)")
        else:
            self._msgid_dedup.appendleft(msg_id)
            intent = sender.intent_for(self)
            event_id = None
            if message.attachment:
                content = await self._handle_twitter_attachment(source, sender, message)
                if content:
                    event_id = await self._send_message(intent, content, timestamp=message.time)
            if message.text and not message.text.isspace():
                message.text = html.unescape(message.text)
                content = TextMessageEventContent(msgtype=MessageType.TEXT, body=message.text)
                event_id = await self._send_message(intent, content, timestamp=message.time)
            if event_id:
                msg = DBMessage(mxid=event_id, mx_room=self.mxid, twid=msg_id,
                                receiver=self.receiver)
                await msg.insert()
                await self._send_delivery_receipt(event_id)
                self.log.debug(f"Handled Twitter message {msg_id} -> {event_id}")

    async def _handle_twitter_attachment(self, source: 'u.User', sender: 'p.Puppet',
                                         message: MessageData
                                         ) -> Optional[MediaMessageEventContent]:
        content = None
        intent = sender.intent_for(self)
        media = message.attachment.media
        if media:
            # TODO this doesn't actually work for gifs and videos
            #      The actual url is in media.video_info.variants[0].url
            #      Gifs are also videos
            reuploaded_info = await self._reupload_twitter_media(source, media, intent)
            content = MediaMessageEventContent(body=reuploaded_info.file_name,
                                               url=reuploaded_info.mxc,
                                               file=reuploaded_info.decryption_info,
                                               external_url=media.media_url_https)
            if message.attachment.video:
                content.msgtype = MessageType.VIDEO
                content.info = VideoInfo(mimetype=reuploaded_info.mime_type,
                                         size=reuploaded_info.size,
                                         width=media.original_info.width,
                                         height=media.original_info.height,
                                         duration=media.video_info.duration_millis // 1000)
            elif message.attachment.photo or message.attachment.animated_gif:
                content.msgtype = MessageType.IMAGE
                content.info = ImageInfo(mimetype=reuploaded_info.mime_type,
                                         size=reuploaded_info.size,
                                         width=media.original_info.width,
                                         height=media.original_info.height)
            # Remove the attachment link from message.text
            start, end = media.indices
            message.text = message.text[:start] + message.text[end:]
        elif message.attachment.tweet:
            # TODO special handling for tweets?
            pass
        return content

    async def _reupload_twitter_media(self, source: 'u.User', attachment: MessageAttachmentMedia,
                                      intent: IntentAPI) -> ReuploadedMediaInfo:
        file_name = path.basename(attachment.media_url_https)
        data, mime_type = await source.client.download_media(attachment)

        upload_mime_type = mime_type
        upload_file_name = file_name
        decryption_info = None
        if self.encrypted and encrypt_attachment:
            data, decryption_info = encrypt_attachment(data)
            upload_mime_type = "application/octet-stream"
            upload_file_name = None

        mxc = await intent.upload_media(data, mime_type=upload_mime_type,
                                        filename=upload_file_name)

        if decryption_info:
            decryption_info.url = mxc
            mxc = None

        return ReuploadedMediaInfo(mxc, decryption_info, mime_type, file_name, len(data))

    async def handle_twitter_reaction_add(self, sender: 'p.Puppet', msg_id: int,
                                          reaction: ReactionKey, time: datetime) -> None:
        async with self._reaction_lock:
            dedup_id = (msg_id, sender.twid, reaction)
            if dedup_id in self._reaction_dedup:
                return
            self._reaction_dedup.appendleft(dedup_id)

        existing = await DBReaction.get_by_twid(msg_id, self.receiver, sender.twid)
        if existing and existing.reaction == reaction:
            return

        message = await DBMessage.get_by_twid(msg_id, self.receiver)
        if not message:
            self.log.debug(f"Ignoring reaction to unknown message {msg_id}")
            return

        intent = sender.intent_for(self)
        mxid = await intent.react(message.mx_room, message.mxid, reaction.emoji, timestamp=time)
        self.log.debug(f"{sender.twid} reacted to {message.mxid} -> {mxid}")
        await self._upsert_reaction(existing, intent, mxid, message, sender, reaction)

    async def handle_twitter_reaction_remove(self, sender: 'p.Puppet', msg_id: int,
                                             key: ReactionKey) -> None:
        reaction = await DBReaction.get_by_twid(msg_id, self.receiver, sender.twid)
        if reaction and reaction.reaction == key:
            try:
                await sender.intent_for(self).redact(reaction.mx_room, reaction.mxid)
            except MForbidden:
                await self.main_intent.redact(reaction.mx_room, reaction.mxid)
            await reaction.delete()
            self.log.trace(f"Removed {reaction} after Twitter removal")

    async def handle_twitter_receipt(self, sender: 'p.Puppet', read_up_to: int) -> None:
        message = await DBMessage.get_by_twid(read_up_to, self.receiver)
        if not message:
            self.log.trace(f"Ignoring read receipt from {sender.twid} "
                           f"up to unknown message {read_up_to}")
            return

        self.log.trace(f"{sender.twid} read messages up to {read_up_to} ({message.mxid})")
        await sender.intent_for(self).mark_read(message.mx_room, message.mxid)

    # endregion
    # region Updating portal info

    async def update_info(self, conv: Conversation) -> None:
        if self.conv_type == ConversationType.ONE_TO_ONE:
            if not self.other_user:
                participant = next(pcp for pcp in conv.participants
                                   if int(pcp.user_id) != self.receiver)
                self.other_user = int(participant.user_id)
                await self.update()
            puppet = await p.Puppet.get_by_twid(self.other_user)
            if not self._main_intent:
                self._main_intent = puppet.default_mxid_intent
            changed = await self._update_name(puppet.name)
        else:
            changed = await self._update_name(conv.name)
        if changed:
            await self.update_bridge_info()
            await self.update()
        await self._update_participants(conv.participants)

    async def _update_name(self, name: str) -> bool:
        if self.name != name:
            self.name = name
            if self.mxid:
                await self.main_intent.set_room_name(self.mxid, name)
            return True
        return False

    async def _update_participants(self, participants: List[Participant]) -> None:
        if not self.mxid:
            return

        # Store the current member list to prevent unnecessary updates
        current_members = {int(participant.user_id) for participant in participants}
        if current_members == self._last_participant_update:
            self.log.trace("Not updating participants: list matches cached list")
            return
        self._last_participant_update = current_members

        # Make sure puppets who should be here are here
        for participant in participants:
            twid = int(participant.user_id)
            puppet = await p.Puppet.get_by_twid(twid)
            await puppet.intent_for(self).ensure_joined(self.mxid)
            if participant.last_read_event_id:
                await self.handle_twitter_receipt(puppet, int(participant.last_read_event_id))

        # Kick puppets who shouldn't be here
        for user_id in await self.main_intent.get_room_members(self.mxid):
            twid = p.Puppet.get_id_from_mxid(user_id)
            if twid and twid not in current_members:
                await self.main_intent.kick_user(self.mxid, p.Puppet.get_mxid_from_id(twid),
                                                 reason="User had left this Twitter chat")

    # endregion
    # region Backfilling

    async def backfill(self, source: 'u.User', is_initial: bool = False) -> None:
        if not is_initial:
            raise RuntimeError("Non-initial backfilling is not supported")
        limit = self.config["bridge.backfill.initial_limit"]
        if limit == 0:
            return
        elif limit < 0:
            limit = None
        with self.backfill_lock:
            await self._backfill(source, limit)

    async def _backfill(self, source: 'u.User', limit: int) -> None:
        self.log.debug("Backfilling history through %s", source.mxid)

        entries = await self._fetch_backfill_entries(source, limit)
        if not entries:
            self.log.debug("Didn't get any entries from server")
            return

        self.log.debug("Got %d entries from server", len(entries))

        backfill_leave = await self._invite_own_puppet_backfill(source)
        async with NotificationDisabler(self.mxid, source):
            for entry in reversed(entries):
                await self._handle_backfill_entry(source, entry)
        for intent in backfill_leave:
            self.log.trace("Leaving room with %s post-backfill", intent.mxid)
            await intent.leave_room(self.mxid)
        self.log.info("Backfilled %d messages through %s", len(entries), source.mxid)

    async def _fetch_backfill_entries(self, source: 'u.User', limit: int
                                      ) -> List[BackfillEntryTypes]:
        conv = source.client.conversation(self.twid)
        entries = []
        self.log.debug("Fetching up to %d messages through %s", limit, source.twid)
        try:
            max_id = None
            while True:
                resp = await conv.fetch(max_id=max_id, include_info=False)
                max_id = resp.min_entry_id
                for entry in resp.entries:
                    if entry:
                        entries += entry.all_types
                    if len(entries) >= limit:
                        break
                if len(entries) >= limit:
                    self.log.debug("Got more messages than limit")
                    break
                elif resp.status == TimelineStatus.AT_END:
                    self.log.debug("Got all messages in conversation")
                    break
        except Exception:
            self.log.warning("Exception while fetching messages", exc_info=True)

        return [entry for entry in entries
                if isinstance(entry, (MessageEntry, ReactionCreateEntry, ReactionDeleteEntry))]

    async def _invite_own_puppet_backfill(self, source: 'u.User') -> Set[IntentAPI]:
        backfill_leave = set()
        # TODO we should probably only invite the puppet when needed
        if self.config["bridge.backfill.invite_own_puppet"]:
            self.log.debug("Adding %s's default puppet to room for backfilling", source.mxid)
            sender = await p.Puppet.get_by_twid(source.twid)
            await self.main_intent.invite_user(self.mxid, sender.default_mxid)
            await sender.default_mxid_intent.join_room_by_id(self.mxid)
            backfill_leave.add(sender.default_mxid_intent)
        return backfill_leave

    async def _handle_backfill_entry(self, source: 'u.User', entry: BackfillEntryTypes) -> None:
        sender = await p.Puppet.get_by_twid(int(entry.sender_id))
        if isinstance(entry, MessageEntry):
            await self.handle_twitter_message(source, sender, entry.message_data, entry.request_id)
        if isinstance(entry, ReactionCreateEntry):
            await self.handle_twitter_reaction_add(sender, int(entry.message_id),
                                                   entry.reaction_key, entry.time)
        elif isinstance(entry, ReactionDeleteEntry):
            await self.handle_twitter_reaction_remove(sender, int(entry.message_id),
                                                      entry.reaction_key)

    # endregion
    # region Bridge info state event

    @property
    def bridge_info_state_key(self) -> str:
        return f"net.maunium.twitter://twitter/{self.twid}"

    @property
    def bridge_info(self) -> Dict[str, Any]:
        return {
            "bridgebot": self.az.bot_mxid,
            "creator": self.main_intent.mxid,
            "protocol": {
                "id": "twitter",
                "displayname": "Twitter DM",
                "avatar_url": self.config["appservice.bot_avatar"],
            },
            "channel": {
                "id": self.twid,
                "displayname": self.name,
            }
        }

    async def update_bridge_info(self) -> None:
        if not self.mxid:
            self.log.debug("Not updating bridge info: no Matrix room created")
            return
        try:
            self.log.debug("Updating bridge info...")
            await self.main_intent.send_state_event(self.mxid, StateBridge,
                                                    self.bridge_info, self.bridge_info_state_key)
            # TODO remove this once https://github.com/matrix-org/matrix-doc/pull/2346 is in spec
            await self.main_intent.send_state_event(self.mxid, StateHalfShotBridge,
                                                    self.bridge_info, self.bridge_info_state_key)
        except Exception:
            self.log.warning("Failed to update bridge info", exc_info=True)

    # endregion
    # region Creating Matrix rooms

    async def create_matrix_room(self, source: 'u.User', info: Conversation) -> Optional[RoomID]:
        if self.mxid:
            try:
                await self._update_matrix_room(source, info)
            except Exception:
                self.log.exception("Failed to update portal")
            return self.mxid
        async with self._create_room_lock:
            return await self._create_matrix_room(source, info)

    async def _update_matrix_room(self, source: 'u.User', info: Conversation) -> None:
        await self.main_intent.invite_user(self.mxid, source.mxid, check_cache=True)
        puppet = await p.Puppet.get_by_custom_mxid(source.mxid)
        if puppet:
            did_join = await puppet.intent.ensure_joined(self.mxid)
            if did_join and self.conv_type == ConversationType.ONE_TO_ONE:
                await source.update_direct_chats({self.main_intent.mxid: [self.mxid]})

        await self.update_info(info)

        # TODO
        # up = DBUserPortal.get(source.fbid, self.fbid, self.fb_receiver)
        # if not up:
        #     in_community = await source._community_helper.add_room(source._community_id, self.mxid)
        #     DBUserPortal(user=source.fbid, portal=self.fbid, portal_receiver=self.fb_receiver,
        #                  in_community=in_community).insert()
        # elif not up.in_community:
        #     in_community = await source._community_helper.add_room(source._community_id, self.mxid)
        #     up.edit(in_community=in_community)

    async def _create_matrix_room(self, source: 'u.User', info: Conversation) -> Optional[RoomID]:
        if self.mxid:
            await self._update_matrix_room(source, info)
            return self.mxid
        await self.update_info(info)
        self.log.debug("Creating Matrix room")
        name: Optional[str] = None
        initial_state = [{
            "type": str(StateBridge),
            "state_key": self.bridge_info_state_key,
            "content": self.bridge_info,
        }, {
            # TODO remove this once https://github.com/matrix-org/matrix-doc/pull/2346 is in spec
            "type": str(StateHalfShotBridge),
            "state_key": self.bridge_info_state_key,
            "content": self.bridge_info,
        }]
        invites = [source.mxid]
        if self.config["bridge.encryption.default"] and self.matrix.e2ee:
            self.encrypted = True
            initial_state.append({
                "type": "m.room.encryption",
                "content": {"algorithm": "m.megolm.v1.aes-sha2"},
            })
            if self.is_direct:
                invites.append(self.az.bot_mxid)
        if self.encrypted or not self.is_direct:
            name = self.name
        if self.config["appservice.community_id"]:
            initial_state.append({
                "type": "m.room.related_groups",
                "content": {"groups": [self.config["appservice.community_id"]]},
            })

        # We lock backfill lock here so any messages that come between the room being created
        # and the initial backfill finishing wouldn't be bridged before the backfill messages.
        with self.backfill_lock:
            self.mxid = await self.main_intent.create_room(name=name, is_direct=self.is_direct,
                                                           initial_state=initial_state,
                                                           invitees=invites)
            if not self.mxid:
                raise Exception("Failed to create room: no mxid returned")

            if self.encrypted and self.matrix.e2ee and self.is_direct:
                try:
                    await self.az.intent.ensure_joined(self.mxid)
                except Exception:
                    self.log.warning("Failed to add bridge bot "
                                     f"to new private chat {self.mxid}")

            await self.update()
            self.log.debug(f"Matrix room created: {self.mxid}")
            self.by_mxid[self.mxid] = self
            if not self.is_direct:
                await self._update_participants(info.participants)
            else:
                puppet = await p.Puppet.get_by_custom_mxid(source.mxid)
                if puppet:
                    try:
                        await puppet.intent.join_room_by_id(self.mxid)
                        await source.update_direct_chats({self.main_intent.mxid: [self.mxid]})
                    except MatrixError:
                        self.log.debug("Failed to join custom puppet into newly created portal",
                                       exc_info=True)

            if not info.trusted:
                msg = "This is a message request. Replying here will accept the request."
                if info.low_quality:
                    msg += " Note: Twitter has marked this as a \"low quality\" message."
                await self.main_intent.send_notice(self.mxid, msg)

            # TODO
            # in_community = await source._community_helper.add_room(source._community_id, self.mxid)
            # DBUserPortal(user=source.fbid, portal=self.fbid, portal_receiver=self.fb_receiver,
            #              in_community=in_community).upsert()

            try:
                await self.backfill(source, is_initial=True)
            except Exception:
                self.log.exception("Failed to backfill new portal")

        return self.mxid

    # endregion
    # region Database getters

    async def postinit(self) -> None:
        self.by_twid[(self.twid, self.receiver)] = self
        if self.mxid:
            self.by_mxid[self.mxid] = self
        if self.other_user and self.is_direct:
            self._main_intent = (await p.Puppet.get_by_twid(self.other_user)).default_mxid_intent
        elif not self.is_direct:
            self._main_intent = self.az.intent

    async def delete(self) -> None:
        await DBMessage.delete_all(self.mxid)
        self.by_mxid.pop(self.mxid, None)
        self.mxid = None
        self.encrypted = False
        await self.update()

    async def save(self) -> None:
        await self.update()

    @classmethod
    async def all_with_room(cls) -> AsyncGenerator['Portal', None]:
        portals = await super().all_with_room()
        portal: cls
        for index, portal in enumerate(portals):
            try:
                yield cls.by_twid[(portal.twid, portal.receiver)]
            except KeyError:
                await portal.postinit()
                yield portal

    @classmethod
    async def get_by_mxid(cls, mxid: RoomID) -> Optional['Portal']:
        try:
            return cls.by_mxid[mxid]
        except KeyError:
            pass

        portal = cast(cls, await super().get_by_mxid(mxid))
        if portal is not None:
            await portal.postinit()
            return portal

        return None

    @classmethod
    async def get_by_twid(cls, twid: str, receiver: int = 0,
                          conv_type: Optional[ConversationType] = None) -> Optional['Portal']:
        if conv_type == ConversationType.GROUP_DM and receiver != 0:
            receiver = 0
        try:
            return cls.by_twid[(twid, receiver)]
        except KeyError:
            pass

        portal = cast(cls, await super().get_by_twid(twid, receiver))
        if portal is not None:
            await portal.postinit()
            return portal

        if conv_type is not None:
            portal = cls(twid, receiver, conv_type)
            await portal.insert()
            await portal.postinit()
            return portal

        return None

    # endregion
