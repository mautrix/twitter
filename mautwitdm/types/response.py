# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Dict, List, Optional

from attr import dataclass
import attr

from mautrix.types import SerializableAttrs

from .conversation import Conversation, TimelineStatus
from .entry import Entry
from .user import User


@dataclass
class SendResponse(SerializableAttrs):
    entries: List[Entry]
    users: Dict[str, User]
    conversations: Optional[Dict[str, Conversation]] = None


@dataclass
class FetchConversationResponse(SerializableAttrs):
    status: TimelineStatus
    min_entry_id: str
    max_entry_id: str
    entries: List[Entry]
    users: Dict[str, User]
    # Only present if include_info=True
    conversations: Optional[Dict[str, Conversation]] = None


@dataclass
class PollResponse(SerializableAttrs):
    cursor: str
    last_seen_event_id: str
    trusted_last_seen_event_id: str
    untrusted_last_seen_event_id: str

    min_entry_id: Optional[str] = None
    max_entry_id: Optional[str] = None
    entries: Optional[List[Entry]] = None
    users: Optional[Dict[str, User]] = None
    conversations: Optional[Dict[str, Conversation]] = None


@dataclass
class InboxTimeline(SerializableAttrs):
    status: TimelineStatus
    min_entry_id: Optional[str] = None


@dataclass
class InboxTimelines(SerializableAttrs):
    trusted: InboxTimeline
    untrusted: InboxTimeline
    untrusted_low_quality: InboxTimeline


@dataclass
class InitialStateResponse(SerializableAttrs):
    cursor: str
    last_seen_event_id: str
    trusted_last_seen_event_id: str
    untrusted_last_seen_event_id: str

    inbox_timelines: InboxTimelines
    entries: List[Entry] = attr.ib(factory=lambda: [])
    users: Dict[str, User] = attr.ib(factory=lambda: {})
    conversations: Dict[str, Conversation] = attr.ib(factory=lambda: {})


@dataclass
class ImageInfo(SerializableAttrs):
    image_type: str
    w: int
    h: int


@dataclass
class VideoInfo(SerializableAttrs):
    video_type: str


@dataclass
class MediaUploadResponse(SerializableAttrs):
    media_id: int
    media_id_string: str
    media_key: str
    size: int
    expires_after_secs: int
    image: Optional[ImageInfo] = None
    video: Optional[VideoInfo] = None
