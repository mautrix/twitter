# Copyright (c) 2021 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import List, Tuple, Union

from attr import dataclass

from mautrix.types import SerializableAttrs


@dataclass
class MessageEntityURL(SerializableAttrs):
    url: str
    expanded_url: str
    display_url: str
    indices: Tuple[int, int]


@dataclass
class MessageEntityUserMention(SerializableAttrs):
    screen_name: str
    name: str
    id: int
    id_str: str
    indices: Tuple[int, int]


@dataclass
class MessageEntitySimple(SerializableAttrs):
    text: str
    indices: Tuple[int, int]


@dataclass
class MessageEntities(SerializableAttrs):
    hashtags: List[MessageEntitySimple]
    symbols: List[MessageEntitySimple]
    user_mentions: List[MessageEntityUserMention]
    urls: List[MessageEntityURL]

    @property
    def all(
        self,
    ) -> List[Union[MessageEntitySimple, MessageEntityUserMention, MessageEntityURL]]:
        entities = self.hashtags + self.symbols + self.user_mentions + self.urls
        return sorted(entities, key=lambda entity: entity.indices[0])
