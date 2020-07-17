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
from typing import Dict, Tuple, Optional, TYPE_CHECKING, cast
import asyncio

from mautwitdm import ConversationType
from mautrix.appservice import IntentAPI
from mautrix.bridge import BasePortal
from mautrix.types import RoomID

from .db import Portal as DBPortal
from .config import Config

if TYPE_CHECKING:
    from .__main__ import TwitterBridge


class Portal(DBPortal, BasePortal):
    by_mxid: Dict[RoomID, 'Portal'] = {}
    by_twid: Dict[Tuple[str, int], 'Portal'] = {}
    config: Config

    _main_intent: Optional[IntentAPI]
    _create_room_lock: asyncio.Lock

    def __init__(self, twid: str, receiver: int, conv_type: ConversationType,
                 other_user: Optional[int] = None, mxid: Optional[RoomID] = None,
                 name: Optional[str] = None, encrypted: bool = False) -> None:
        super().__init__(twid, receiver, conv_type, other_user, mxid, name, encrypted)
        self._create_room_lock = asyncio.Lock()
        self.log = self.log.getChild(twid)

    @property
    def is_direct(self) -> bool:
        return self.conv_type == ConversationType.ONE_TO_ONE

    @classmethod
    def init_cls(cls, bridge: 'TwitterBridge') -> None:
        cls.config = bridge.config

    # region Database getters

    def _add_to_cache(self) -> None:
        self.by_twid[(self.twid, self.receiver)] = self
        if self.mxid:
            self.by_mxid[self.mxid] = self

    @classmethod
    async def get_by_mxid(cls, mxid: RoomID) -> Optional['Portal']:
        try:
            return cls.by_mxid[mxid]
        except KeyError:
            pass

        portal = cast(Portal, await super().get_by_mxid(mxid))
        if portal is not None:
            portal._add_to_cache()
            return portal

        return None

    @classmethod
    async def get_by_twid(cls, twid: str, receiver: int = 0,
                          conv_type: Optional[ConversationType] = None) -> Optional['Portal']:
        if conv_type == ConversationType.GROUP_DM and receiver != 0:
            raise ValueError("receiver must be 0 when conv_type is GROUP_DM")
        try:
            return cls.by_twid[(twid, receiver)]
        except KeyError:
            pass

        portal = cast(Portal, await super().get_by_twid(twid, receiver))
        if portal is not None:
            portal._add_to_cache()
            return portal

        if conv_type is not None:
            portal = Portal(twid, receiver, conv_type)
            await portal.insert()
            portal._add_to_cache()
            return portal

        return None

    # endregion
