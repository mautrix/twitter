# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import List, Tuple

from attr import dataclass

from mautrix.types import SerializableAttrs


@dataclass
class MessageEntityURL(SerializableAttrs['MessageEntityURL']):
    url: str
    expanded_url: str
    display_url: str
    indices: Tuple[int, int]


@dataclass
class MessageEntityUserMention(SerializableAttrs['MessageEntityUserMention']):
    screen_name: str
    name: str
    id: int
    id_str: str
    indices: Tuple[int, int]


@dataclass
class MessageEntitySimple(SerializableAttrs['MessageEntitySimple']):
    text: str
    indices: Tuple[int, int]


@dataclass
class MessageEntities(SerializableAttrs['MessageEntities']):
    hashtags: List[MessageEntitySimple]
    symbols: List[MessageEntitySimple]
    user_mentions: List[MessageEntityUserMention]
    urls: List[MessageEntityURL]
