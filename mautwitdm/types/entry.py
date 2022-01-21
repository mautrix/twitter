# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import List, Optional, Union

from attr import dataclass

from mautrix.types import SerializableAttrs

from .conversation import Conversation
from .message import MessageEntry
from .reaction import ReactionCreateEntry, ReactionDeleteEntry
from .util import StringTimestamp


@dataclass
class TrustConversationEntry(SerializableAttrs):
    id: str
    time: StringTimestamp
    conversation_id: str
    reason: str
    affects_sort: Optional[bool] = None

    conversation: Optional[Conversation] = None


@dataclass
class ConversationCreateEntry(SerializableAttrs):
    id: str
    time: StringTimestamp
    conversation_id: str
    affects_sort: Optional[bool] = None
    request_id: Optional[str] = None

    conversation: Optional[Conversation] = None


@dataclass
class ConversationNameUpdateEntry(SerializableAttrs):
    id: str
    time: StringTimestamp
    conversation_id: str
    conversation_name: str
    by_user_id: str
    affects_sort: Optional[bool] = None

    conversation: Optional[Conversation] = None


@dataclass
class ConversationReadEntry(SerializableAttrs):
    id: str
    last_read_event_id: str
    time: StringTimestamp
    conversation_id: str
    affects_sort: Optional[bool] = None

    conversation: Optional[Conversation] = None


EntryType = Union[
    MessageEntry,
    TrustConversationEntry,
    ConversationCreateEntry,
    ConversationNameUpdateEntry,
    ReactionCreateEntry,
    ReactionCreateEntry,
    ConversationReadEntry,
]


@dataclass
class Entry(SerializableAttrs):
    message: Optional[MessageEntry] = None
    trust_conversation: Optional[TrustConversationEntry] = None
    conversation_create: Optional[ConversationCreateEntry] = None
    conversation_name_update: Optional[ConversationNameUpdateEntry] = None
    reaction_create: Optional[ReactionCreateEntry] = None
    reaction_delete: Optional[ReactionDeleteEntry] = None
    conversation_read: Optional[ConversationReadEntry] = None

    @property
    def all_types(self) -> List[EntryType]:
        items = (
            self.conversation_create,
            self.conversation_name_update,
            self.trust_conversation,
            self.message,
            self.reaction_create,
            self.reaction_delete,
            self.conversation_read,
        )
        return [item for item in items if item is not None]
