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

from mautrix.types import EventID, RoomID
from mautrix.util.async_db import Database
from mautwitdm.types import ReactionKey

fake_db = Database.create("") if TYPE_CHECKING else None


@dataclass
class Reaction:
    db: ClassVar[Database] = fake_db

    mxid: EventID
    mx_room: RoomID
    tw_msgid: int
    tw_receiver: int
    tw_sender: int
    reaction: ReactionKey
    tw_reaction_id: int

    async def insert(self) -> None:
        q = (
            "INSERT INTO reaction (mxid, mx_room, tw_msgid, tw_receiver, tw_sender, reaction, tw_reaction_id) "
            "VALUES ($1, $2, $3, $4, $5, $6, $7)"
        )
        await self.db.execute(
            q,
            self.mxid,
            self.mx_room,
            self.tw_msgid,
            self.tw_receiver,
            self.tw_sender,
            self.reaction.value,
            self.tw_reaction_id,
        )

    async def edit(
        self, mx_room: RoomID, mxid: EventID, reaction: ReactionKey, tw_reaction_id: int | None
    ) -> None:
        q = (
            "UPDATE reaction SET mxid=$1, mx_room=$2, reaction=$3, tw_reaction_id=$7 "
            "WHERE tw_msgid=$4 AND tw_receiver=$5 AND tw_sender=$6"
        )
        await self.db.execute(
            q,
            mxid,
            mx_room,
            reaction.value,
            self.tw_msgid,
            self.tw_receiver,
            self.tw_sender,
            tw_reaction_id,
        )

    async def delete(self) -> None:
        q = "DELETE FROM reaction WHERE tw_msgid=$1 AND tw_receiver=$2 AND tw_sender=$3"
        await self.db.execute(q, self.tw_msgid, self.tw_receiver, self.tw_sender)

    @classmethod
    def _from_row(cls, row: asyncpg.Record) -> Reaction | None:
        if not row:
            return None
        data = {**row}
        reaction = ReactionKey(data.pop("reaction"))
        return cls(reaction=reaction, **data)

    @classmethod
    async def get_by_mxid(cls, mxid: EventID, mx_room: RoomID) -> Reaction | None:
        q = (
            "SELECT mxid, mx_room, tw_msgid, tw_receiver, tw_sender, reaction, tw_reaction_id "
            "FROM reaction WHERE mxid=$1 AND mx_room=$2"
        )
        return cls._from_row(await cls.db.fetchrow(q, mxid, mx_room))

    @classmethod
    async def get_last(cls, mx_room: RoomID) -> Reaction | None:
        q = (
            "SELECT mxid, mx_room, tw_msgid, tw_receiver, tw_sender, reaction, tw_reaction_id "
            "FROM reaction WHERE mx_room=$1 and tw_reaction_id IS NOT NULL ORDER BY tw_reaction_id DESC LIMIT 1"
        )
        return cls._from_row(await cls.db.fetchrow(q, mx_room))

    @classmethod
    async def get_by_twid(
        cls,
        tw_msgid: int,
        tw_receiver: int,
        tw_sender: int,
    ) -> Reaction | None:
        q = (
            "SELECT mxid, mx_room, tw_msgid, tw_receiver, tw_sender, reaction, tw_reaction_id "
            "FROM reaction WHERE tw_msgid=$1 AND tw_sender=$2 AND tw_receiver=$3"
        )
        return cls._from_row(await cls.db.fetchrow(q, tw_msgid, tw_sender, tw_receiver))
