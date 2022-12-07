# mautrix-twitter - A Matrix-Twitter DM puppeting bridge
# Copyright (C) 2022 Tulir Asokan, Max Sandholm
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

from mautrix.util.async_db import Database

fake_db = Database.create("") if TYPE_CHECKING else None


@dataclass
class BackfillStatus:
    db: ClassVar[Database] = fake_db

    twid: str
    receiver: int | None
    backfill_user: int
    dispatched: bool
    message_count: int
    state: int

    @property
    def _values(self):
        return (
            self.twid,
            self.receiver,
            self.backfill_user,
            self.dispatched,
            self.message_count,
            self.state,
        )

    async def insert(self) -> None:
        q = """INSERT INTO backfill_status (twid, receiver, backfill_user, dispatched, message_count, state)
            VALUES ($1, $2, $3, $4, $5, $6)"""
        await self.db.execute(q, *self._values)

    async def update(self) -> None:
        q = """UPDATE backfill_status SET backfill_user=$3, dispatched=$4, message_count=$5, state=$6
            WHERE twid=$1 AND receiver=$2"""
        await self.db.execute(q, *self._values)

    @classmethod
    async def get_by_twid(cls, twid: int, receiver: int = 0) -> BackfillStatus | None:
        q = """SELECT twid, receiver, backfill_user, dispatched, message_count, state FROM backfill_status
        WHERE twid =$1 AND receiver =$2"""

        row = await cls.db.fetchrow(q, twid, receiver)
        if not row:
            return None
        return cls(**row)

    @classmethod
    async def get_next_unfinished_status(cls) -> BackfillStatus | None:
        q = """SELECT twid, receiver, backfill_user, dispatched, message_count, state FROM backfill_status
        WHERE dispatched IS FALSE
        AND state < 2
        ORDER BY state ASC, message_count ASC
        LIMIT 1"""
        row = await cls.db.fetchrow(q)
        if not row:
            return None
        return cls(**row)
