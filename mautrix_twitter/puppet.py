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
from typing import Optional, Dict, List, Awaitable, TYPE_CHECKING, cast

from mautrix.bridge import BasePuppet, CustomPuppetMixin
from mautrix.types import ContentURI, UserID, SyncToken
from mautrix.util.simple_template import SimpleTemplate

from .db import Puppet as DBPuppet
from .config import Config

if TYPE_CHECKING:
    from .__main__ import TwitterBridge


class Puppet(DBPuppet, BasePuppet, CustomPuppetMixin):
    by_twid: Dict[int, 'Puppet'] = {}
    by_custom_mxid: Dict[UserID, 'Puppet'] = {}
    hs_domain: str
    mxid_template: SimpleTemplate

    config: Config

    def __init__(self, twid: int, name: Optional[str] = None, photo_url: Optional[str] = None,
                 photo_mxc: Optional[ContentURI] = None, is_registered: bool = False,
                 custom_mxid: Optional[UserID] = None, access_token: Optional[str] = None,
                 next_batch: Optional[SyncToken] = None) -> None:
        super().__init__(twid, name, photo_url, photo_mxc, is_registered, custom_mxid,
                         access_token, next_batch)
        self.log = self.log.getChild(str(twid))

    @classmethod
    async def init_cls(cls, bridge: 'TwitterBridge') -> List[Awaitable[None]]:
        cls.config = bridge.config
        cls.hs_domain = cls.config["homeserver.domain"]
        cls.mxid_template = SimpleTemplate(cls.config["bridge.username_template"], "userid",
                                           prefix="@", suffix=f":{cls.hs_domain}", type=str)
        cls.sync_with_custom_puppets = cls.config["bridge.sync_with_custom_puppets"]
        secret = cls.config["bridge.login_shared_secret"]
        cls.login_shared_secret = secret.encode("utf-8") if secret else None
        cls.login_device_name = "Twitter DM Bridge"
        return [puppet.try_start() for puppet in await cls.all_with_custom_mxid()]

    # region Database getters

    def _add_to_cache(self) -> None:
        self.by_twid[self.twid] = self
        if self.custom_mxid:
            self.by_custom_mxid[self.custom_mxid] = self

    @classmethod
    async def get_by_twid(cls, twid: int, create: bool = True) -> Optional['Puppet']:
        try:
            return cls.by_twid[twid]
        except KeyError:
            pass

        puppet = cast(Puppet, await super().get_by_twid(twid))
        if puppet is not None:
            puppet._add_to_cache()
            return puppet

        if create:
            puppet = Puppet(twid)
            await puppet.insert()
            puppet._add_to_cache()
            return puppet

        return None

    @classmethod
    async def all_with_custom_mxid(cls) -> List['Puppet']:
        puppets = await super().all_with_custom_mxid()
        puppet: Puppet
        for index, puppet in enumerate(puppets):
            try:
                puppets[index] = cls.by_twid[puppet.twid]
            except KeyError:
                puppet._add_to_cache()
        return puppets

    # endregion
