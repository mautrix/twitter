# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Optional, Union, List

from attr import dataclass

from mautrix.types import SerializableAttrs

from .message import MessageEntry
from .reaction import ReactionDeleteEntry, ReactionCreateEntry
from .util import StringTimestamp


@dataclass
class TrustConversationEntry(SerializableAttrs['TrustConversationEntry']):
    id: str
    time: StringTimestamp
    affects_sort: bool
    conversation_id: str
    reason: str


@dataclass
class ConversationCreateEntry(SerializableAttrs['ConversationCreateEntry']):
    id: str
    time: StringTimestamp
    affects_sort: bool
    request_id: str
    conversation_id: str


@dataclass
class ConversationNameUpdateEntry(SerializableAttrs['ConversationNameUpdateEntry']):
    id: str
    time: StringTimestamp
    affects_sort: bool
    conversation_id: str
    conversation_name: str
    by_user_id: str


EntryType = Union[MessageEntry, TrustConversationEntry, ConversationCreateEntry,
                  ConversationNameUpdateEntry, ReactionCreateEntry, ReactionCreateEntry]


@dataclass
class Entry(SerializableAttrs['Entry']):
    message: Optional[MessageEntry] = None
    trust_conversation: Optional[TrustConversationEntry] = None
    conversation_create: Optional[ConversationCreateEntry] = None
    conversation_name_update: Optional[ConversationNameUpdateEntry] = None
    reaction_create: Optional[ReactionCreateEntry] = None
    reaction_delete: Optional[ReactionDeleteEntry] = None

    @property
    def all_types(self) -> List[EntryType]:
        items = (self.conversation_create, self.conversation_name_update, self.trust_conversation,
                 self.message, self.reaction_create, self.reaction_delete)
        return [item for item in items
                if item is not None]
