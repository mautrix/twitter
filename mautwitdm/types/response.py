# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Optional, List, Dict

from attr import dataclass

from mautrix.types import SerializableAttrs

from .user import User
from .entry import Entry
from .conversation import Conversation, TimelineStatus


@dataclass
class SendResponse(SerializableAttrs['SendResponse']):
    entries: List[Entry]
    users: Dict[str, User]
    conversations: Optional[Dict[str, Conversation]] = None


@dataclass
class FetchConversationResponse(SerializableAttrs['FetchConversationResponse']):
    status: TimelineStatus
    min_entry_id: str
    max_entry_id: str
    entries: List[Entry]
    users: Dict[str, User]
    # Only present if include_info=True
    conversations: Optional[Dict[str, Conversation]] = None


@dataclass
class PollResponse(SerializableAttrs['PollResponse']):
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
class InboxTimeline(SerializableAttrs['InboxTimeline']):
    status: TimelineStatus
    min_entry_id: Optional[str] = None


@dataclass
class InboxTimelines(SerializableAttrs['InboxTimelines']):
    trusted: InboxTimeline
    untrusted: InboxTimeline
    untrusted_low_quality: InboxTimeline


@dataclass
class InitialStateResponse(SerializableAttrs['InitialStateResponse']):
    cursor: str
    last_seen_event_id: str
    trusted_last_seen_event_id: str
    untrusted_last_seen_event_id: str

    inbox_timelines: InboxTimelines
    entries: List[Entry]
    users: Dict[str, User]
    conversations: Dict[str, Conversation]


@dataclass
class ImageInfo(SerializableAttrs['ImageInfo']):
    image_type: str
    w: int
    h: int


@dataclass
class VideoInfo(SerializableAttrs['VideoInfo']):
    video_type: str


@dataclass
class MediaUploadResponse(SerializableAttrs['MediaUploadResponse']):
    media_id: int
    media_id_string: str
    media_key: str
    size: int
    expires_after_secs: int
    image: Optional[ImageInfo] = None
    video: Optional[VideoInfo] = None
