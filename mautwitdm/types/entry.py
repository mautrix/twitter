# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Optional

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
class Entry(SerializableAttrs['Entry']):
    message: Optional[MessageEntry] = None
    trust_conversation: Optional[TrustConversationEntry] = None
    reaction_create: Optional[ReactionCreateEntry] = None
    reaction_delete: Optional[ReactionDeleteEntry] = None
