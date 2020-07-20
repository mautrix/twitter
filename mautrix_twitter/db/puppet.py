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

from mautrix.types import UserID, SyncToken, ContentURI
from mautrix.util.async_db import Database

fake_db = Database("") if TYPE_CHECKING else None


@dataclass
class Puppet:
    db: ClassVar[Database] = fake_db

    twid: int
    name: Optional[str]
    photo_url: Optional[str]
    photo_mxc: Optional[ContentURI]

    is_registered: bool

    custom_mxid: Optional[UserID]
    access_token: Optional[str]
    next_batch: Optional[SyncToken]

    async def insert(self) -> None:
        q = ("INSERT INTO puppet (twid, name, photo_url, photo_mxc, is_registered, custom_mxid,"
             "                    access_token, next_batch) "
             "VALUES ($1, $2, $3, $4, $5, $6, $7, $8)")
        await self.db.execute(q, self.twid, self.name, self.photo_url, self.photo_mxc,
                              self.is_registered, self.custom_mxid, self.access_token,
                              self.next_batch)

    async def update(self) -> None:
        q = ("UPDATE puppet SET name=$2, photo_url=$3, photo_mxc=$4, is_registered=$5,"
             "                  custom_mxid=$6, access_token=$7, next_batch=$8 WHERE twid=$1")
        await self.db.execute(q, self.twid, self.name, self.photo_url, self.photo_mxc,
                              self.is_registered, self.custom_mxid, self.access_token,
                              self.next_batch)

    @classmethod
    async def get_by_twid(cls, twid: int) -> Optional['Puppet']:
        row = await cls.db.fetchrow("SELECT twid, name, photo_url, photo_mxc, is_registered,"
                                    "       custom_mxid, access_token, next_batch "
                                    "FROM puppet WHERE twid=$1", twid)
        if not row:
            return None
        return cls(**row)

    @classmethod
    async def get_by_custom_mxid(cls, mxid: UserID) -> Optional['Puppet']:
        row = await cls.db.fetchrow("SELECT twid, name, photo_url, photo_mxc, is_registered,"
                                    "       custom_mxid, access_token, next_batch "
                                    "FROM puppet WHERE custom_mxid=$1", mxid)
        if not row:
            return None
        return cls(**row)

    @classmethod
    async def all_with_custom_mxid(cls) -> List['Puppet']:
        rows = await cls.db.fetch("SELECT twid, name, photo_url, photo_mxc, is_registered,"
                                  "       custom_mxid, access_token, next_batch "
                                  "FROM puppet WHERE custom_mxid IS NOT NULL")
        return [cls(**row) for row in rows]
