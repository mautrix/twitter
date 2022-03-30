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

from typing import TYPE_CHECKING, AsyncGenerator, AsyncIterable, Awaitable, cast
import asyncio
import logging

from mautrix.bridge import BaseUser, async_getter_lock
from mautrix.types import EventID, MessageType, RoomID, TextMessageEventContent, UserID
from mautrix.util.bridge_state import BridgeState, BridgeStateEvent
from mautrix.util.opt_prometheus import Gauge, Summary, async_time
from mautwitdm import TwitterAPI
from mautwitdm.errors import TwitterAuthError, TwitterError
from mautwitdm.poller import PollingErrored, PollingErrorResolved, PollingStarted, PollingStopped
from mautwitdm.types import (
    Conversation,
    ConversationReadEntry,
    MessageEntry,
    ReactionCreateEntry,
    ReactionDeleteEntry,
    User as TwitterUser,
)
from mautwitdm.types.conversation import ConversationType

from . import portal as po, puppet as pu
from .config import Config
from .db import Portal as DBPortal, User as DBUser

if TYPE_CHECKING:
    from .__main__ import TwitterBridge

METRIC_CONVERSATION_UPDATE = Summary(
    "bridge_on_conversation_update", "calls to handle_conversation_update"
)
METRIC_USER_UPDATE = Summary("bridge_on_user_update", "calls to handle_user_update")
METRIC_MESSAGE = Summary("bridge_on_message", "calls to handle_message")
METRIC_REACTION = Summary("bridge_on_reaction", "calls to handle_reaction")
METRIC_RECEIPT = Summary("bridge_on_receipt", "calls to handle_receipt")
METRIC_LOGGED_IN = Gauge("bridge_logged_in", "Users logged into the bridge")
METRIC_CONNECTED = Gauge("bridge_connected", "Bridged users connected to Twitter")

BridgeState.human_readable_errors.update(
    {
        "logged-out": "You're not logged into Twitter",
        "twitter-connection-failed": "Failed to connect to Twitter",
        "twitter-not-connected": None,
        "twitter-connection-error": "An error occurred while polling Twitter for new messages",
    }
)


