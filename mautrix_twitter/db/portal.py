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

from typing import TYPE_CHECKING, ClassVar

from attr import dataclass
import asyncpg

from mautrix.types import RoomID
from mautrix.util.async_db import Database
from mautwitdm.types import ConversationType

fake_db = Database.create("") if TYPE_CHECKING else None


@dataclass
class Portal:
    db: ClassVar[Database] = fake_db

    twid: str
    receiver: int
    conv_type: ConversationType
    other_user: int | None
    mxid: RoomID | None
    name: str | None
    encrypted: bool

    @property
    def _values(self):
        return (
            self.twid,
            self.receiver,
            self.conv_type.value,
            self.other_user,
            self.mxid,
            self.name,
            self.encrypted,
        )

    async def insert(self) -> None:
        q = (
            "INSERT INTO portal (twid, receiver, conv_type, other_user, mxid, name, encrypted) "
            "VALUES ($1, $2, $3, $4, $5, $6, $7)"
        )
        await self.db.execute(q, *self._values)

    async def update(self) -> None:
        q = (
            "UPDATE portal SET conv_type=$3, other_user=$4, mxid=$5, name=$6, encrypted=$7 "
            "WHERE twid=$1 AND receiver=$2"
        )
        await self.db.execute(q, *self._values)

    @classmethod
    def _from_row(cls, row: asyncpg.Record) -> "Portal":
        data = {**row}
        return cls(conv_type=ConversationType(data.pop("conv_type")), **data)

    @classmethod
    async def get_by_mxid(cls, mxid: RoomID) -> Portal | None:
        q = (
            "SELECT twid, receiver, conv_type, other_user, mxid, name, encrypted "
            "FROM portal WHERE mxid=$1"
        )
        row = await cls.db.fetchrow(q, mxid)
        if not row:
            return None
        return cls._from_row(row)

    @classmethod
    async def get_by_twid(cls, twid: str, receiver: int = 0) -> Portal | None:
        q = (
            "SELECT twid, receiver, conv_type, other_user, mxid, name, encrypted "
            "FROM portal WHERE twid=$1 AND receiver=$2"
        )
        row = await cls.db.fetchrow(q, twid, receiver)
        if not row:
            return None
        return cls._from_row(row)

    @classmethod
    async def find_private_chats_of(cls, receiver: int) -> list[Portal]:
        q = (
            "SELECT twid, receiver, conv_type, other_user, mxid, name, encrypted FROM portal "
            "WHERE receiver=$1 AND conv_type='ONE_TO_ONE'"
        )
        rows = await cls.db.fetch(q, receiver)
        return [cls._from_row(row) for row in rows]

    @classmethod
    async def find_private_chats_with(cls, other_user: int) -> list[Portal]:
        q = (
            "SELECT twid, receiver, conv_type, other_user, mxid, name, encrypted FROM portal "
            "WHERE other_user=$1 AND conv_type='ONE_TO_ONE'"
        )
        rows = await cls.db.fetch(q, other_user)
        return [cls._from_row(row) for row in rows]

    @classmethod
    async def all_with_room(cls) -> list[Portal]:
        q = (
            "SELECT twid, receiver, conv_type, other_user, mxid, name, encrypted FROM portal "
            "WHERE mxid IS NOT NULL"
        )
        rows = await cls.db.fetch(q)
        return [cls._from_row(row) for row in rows]
