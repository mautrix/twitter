# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Dict, Optional

from attr import dataclass

from mautrix.types import SerializableAttrs, SerializableEnum

from .util import StringTimestamp
from .conversation import Conversation


class ReactionKey(SerializableEnum):
    FUNNY = "funny"  # 😂
    SURPRISED = "surprised"  # 😲
    SAD = "sad"  # 😢
    LIKE = "like"  # ❤
    EXCITED = "excited"  # 🔥
    AGREE = "agree"  # 👍
    DISAGREE = "disagree"  # 👎

    @property
    def emoji(self) -> str:
        return _key_to_emoji[self]

    @classmethod
    def from_emoji(cls, emoji: str) -> 'ReactionKey':
        try:
            return _emoji_to_key[emoji.rstrip("\uFE0F")]
        except KeyError:
            raise ValueError(f"Unsupported reaction emoji {emoji}")


_key_to_emoji: Dict[ReactionKey, str] = {
    ReactionKey.FUNNY: "\U0001F602",
    ReactionKey.SURPRISED: "\U0001F632",
    ReactionKey.SAD: "\U0001F622",
    ReactionKey.LIKE: "\u2764\uFE0F",
    ReactionKey.EXCITED: "\U0001F525",
    ReactionKey.AGREE: "\U0001F44D\uFE0F",
    ReactionKey.DISAGREE: "\U0001F44E\uFE0F",
}

_emoji_to_key: Dict[str, ReactionKey] = {
    "\U0001F44D": ReactionKey.AGREE,
    "\U0001F44E": ReactionKey.DISAGREE,
    "\u2764": ReactionKey.LIKE,
    "\U0001F525": ReactionKey.EXCITED,
    "\U0001F622": ReactionKey.SAD,
    "\U0001F632": ReactionKey.SURPRISED,
    "\U0001F62E": ReactionKey.SURPRISED,  # This is 😮, which is commonly used for surprised
    "\U0001F602": ReactionKey.FUNNY,
}


@dataclass
class ReactionCreateEntry(SerializableAttrs['ReactionCreateEntry']):
    id: str
    time: StringTimestamp
    conversation_id: str
    message_id: str
    reaction_key: ReactionKey
    sender_id: str
    affects_sort: Optional[bool] = None

    conversation: Optional[Conversation] = None


@dataclass
class ReactionDeleteEntry(SerializableAttrs['ReactionDeleteEntry']):
    id: str
    time: StringTimestamp
    conversation_id: str
    message_id: str
    reaction_key: ReactionKey
    sender_id: str
    affects_sort: Optional[bool] = None

    conversation: Optional[Conversation] = None
