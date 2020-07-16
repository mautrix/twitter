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
from mautrix.types import RoomID, UserID

from .version import version, linkified_version
from .config import Config


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

    def prepare_db(self) -> None:
        pass

    async def get_user(self, user_id: UserID) -> 'BaseUser':
        pass

    async def get_portal(self, room_id: RoomID) -> 'BasePortal':
        pass

    async def get_puppet(self, user_id: UserID, create: bool = False) -> 'BasePuppet':
        pass

    async def get_double_puppet(self, user_id: UserID) -> 'BasePuppet':
        pass


TwitterBridge().run()
