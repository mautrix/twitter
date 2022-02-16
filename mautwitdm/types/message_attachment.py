# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Dict, List, Optional, Tuple, Union

from attr import dataclass

from mautrix.types import SerializableAttrs

from .message_entity import MessageEntities
from .user import User
from .util import StringDateTime


@dataclass
class ImageSizeInfo(SerializableAttrs):
    resize: str
    w: int
    h: int


@dataclass
class OriginalImageSizeInfo(SerializableAttrs):
    width: int
    height: int


@dataclass
class VideoVariant(SerializableAttrs):
    content_type: str
    url: str
    bitrate: Optional[int] = None


@dataclass
class VideoInfo(SerializableAttrs):
    aspect_ratio: Tuple[int, int]
    variants: List[VideoVariant]
    duration_millis: Optional[int] = None


@dataclass
class MessageAttachmentMedia(SerializableAttrs):
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
class ImageBindingValue(SerializableAttrs):
    url: str
    width: int
    height: int
    alt: Optional[str]


@dataclass
class RGB(SerializableAttrs):
    red: int
    green: int
    blue: int


@dataclass
class ImageColor(SerializableAttrs):
    percentage: float
    rgb: RGB


@dataclass
class ImageColorBindingValue(SerializableAttrs):
    palette: List[ImageColor]


@dataclass
class CardBindingValue(SerializableAttrs):
    type: str
    string_value: Optional[str] = None
    image_value: Optional[ImageBindingValue] = None
    image_color_value: Optional[ImageColorBindingValue] = None
    scribe_key: Optional[str] = None


@dataclass
class MessageAttachmentCard(SerializableAttrs):
    name: str
    url: str
    binding_values: Dict[str, CardBindingValue]
    users: Optional[Dict[str, User]] = None

    def string(self, key: str, default: Optional[str] = None) -> Optional[str]:
        try:
            val = self.binding_values[key].string_value
            return val if val is not None else default
        except KeyError:
            return default

    def image(self, key: str) -> Optional[ImageBindingValue]:
        try:
            return self.binding_values[key].image_value
        except KeyError:
            return None


@dataclass
class MessageEntityMedia(MessageAttachmentMedia, SerializableAttrs):
    indices: Optional[Tuple[int, int]] = None


@dataclass
class ExtendedMessageEntities(SerializableAttrs):
    media: List[MessageEntityMedia]


@dataclass
class TweetAttachmentStatus(SerializableAttrs):
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
    card: Optional[MessageAttachmentCard] = None


@dataclass
class MessageAttachmentTweet(SerializableAttrs):
    id: str
    url: str
    display_url: str
    expanded_url: str
    status: TweetAttachmentStatus
    indices: Optional[Tuple[int, int]] = None


@dataclass
class MessageAttachment(SerializableAttrs):
    photo: Optional[MessageAttachmentMedia] = None
    video: Optional[MessageAttachmentMedia] = None
    animated_gif: Optional[MessageAttachmentMedia] = None
    tweet: Optional[MessageAttachmentTweet] = None
    card: Optional[MessageAttachmentCard] = None

    @property
    def media(self) -> Optional[MessageAttachmentMedia]:
        if self.video:
            return self.video
        if self.animated_gif:
            return self.animated_gif
        if self.photo:
            return self.photo
        return None

    @property
    def url_preview(self) -> Optional[Union[MessageAttachmentTweet, MessageAttachmentCard]]:
        if self.tweet:
            return self.tweet
        if self.card:
            return self.card
        return None
