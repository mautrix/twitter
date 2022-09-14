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
from asyncpg import Connection

from mautrix.util.async_db import Scheme, UpgradeTable

upgrade_table = UpgradeTable()


@upgrade_table.register(description="Latest revision", upgrades_to=4)
async def upgrade_latest(conn: Connection, scheme: Scheme) -> None:
    if scheme in (Scheme.POSTGRES, Scheme.COCKROACH):
        await conn.execute("CREATE TYPE twitter_conv_type AS ENUM ('ONE_TO_ONE', 'GROUP_DM')")
        await conn.execute(
            """CREATE TYPE twitter_reaction_key AS ENUM (
                'funny', 'surprised', 'sad', 'like', 'excited', 'agree', 'disagree'
            )"""
        )
    await conn.execute(
        """CREATE TABLE portal (
            twid        TEXT,
            receiver    BIGINT,
            conv_type   twitter_conv_type NOT NULL,
            other_user  BIGINT,
            mxid        TEXT,
            name        TEXT,
            encrypted   BOOLEAN NOT NULL DEFAULT false,

            PRIMARY KEY (twid, receiver)
        )"""
    )
    await conn.execute(
        """CREATE TABLE "user" (
            mxid        TEXT PRIMARY KEY,
            twid        BIGINT,
            auth_token  TEXT,
            csrf_token  TEXT,
            poll_cursor TEXT,
            notice_room TEXT
        )"""
    )
    await conn.execute(
        """CREATE TABLE puppet (
            twid      BIGINT PRIMARY KEY,
            name      TEXT,
            photo_url TEXT,
            photo_mxc TEXT,

            is_registered BOOLEAN NOT NULL DEFAULT false,

            custom_mxid  TEXT,
            access_token TEXT,
            next_batch   TEXT,
            base_url     TEXT
        )"""
    )
    await conn.execute(
        """CREATE TABLE message (
            mxid     TEXT NOT NULL,
            mx_room  TEXT NOT NULL,
            twid     BIGINT,
            receiver BIGINT,

            PRIMARY KEY (twid, receiver),
            UNIQUE (mxid, mx_room)
        )"""
    )
    await conn.execute(
        """CREATE TABLE reaction (
            mxid        TEXT NOT NULL,
            mx_room     TEXT NOT NULL,
            tw_msgid    BIGINT,
            tw_receiver BIGINT,
            tw_sender   BIGINT,
            reaction    twitter_reaction_key NOT NULL,

            tw_reaction_id BIGINT,

            PRIMARY KEY (tw_msgid, tw_receiver, tw_sender),
            FOREIGN KEY (tw_msgid, tw_receiver) REFERENCES message(twid, receiver)
                ON DELETE CASCADE ON UPDATE CASCADE,
            UNIQUE (mxid, mx_room)
        )"""
    )


@upgrade_table.register(description="Add double-puppeting base_url to puppet table")
async def upgrade_v2(conn: Connection) -> None:
    await conn.execute("ALTER TABLE puppet ADD COLUMN base_url TEXT")


@upgrade_table.register(description="Store Twitter reaction IDs for marking things read")
async def upgrade_v3(conn: Connection) -> None:
    await conn.execute("ALTER TABLE reaction ADD COLUMN tw_reaction_id BIGINT")


@upgrade_table.register(description="Replace VARCHAR(255) with TEXT")
async def upgrade_v4(conn: Connection) -> None:
    tables = {
        "portal": ("twid", "mxid", "name"),
        "user": ("mxid", "auth_token", "csrf_token", "poll_cursor", "notice_room"),
        "puppet": ("name", "photo_url", "photo_mxc", "custom_mxid", "next_batch"),
        "message": ("mxid", "mx_room"),
        "reaction": ("mxid", "mx_room"),
    }
    for table, columns in tables.items():
        for column in columns:
            await conn.execute(f'ALTER TABLE "{table}" ALTER COLUMN "{column}" TYPE TEXT')
    await conn.execute("DROP TABLE user_portal")