class User(DBUser, BaseUser):
    by_mxid: dict[UserID, User] = {}
    by_twid: dict[int, User] = {}
    config: Config

    client: TwitterAPI | None

    permission_level: str
    username: str | None

    _notice_room_lock: asyncio.Lock
    _notice_send_lock: asyncio.Lock
    _is_logged_in: bool | None
    _connected: bool
    _intentional_stop: bool

    def __init__(
        self,
        mxid: UserID,
        twid: int | None = None,
        auth_token: str | None = None,
        csrf_token: str | None = None,
        poll_cursor: str | None = None,
        notice_room: RoomID | None = None,
    ) -> None:
        super().__init__(
            mxid=mxid,
            twid=twid,
            auth_token=auth_token,
            csrf_token=csrf_token,
            poll_cursor=poll_cursor,
            notice_room=notice_room,
        )
        BaseUser.__init__(self)
        self._notice_room_lock = asyncio.Lock()
        self._notice_send_lock = asyncio.Lock()
        perms = self.config.get_permissions(mxid)
        self.is_whitelisted, self.is_admin, self.permission_level = perms
        self.client = None
        self.username = None
        self._is_logged_in = None
        self._connected = False
        self._intentional_stop = False

    @classmethod
    def init_cls(cls, bridge: "TwitterBridge") -> AsyncIterable[Awaitable[None]]:
        cls.bridge = bridge
        cls.config = bridge.config
        cls.az = bridge.az
        cls.loop = bridge.loop
        TwitterAPI.error_sleep = cls.config["bridge.error_sleep"]
        TwitterAPI.max_poll_errors = cls.config["bridge.max_poll_errors"]
        return (user.try_connect() async for user in cls.all_logged_in())

    async def update(self) -> None:
        if self.client:
            self.auth_token, self.csrf_token = self.client.tokens
            self.poll_cursor = self.client.poll_cursor
        await super().update()

    # region Connection management

    async def is_logged_in(self, ignore_cache: bool = False) -> bool:
        if not self.client:
            return False
        if self._is_logged_in is None:
            try:
                self._is_logged_in = await self.client.get_user_identifier() is not None
            except Exception:
                self._is_logged_in = False
        return self.client and self._is_logged_in

    async def get_puppet(self) -> pu.Puppet | None:
        if not self.twid:
            return None
        return await pu.Puppet.get_by_twid(self.twid)

    async def try_connect(self) -> None:
        try:
            await self.connect()
        except TwitterAuthError as e:
            self.log.exception("Auth error while connecting to Twitter")
            await self.push_bridge_state(
                BridgeStateEvent.BAD_CREDENTIALS,
                error="twitter-auth-error",
                message=e.message,
            )
        except TwitterError as e:
            self.log.exception("Error while connecting to Twitter")
            await self.push_bridge_state(
                BridgeStateEvent.UNKNOWN_ERROR,
                error="twitter-unknown-error",
                message=e.message,
            )
        except Exception:
            self.log.exception("Unknown exception while connecting to Twitter")
            await self.push_bridge_state(
                BridgeStateEvent.UNKNOWN_ERROR, error="twitter-connection-failed"
            )

    async def connect(self, auth_token: str | None = None, csrf_token: str | None = None) -> None:
        client = TwitterAPI(
            log=logging.getLogger("mau.twitter.api").getChild(self.mxid),
            loop=self.loop,
            node_id=hash(self.mxid) % (2**48),
        )
        client.poll_cursor = self.poll_cursor
        client.set_tokens(auth_token or self.auth_token, csrf_token or self.csrf_token)

        # Initial ping to make sure auth works
        await client.get_user_identifier()

        self.client = client
        self.client.add_handler(Conversation, self.handle_conversation_update)
        self.client.add_handler(TwitterUser, self.handle_user_update)
        self.client.add_handler(MessageEntry, self.handle_message)
        self.client.add_handler(ReactionCreateEntry, self.handle_reaction)
        self.client.add_handler(ReactionDeleteEntry, self.handle_reaction)
        self.client.add_handler(ConversationReadEntry, self.handle_receipt)
        self.client.add_handler(PollingStarted, self.on_connect)
        self.client.add_handler(PollingErrorResolved, self.on_connect)
        self.client.add_handler(PollingStopped, self.on_disconnect)
        self.client.add_handler(PollingErrored, self.on_disconnect)
        self.client.add_handler(PollingErrored, self.on_error)
        self.client.add_handler(PollingErrorResolved, self.on_error_resolved)

        user_info = await self.get_info()
        self.twid = user_info.id
        self._track_metric(METRIC_LOGGED_IN, True)
        self.by_twid[self.twid] = self

        await self.update()

        self._intentional_stop = False
        if self.poll_cursor:
            self.log.debug("Poll cursor set, starting polling right away (not initial syncing)")
            self.client.start_polling()
        else:
            asyncio.create_task(self._try_initial_sync())
        asyncio.create_task(self._try_sync_puppet(user_info))

    async def fill_bridge_state(self, state: BridgeState) -> None:
        await super().fill_bridge_state(state)
        if self.twid:
            state.remote_id = str(self.twid)
            puppet = await pu.Puppet.get_by_twid(self.twid)
            state.remote_name = puppet.name

    async def get_bridge_states(self) -> list[BridgeState]:
        if not self.twid:
            return []
        state = BridgeState(state_event=BridgeStateEvent.UNKNOWN_ERROR)
        if self._connected:
            state.state_event = BridgeStateEvent.CONNECTED
        return [state]

    async def on_connect(self, evt: PollingStarted | PollingErrorResolved) -> None:
        self._track_metric(METRIC_CONNECTED, True)
        self._connected = True
        await self.push_bridge_state(BridgeStateEvent.CONNECTED)

    async def on_disconnect(self, evt: PollingStopped | PollingErrored) -> None:
        self._track_metric(METRIC_CONNECTED, False)
        self._connected = False
        if isinstance(evt, PollingStopped) and not self._intentional_stop:
            await self.push_bridge_state(
                BridgeStateEvent.UNKNOWN_ERROR, error="twitter-not-connected"
            )
        self._intentional_stop = False

    # TODO this stuff could probably be moved to mautrix-python
    async def get_notice_room(self) -> RoomID:
        if not self.notice_room:
            async with self._notice_room_lock:
                # If someone already created the room while this call was waiting,
                # don't make a new room
                if self.notice_room:
                    return self.notice_room
                creation_content = {}
                if not self.config["bridge.federate_rooms"]:
                    creation_content["m.federate"] = False
                self.notice_room = await self.az.intent.create_room(
                    is_direct=True,
                    invitees=[self.mxid],
                    topic="Twitter DM bridge notices",
                    creation_content=creation_content,
                )
                await self.update()
        return self.notice_room

    async def send_bridge_notice(
        self,
        text: str,
        edit: EventID | None = None,
        state_event: BridgeStateEvent | None = None,
        important: bool = False,
    ) -> EventID | None:
        if state_event:
            await self.push_bridge_state(state_event, message=text)
        if self.config["bridge.disable_bridge_notices"]:
            return None
        event_id = None
        try:
            self.log.debug("Sending bridge notice: %s", text)
            content = TextMessageEventContent(
                body=text,
                msgtype=(MessageType.TEXT if important else MessageType.NOTICE),
            )
            if edit:
                content.set_edit(edit)
            # This is locked to prevent notices going out in the wrong order
            async with self._notice_send_lock:
                event_id = await self.az.intent.send_message(await self.get_notice_room(), content)
        except Exception:
            self.log.warning("Failed to send bridge notice", exc_info=True)
        return edit or event_id

    async def on_error(self, evt: PollingErrored) -> None:
        if evt.fatal:
            await self.send_bridge_notice(
                f"Fatal error while polling Twitter: {evt.error}",
                state_event=BridgeStateEvent.UNKNOWN_ERROR,
                important=evt.fatal,
            )
        elif evt.count == 1 and self.config["bridge.temporary_disconnect_notices"]:
            await self.send_bridge_notice(
                f"Error while polling Twitter: {evt.error}\nThe bridge will keep retrying.",
                state_event=BridgeStateEvent.TRANSIENT_DISCONNECT,
            )
        else:
            state_event = (
                BridgeStateEvent.TRANSIENT_DISCONNECT
                if evt.count < 5
                else BridgeStateEvent.UNKNOWN_ERROR
            )
            await self.push_bridge_state(state_event, f"Error while polling Twitter: {evt.error}")

    async def on_error_resolved(self, evt: PollingErrorResolved) -> None:
        if self.config["bridge.temporary_disconnect_notices"]:
            await self.send_bridge_notice(
                "Twitter polling error resolved", state_event=BridgeStateEvent.CONNECTED
            )

    async def _try_sync_puppet(self, user_info: TwitterUser) -> None:
        puppet = await pu.Puppet.get_by_twid(self.twid)
        try:
            await puppet.update_info(user_info)
        except Exception:
            self.log.exception("Failed to update own puppet info")
        try:
            if puppet.custom_mxid != self.mxid and puppet.can_auto_login(self.mxid):
                self.log.info(f"Automatically enabling custom puppet")
                await puppet.switch_mxid(access_token="auto", mxid=self.mxid)
        except Exception:
            self.log.exception("Failed to automatically enable custom puppet")

    async def _try_initial_sync(self) -> None:
        try:
            await self.sync()
        except Exception:
            self.log.exception("Exception while syncing conversations")
        self.log.debug("Initial sync completed, starting polling")
        self.client.start_polling()

    async def get_direct_chats(self) -> dict[UserID, list[RoomID]]:
        return {
            pu.Puppet.get_mxid_from_id(portal.other_user): [portal.mxid]
            for portal in await DBPortal.find_private_chats_of(self.twid)
            if portal.mxid
        }

    async def get_portal_with(self, puppet: pu.Puppet, create: bool = True) -> po.Portal | None:
        # We should probably make this work eventually, but for now, creating chats will just not
        # work.
        return None

    async def sync(self) -> None:
        await self.push_bridge_state(BridgeStateEvent.BACKFILLING)
        resp = await self.client.inbox_initial_state(set_poll_cursor=False)
        if not self.poll_cursor:
            self.poll_cursor = resp.cursor
        self.client.poll_cursor = self.poll_cursor
        limit = self.config["bridge.initial_conversation_sync"]
        conversations = sorted(
            resp.conversations.values(), key=lambda conv: conv.sort_timestamp, reverse=True
        )
        if limit < 0:
            limit = len(conversations)
        for user in resp.users.values():
            await self.handle_user_update(user)
        index = 0
        for conversation in conversations:
            create_portal = index < limit and conversation.trusted
            if create_portal:
                index += 1
            await self.handle_conversation_update(conversation, create_portal=create_portal)
        await self.update_direct_chats()

    async def get_info(self) -> TwitterUser:
        settings = await self.client.get_settings()
        self.username = settings["screen_name"]
        return (await self.client.lookup_users(usernames=[self.username]))[0]

    async def stop(self) -> None:
        if self.client:
            self._intentional_stop = True
            self.client.stop_polling()
        self._track_metric(METRIC_CONNECTED, False)
        await self.update()

    async def logout(self) -> None:
        if self.client:
            self._intentional_stop = True
            self.client.stop_polling()
        self._track_metric(METRIC_CONNECTED, False)
        self._track_metric(METRIC_LOGGED_IN, False)
        puppet = await pu.Puppet.get_by_twid(self.twid, create=False)
        if puppet and puppet.is_real_user:
            await puppet.switch_mxid(None, None)
        try:
            del self.by_twid[self.twid]
        except KeyError:
            pass
        self.client = None
        self._is_logged_in = None
        self.twid = None
        self.poll_cursor = None
        self.auth_token = None
        self.csrf_token = None
        await self.update()

    # endregion
    # region Event handlers

    @async_time(METRIC_CONVERSATION_UPDATE)
    async def handle_conversation_update(
        self, evt: Conversation, create_portal: bool = False
    ) -> None:
        portal = await po.Portal.get_by_twid(
            evt.conversation_id, receiver=self.twid, conv_type=evt.type
        )
        if not portal.mxid:
            if create_portal:
                await portal.create_matrix_room(self, evt)
        else:
            # We don't want to do the invite_user and such things each time conversation info
            # comes down polling, so if the room already exists, only call .update_info()
            await portal.update_info(evt)

    @async_time(METRIC_USER_UPDATE)
    async def handle_user_update(self, user: TwitterUser) -> None:
        puppet = await pu.Puppet.get_by_twid(user.id)
        await puppet.update_info(user)

    @async_time(METRIC_MESSAGE)
    async def handle_message(self, evt: MessageEntry) -> None:
        portal = await po.Portal.get_by_twid(
            evt.conversation_id, receiver=self.twid, conv_type=evt.conversation.type
        )
        if not portal.mxid:
            await portal.create_matrix_room(self, evt.conversation)
        sender = await pu.Puppet.get_by_twid(int(evt.message_data.sender_id))
        await portal.handle_twitter_message(self, sender, evt.message_data, evt.request_id)

    @async_time(METRIC_REACTION)
    async def handle_reaction(self, evt: ReactionCreateEntry | ReactionDeleteEntry) -> None:
        portal = await po.Portal.get_by_twid(
            evt.conversation_id, receiver=self.twid, conv_type=evt.conversation.type
        )
        if not portal.mxid:
            self.log.debug(f"Ignoring reaction in conversation {evt.conversation_id} with no room")
            return
        puppet = await pu.Puppet.get_by_twid(int(evt.sender_id))
        if isinstance(evt, ReactionCreateEntry):
            await portal.handle_twitter_reaction_add(
                puppet, int(evt.message_id), evt.reaction_key, evt.time, int(evt.id)
            )
        else:
            await portal.handle_twitter_reaction_remove(
                puppet, int(evt.message_id), evt.reaction_key
            )

    @async_time(METRIC_RECEIPT)
    async def handle_receipt(self, evt: ConversationReadEntry) -> None:
        portal = await po.Portal.get_by_twid(
            evt.conversation_id, receiver=self.twid, conv_type=evt.conversation.type
        )
        if not portal.mxid:
            return
        sender = await pu.Puppet.get_by_twid(self.twid)
        await portal.handle_twitter_receipt(sender, int(evt.last_read_event_id), historical=False)

    # endregion
    # region Database getters

    def _add_to_cache(self) -> None:
        self.by_mxid[self.mxid] = self
        if self.twid:
            self.by_twid[self.twid] = self

    @classmethod
    @async_getter_lock
    async def get_by_mxid(cls, mxid: UserID, *, create: bool = True) -> User | None:
        # Never allow ghosts to be users
        if pu.Puppet.get_id_from_mxid(mxid):
            return None
        try:
            return cls.by_mxid[mxid]
        except KeyError:
            pass

        user = cast(cls, await super().get_by_mxid(mxid))
        if user is not None:
            user._add_to_cache()
            return user

        if create:
            user = cls(mxid)
            await user.insert()
            user._add_to_cache()
            return user

        return None

    @classmethod
    @async_getter_lock
    async def get_by_twid(cls, twid: int) -> User | None:
        try:
            return cls.by_twid[twid]
        except KeyError:
            pass

        user = cast(cls, await super().get_by_twid(twid))
        if user is not None:
            user._add_to_cache()
            return user

        return None

    @classmethod
    async def all_logged_in(cls) -> AsyncGenerator[User, None]:
        users = await super().all_logged_in()
        user: cls
        for index, user in enumerate(users):
            try:
                yield cls.by_mxid[user.mxid]
            except KeyError:
                user._add_to_cache()
                yield user

    # endregion
