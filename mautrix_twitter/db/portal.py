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
from typing import Optional, ClassVar, List, TYPE_CHECKING

from attr import dataclass

from mautrix.types import RoomID
from mautrix.util.async_db import Database

from mautwitdm.types import ConversationType

fake_db = Database("") if TYPE_CHECKING else None


@dataclass
class Portal:
    db: ClassVar[Database] = fake_db

    twid: str
    receiver: int
    conv_type: ConversationType
    other_user: Optional[int]
    mxid: Optional[RoomID]
    name: Optional[str]
    encrypted: bool

    async def insert(self) -> None:
        q = ("INSERT INTO portal (twid, receiver, conv_type, other_user, mxid, name, encrypted) "
             "VALUES ($1, $2, $3, $4, $5, $6, $7)")
        await self.db.execute(q, self.twid, self.receiver, self.conv_type.value, self.other_user,
                              self.mxid, self.name, self.encrypted)

    async def update(self) -> None:
        q = ("UPDATE portal SET conv_type=$3, other_user=$4, mxid=$5, name=$6, encrypted=$7 "
             "WHERE twid=$1 AND receiver=$2")
        await self.db.execute(q, self.twid, self.receiver, self.conv_type.value, self.other_user,
                              self.mxid, self.name, self.encrypted)

    @classmethod
    async def get_by_mxid(cls, mxid: RoomID) -> Optional['Portal']:
        q = ("SELECT twid, receiver, conv_type, other_user, mxid, name, encrypted "
             "FROM portal WHERE mxid=$1")
        row = await cls.db.fetchrow(q, mxid)
        if not row:
            return None
        data = {**row}
        return cls(conv_type=ConversationType(data.pop("conv_type")), **data)

    @classmethod
    async def get_by_twid(cls, twid: str, receiver: int = 0) -> Optional['Portal']:
        q = ("SELECT twid, receiver, conv_type, other_user, mxid, name, encrypted "
             "FROM portal WHERE twid=$1 AND receiver=$2")
        row = await cls.db.fetchrow(q, twid, receiver)
        if not row:
            return None
        data = {**row}
        return cls(conv_type=ConversationType(data.pop("conv_type")), **data)

    @classmethod
    async def find_private_chats(cls, receiver: int) -> List['Portal']:
        q = ("SELECT twid, receiver, conv_type, other_user, mxid, name, encrypted FROM portal "
             "WHERE receiver=$1 AND conv_type='ONE_TO_ONE'")
        rows = await cls.db.fetch(q, receiver)
        return [cls(**row) for row in rows]

    @classmethod
    async def all_with_room(cls) -> List['Portal']:
        q = ("SELECT twid, receiver, conv_type, other_user, mxid, name, encrypted FROM portal "
             'WHERE mxid IS NOT NULL')
        rows = await cls.db.fetch(q)
        return [cls(**row) for row in rows]
