# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Dict, Optional, Union

from attr import dataclass

from mautrix.types import SerializableAttrs, SerializableEnum
from mautrix.util import variation_selector

from .conversation import Conversation
from .util import StringTimestamp


class ReactionKey(SerializableEnum):
    FUNNY = "funny"  # 😂
    SURPRISED = "surprised"  # 😲
    SAD = "sad"  # 😢
    LIKE = "like"  # ❤
    EXCITED = "excited"  # 🔥
    AGREE = "agree"  # 👍
    DISAGREE = "disagree"  # 👎
    EMOJI = "emoji"  # arbitrary emoji

    @property
    def emoji(self) -> str:
        return _key_to_emoji[self]

    @classmethod
    def from_emoji(cls, emoji: str) -> Union["ReactionKey", str]:
        try:
            return _emoji_to_key[emoji.rstrip("\uFE0F")]
        except KeyError:
            return ReactionKey.EMOJI


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
class ReactionCreateEntry(SerializableAttrs):
    id: str
    time: StringTimestamp
    conversation_id: str
    message_id: str
    reaction_key: ReactionKey
    sender_id: str
    emoji_reaction: Optional[str] = None
    affects_sort: Optional[bool] = None
    request_id: Optional[str] = None

    conversation: Optional[Conversation] = None

    @property
    def reaction_emoji(self) -> str:
        if self.reaction_key == ReactionKey.EMOJI:
            return self.emoji_reaction
        return self.reaction_key.emoji


@dataclass
class ReactionDeleteEntry(SerializableAttrs):
    id: str
    time: StringTimestamp
    conversation_id: str
    message_id: str
    reaction_key: ReactionKey
    sender_id: str
    emoji_reaction: Optional[str] = None
    affects_sort: Optional[bool] = None

    conversation: Optional[Conversation] = None

    @property
    def reaction_emoji(self) -> str:
        if self.reaction_key == ReactionKey.EMOJI:
            emoji = self.emoji_reaction
        else:
            emoji = self.reaction_key.emoji
        return variation_selector.remove(emoji)
