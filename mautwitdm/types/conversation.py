# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import List, Optional

from attr import dataclass

from mautrix.types import SerializableAttrs, SerializableEnum


class ConversationType(SerializableEnum):
    ONE_TO_ONE = "ONE_TO_ONE"


@dataclass
class Participant(SerializableAttrs['Participant']):
    user_id: str
    last_read_event_id: Optional[str] = None


@dataclass
class Conversation(SerializableAttrs['Conversation']):
    conversation_id: str
    type: ConversationType
    sort_event_id: str
    sort_timestamp: str
    participants: List[Participant]
    notifications_disabled: bool
    mention_notifications_disabled: bool
    read_only: bool
    trusted: bool
    low_quality: bool
