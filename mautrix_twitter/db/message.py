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

from mautrix.types import EventID, RoomID
from mautrix.util.async_db import Database

fake_db = Database.create("") if TYPE_CHECKING else None


@dataclass
class Message:
    db: ClassVar[Database] = fake_db

    mxid: EventID
    mx_room: RoomID
    twid: int
    receiver: int

    async def insert(self) -> None:
        q = "INSERT INTO message (mxid, mx_room, twid, receiver) VALUES ($1, $2, $3, $4)"
        await self.db.execute(q, self.mxid, self.mx_room, self.twid, self.receiver)

    async def delete(self) -> None:
        q = "DELETE FROM message WHERE twid=$1 AND receiver=$2"
        await self.db.execute(q, self.twid, self.receiver)

    @classmethod
    async def delete_all(cls, room_id: RoomID) -> None:
        await cls.db.execute("DELETE FROM message WHERE mx_room=$1", room_id)

    @classmethod
    async def get_by_mxid(cls, mxid: EventID, mx_room: RoomID) -> Message | None:
        q = "SELECT mxid, mx_room, twid, receiver FROM message WHERE mxid=$1 AND mx_room=$2"
        row = await cls.db.fetchrow(q, mxid, mx_room)
        if not row:
            return None
        return cls(**row)

    @classmethod
    async def get_last(cls, mx_room: RoomID) -> Message | None:
        q = "SELECT mxid, mx_room, twid, receiver FROM message WHERE mx_room=$1 ORDER BY twid DESC LIMIT 1"
        row = await cls.db.fetchrow(q, mx_room)
        if not row:
            return None
        return cls(**row)

    @classmethod
    async def get_by_twid(cls, twid: int, receiver: int = 0) -> Message | None:
        q = "SELECT mxid, mx_room, twid, receiver FROM message WHERE twid=$1 AND receiver=$2"
        row = await cls.db.fetchrow(q, twid, receiver)
        if not row:
            return None
        return cls(**row)
