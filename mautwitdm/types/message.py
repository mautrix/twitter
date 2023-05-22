# Copyright (c) 2022 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import List, Optional

from attr import dataclass
import attr

from mautrix.types import SerializableAttrs

from .conversation import Conversation
from .message_attachment import MessageAttachment
from .message_entity import MessageEntities
from .reaction import ReactionCreateEntry
from .util import StringTimestamp


@dataclass
class MessageData(SerializableAttrs):
    id: str
    time: StringTimestamp
    sender_id: str
    text: str
    recipient_id: Optional[str] = None
    conversation_id: Optional[str] = None
    entities: Optional[MessageEntities] = None
    attachment: Optional[MessageAttachment] = None
    reply_data: Optional["MessageData"] = None


# Resolve reply_data field
attr.resolve_types(MessageData)


@dataclass
class MessageEntry(SerializableAttrs):
    id: str
    time: StringTimestamp
    conversation_id: str
    message_data: MessageData
    message_reactions: Optional[List[ReactionCreateEntry]] = None
    affects_sort: Optional[bool] = None
    request_id: Optional[str] = None

    @property
    def sender_id(self) -> str:
        return self.message_data.sender_id

    conversation: Optional[Conversation] = None
