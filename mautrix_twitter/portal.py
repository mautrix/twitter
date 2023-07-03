# mautrix-twitter - A Matrix-Twitter DM puppeting bridge
# Copyright (C) 2022 Tulir Asokan
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
from __future__ import annotations

from typing import TYPE_CHECKING, Any, AsyncGenerator, Awaitable, Literal, NamedTuple, Tuple, cast
from collections import deque
from datetime import datetime, timedelta
import asyncio
import base64
import hashlib
import time

from yarl import URL
import magic

from mautrix.appservice import DOUBLE_PUPPET_SOURCE_KEY, AppService, IntentAPI
from mautrix.bridge import BasePortal, NotificationDisabler, async_getter_lock
from mautrix.errors import MatrixError, MForbidden
from mautrix.types import (
    AudioInfo,
    BatchSendEvent,
    BatchSendStateEvent,
    BeeperMessageStatusEventContent,
    ContentURI,
    EncryptedFile,
    EventID,
    EventType,
    ImageInfo,
    MediaMessageEventContent,
    Membership,
    MemberStateEventContent,
    MessageEventContent,
    MessageStatus,
    MessageStatusReason,
    MessageType,
    ReactionEventContent,
    RelatesTo,
    RelationType,
    RoomID,
    TextMessageEventContent,
    ThumbnailInfo,
    VideoInfo,
)
from mautrix.util import background_task, ffmpeg, variation_selector
from mautrix.util.message_send_checkpoint import MessageSendCheckpointStatus
from mautrix.util.simple_lock import SimpleLock
from mautwitdm.errors import UnsupportedAttachmentError
from mautwitdm.types import (
    Conversation,
    ConversationType,
    MessageData,
    MessageEntry,
    Participant,
    ReactionCreateEntry,
    VideoVariant,
)

from . import backfill as b, matrix as m, puppet as p, user as u
from .config import Config
from .db import (
    BackfillStatus as DBBackfillStatus,
    Message as DBMessage,
    Portal as DBPortal,
    Reaction as DBReaction,
)
from .formatter import twitter_to_matrix

if TYPE_CHECKING:
    from .__main__ import TwitterBridge

try:
    from mautrix.crypto.attachments import decrypt_attachment, encrypt_attachment
except ImportError:
    encrypt_attachment = decrypt_attachment = None

StateBridge = EventType.find("m.bridge", EventType.Class.STATE)
StateHalfShotBridge = EventType.find("uk.half-shot.bridge", EventType.Class.STATE)
StateMarker = EventType.find("org.matrix.msc2716.marker", EventType.Class.STATE)


class ReuploadedMediaInfo(NamedTuple):
    mxc: ContentURI | None
    decryption_info: EncryptedFile | None
    mime_type: str
    file_name: str
    size: int


