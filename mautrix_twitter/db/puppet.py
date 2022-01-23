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
from yarl import URL
import asyncpg

from mautrix.types import ContentURI, SyncToken, UserID
from mautrix.util.async_db import Database

fake_db = Database.create("") if TYPE_CHECKING else None


@dataclass
class Puppet:
    db: ClassVar[Database] = fake_db

    twid: int
    name: str | None
    photo_url: str | None
    photo_mxc: ContentURI | None

    is_registered: bool

    custom_mxid: UserID | None
    access_token: str | None
    next_batch: SyncToken | None
    base_url: URL | None

    @property
    def _values(self):
        return (
            self.twid,
            self.name,
            self.photo_url,
            self.photo_mxc,
            self.is_registered,
            self.custom_mxid,
            self.access_token,
            self.next_batch,
            str(self.base_url) if self.base_url else None,
        )

    async def insert(self) -> None:
        q = (
            "INSERT INTO puppet (twid, name, photo_url, photo_mxc, is_registered, custom_mxid,"
            "                    access_token, next_batch, base_url) "
            "VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
        )
        await self.db.execute(q, *self._values)

    async def update(self) -> None:
        q = (
            "UPDATE puppet SET name=$2, photo_url=$3, photo_mxc=$4, is_registered=$5,"
            "                  custom_mxid=$6, access_token=$7, next_batch=$8, base_url=$9 "
            "WHERE twid=$1"
        )
        await self.db.execute(q, *self._values)

    @classmethod
    def _from_row(cls, row: asyncpg.Record) -> Puppet | None:
        if not row:
            return None
        data = {**row}
        base_url_str = data.pop("base_url")
        base_url = URL(base_url_str) if base_url_str is not None else None
        return cls(base_url=base_url, **data)

    @classmethod
    async def get_by_twid(cls, twid: int) -> Puppet | None:
        q = (
            "SELECT twid, name, photo_url, photo_mxc, is_registered,"
            "       custom_mxid, access_token, next_batch, base_url "
            "FROM puppet WHERE twid=$1"
        )
        return cls._from_row(await cls.db.fetchrow(q, twid))

    @classmethod
    async def get_by_custom_mxid(cls, mxid: UserID) -> Puppet | None:
        q = (
            "SELECT twid, name, photo_url, photo_mxc, is_registered,"
            "       custom_mxid, access_token, next_batch, base_url "
            "FROM puppet WHERE custom_mxid=$1"
        )
        return cls._from_row(await cls.db.fetchrow(q, mxid))

    @classmethod
    async def all_with_custom_mxid(cls) -> list[Puppet]:
        q = (
            "SELECT twid, name, photo_url, photo_mxc, is_registered,"
            "       custom_mxid, access_token, next_batch, base_url "
            "FROM puppet WHERE custom_mxid IS NOT NULL"
        )
        return [cls._from_row(row) for row in await cls.db.fetch(q)]
