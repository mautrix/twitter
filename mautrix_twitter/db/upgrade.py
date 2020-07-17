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

from mautrix.util.async_db import UpgradeTable

upgrade_table = UpgradeTable()


@upgrade_table.register(description="Initial revision")
async def upgrade_v1(conn: Connection) -> None:
    await conn.execute("CREATE TYPE twitter_conv_type AS ENUM ('ONE_TO_ONE', 'GROUP_DM')")
    await conn.execute("""CREATE TABLE portal (
        twid        VARCHAR(255),
        receiver    BIGINT,
        conv_type   twitter_conv_type NOT NULL,
        other_user  BIGINT,
        mxid        VARCHAR(255),
        name        VARCHAR(255),
        encrypted   BOOLEAN NOT NULL DEFAULT false,

        PRIMARY KEY (twid, receiver)
    )""")
    await conn.execute("""CREATE TABLE "user" (
        mxid        VARCHAR(255) PRIMARY KEY,
        twid        BIGINT,
        auth_token  VARCHAR(255),
        csrf_token  VARCHAR(255),
        poll_cursor VARCHAR(255),
        notice_room VARCHAR(255)
    )""")
    await conn.execute("""CREATE TABLE puppet (
        twid      BIGINT PRIMARY KEY,
        name      VARCHAR(255),
        photo_url VARCHAR(255),
        photo_mxc VARCHAR(255),

        is_registered BOOLEAN NOT NULL DEFAULT false,

        custom_mxid  VARCHAR(255),
        access_token TEXT,
        next_batch   VARCHAR(255)
    )""")
    await conn.execute("""CREATE TABLE user_portal (
        "user"          BIGINT,
        portal          VARCHAR(255),
        portal_receiver BIGINT,
        in_community    BOOLEAN NOT NULL DEFAULT false,

        FOREIGN KEY (portal, portal_receiver) REFERENCES portal(twid, receiver)
            ON UPDATE CASCADE ON DELETE CASCADE
    )""")
    await conn.execute("""CREATE TABLE message (
        mxid     VARCHAR(255) NOT NULL,
        mx_room  VARCHAR(255) NOT NULL,
        twid     BIGINT,
        receiver BIGINT,

        PRIMARY KEY (twid, receiver),
        UNIQUE (mxid, mx_room)
    )""")
    await conn.execute("CREATE TYPE twitter_reaction_key AS ENUM ('funny', 'surprised', 'sad',"
                       "                                          'like', 'excited', 'agree',"
                       "                                          'disagree')")
    await conn.execute("""CREATE TABLE reaction (
        mxid        VARCHAR(255) NOT NULL,
        mx_room     VARCHAR(255) NOT NULL,
        tw_msgid    BIGINT,
        tw_receiver BIGINT,
        tw_sender   BIGINT,
        reaction    twitter_reaction_key NOT NULL,

        PRIMARY KEY (tw_msgid, tw_receiver, tw_sender),
        FOREIGN KEY (tw_msgid, tw_receiver) REFERENCES message(twid, receiver)
            ON DELETE CASCADE ON UPDATE CASCADE,
        UNIQUE (mxid, mx_room)
    )""")
