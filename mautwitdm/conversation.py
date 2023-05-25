# Copyright (c) 2022 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from __future__ import annotations

from uuid import UUID

from yarl import URL

from . import twitter as tw
from .errors import check_error
from .types import FetchConversationResponse, ReactionKey, SendResponse


class Conversation:
    api: tw.TwitterAPI
    id: str

    def __init__(self, api: tw.TwitterAPI, id: str) -> None:
        self.api = api
        self.id = id

    @property
    def api_url(self) -> URL:
        """The base URL for API requests related to this conversation."""
        return self.api.dm_url / "conversation" / self.id

    async def mark_read(self, last_read_event_id: int | str) -> None:
        """Mark the conversation as read, up to the given event ID."""
        req = {"conversationId": self.id, "last_read_event_id": str(last_read_event_id)}
        url = self.api_url / "mark_read.json"
        async with self.api.http.post(url, headers=self.api.headers, data=req) as resp:
            await check_error(resp)

    async def mark_typing(self) -> None:
        """Send a typing notification. This request should be repeated every 2-3 seconds."""
        async with self.api.http.post(self.api_url / "typing.json", headers=self.api.headers) as r:
            await check_error(r)

    async def accept(self) -> None:
        """Accept the conversation (when DMing users who you don't follow)."""
        async with self.api.http.post(self.api_url / "accept.json", headers=self.api.headers) as r:
            await check_error(r)

    async def send(
        self,
        text: str,
        media_id: str | int | None = None,
        voice_message: bool = False,
        reply_to_id: str | int | None = None,
        request_id: UUID | str | None = None,
    ) -> SendResponse:
        """
        Send a message to this conversation.

        Args:
            text: The text to send. May be an empty string if only sending media.
            media_id: The media ID to send. Use :meth:`TwitterUploader.upload` to upload media and
                get a media ID.
            voice_message: Whether the media message is a voice message.
            reply_to_id: ID of message to reply to.
            request_id: The transaction ID for this request. It will be included in the message
                when polling and can be used for deduplication. If not provided, one will be
                automatically created using :meth:`TwitterAPI.new_request_id`

        Returns:
            The send response from the server.
        """
        data = {
            "cards_platform": "Web-12",
            "conversation_id": self.id,
            "dm_users": False,
            "include_cards": 1,
            "include_quote_count": True,
            "recipient_ids": False,
            "request_id": str(request_id or self.api.new_request_id()),
            "text": text,
        }
        url = (self.api.dm_url / "new2.json").with_query(self.api.poll_query_params)
        if reply_to_id:
            data["reply_to_dm_id"] = str(reply_to_id)
        if media_id:
            data["media_id"] = str(media_id)
            if voice_message:
                data["audio_only_media_attachment"] = True
        async with self.api.http.post(url, json=data, headers=self.api.headers) as resp:
            resp_data = await check_error(resp)
            return SendResponse.deserialize(resp_data)

    @staticmethod
    def _reaction_key_to_params(emoji: str) -> dict[str, str]:
        key = ReactionKey.from_emoji(emoji)
        if key == ReactionKey.EMOJI:
            return {
                "emoji_reaction": emoji,
                "reaction_key": str(key),
            }
        else:
            return {
                "reaction_key": str(key),
            }

    async def react(self, message_id: str | int, emoji: str) -> None:
        """
        React to a message. Reacting to the same message multiple times will override earlier
        reactions.

        Args:
            message_id: The message ID to react to.
            emoji: The reaction itself.
        """
        query = {
            "conversation_id": self.id,
            "dm_id": str(message_id),
            **self._reaction_key_to_params(emoji),
        }
        url = (self.api.dm_url / "reaction" / "new.json").with_query(query)
        async with self.api.http.post(url, headers=self.api.headers) as resp:
            await check_error(resp)

    async def delete_reaction(self, message_id: str | int, emoji: str) -> None:
        """
        Delete an earlier reaction.

        Args:
            message_id: The message ID to react to.
            emoji: The reaction itself.
        """
        url = (self.api.dm_url / "reaction" / "delete.json").with_query(
            {
                "conversation_id": self.id,
                "dm_id": str(message_id),
                **self._reaction_key_to_params(emoji),
            }
        )
        async with self.api.http.post(url, headers=self.api.headers) as resp:
            await check_error(resp)

    async def fetch(self, max_id: str | None = None) -> FetchConversationResponse:
        """
        Fetch the conversation metadata and message history.

        Args:
            max_id: The maximum message ID to fetch.

        Returns:
            The requested metadata and message history.
        """
        query = {
            **self.api.full_state_params,
            "context": "FETCH_DM_CONVERSATION",
            "include_conversation_info": "true",
        }
        req = (self.api.dm_url / "conversation" / f"{self.id}.json").with_query(query)
        if max_id:
            req = req.update_query({"max_id": max_id})
        async with self.api.http.get(req, headers=self.api.headers) as resp:
            resp_data = await check_error(resp)
        data = resp_data["conversation_timeline"]
        if "entries" not in data:
            data["entries"] = None
        if "min_entry_id" not in data:
            data["min_entry_id"] = None
        if "max_entry_id" not in data:
            data["max_entry_id"] = None
        return FetchConversationResponse.deserialize(data)
