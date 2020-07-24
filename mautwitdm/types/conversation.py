# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import List, Optional

from attr import dataclass

from mautrix.types import SerializableAttrs, SerializableEnum

from .util import StringTimestamp


class ConversationType(SerializableEnum):
    ONE_TO_ONE = "ONE_TO_ONE"
    GROUP_DM = "GROUP_DM"


class TimelineStatus(SerializableEnum):
    AT_END = "AT_END"
    HAS_MORE = "HAS_MORE"


@dataclass
class Participant(SerializableAttrs['Participant']):
    user_id: str

    # This seems to be only for one-to-one chats
    last_read_event_id: Optional[str] = None

    # These seem to be only for group chats
    join_time: Optional[StringTimestamp] = None
    join_conversation_event_id: Optional[str] = None
    is_admin: Optional[bool] = None


@dataclass
class Conversation(SerializableAttrs['Conversation']):
    conversation_id: str
    type: ConversationType
    sort_event_id: str
    sort_timestamp: StringTimestamp
    participants: List[Participant]
    notifications_disabled: bool
    mention_notifications_disabled: bool
    trusted: bool
    low_quality: bool

    # These are present in some responses
    min_entry_id: Optional[str] = None
    max_entry_id: Optional[str] = None
    status: Optional[TimelineStatus] = None

    # This seems to be only for one-to-one chats
    read_only: Optional[bool] = None

    # These seem to be only for group chats
    create_time: Optional[StringTimestamp] = None
    created_by_user_id: Optional[str] = None
    name: Optional[str] = None
    last_read_event_id: Optional[str] = None
