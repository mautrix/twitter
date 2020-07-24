# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Optional, Union, TYPE_CHECKING
from uuid import UUID

from yarl import URL

from .types import ReactionKey, SendResponse, FetchConversationResponse
from .errors import check_error

if TYPE_CHECKING:
    from .twitter import TwitterAPI


class Conversation:
    api: 'TwitterAPI'
    id: str

    def __init__(self, api: 'TwitterAPI', id: str) -> None:
        self.api = api
        self.id = id

    @property
    def api_url(self) -> URL:
        """The base URL for API requests related to this conversation."""
        return self.api.dm_url / "conversation" / self.id

    async def mark_read(self, last_read_event_id: Union[int, str]) -> None:
        """Mark the conversation as read, up to the given event ID."""
        req = {"conversationId": self.id,
               "last_read_event_id": str(last_read_event_id)}
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

    async def send(self, text: str, media_id: Optional[Union[str, int]] = None,
                   request_id: Optional[Union[UUID, str]] = None) -> SendResponse:
        """
        Send a message to this conversation.

        Args:
            text: The text to send. May be an empty string if only sending media.
            media_id: The media ID to send. Use :meth:`TwitterUploader.upload` to upload media and
                get a media ID.
            request_id: The transaction ID for this request. It will be included in the message
                when polling and can be used for deduplication. If not provided, one will be
                automatically created using :meth:`TwitterAPI.new_request_id`

        Returns:
            The send response from the server.
        """
        data = {
            **self.api.poll_params,
            "text": text,
            "conversation_id": self.id,
            "recipient_ids": "false",
            "request_id": str(request_id or self.api.new_request_id()),
        }
        url = self.api.dm_url / "new.json"
        if media_id:
            data["media_id"] = str(media_id)
        async with self.api.http.post(url, data=data, headers=self.api.headers) as resp:
            resp_data = await check_error(resp)
            return SendResponse.deserialize(resp_data)

    async def react(self, message_id: Union[str, int], key: ReactionKey) -> None:
        """
        React to a message. Reacting to the same message multiple times will override earlier
        reactions.

        Args:
            message_id: The message ID to react to.
            key: The reaction itself.
        """
        url = (self.api.dm_url / "reaction" / "new.json").with_query({
            "reaction_key": str(key),
            "conversation_id": self.id,
            "dm_id": str(message_id),
        })
        async with self.api.http.post(url, headers=self.api.headers) as resp:
            await check_error(resp)

    async def delete_reaction(self, message_id: Union[str, int], key: ReactionKey) -> None:
        """
        Delete an earlier reaction.

        Args:
            message_id: The message ID to react to.
            key: The reaction itself.
        """
        url = (self.api.dm_url / "reaction" / "delete.json").with_query({
            "reaction_key": str(key),
            "conversation_id": self.id,
            "dm_id": str(message_id),
        })
        async with self.api.http.post(url, headers=self.api.headers) as resp:
            await check_error(resp)

    async def fetch(self, max_id: Optional[str] = None, include_info: bool = True
                    ) -> FetchConversationResponse:
        """
        Fetch the conversation metadata and message history.

        Args:
            max_id: The maximum message ID to fetch.
            include_info: Whether or not to include conversation metadata in the response.

        Returns:
            The requested metadata and message history.
        """
        query = {
            **self.api.full_state_params,
            "context": "FETCH_DM_CONVERSATION",
        }
        if include_info:
            query["include_conversation_info"] = "true"
        req = (self.api.dm_url / "conversation" / f"{self.id}.json").with_query(query)
        if max_id:
            req = req.update_query({"max_id": max_id})
        async with self.api.http.get(req, headers=self.api.headers) as resp:
            resp_data = await check_error(resp)
        return FetchConversationResponse.deserialize(resp_data["conversation_timeline"])
