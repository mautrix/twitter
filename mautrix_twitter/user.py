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
from typing import (Dict, Optional, AsyncIterable, Awaitable, AsyncGenerator, Union, List,
                    TYPE_CHECKING, cast)
import asyncio
import logging

from mautwitdm import TwitterAPI
from mautwitdm.types import (MessageEntry, ReactionCreateEntry, ReactionDeleteEntry, Conversation,
                             User as TwitterUser, ConversationReadEntry)
from mautrix.bridge import BaseUser
from mautrix.types import UserID, RoomID
from mautrix.appservice import AppService

from .db import User as DBUser, Portal as DBPortal
from .config import Config
from . import puppet as pu, portal as po

if TYPE_CHECKING:
    from .__main__ import TwitterBridge


class User(DBUser, BaseUser):
    by_mxid: Dict[UserID, 'User'] = {}
    by_twid: Dict[int, 'User'] = {}
    config: Config
    az: AppService
    loop: asyncio.AbstractEventLoop

    client: Optional[TwitterAPI]

    is_admin: bool
    permission_level: str
    username: Optional[str]

    _notice_room_lock: asyncio.Lock

    def __init__(self, mxid: UserID, twid: Optional[int] = None, auth_token: Optional[str] = None,
                 csrf_token: Optional[str] = None, poll_cursor: Optional[str] = None,
                 notice_room: Optional[RoomID] = None) -> None:
        super().__init__(mxid=mxid, twid=twid, auth_token=auth_token, csrf_token=csrf_token,
                         poll_cursor=poll_cursor, notice_room=notice_room)
        self._notice_room_lock = asyncio.Lock()
        perms = self.config.get_permissions(mxid)
        self.is_whitelisted, self.is_admin, self.permission_level = perms
        self.log = self.log.getChild(self.mxid)
        self.client = None
        self.username = None
        self.dm_update_lock = asyncio.Lock()

    @classmethod
    def init_cls(cls, bridge: 'TwitterBridge') -> AsyncIterable[Awaitable[None]]:
        cls.bridge = bridge
        cls.config = bridge.config
        cls.az = bridge.az
        cls.loop = bridge.loop
        return (user.try_connect() async for user in cls.all_logged_in())

    async def update(self) -> None:
        if self.client:
            self.auth_token, self.csrf_token = self.client.tokens
            self.poll_cursor = self.client.poll_cursor
        await super().update()

    # region Connection management

    async def is_logged_in(self) -> bool:
        try:
            return self.client and await self.client.get_user_identifier() is not None
        except Exception:
            return False

    async def try_connect(self) -> None:
        try:
            await self.connect()
        except Exception:
            self.log.exception("Error while connecting to Twitter")

    async def connect(self, auth_token: Optional[str] = None, csrf_token: Optional[str] = None
                      ) -> None:
        client = TwitterAPI(log=logging.getLogger("mau.twitter.api").getChild(self.mxid),
                            loop=self.loop, node_id=hash(self.mxid) % (2 ** 48))
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

        user_info = await self.get_info()
        self.twid = user_info.id
        self.by_twid[self.twid] = self

        await self.update()

        if self.poll_cursor:
            self.log.debug("Poll cursor set, starting polling right away (not initial syncing)")
            self.client.start_polling()
        else:
            self.loop.create_task(self._try_initial_sync())
        self.loop.create_task(self._try_sync_puppet(user_info))

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

    async def get_direct_chats(self) -> Dict[UserID, List[RoomID]]:
        return {
            pu.Puppet.get_mxid_from_id(portal.other_user): [portal.mxid]
            for portal in await DBPortal.find_private_chats(self.twid)
            if portal.mxid
        }

    async def sync(self) -> None:
        resp = await self.client.inbox_initial_state(set_poll_cursor=False)
        if not self.poll_cursor:
            self.poll_cursor = resp.cursor
        self.client.poll_cursor = self.poll_cursor
        limit = self.config["bridge.initial_conversation_sync"]
        conversations = sorted(resp.conversations.values(), key=lambda conv: conv.sort_timestamp)
        if limit < 0:
            limit = len(conversations)
        for i, conversation in enumerate(conversations):
            await self.handle_conversation_update(conversation, create_portal=i < limit)
        for user in resp.users.values():
            await self.handle_user_update(user)
        await self.update_direct_chats()

    async def get_info(self) -> TwitterUser:
        settings = await self.client.get_settings()
        self.username = settings["screen_name"]
        return (await self.client.lookup_users(usernames=[self.username]))[0]

    async def stop(self) -> None:
        if self.client:
            self.client.stop_polling()
        await self.update()

    async def logout(self) -> None:
        if self.client:
            self.client.stop_polling()
        puppet = await pu.Puppet.get_by_twid(self.twid, create=False)
        if puppet and puppet.is_real_user:
            await puppet.switch_mxid(None, None)
        try:
            del self.by_twid[self.twid]
        except KeyError:
            pass
        self.client = None
        self.twid = None
        self.poll_cursor = None
        self.auth_token = None
        self.csrf_token = None
        await self.update()

    # endregion
    # region Event handlers

    async def handle_conversation_update(self, evt: Conversation, create_portal: bool = False
                                         ) -> None:
        portal = await po.Portal.get_by_twid(evt.conversation_id, self.twid, conv_type=evt.type)
        if not portal.mxid:
            if create_portal:
                await portal.create_matrix_room(self, evt)
        else:
            # We don't want to do the invite_user and such things each time conversation info
            # comes down polling, so if the room already exists, only call .update_info()
            await portal.update_info(evt)

    async def handle_user_update(self, user: TwitterUser) -> None:
        puppet = await pu.Puppet.get_by_twid(user.id)
        await puppet.update_info(user)

    async def handle_message(self, evt: MessageEntry) -> None:
        portal = await po.Portal.get_by_twid(evt.conversation_id, self.twid,
                                             conv_type=evt.conversation.type)
        if not portal.mxid:
            await portal.create_matrix_room(self, evt.conversation)
        sender = await pu.Puppet.get_by_twid(int(evt.message_data.sender_id))
        await portal.handle_twitter_message(self, sender, evt.message_data, evt.request_id)

    async def handle_reaction(self, evt: Union[ReactionCreateEntry, ReactionDeleteEntry]) -> None:
        portal = await po.Portal.get_by_twid(evt.conversation_id, self.twid,
                                             conv_type=evt.conversation.type)
        if not portal.mxid:
            self.log.debug(f"Ignoring reaction in conversation {evt.conversation_id} with no room")
            return
        puppet = await pu.Puppet.get_by_twid(int(evt.sender_id))
        if isinstance(evt, ReactionCreateEntry):
            await portal.handle_twitter_reaction_add(puppet, int(evt.message_id),
                                                     evt.reaction_key, evt.time)
        else:
            await portal.handle_twitter_reaction_remove(puppet, int(evt.message_id),
                                                        evt.reaction_key)

    async def handle_receipt(self, evt: ConversationReadEntry) -> None:
        portal = await po.Portal.get_by_twid(evt.conversation_id, self.twid,
                                             conv_type=evt.conversation.type)
        if not portal.mxid:
            return
        sender = await pu.Puppet.get_by_twid(self.twid)
        await portal.handle_twitter_receipt(sender, int(evt.last_read_event_id))

    # endregion
    # region Database getters

    def _add_to_cache(self) -> None:
        self.by_mxid[self.mxid] = self
        if self.twid:
            self.by_twid[self.twid] = self

    @classmethod
    async def get_by_mxid(cls, mxid: UserID, create: bool = True) -> Optional['User']:
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
    async def get_by_twid(cls, twid: int) -> Optional['User']:
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
    async def all_logged_in(cls) -> AsyncGenerator['User', None]:
        users = await super().all_logged_in()
        user: cls
        for index, user in enumerate(users):
            try:
                yield cls.by_mxid[user.mxid]
            except KeyError:
                user._add_to_cache()
                yield user

    # endregion
