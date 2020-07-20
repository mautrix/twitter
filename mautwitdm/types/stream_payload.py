# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Optional, Union, List

from attr import dataclass

from mautrix.types import SerializableAttrs


@dataclass
class DMTypingEvent(SerializableAttrs['DMTypingEvent']):
    user_id: str
    conversation_id: str


@dataclass
class DMUpdateEvent(SerializableAttrs['DMUpdateEvent']):
    user_id: str
    conversation_id: str


@dataclass
class SubscriptionError(SerializableAttrs['SubscriptionError']):
    topic: str
    code: int
    message: str


@dataclass
class SubscriptionsEvent(SerializableAttrs['SubscriptionsEvent']):
    errors: List[SubscriptionError]


StreamEventType = Union[SubscriptionsEvent, DMTypingEvent, DMUpdateEvent]


@dataclass
class StreamEvent(SerializableAttrs['Event']):
    dm_typing: Optional[DMTypingEvent] = None
    dm_update: Optional[DMUpdateEvent] = None
    subscriptions: Optional[SubscriptionsEvent] = None

    @property
    def all_types(self) -> List[StreamEventType]:
        items = (self.dm_typing, self.dm_update)
        return [item for item in items
                if item is not None]
