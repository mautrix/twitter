# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Optional

from attr import dataclass

from mautrix.types import SerializableAttrs

from .util import StringDateTime


@dataclass
class User(SerializableAttrs['User']):
    id: int
    id_str: str
    name: str
    screen_name: str
    profile_image_url: Optional[str]
    profile_image_url_https: Optional[str]
    description: str
    created_at: StringDateTime
    verified: bool
    protected: bool

    can_media_tag: bool
    following: bool
    follow_request_sent: bool
    blocking: bool

    friends_count: int
    followers_count: int

    # Note: the DM /new.json response users have more fields that /user_updates.json doesn't.
    # Those aren't that important so they're listed here. If they're needed, a new FullUser class
    # could be added.
