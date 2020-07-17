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

    friends_count: int
    followers_count: int

    statuses_count: Optional[int] = None

    followed_by: Optional[bool] = None
    blocking: Optional[bool] = None
    suspended: Optional[bool] = None
    bocked_by: Optional[bool] = None

    url: Optional[str] = None
    location: Optional[str] = None
