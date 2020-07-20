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
from typing import Optional, ClassVar, TYPE_CHECKING

from attr import dataclass

from mautwitdm.types import ReactionKey
from mautrix.types import RoomID, EventID
from mautrix.util.async_db import Database

fake_db = Database("") if TYPE_CHECKING else None


@dataclass
class Reaction:
    db: ClassVar[Database] = fake_db

    mxid: EventID
    mx_room: RoomID
    tw_msgid: int
    tw_receiver: int
    tw_sender: int
    reaction: ReactionKey

    async def insert(self) -> None:
        q = ("INSERT INTO reaction (mxid, mx_room, tw_msgid, tw_receiver, tw_sender, reaction) "
             "VALUES ($1, $2, $3, $4, $5, $6)")
        await self.db.execute(q, self.mxid, self.mx_room, self.tw_msgid, self.tw_receiver,
                              self.tw_sender, self.reaction.value)

    async def edit(self, mx_room: RoomID, mxid: EventID, reaction: ReactionKey) -> None:
        await self.db.execute("UPDATE reaction SET mxid=$1, mx_room=$2, reaction=$3 "
                              "WHERE tw_msgid=$4 AND tw_receiver=$5 AND tw_sender=$6",
                              mxid, mx_room, reaction.value, self.tw_msgid, self.tw_receiver,
                              self.tw_sender)

    async def delete(self) -> None:
        q = "DELETE FROM reaction WHERE tw_msgid=$1 AND tw_receiver=$2 AND tw_sender=$3"
        await self.db.execute(q, self.tw_msgid, self.tw_receiver, self.tw_sender)

    @classmethod
    async def get_by_mxid(cls, mxid: EventID, mx_room: RoomID) -> Optional['Reaction']:
        q = ("SELECT mxid, mx_room, tw_msgid, tw_receiver, tw_sender, reaction "
             "FROM reaction WHERE mxid=$1 AND mx_room=$2")
        row = await cls.db.fetchrow(q, mxid, mx_room)
        if not row:
            return None
        data = {**row}
        return cls(reaction=ReactionKey(data.pop("reaction")), **data)

    @classmethod
    async def get_by_twid(cls, tw_msgid: int, tw_receiver: int, tw_sender: int,
                          ) -> Optional['Reaction']:
        q = ("SELECT mxid, mx_room, tw_msgid, tw_receiver, tw_sender, reaction "
             "FROM reaction WHERE tw_msgid=$1 AND tw_sender=$2 AND tw_receiver=$3")
        row = await cls.db.fetchrow(q, tw_msgid, tw_sender, tw_receiver)
        if not row:
            return None
        data = {**row}
        return cls(reaction=ReactionKey(data.pop("reaction")), **data)