class Portal(DBPortal, BasePortal):
    by_mxid: dict[RoomID, Portal] = {}
    by_twid: dict[tuple[str, int], Portal] = {}
    config: Config
    matrix: "m.MatrixHandler"
    az: AppService
    private_chat_portal_meta: Literal["default", "always", "never"]

    _main_intent: IntentAPI | None
    _create_room_lock: asyncio.Lock
    backfill_lock: SimpleLock
    _msgid_dedup: deque[int]
    _reqid_dedup: set[str]
    _reaction_dedup: deque[tuple[int, int, str]]

    _main_intent: IntentAPI
    _last_participant_update: set[int]
    _reaction_lock: asyncio.Lock

    def __init__(
        self,
        twid: str,
        receiver: int,
        conv_type: ConversationType,
        other_user: int | None = None,
        mxid: RoomID | None = None,
        name: str | None = None,
        encrypted: bool = False,
        next_batch_id: str | None = None,
    ) -> None:
        super().__init__(
            twid,
            receiver,
            conv_type,
            other_user,
            mxid,
            name,
            encrypted,
            next_batch_id,
        )
        self._create_room_lock = asyncio.Lock()
        self.log = self.log.getChild(twid)
        self._msgid_dedup = deque(maxlen=100)
        self._reaction_dedup = deque(maxlen=100)
        self._reqid_dedup = set()
        self._last_participant_update = set()

        self.backfill_lock = SimpleLock(
            "Waiting for backfilling to finish before handling %s", log=self.log
        )
        self._main_intent = None
        self._reaction_lock = asyncio.Lock()

    @property
    def is_direct(self) -> bool:
        return self.conv_type == ConversationType.ONE_TO_ONE

    @property
    def set_dm_room_metadata(self) -> bool:
        return (
            not self.is_direct
            or self.private_chat_portal_meta == "always"
            or (self.encrypted and self.private_chat_portal_meta != "never")
        )

    @property
    def main_intent(self) -> IntentAPI:
        if not self._main_intent:
            raise ValueError("Portal must be postinit()ed before main_intent can be used")
        return self._main_intent

    @classmethod
    def init_cls(cls, bridge: "TwitterBridge") -> None:
        cls.config = bridge.config
        cls.matrix = bridge.matrix
        cls.az = bridge.az
        cls.loop = bridge.loop
        BasePortal.bridge = bridge
        cls.private_chat_portal_meta = cls.config["bridge.private_chat_portal_meta"]
        NotificationDisabler.puppet_cls = p.Puppet
        NotificationDisabler.config_enabled = cls.config["bridge.backfill.disable_notifications"]

    # region Misc

    async def _send_delivery_receipt(self, event_id: EventID) -> None:
        if event_id and self.config["bridge.delivery_receipts"]:
            try:
                await self.az.intent.mark_read(self.mxid, event_id)
            except Exception:
                self.log.exception("Failed to send delivery receipt for %s", event_id)

    async def _upsert_reaction(
        self,
        existing: DBReaction | None,
        intent: IntentAPI,
        mxid: EventID,
        message: DBMessage,
        sender: u.User | p.Puppet,
        reaction: str,
        reaction_id: int | None = None,
    ) -> None:
        if existing:
            self.log.debug(
                f"_upsert_reaction redacting {existing.mxid} and inserting {mxid}"
                f" (message: {message.mxid})"
            )
            await intent.redact(existing.mx_room, existing.mxid)
            await existing.edit(
                reaction=reaction, mxid=mxid, mx_room=message.mx_room, tw_reaction_id=reaction_id
            )
        else:
            self.log.debug(f"_upsert_reaction inserting {mxid} (message: {message.mxid})")
            await DBReaction(
                mxid=mxid,
                mx_room=message.mx_room,
                tw_msgid=message.twid,
                tw_receiver=self.receiver,
                tw_sender=sender.twid,
                reaction=reaction,
                tw_reaction_id=reaction_id,
            ).insert()

    # endregion
    # region Matrix event handling

    async def _send_bridge_success(
        self,
        sender: u.User,
        event_id: EventID,
        event_type: EventType,
        msgtype: MessageType | None = None,
    ) -> None:
        sender.send_remote_checkpoint(
            status=MessageSendCheckpointStatus.SUCCESS,
            event_id=event_id,
            room_id=self.mxid,
            event_type=event_type,
            message_type=msgtype,
        )
        background_task.create(self._send_message_status(event_id, err=None))
        await self._send_delivery_receipt(event_id)

    async def _send_bridge_error(
        self,
        sender: u.User,
        err: Exception,
        event_id: EventID,
        event_type: EventType,
        message_type: MessageType | None = None,
    ) -> None:
        sender.send_remote_checkpoint(
            self._status_from_exception(err),
            event_id,
            self.mxid,
            event_type,
            message_type=message_type,
            error=err,
        )

        send_notice = not isinstance(err, NotImplementedError)
        if self.config["bridge.delivery_error_reports"] and send_notice:
            event_type_str = {
                EventType.REACTION: "reaction",
                EventType.ROOM_REDACTION: "redaction",
            }.get(event_type, "message")
            await self._send_message(
                self.main_intent,
                TextMessageEventContent(
                    msgtype=MessageType.NOTICE,
                    body=f"\u26a0 Your {event_type_str} may not have been bridged: {str(err)}",
                ),
            )
        background_task.create(self._send_message_status(event_id, err))

    async def _send_message_status(self, event_id: EventID, err: Exception | None) -> None:
        if not self.config["bridge.message_status_events"]:
            return
        intent = self.az.intent if self.encrypted else self.main_intent
        status = BeeperMessageStatusEventContent(
            network=self.bridge_info_state_key,
            relates_to=RelatesTo(
                rel_type=RelationType.REFERENCE,
                event_id=event_id,
            ),
        )
        if err:
            status.status = MessageStatus.RETRIABLE
            status.reason = MessageStatusReason.GENERIC_ERROR
            status.error = str(err)
            if isinstance(err, NotImplementedError):
                if isinstance(err, UnsupportedAttachmentError):
                    status.message = str(err)
                status.status = MessageStatus.FAIL
                status.reason = MessageStatusReason.UNSUPPORTED
        else:
            status.status = MessageStatus.SUCCESS

        await intent.send_message_event(
            room_id=self.mxid,
            event_type=EventType.BEEPER_MESSAGE_STATUS,
            content=status,
        )

    @staticmethod
    def _status_from_exception(e: Exception) -> MessageSendCheckpointStatus:
        if isinstance(e, NotImplementedError):
            return MessageSendCheckpointStatus.UNSUPPORTED
        return MessageSendCheckpointStatus.PERM_FAILURE

    async def handle_matrix_message(
        self, sender: u.User, message: MessageEventContent, event_id: EventID
    ) -> None:
        try:
            await self._handle_matrix_message(sender, message, event_id)
        except Exception as e:
            self.log.exception(f"Error handling Matrix event {event_id}")
            await self._send_bridge_error(
                sender, e, event_id, EventType.ROOM_MESSAGE, message.msgtype
            )
        else:
            await self._send_bridge_success(
                sender, event_id, EventType.ROOM_MESSAGE, message.msgtype
            )

    async def _handle_matrix_message(
        self, sender: u.User, message: MessageEventContent, event_id: EventID
    ) -> None:
        if not sender.client:
            raise NotImplementedError("user is not connected")
        request_id = str(sender.client.new_request_id())
        self._reqid_dedup.add(request_id)
        media_id = None
        if message.msgtype == MessageType.TEXT:
            text = message.body
        elif message.msgtype == MessageType.EMOTE:
            text = f"/me {message.body}"
        elif message.msgtype.is_media:
            if message.msgtype != MessageType.IMAGE and message.msgtype != MessageType.VIDEO:
                raise UnsupportedAttachmentError(
                    "Non-image/video files are not supported by Twitter"
                )
            if message.file and decrypt_attachment:
                data = await self.main_intent.download_media(message.file.url)
                data = decrypt_attachment(
                    data,
                    message.file.key.key,
                    message.file.hashes.get("sha256"),
                    message.file.iv,
                )
            else:
                data = await self.main_intent.download_media(message.url)
            mime_type = message.info.mimetype or magic.from_buffer(data, mime=True)
            upload_resp = await sender.client.upload(data, mime_type=mime_type)
            media_id = upload_resp.media_id
            filename = message.get("filename", None)
            if filename and filename != message.body:
                text = message.body
            else:
                text = ""
        else:
            raise NotImplementedError(f"unsupported msgtype '{message.msgtype.value}'")
        reply_to = None
        if message.get_reply_to():
            reply_to_msg = await DBMessage.get_by_mxid(message.get_reply_to(), self.mxid)
            if reply_to_msg:
                reply_to = reply_to_msg.twid
        resp = await sender.client.conversation(self.twid).send(
            text, media_id=media_id, request_id=request_id, reply_to_id=reply_to
        )
        resp_msg_id = int(resp.entries[0].message.id)
        self._msgid_dedup.appendleft(resp_msg_id)
        msg = DBMessage(mxid=event_id, mx_room=self.mxid, twid=resp_msg_id, receiver=self.receiver)
        await msg.insert()
        self._reqid_dedup.remove(request_id)
        self.log.debug(f"Handled Matrix message {event_id} -> {resp_msg_id}")

    async def handle_matrix_reaction(
        self, sender: u.User, event_id: EventID, reacting_to: EventID, reaction: str
    ) -> None:
        try:
            await self._handle_matrix_reaction(sender, event_id, reacting_to, reaction)
        except Exception as e:
            self.log.exception(f"Failed to react to {event_id}")
            await self._send_bridge_error(sender, e, event_id, EventType.REACTION)
        else:
            await self._send_bridge_success(sender, event_id, EventType.REACTION)

    async def _handle_matrix_reaction(
        self, sender: u.User, event_id: EventID, reacting_to: EventID, reaction: str
    ) -> None:
        reaction = variation_selector.remove(reaction)
        message = await DBMessage.get_by_mxid(reacting_to, self.mxid)
        if not message:
            raise NotImplementedError(f"Unknown reaction target event")

        async with self._reaction_lock:
            dedup_id = (message.twid, sender.twid, reaction)
            self._reaction_dedup.appendleft(dedup_id)

            existing = await DBReaction.get_by_message_twid(
                message.twid, message.receiver, sender.twid
            )
            if existing and existing.reaction == reaction:
                return

            await sender.client.conversation(self.twid).react(message.twid, reaction)
            await self._upsert_reaction(
                existing, self.main_intent, event_id, message, sender, reaction
            )
            self.log.debug(f"{sender.mxid} reacted to {message.twid} with {reaction}")

    async def handle_matrix_redaction(
        self, sender: u.User, event_id: EventID, redaction_event_id: EventID
    ) -> None:
        try:
            await self._handle_matrix_redaction(sender, event_id)
        except Exception as e:
            self.log.exception(f"Error handling redaction of {event_id}")
            await self._send_bridge_error(sender, e, redaction_event_id, EventType.ROOM_REDACTION)
        else:
            await self._send_bridge_success(sender, redaction_event_id, EventType.ROOM_REDACTION)

    async def _handle_matrix_redaction(self, sender: u.User, event_id: EventID) -> None:
        async with self._reaction_lock:
            reaction = await DBReaction.get_by_mxid(event_id, self.mxid)
            if reaction:
                await reaction.delete()
                await sender.client.conversation(self.twid).delete_reaction(
                    reaction.tw_msgid, reaction.reaction
                )
                self.log.trace(f"Removed {reaction} after Matrix redaction")
                return

        raise NotImplementedError("Message redactions are not supported")

    async def handle_matrix_leave(self, user: u.User) -> None:
        if self.is_direct:
            self.log.info(f"{user.mxid} left private chat portal with {self.twid}")
            if user.twid == self.receiver:
                self.log.info(
                    f"{user.mxid} was the recipient of this portal. Cleaning up and deleting..."
                )
                await self.cleanup_and_delete()
        else:
            self.log.debug(f"{user.mxid} left portal to {self.twid}")
            # TODO cleanup if empty

    # endregion
    # region Twitter event handling

    async def handle_twitter_message(
        self, source: u.User, sender: p.Puppet, message: MessageData, request_id: str
    ) -> None:
        if await self._twitter_message_dedupe(request_id, int(message.id), message.sender_id):
            await self._handle_deduplicated_twitter_message(
                source, sender, message, int(message.id)
            )

    async def _twitter_message_dedupe(self, request_id: str, msg_id: int, sender_id: str) -> bool:
        if request_id in self._reqid_dedup:
            self.log.debug(
                f"Ignoring message {msg_id} by {sender_id}"
                " as it was sent by us (request_id in dedup queue)"
            )
            return False
        if msg_id in self._msgid_dedup:
            self.log.debug(
                f"Ignoring message {msg_id} by {sender_id}"
                " as it was already handled (message.id in dedup queue)"
            )
            return False
        self._msgid_dedup.appendleft(msg_id)

        if await DBMessage.get_by_twid(msg_id, self.receiver) is not None:
            self.log.debug(
                f"Ignoring message {msg_id} by {sender_id}"
                " as it was already handled (message.id found in database)"
            )
            return False

        return True

    def deterministic_event_id(self, msg_id: str, part: str) -> EventID:
        hash_content = f"{self.mxid}/twitter/{msg_id}/{part}"
        hashed = hashlib.sha256(hash_content.encode("utf-8")).digest()
        b64hash = base64.urlsafe_b64encode(hashed).decode("utf-8").rstrip("=")
        return EventID(f"${b64hash}:twitter.com")

    async def _convert_twitter_message(
        self, source: u.User, sender: p.Puppet, message: MessageData
    ) -> list[MessageEventContent]:
        converted = []
        intent = sender.intent_for(self)
        if message.reply_data:
            reply_to_msg = await DBMessage.get_by_twid(int(message.reply_data.id), self.receiver)
        else:
            reply_to_msg = None
        media_content = None
        if message.attachment and message.attachment.media:
            media_content = await self._handle_twitter_media(source, intent, message)
            if media_content:
                if reply_to_msg:
                    media_content.set_reply(reply_to_msg.mxid)
                converted.append(media_content)
        if message.text and not message.text.isspace():
            text_content = await twitter_to_matrix(message)
            text_content["com.beeper.linkpreviews"] = await self._twitter_preview_to_beeper(
                source, intent, message
            )
            if reply_to_msg:
                text_content.set_reply(reply_to_msg.mxid)
            if media_content and self.config["bridge.caption_in_message"]:
                media_content["filename"] = media_content.body
                media_content.body = text_content.body
                if text_content.formatted_body:
                    media_content["format"] = str(text_content.format)
                    media_content["formatted_body"] = text_content.formatted_body
            else:
                converted.append(text_content)
        return converted

    async def _handle_deduplicated_twitter_message(
        self, source: u.User, sender: p.Puppet, message: MessageData, msg_id: int
    ) -> None:
        intent = sender.intent_for(self)
        event_id = None
        converted = await self._convert_twitter_message(source, sender, message)
        for content in converted:
            event_id = await self._send_message(intent, content, timestamp=message.time)
        if event_id:
            msg = DBMessage(
                mxid=event_id,
                mx_room=self.mxid,
                twid=msg_id,
                receiver=self.receiver,
            )
            await msg.insert()
            await self._send_delivery_receipt(event_id)
            self.log.debug(f"Handled Twitter message {msg_id} -> {event_id}")

    @staticmethod
    def _is_better_mime(best: VideoVariant, current: VideoVariant) -> bool:
        order = ["video/mp4"]
        try:
            best_quality = order.index(best.content_type)
        except (IndexError, ValueError):
            best_quality = -1
        try:
            current_quality = order.index(current.content_type)
        except (IndexError, ValueError):
            current_quality = -1
        return current_quality > best_quality

    async def _twitter_preview_to_beeper(
        self, source: u.User, intent: IntentAPI, message: MessageData
    ) -> list[dict[str, Any]]:
        if not message.attachment or not message.attachment.url_preview:
            return []
        if message.attachment.card:
            preview = await self._twitter_card_to_beeper(source, intent, message)
        elif message.attachment.tweet:
            preview = await self._twitter_tweet_to_beeper(source, intent, message)
        else:
            return []
        return [{k: v for k, v in preview.items() if v is not None}]

    async def _twitter_card_to_beeper(
        self, source: u.User, intent: IntentAPI, message: MessageData
    ) -> dict[str, Any]:
        card = message.attachment.card
        card_url = card.string("card_url", "")
        if card_url.startswith("https://t.co") and message.entities:
            try:
                # Try to find the actual expanded URL from the message entities
                card_url = next(e for e in message.entities.urls if e.url == card_url).expanded_url
            except StopIteration:
                pass
        preview = {
            "matched_url": card_url,
            "og:title": card.string("title"),
            "og:description": card.string("description"),
        }
        image = (
            card.image("photo_image_full_size_original")
            or card.image("summary_photo_image_original")
            or card.image("thumbnail_image_original")
        )
        if image:
            preview.update(
                {
                    **await self._twitter_image_to_beeper(source, intent, image.url),
                    "og:image:width": image.width,
                    "og:image:height": image.height,
                }
            )
        return preview

    async def _twitter_tweet_to_beeper(
        self, source: u.User, intent: IntentAPI, message: MessageData
    ) -> dict[str, Any]:
        tweet = message.attachment.tweet.status
        preview = {
            "matched_url": message.attachment.tweet.expanded_url,
            "og:url": message.attachment.tweet.expanded_url,
            "og:title": f"{tweet.user.name} on Twitter",
            "og:description": tweet.full_text,
        }
        if tweet.extended_entities and tweet.extended_entities.media:
            media = tweet.extended_entities.media[0]
            if "medium" in media.sizes:
                size_name = "medium"
            else:
                size_name = next(media.sizes.keys().__iter__())
            url = str(URL(media.media_url_https).update_query({"name": size_name}))
            preview.update(
                {
                    **await self._twitter_image_to_beeper(source, intent, url),
                    "og:image:width": media.sizes[size_name].w,
                    "og:image:height": media.sizes[size_name].h,
                }
            )

        return preview

    async def _twitter_image_to_beeper(
        self, source: u.User, intent: IntentAPI, url: str
    ) -> dict[str, Any]:
        info = await self._reupload_twitter_media(source, url, intent)
        output = {
            "og:image:type": info.mime_type,
            "matrix:image:size": info.size,
        }
        if info.decryption_info:
            output["beeper:image:encryption"] = info.decryption_info.serialize()
        else:
            output["og:image"] = info.mxc
        return output

    async def _handle_twitter_media(
        self, source: u.User, intent: IntentAPI, message: MessageData
    ) -> MediaMessageEventContent | None:
        media = message.attachment.media
        reuploaded_info = await self._reupload_twitter_media(source, media.media_url_https, intent)
        thumbnail_info = None
        if media.video_info:
            thumbnail_info = reuploaded_info
            best_variant = None
            for variant in media.video_info.variants:
                if (
                    not best_variant
                    or (variant.bitrate or 0) > (best_variant.bitrate or 0)
                    or self._is_better_mime(best_variant, variant)
                ):
                    best_variant = variant
            reuploaded_info = await self._reupload_twitter_media(
                source, best_variant.url, intent, convert_to_audio=media.audio_only
            )
        content = MediaMessageEventContent(
            body=reuploaded_info.file_name,
            url=reuploaded_info.mxc,
            file=reuploaded_info.decryption_info,
        )
        if message.attachment.video and message.attachment.video.audio_only:
            content.msgtype = MessageType.AUDIO
            content.info = AudioInfo(
                mimetype=reuploaded_info.mime_type,
                size=reuploaded_info.size,
                duration=media.video_info.duration_millis or None,
            )
            content["org.matrix.msc1767.audio"] = (
                {
                    "duration": content.info.duration,
                }
                if content.info.duration
                else {}
            )
            content["org.matrix.msc3245.voice"] = {}
        elif message.attachment.video or (
            message.attachment.animated_gif and reuploaded_info.mime_type.startswith("video/")
        ):
            content.msgtype = MessageType.VIDEO
            content.info = VideoInfo(
                mimetype=reuploaded_info.mime_type,
                size=reuploaded_info.size,
                width=media.original_info.width,
                height=media.original_info.height,
                duration=media.video_info.duration_millis or None,
            )
            if message.attachment.animated_gif:
                content.info["fi.mau.loop"] = True
                content.info["fi.mau.autoplay"] = True
                content.info["fi.mau.hide_controls"] = True
                content.info["fi.mau.no_audio"] = True
        elif message.attachment.photo or message.attachment.animated_gif:
            content.msgtype = MessageType.IMAGE
            content.info = ImageInfo(
                mimetype=reuploaded_info.mime_type,
                size=reuploaded_info.size,
                width=media.original_info.width,
                height=media.original_info.height,
            )
        if thumbnail_info:
            content.info.thumbnail_url = thumbnail_info.mxc
            content.info.thumbnail_file = thumbnail_info.decryption_info
            content.info.thumbnail_info = ThumbnailInfo(
                mimetype=thumbnail_info.mime_type,
                size=thumbnail_info.size,
                width=media.original_info.width,
                height=media.original_info.height,
            )
        # Remove the attachment link from message.text
        start, end = media.indices
        message.text = message.text[:start] + message.text[end:]
        if message.entities and message.entities.urls:
            message.entities.urls = [u for u in message.entities.urls if u.url != media.url]
        return content

    async def _reupload_twitter_media(
        self, source: u.User, url: str, intent: IntentAPI, convert_to_audio: bool = False
    ) -> ReuploadedMediaInfo:
        file_name = URL(url).name
        data, mime_type = await source.client.download_media(url)

        if convert_to_audio and (mime_type.startswith("video/") or mime_type.startswith("audio/")):
            data = await ffmpeg.convert_bytes(
                data,
                ".ogg",
                output_args=("-c:a", "libopus"),
                input_mime=mime_type,
            )
            mime_type = "audio/ogg"
            file_name += ".ogg"

        upload_mime_type = mime_type
        upload_file_name = file_name
        decryption_info = None
        if self.encrypted and encrypt_attachment:
            data, decryption_info = encrypt_attachment(data)
            upload_mime_type = "application/octet-stream"
            upload_file_name = None

        mxc = await intent.upload_media(
            data,
            mime_type=upload_mime_type,
            filename=upload_file_name,
            async_upload=self.config["homeserver.async_media"],
        )

        if decryption_info:
            decryption_info.url = mxc
            mxc = None

        return ReuploadedMediaInfo(mxc, decryption_info, mime_type, file_name, len(data))

    async def handle_twitter_reaction_add(
        self,
        sender: p.Puppet,
        msg_id: int,
        emoji: str,
        time: datetime,
        reaction_id: int,
    ) -> None:
        async with self._reaction_lock:
            # TODO update the database with the reaction_id of outgoing reactions
            dedup_id = (msg_id, sender.twid, emoji)
            if dedup_id in self._reaction_dedup:
                self.log.debug(
                    f"Ignoring duplicate reaction from {sender.twid} to {msg_id} (dedup queue)"
                )
                return
            self._reaction_dedup.appendleft(dedup_id)

        existing = await DBReaction.get_by_message_twid(msg_id, self.receiver, sender.twid)
        if existing and existing.reaction == emoji:
            if not existing.tw_reaction_id:
                await existing.update_id(reaction_id)
            self.log.debug(
                f"Ignoring duplicate reaction from {sender.twid} to {msg_id} (database)"
            )
            return

        message = await DBMessage.get_by_twid(msg_id, self.receiver)
        if not message:
            self.log.debug(f"Ignoring reaction to unknown message {msg_id}")
            return

        intent = sender.intent_for(self)
        mxid = await intent.react(
            message.mx_room, message.mxid, variation_selector.add(emoji), timestamp=time
        )
        self.log.debug(f"{sender.twid} reacted to {msg_id}/{message.mxid} -> {mxid}")
        await self._upsert_reaction(existing, intent, mxid, message, sender, emoji, reaction_id)

    async def handle_twitter_reaction_remove(
        self, sender: p.Puppet, msg_id: int, emoji: str
    ) -> None:
        reaction = await DBReaction.get_by_message_twid(msg_id, self.receiver, sender.twid)
        if reaction and (reaction.reaction == emoji or not emoji):
            try:
                self._reaction_dedup.remove((msg_id, sender.twid, emoji))
            except ValueError:
                pass
            try:
                await sender.intent_for(self).redact(reaction.mx_room, reaction.mxid)
            except MForbidden:
                await self.main_intent.redact(reaction.mx_room, reaction.mxid)
            await reaction.delete()
            self.log.debug(f"Removed {reaction} after Twitter removal")

    async def handle_twitter_receipt(
        self, sender: p.Puppet, read_up_to: int, historical: bool = False
    ) -> None:
        message = await DBMessage.get_by_twid(read_up_to, self.receiver)
        if not message:
            message = await DBReaction.get_by_reaction_twid(read_up_to, self.receiver)
            if not message:
                self.log.debug(
                    f"Ignoring read receipt from {sender.twid} "
                    f"up to unknown message {read_up_to} ({historical=})"
                )
                return

        self.log.debug(
            f"{sender.twid} read messages up to {read_up_to} ({message.mxid}, {historical=})"
        )
        await sender.intent_for(self).mark_read(message.mx_room, message.mxid)

    # endregion
    # region Updating portal info

    async def update_info(self, conv: Conversation) -> None:
        if self.conv_type == ConversationType.ONE_TO_ONE:
            if not self.other_user:
                if len(conv.participants) == 1:
                    self.other_user = int(conv.participants[0].user_id)
                else:
                    participant = next(
                        pcp for pcp in conv.participants if int(pcp.user_id) != self.receiver
                    )
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

    async def update_name(self, name: str) -> None:
        changed = await self._update_name(name)

        if changed:
            await self.update_bridge_info()
            await self.update()

    async def _update_name(self, name: str) -> bool:
        if self.name != name:
            self.name = name
            if self.mxid and self.set_dm_room_metadata:
                await self.main_intent.set_room_name(self.mxid, name)
            return True
        return False

    async def _update_participants(self, participants: list[Participant]) -> None:
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
            await puppet.intent_for(self).ensure_joined(self.mxid, bot=self.main_intent)
            if participant.last_read_event_id:
                await self.handle_twitter_receipt(
                    puppet, int(participant.last_read_event_id), historical=True
                )

        # Kick puppets who shouldn't be here
        for user_id in await self.main_intent.get_room_members(self.mxid):
            twid = p.Puppet.get_id_from_mxid(user_id)
            if twid and twid not in current_members:
                await self.main_intent.kick_user(
                    self.mxid,
                    p.Puppet.get_mxid_from_id(twid),
                    reason="User had left this Twitter chat",
                )

    # endregion
    # region Backfilling

    async def backfill(self, source: u.User, is_initial: bool = False) -> int:
        limit = self.config["bridge.backfill.initial_limit"]
        if limit == 0:
            return 0
        elif limit < 0:
            limit = None
        with self.backfill_lock:
            return await self._backfill(source, limit, is_initial)

    async def _backfill(self, source: u.User, limit: int, is_initial: bool = False) -> int:
        if is_initial:
            self.log.debug("Backfilling initial batch through %s", source.mxid)
        elif not self.config["bridge.backfill.backwards"]:
            self.log.debug("Not backfilling history, disabled in config")
            return 0
        else:
            self.log.debug("Backfilling history through %s", source.mxid)

        first_message = await DBMessage.get_first(self.mxid)
        max_id = None
        if not is_initial:
            if first_message is not None:
                max_id = first_message.twid
            else:
                self.log.warning("Can't backfill without a first bridged message")
                raise b.NoFirstMessageException()

        mark_read, entries = await self._fetch_backfill_entries(source, limit, max_id)
        if not entries:
            self.log.debug("Didn't get any entries from server")
            return 0

        self.log.debug("Got %d entries from server", len(entries))

        filled = 0
        if self.config["bridge.backfill.backwards"]:
            filled = await self._batch_handle_backfill(
                source, reversed(entries), is_initial, mark_read
            )
        else:
            backfill_leave = await self._invite_own_puppet_backfill(source)
            async with NotificationDisabler(self.mxid, source):
                for entry in reversed(entries):
                    await self._handle_backfill_entry(source, entry)
                    filled += 1
            for intent in backfill_leave:
                self.log.trace("Leaving room with %s post-backfill", intent.mxid)
                await intent.leave_room(self.mxid)
        self.log.info(
            "Backfilled %d messages (%d events) through %s", len(entries), filled, source.mxid
        )
        return filled

    async def _fetch_backfill_entries(
        self, source: u.User, limit: int, max_id: int | None = None
    ) -> tuple[bool, list[MessageEntry | ReactionCreateEntry]]:
        conv = source.client.conversation(self.twid)
        entries: list[MessageEntry | ReactionCreateEntry] = []
        message_count = 0
        self.log.debug("Fetching up to %d messages through %s", limit, source.twid)
        try:
            mark_read = False
            self.log.debug("Fetching with max_id %s", max_id)
            resp = await conv.fetch(max_id=max_id)
            resp_entries = (
                list(filter(lambda x: x is not None, resp.entries))
                if resp.entries is not None
                else None
            )
            try:
                if datetime.now() - timedelta(days=30) > datetime.fromtimestamp(
                    resp.conversations[self.twid].sort_timestamp
                ) or (
                    resp_entries is not None
                    and len(resp_entries) != 0
                    and int(resp.conversations[self.twid].last_read_event_id)
                    >= int(resp_entries[0].message.id)
                ):
                    mark_read = True
            except:
                mark_read = True

            while True:
                if resp_entries is None or len(resp_entries) == 0:
                    break
                for entry in resp_entries:
                    if entry and entry.message:
                        entries.append(entry.message)
                        message_count += 1
                        if entry.message.message_reactions:
                            entries += entry.message.message_reactions
                    if message_count >= limit:
                        break
                if message_count >= limit:
                    self.log.debug("Got more messages than limit")
                    break

                if resp.min_entry_id is None or resp.min_entry_id == max_id:
                    break
                max_id = resp.min_entry_id
                self.log.debug("Fetching more entries with max_id %s", max_id)
                resp = await conv.fetch(max_id=max_id)

            if len(entries) == 0:
                return None, None
            entries.sort(key=lambda ent: (ent.time, ent.id), reverse=True)
            self.log.debug("Finished fetching entries")
            return mark_read, entries
        except Exception:
            self.log.warning("Exception while fetching messages", exc_info=True)
            return None, None

    async def _invite_own_puppet_backfill(self, source: u.User) -> set[IntentAPI]:
        backfill_leave = set()
        # TODO we should probably only invite the puppet when needed
        if self.config["bridge.backfill.invite_own_puppet"]:
            self.log.debug("Adding %s's default puppet to room for backfilling", source.mxid)
            sender = await p.Puppet.get_by_twid(source.twid)
            await self.main_intent.invite_user(self.mxid, sender.default_mxid)
            await sender.default_mxid_intent.join_room_by_id(self.mxid)
            backfill_leave.add(sender.default_mxid_intent)
        return backfill_leave

    async def _handle_backfill_entry(
        self, source: u.User, entry: MessageEntry | ReactionCreateEntry
    ) -> None:
        sender = await p.Puppet.get_by_twid(int(entry.sender_id))
        if isinstance(entry, MessageEntry):
            await self.handle_twitter_message(source, sender, entry.message_data, entry.request_id)
        if isinstance(entry, ReactionCreateEntry):
            await self.handle_twitter_reaction_add(
                sender, int(entry.message_id), entry.reaction_emoji, entry.time, int(entry.id)
            )

    async def _batch_handle_backfill(
        self,
        source: u.User,
        entries: list[MessageEntry | ReactionCreateEntry],
        is_forward: bool,
        mark_read: bool,
    ) -> int:
        events = []
        twids = []
        users_in_batch = set()
        self.log.debug("Converting Twitter messages in batch")
        for i, entry in enumerate(entries):
            sender = await p.Puppet.get_by_twid(int(entry.sender_id))
            users_in_batch.add(sender.mxid)
            if isinstance(entry, MessageEntry):
                msg_id = int(entry.message_data.id)
                if await self._twitter_message_dedupe(
                    entry.request_id, msg_id, entry.message_data.sender_id
                ):
                    converted = await self._convert_twitter_message(
                        source, sender, entry.message_data
                    )
                    for i, content in enumerate(converted):
                        event_type = EventType.ROOM_MESSAGE
                        if self.encrypted and self.matrix.e2ee:
                            event_type, content = await self.matrix.e2ee.encrypt(
                                self.mxid, event_type, content
                            )
                        content[DOUBLE_PUPPET_SOURCE_KEY] = self.bridge.name
                        e = BatchSendEvent(
                            type=event_type,
                            content=content,
                            sender=sender.mxid,
                            timestamp=int(entry.message_data.time.timestamp() * 1000),
                        )
                        if self.bridge.homeserver_software.is_hungry:
                            e.event_id = self.deterministic_event_id(
                                entry.id, "main" if i + 1 == len(converted) else str(i)
                            )
                        events.append(e)
                        twids.append((msg_id, "message"))
            elif (
                isinstance(entry, ReactionCreateEntry)
                and self.bridge.homeserver_software.is_hungry
            ):
                content = ReactionEventContent(
                    relates_to=RelatesTo(
                        rel_type=RelationType.ANNOTATION,
                        event_id=self.deterministic_event_id(entry.message_id, "main"),
                        key=entry.reaction_emoji,
                    )
                )
                content[DOUBLE_PUPPET_SOURCE_KEY] = self.bridge.name
                e = BatchSendEvent(
                    type=EventType.REACTION,
                    content=content,
                    sender=sender.mxid,
                    event_id=self.deterministic_event_id(entry.id, f"reaction{i}"),
                    timestamp=int(entry.time.timestamp() * 1000),
                )
                events.append(e)
                twids.append(
                    (entry.id, "reaction", entry.message_id, entry.sender_id, entry.reaction_emoji)
                )
        if len(events) == 0:
            self.log.warning("No bridgeable messages in backfill batch")
            return 0

        intent = self.main_intent

        state_events_at_start = []
        if not self.bridge.homeserver_software.is_hungry:
            before_first_message_timestamp = events[0].timestamp - 1
            self.log.debug("Adding member state events to batch")
            for mxid in users_in_batch:
                puppet = await p.Puppet.get_by_mxid(mxid)
                if puppet is None:
                    self.log.warning(f"No puppet found for user {mxid} while backfilling")
                    continue
                state_events_at_start.append(
                    BatchSendStateEvent(
                        type=EventType.ROOM_MEMBER,
                        content=MemberStateEventContent(
                            Membership.INVITE, avatar_url=puppet.photo_mxc, displayname=puppet.name
                        ),
                        sender=intent.mxid,
                        timestamp=before_first_message_timestamp,
                        state_key=mxid,
                    )
                )
                state_events_at_start.append(
                    BatchSendStateEvent(
                        type=EventType.ROOM_MEMBER,
                        content=MemberStateEventContent(
                            Membership.JOIN, avatar_url=puppet.photo_mxc, displayname=puppet.name
                        ),
                        sender=mxid,
                        timestamp=before_first_message_timestamp,
                        state_key=mxid,
                    )
                )

        first_event = None
        if not is_forward:
            first_event = await DBMessage.get_first(self.mxid)
            if first_event is None:
                raise b.NoFirstMessageException
            else:
                first_event = first_event.mxid
        self.log.debug("Sending batch send request")
        resp = await intent.batch_send(
            self.mxid,
            first_event,
            batch_id=None if is_forward else self.next_batch_id,
            events=events,
            state_events_at_start=state_events_at_start,
            beeper_new_messages=is_forward,
            beeper_mark_read_by=source.mxid if mark_read else None,
        )
        if resp.base_insertion_event_id is not None:
            self.log.debug("Sending msc2716 insertion marker event")
            await self.main_intent.send_state_event(
                self.mxid,
                StateMarker,
                {
                    "org.matrix.msc2716.marker.insertion": resp.base_insertion_event_id,
                    "com.beeper.timestamp": int(time.time() * 1000),
                },
                state_key=resp.base_insertion_event_id,
            )

        for i, event_id in enumerate(resp.event_ids):
            if twids[i][1] == "message":
                msg = DBMessage(event_id, self.mxid, twids[i][0], self.receiver)
                await msg.upsert()
            elif twids[i][1] == "reaction" and self.bridge.homeserver_software.is_hungry:
                reaction = DBReaction(
                    mxid=event_id,
                    mx_room=self.mxid,
                    tw_msgid=twids[i][2],
                    tw_receiver=self.receiver,
                    tw_sender=twids[i][3],
                    reaction=twids[i][4],
                    tw_reaction_id=twids[i][0],
                )
                await reaction.insert()
        self.next_batch_id = resp.next_batch_id
        await self.update()

        return len(events)

    # endregion
    # region Bridge info state event

    @property
    def bridge_info_state_key(self) -> str:
        return f"net.maunium.twitter://twitter/{self.twid}"

    @property
    def bridge_info(self) -> dict[str, Any]:
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
            },
        }

    async def update_bridge_info(self) -> None:
        if not self.mxid:
            self.log.debug("Not updating bridge info: no Matrix room created")
            return
        try:
            self.log.debug("Updating bridge info...")
            await self.main_intent.send_state_event(
                self.mxid, StateBridge, self.bridge_info, self.bridge_info_state_key
            )
            # TODO remove this once https://github.com/matrix-org/matrix-doc/pull/2346 is in spec
            await self.main_intent.send_state_event(
                self.mxid,
                StateHalfShotBridge,
                self.bridge_info,
                self.bridge_info_state_key,
            )
        except Exception:
            self.log.warning("Failed to update bridge info", exc_info=True)

    # endregion
    # region Creating Matrix rooms

    async def create_matrix_room(self, source: u.User, info: Conversation) -> RoomID | None:
        if self.mxid:
            try:
                await self._update_matrix_room(source, info)
            except Exception:
                self.log.exception("Failed to update portal")
            return self.mxid
        async with self._create_room_lock:
            return await self._create_matrix_room(source, info)

    def _get_invite_content(self, double_puppet: p.Puppet | None) -> dict[str, Any]:
        invite_content = {}
        if double_puppet:
            invite_content["fi.mau.will_auto_accept"] = True
        if self.is_direct:
            invite_content["is_direct"] = True
        return invite_content

    async def _add_user(self, user: u.User) -> None:
        puppet = await p.Puppet.get_by_custom_mxid(user.mxid)
        await self.main_intent.invite_user(
            self.mxid,
            user.mxid,
            check_cache=True,
            extra_content=self._get_invite_content(puppet),
        )
        if puppet:
            did_join = await puppet.intent.ensure_joined(self.mxid)
            if did_join and self.is_direct:
                await user.update_direct_chats({self.main_intent.mxid: [self.mxid]})

    async def _update_matrix_room(self, source: u.User, info: Conversation) -> None:
        await self._add_user(source)
        await self.update_info(info)

    async def _create_matrix_room(self, source: u.User, info: Conversation) -> RoomID | None:
        if self.mxid:
            await self._update_matrix_room(source, info)
            return self.mxid
        await self.update_info(info)
        self.log.debug("Creating Matrix room")
        initial_state = [
            {
                "type": str(StateBridge),
                "state_key": self.bridge_info_state_key,
                "content": self.bridge_info,
            },
            {
                # TODO remove this once https://github.com/matrix-org/matrix-doc/pull/2346 is in spec
                "type": str(StateHalfShotBridge),
                "state_key": self.bridge_info_state_key,
                "content": self.bridge_info,
            },
        ]
        invites = []
        if self.config["bridge.encryption.default"] and self.matrix.e2ee:
            self.encrypted = True
            initial_state.append(
                {
                    "type": "m.room.encryption",
                    "content": self.get_encryption_state_event_json(),
                }
            )
            if self.is_direct:
                invites.append(self.az.bot_mxid)

        # We lock backfill lock here so any messages that come between the room being created
        # and the initial backfill finishing wouldn't be bridged before the backfill messages.
        with self.backfill_lock:
            creation_content = {}
            if not self.config["bridge.federate_rooms"]:
                creation_content["m.federate"] = False
            self.mxid = await self.main_intent.create_room(
                name=self.name if self.set_dm_room_metadata else None,
                is_direct=self.is_direct,
                initial_state=initial_state,
                invitees=invites,
                creation_content=creation_content,
            )
            if not self.mxid:
                raise Exception("Failed to create room: no mxid returned")

            if self.encrypted and self.matrix.e2ee and self.is_direct:
                try:
                    await self.az.intent.ensure_joined(self.mxid)
                except Exception:
                    self.log.warning(f"Failed to add bridge bot to new private chat {self.mxid}")

            await self.update()
            self.log.debug(f"Matrix room created: {self.mxid}")
            self.by_mxid[self.mxid] = self
            await self._add_user(source)
            await self._update_participants(info.participants)

            puppet = await p.Puppet.get_by_custom_mxid(source.mxid)
            if puppet:
                try:
                    await puppet.intent.join_room_by_id(self.mxid)
                    if self.is_direct:
                        await source.update_direct_chats({self.main_intent.mxid: [self.mxid]})
                    if self.config["bridge.low_quality_tag"]:
                        await source.tag_room(
                            puppet,
                            self,
                            self.config["bridge.low_quality_tag"],
                            info.low_quality == True,
                        )
                    if self.config["bridge.low_quality_mute"]:
                        await source.set_muted(puppet, self, info.low_quality == True)
                except MatrixError:
                    self.log.debug(
                        "Failed to join custom puppet into newly created portal", exc_info=True
                    )

            if not info.trusted:
                msg = "This is a message request. Replying here will accept the request."
                if info.low_quality:
                    msg += ' Note: Twitter has marked this as a "low quality" message.'
                await self.main_intent.send_notice(self.mxid, msg)

            try:
                await self._enqueue_backfills(source)
            except Exception:
                self.log.exception("Failed to backfill new portal")

            # Update participants again after backfill to sync read receipts
            self._last_participant_update = set()
            await self._update_participants(info.participants)

        return self.mxid

    async def _enqueue_backfills(self, source: u.User) -> None:
        state = DBBackfillStatus(self.twid, self.receiver, source.twid, False, 0, 0)
        await state.insert()
        await b.BackfillStatus.recheck()

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
    def all_with_room(cls) -> AsyncGenerator[Portal, None]:
        return cls._db_to_portals(super().all_with_room())

    @classmethod
    def find_private_chats_with(cls, other_user: int) -> AsyncGenerator[Portal, None]:
        return cls._db_to_portals(super().find_private_chats_with(other_user))

    @classmethod
    async def _db_to_portals(cls, query: Awaitable[list[Portal]]) -> AsyncGenerator[Portal, None]:
        portals = await query
        for index, portal in enumerate(portals):
            try:
                yield cls.by_twid[(portal.twid, portal.receiver)]
            except KeyError:
                await portal.postinit()
                yield portal

    @classmethod
    @async_getter_lock
    async def get_by_mxid(cls, mxid: RoomID) -> Portal | None:
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
    @async_getter_lock
    async def get_by_twid(
        cls,
        twid: str,
        *,
        receiver: int = 0,
        conv_type: ConversationType | None = None,
    ) -> Portal | None:
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

    async def get_dm_puppet(self) -> p.Puppet | None:
        if not self.is_direct:
            return None
        return await p.Puppet.get_by_twid(self.twid)

    # endregion
