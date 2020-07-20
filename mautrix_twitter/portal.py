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
                    TYPE_CHECKING, cast)
from collections import deque
from datetime import datetime
import asyncio

from mautwitdm.types import ConversationType, Conversation, MessageData, Participant, ReactionKey
from mautrix.appservice import AppService, IntentAPI
from mautrix.bridge import BasePortal
from mautrix.types import (EventID, MessageEventContent, RoomID, EventType, MessageType,
                           TextMessageEventContent)
from mautrix.errors import MatrixError, MForbidden
from mautrix.util.simple_lock import SimpleLock

from .db import Portal as DBPortal, Message as DBMessage, Reaction as DBReaction
from .config import Config
from . import user as u, puppet as p, matrix as m

if TYPE_CHECKING:
    from .__main__ import TwitterBridge

StateBridge = EventType.find("m.bridge", EventType.Class.STATE)
StateHalfShotBridge = EventType.find("uk.half-shot.bridge", EventType.Class.STATE)


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
        request_id = str(sender.client.new_request_id())
        self._reqid_dedup.add(request_id)
        resp = await sender.client.conversation(self.twid).send(message.body,
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
            content = TextMessageEventContent(msgtype=MessageType.TEXT, body=message.text)
            event_id = await self._send_message(intent, content, timestamp=message.time)
            msg = DBMessage(mxid=event_id, mx_room=self.mxid, twid=msg_id, receiver=self.receiver)
            await msg.insert()
            await self._send_delivery_receipt(event_id)
            self.log.debug(f"Handled Twitter message {msg_id} -> {event_id}")

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

    async def backfill(self, user: 'u.User', is_initial: bool = False) -> None:
        pass

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
            await puppet.intent.ensure_joined(self.mxid)

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
                    except MatrixError:
                        self.log.debug("Failed to join custom puppet into newly created portal",
                                       exc_info=True)

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
        self.by_mxid.pop(self.mxid, None)
        self.by_twid.pop((self.twid, self.receiver), None)
        await super().delete()
        await DBMessage.delete_all(self.mxid)

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
