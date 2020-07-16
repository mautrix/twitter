# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Optional, Union, TYPE_CHECKING
from uuid import UUID

from yarl import URL

from .types import ReactionKey, SendResponse, FetchConversationResponse

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
        return self.api.dm_url / "conversation" / self.id

    async def mark_read(self, last_read_event_id: str) -> None:
        """Mark the conversation as read, up to the given event ID."""
        await self.api.http.post(self.api_url / "mark_read.json",
                                 headers=self.api.headers,
                                 json={"conversationId": self.id,
                                       "last_read_event_id": last_read_event_id})

    async def mark_typing(self) -> None:
        """Send a typing notification. This request should be repeated every 2-3 seconds."""
        await self.api.http.post(self.api_url / "typing.json", headers=self.api.headers)

    async def accept(self) -> None:
        """Accept the conversation (when DMing users who you don't follow)."""
        await self.api.http.post(self.api_url / "accept.json", headers=self.api.headers)

    async def send(self, text: str, media_id: Optional[str] = None,
                   request_id: Optional[Union[UUID, str]] = None) -> SendResponse:
        data = {
            **self.api.poll_params,
            "text": text,
            "conversation_id": self.id,
            "recipient_ids": "false",
            "request_id": str(request_id or self.api.new_request_id()),
        }
        url = self.api.dm_url / "new.json"
        if media_id:
            data["media_id"] = media_id
        async with self.api.http.post(url, data=data, headers=self.api.headers) as resp:
            resp_data = await resp.json()
            return SendResponse.deserialize(resp_data)

    async def react(self, key: ReactionKey, message_id: str) -> None:
        url = (self.api.dm_url / "reaction" / "new.json").with_query({
            "reaction_key": str(key),
            "conversation_id": self.id,
            "dm_id": message_id,
        })
        async with self.api.http.post(url, headers=self.api.headers) as resp:
            resp.raise_for_status()

    async def fetch(self, max_id: Optional[str] = None) -> FetchConversationResponse:
        req = (self.api.dm_url / "conversation" / f"{self.id}.json").with_query({
            **self.api.full_state_params,
            "include_conversation_info": "true",
            "context": "FETCH_DM_CONVERSATION",
        })
        if max_id:
            req = req.update_query({"max_id": max_id})
        async with self.api.http.get(req, headers=self.api.headers) as resp:
            resp.raise_for_status()
            resp_data = await resp.json()
        return FetchConversationResponse.deserialize(resp_data)
