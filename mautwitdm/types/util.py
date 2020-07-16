# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import NewType
from datetime import datetime, timezone

from mautrix.types import deserializer, serializer

StringTimestamp = NewType('StringTimestamp', datetime)


@deserializer(StringTimestamp)
def deserialize_str_timestamp(val: str) -> StringTimestamp:
    return StringTimestamp(datetime.utcfromtimestamp(int(val) / 1000))


@serializer(StringTimestamp)
def serialize_str_timestamp(val: StringTimestamp) -> str:
    return str(val.timestamp() * 1000)


StringDateTime = NewType('StringDateTime', datetime)


@deserializer(StringDateTime)
def deserialize_str_datetime(val: str) -> StringDateTime:
    return StringDateTime(datetime.strptime(val, "%a %b %d %H:%M:%S +0000 %Y")
                          .replace(tzinfo=timezone.utc))


@serializer(StringDateTime)
def serialize_str_datetime(val: StringDateTime) -> str:
    return str(val.strftime("%a %b %d %H:%M:%S +0000 %Y"))
