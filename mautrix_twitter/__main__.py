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
from mautrix.bridge import Bridge, BaseUser, BasePuppet, BasePortal
from mautrix.bridge.state_store.asyncpg import PgBridgeStateStore
from mautrix.types import RoomID, UserID
from mautrix.util.async_db import Database

from .version import version, linkified_version
from .config import Config
from .db import upgrade_table, init as init_db
from .matrix import MatrixHandler
from .user import User
from .portal import Portal
from .puppet import Puppet


class TwitterBridge(Bridge):
    module = "mautrix_twitter"
    name = "mautrix-twitter"
    command = "python -m mautrix-twitter"
    description = "A Matrix-Twitter DM puppeting bridge."
    repo_url = "https://github.com/tulir/mautrix-twitter"
    real_user_content_key = "net.maunium.twitter.puppet"
    version = version
    markdown_version = linkified_version
    config_class = Config
    matrix_class = MatrixHandler

    db: Database
    config: Config
    state_store: PgBridgeStateStore

    def make_state_store(self) -> None:
        self.state_store = PgBridgeStateStore(self.db, self.get_puppet, self.get_double_puppet)

    def prepare_db(self) -> None:
        self.db = Database(self.config["appservice.database"], upgrade_table=upgrade_table,
                           loop=self.loop)
        init_db(self.db)

    async def start(self) -> None:
        await self.db.start()
        await self.state_store.upgrade_table.upgrade(self.db.pool)
        self.add_startup_actions(await User.init_cls(self))
        self.add_startup_actions(await Puppet.init_cls(self))
        Portal.init_cls(self)
        await super().start()

    async def get_user(self, user_id: UserID) -> 'BaseUser':
        pass

    async def get_portal(self, room_id: RoomID) -> 'BasePortal':
        pass

    async def get_puppet(self, user_id: UserID, create: bool = False) -> 'BasePuppet':
        pass

    async def get_double_puppet(self, user_id: UserID) -> 'BasePuppet':
        pass


TwitterBridge().run()
