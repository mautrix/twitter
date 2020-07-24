# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Optional

from attr import dataclass

from mautrix.types import SerializableAttrs

from .message_entity import MessageEntities
from .message_attachment import MessageAttachment
from .conversation import Conversation
from .util import StringTimestamp


@dataclass
class MessageData(SerializableAttrs['MessageData']):
    id: str
    time: StringTimestamp
    sender_id: str
    text: str
    recipient_id: Optional[str] = None
    conversation_id: Optional[str] = None
    entities: Optional[MessageEntities] = None
    attachment: Optional[MessageAttachment] = None


@dataclass
class MessageEntry(SerializableAttrs['MessageEntry']):
    id: str
    time: StringTimestamp
    conversation_id: str
    message_data: MessageData
    affects_sort: Optional[bool] = None
    request_id: Optional[str] = None

    @property
    def sender_id(self) -> str:
        return self.message_data.sender_id

    conversation: Optional[Conversation] = None
