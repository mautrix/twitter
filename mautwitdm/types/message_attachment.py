# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Optional, List, Dict, Tuple

from attr import dataclass

from mautrix.types import SerializableAttrs

from .user import User
from .message_entity import MessageEntities
from .util import StringDateTime


@dataclass
class ImageSizeInfo(SerializableAttrs['ImageSizeInfo']):
    resize: str
    w: int
    h: int


@dataclass
class OriginalImageSizeInfo(SerializableAttrs['OriginalImageSizeInfo']):
    width: int
    height: int


@dataclass
class VideoVariant(SerializableAttrs['VideoVariant']):
    content_type: str
    url: str
    bitrate: Optional[int] = None


@dataclass
class VideoInfo(SerializableAttrs['VideoInfo']):
    aspect_ratio: Tuple[int, int]
    variants: List[VideoVariant]
    duration_millis: Optional[int] = None


@dataclass
class MessageAttachmentMedia(SerializableAttrs['MessageAttachmentMedia']):
    id: int
    id_str: str
    indices: Tuple[int, int]
    media_url: str
    media_url_https: str
    url: str
    display_url: str
    expanded_url: str
    type: str
    sizes: Dict[str, ImageSizeInfo]
    original_info: OriginalImageSizeInfo
    video_info: Optional[VideoInfo] = None


@dataclass
class ImageBindingValue(SerializableAttrs['ImageBindingValue']):
    url: str
    width: int
    height: int


@dataclass
class CardBindingValue(SerializableAttrs['CardBindingValue']):
    type: str
    string_value: Optional[str] = None
    image_value: Optional[ImageBindingValue] = None
    scribe_key: Optional[str] = None


@dataclass
class MessageAttachmentCard(SerializableAttrs['MessageAttachmentCard']):
    name: str
    url: str
    binding_values: Dict[str, CardBindingValue]
    users: Optional[Dict[str, User]] = None


@dataclass
class MessageEntityMedia(MessageAttachmentMedia, SerializableAttrs['MessageEntityMedia']):
    indices: Optional[Tuple[int, int]] = None


@dataclass
class ExtendedMessageEntities(SerializableAttrs['ExtendedMessageEntities']):
    media: List[MessageEntityMedia]


@dataclass
class TweetAttachmentStatus(SerializableAttrs['TweetAttachmentStatus']):
    created_at: StringDateTime
    id: int
    id_str: str
    full_text: str
    truncated: bool
    display_text_range: Tuple[int, int]
    source: str
    user: User
    in_reply_to_status_id: Optional[int]
    in_reply_to_status_id_str: Optional[str]
    in_reply_to_user_id: Optional[int]
    in_reply_to_user_id_str: Optional[str]
    in_reply_to_screen_name: Optional[str]
    is_quote_status: bool
    retweet_count: int
    favorite_count: int
    reply_count: int
    quote_count: int
    favorited: bool
    retweeted: bool
    lang: Optional[str] = None
    supplemental_language: Optional[str] = None
    possibly_sensitive: Optional[bool] = None
    possibly_sensitive_editable: Optional[bool] = None
    entities: Optional[MessageEntities] = None
    extended_entities: Optional[ExtendedMessageEntities] = None


@dataclass
class MessageAttachmentTweet(SerializableAttrs['MessageAttachmentTweet']):
    id: str
    url: str
    display_url: str
    expanded_url: str
    status: TweetAttachmentStatus
    indices: Optional[Tuple[int, int]] = None


@dataclass
class MessageAttachment(SerializableAttrs['MessageAttachment']):
    photo: Optional[MessageAttachmentMedia] = None
    video: Optional[MessageAttachmentMedia] = None
    animated_gif: Optional[MessageAttachmentMedia] = None
    tweet: Optional[MessageAttachmentTweet] = None

    @property
    def media(self) -> Optional[MessageAttachmentMedia]:
        if self.video:
            return self.video
        if self.animated_gif:
            return self.animated_gif
        if self.photo:
            return self.photo
