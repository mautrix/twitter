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
from os import path

from aiohttp import ClientSession
from yarl import URL

from mautrix.appservice import IntentAPI
from mautrix.bridge import BasePuppet, async_getter_lock
from mautrix.types import ContentURI, RoomID, SyncToken, UserID
from mautrix.util.simple_template import SimpleTemplate
from mautwitdm.types import User

from . import portal as p
from .config import Config
from .db import Puppet as DBPuppet

if TYPE_CHECKING:
    from .__main__ import TwitterBridge


class Puppet(DBPuppet, BasePuppet):
    by_twid: dict[int, Puppet] = {}
    by_custom_mxid: dict[UserID, Puppet] = {}
    hs_domain: str
    mxid_template: SimpleTemplate[int]

    config: Config

    default_mxid_intent: IntentAPI
    default_mxid: UserID

    def __init__(
        self,
        twid: int,
        name: str | None = None,
        photo_url: str | None = None,
        photo_mxc: ContentURI | None = None,
        is_registered: bool = False,
        custom_mxid: UserID | None = None,
        access_token: str | None = None,
        next_batch: SyncToken | None = None,
        base_url: URL | None = None,
    ) -> None:
        super().__init__(
            twid=twid,
            name=name,
            photo_url=photo_url,
            photo_mxc=photo_mxc,
            is_registered=is_registered,
            custom_mxid=custom_mxid,
            access_token=access_token,
            next_batch=next_batch,
            base_url=base_url,
        )
        self.log = self.log.getChild(str(twid))

        self.default_mxid = self.get_mxid_from_id(twid)
        self.default_mxid_intent = self.az.intent.user(self.default_mxid)
        self.intent = self._fresh_intent()

    @classmethod
    def init_cls(cls, bridge: "TwitterBridge") -> AsyncIterable[Awaitable[None]]:
        cls.config = bridge.config
        cls.loop = bridge.loop
        cls.mx = bridge.matrix
        cls.az = bridge.az
        cls.hs_domain = cls.config["homeserver.domain"]
        cls.mxid_template = SimpleTemplate(
            cls.config["bridge.username_template"],
            "userid",
            prefix="@",
            suffix=f":{cls.hs_domain}",
            type=int,
        )
        cls.sync_with_custom_puppets = cls.config["bridge.sync_with_custom_puppets"]
        cls.homeserver_url_map = {
            server: URL(url)
            for server, url in cls.config["bridge.double_puppet_server_map"].items()
        }
        cls.allow_discover_url = cls.config["bridge.double_puppet_allow_discovery"]
        cls.login_shared_secret_map = {
            server: secret.encode("utf-8")
            for server, secret in cls.config["bridge.login_shared_secret_map"].items()
        }
        cls.login_device_name = "Twitter DM Bridge"
        return (puppet.try_start() async for puppet in cls.all_with_custom_mxid())

    def intent_for(self, portal: p.Portal) -> IntentAPI:
        if portal.other_user == self.twid or (
            self.config["bridge.backfill.invite_own_puppet"] and portal.backfill_lock.locked
        ):
            return self.default_mxid_intent
        return self.intent

    async def update_info(self, info: User) -> None:
        update = False
        update = await self._update_name(info) or update
        update = await self._update_avatar(info.profile_image_url_https) or update
        if update:
            await self.update()

    @classmethod
    def _get_displayname(cls, info: User) -> str:
        return cls.config["bridge.displayname_template"].format(
            displayname=info.name, id=info.id, username=info.screen_name
        )

    async def _update_name(self, info: User) -> bool:
        name = self._get_displayname(info)
        if name != self.name:
            self.name = name
            await self.default_mxid_intent.set_displayname(self.name)
            await self._update_portal_names()
            return True
        return False

    async def _update_portal_names(self) -> None:
        async for portal in p.Portal.find_private_chats_with(self.twid):
            await portal.update_name(self.name)

    async def _update_avatar(self, image_url: str) -> bool:
        if image_url != self.photo_url:
            self.photo_url = image_url
            url = URL(self.photo_url.replace("_normal.", "_400x400."))
            file_name = path.basename(url.path)
            async with ClientSession() as sess, sess.get(url) as resp:
                content_type = resp.headers["Content-Type"]
                resp_data = await resp.read()
            mxc = await self.default_mxid_intent.upload_media(
                data=resp_data,
                filename=file_name,
                mime_type=content_type,
                async_upload=self.config["homeserver.async_media"],
            )
            self.photo_mxc = mxc
            await self.default_mxid_intent.set_avatar_url(mxc)
            return True
        return False

    async def default_puppet_should_leave_room(self, room_id: RoomID) -> bool:
        portal = await p.Portal.get_by_mxid(room_id)
        return portal and portal.other_user != self.twid

    # region Database getters

    def _add_to_cache(self) -> None:
        self.by_twid[self.twid] = self
        if self.custom_mxid:
            self.by_custom_mxid[self.custom_mxid] = self

    async def save(self) -> None:
        await self.update()

    @classmethod
    async def get_by_mxid(cls, mxid: UserID, create: bool = True) -> Puppet | None:
        twid = cls.get_id_from_mxid(mxid)
        if twid:
            return await cls.get_by_twid(twid, create=create)
        return None

    @classmethod
    @async_getter_lock
    async def get_by_custom_mxid(cls, mxid: UserID) -> Puppet | None:
        try:
            return cls.by_custom_mxid[mxid]
        except KeyError:
            pass

        puppet = cast(cls, await super().get_by_custom_mxid(mxid))
        if puppet:
            puppet._add_to_cache()
            return puppet

        return None

    @classmethod
    def get_id_from_mxid(cls, mxid: UserID) -> int | None:
        return cls.mxid_template.parse(mxid)

    @classmethod
    def get_mxid_from_id(cls, twid: int) -> UserID:
        return UserID(cls.mxid_template.format_full(twid))

    @classmethod
    @async_getter_lock
    async def get_by_twid(cls, twid: int, *, create: bool = True) -> Puppet | None:
        try:
            return cls.by_twid[twid]
        except KeyError:
            pass

        puppet = cast(cls, await super().get_by_twid(twid))
        if puppet is not None:
            puppet._add_to_cache()
            return puppet

        if create:
            puppet = cls(twid)
            await puppet.insert()
            puppet._add_to_cache()
            return puppet

        return None

    @classmethod
    async def all_with_custom_mxid(cls) -> AsyncGenerator[Puppet, None]:
        puppets = await super().all_with_custom_mxid()
        puppet: cls
        for index, puppet in enumerate(puppets):
            try:
                yield cls.by_twid[puppet.twid]
            except KeyError:
                puppet._add_to_cache()
                yield puppet

    # endregion
