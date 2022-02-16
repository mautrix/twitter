# mautrix-twitter - A Matrix-Twitter DM puppeting bridge
# Copyright (C) 2022 Tulir Asokan
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.
import html

from mautrix.types import Format, MessageType, TextMessageEventContent
from mautwitdm.types import MessageData, MessageEntityURL, MessageEntityUserMention

from . import puppet as pu


async def twitter_to_matrix(message: MessageData) -> TextMessageEventContent:
    content = TextMessageEventContent(
        msgtype=MessageType.TEXT,
        body=message.text,
        format=Format.HTML,
        formatted_body=message.text,
    )
    for entity in reversed(message.entities.all) if message.entities else []:
        start, end = entity.indices
        if isinstance(entity, MessageEntityURL):
            content.body = content.body[:start] + entity.expanded_url + content.body[end:]
            content.formatted_body = (
                f"{content.formatted_body[:start]}"
                f'<a href="{entity.expanded_url}">{entity.expanded_url}</a>'
                f"{content.formatted_body[end:]}"
            )
        elif isinstance(entity, MessageEntityUserMention):
            puppet = await pu.Puppet.get_by_twid(entity.id, create=False)
            if puppet:
                user_url = f"https://matrix.to/#/{puppet.mxid}"
                content.formatted_body = (
                    f"{content.formatted_body[:start]}"
                    f'<a href="{user_url}">{puppet.name or entity.name}</a>'
                    f"{content.formatted_body[end:]}"
                )
        else:
            # Get the sigil (# or $) from the body
            text = content.formatted_body[start:end][0] + entity.text
            content.formatted_body = (
                f"{content.formatted_body[:start]}"
                f'<font color="#3771bb">{text}</font>'
                f"{content.formatted_body[end:]}"
            )
    if content.formatted_body == content.body:
        content.formatted_body = None
        content.format = None
    content.body = html.unescape(content.body)
    return content
