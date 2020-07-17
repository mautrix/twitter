# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Dict, Optional, Type, TypeVar, List, Callable, Awaitable, Any
import logging
import asyncio

from aiohttp import ClientSession
from yarl import URL

from .types import PollResponse, InitialStateResponse

T = TypeVar('T')
Handler = Callable[[T], Awaitable[Any]]
HandlerMap = Dict[Type[T], List[Handler]]


class TwitterPoller:
    dm_url: URL
    log: logging.Logger
    loop: asyncio.AbstractEventLoop
    http: ClientSession
    headers: Dict[str, str]

    poll_sleep: int = 1
    poll_cursor: Optional[str]
    _poll_task: Optional[asyncio.Task]
    _handlers: HandlerMap

    @property
    def poll_params(self) -> Dict[str, str]:
        return {
            "cards_platform": "Web-12",
            "include_cards": "1",
            "include_ext_alt_text": "true",
            "include_quote_count": "true",
            "include_reply_count": "1",
            "tweet_mode": "extended",
            "dm_users": "false",
            "include_groups": "true",
            "include_inbox_timelines": "true",
            "include_ext_media_color": "true",
            "support_reactions": "true",
            "ext": "mediaColor,altText,mediaStats,highlightedLabel",
        }

    @property
    def full_state_params(self) -> Dict[str, str]:
        return {
            "include_profile_interstitial_type": "1",
            "include_blocking": "1",
            "include_blocked_by": "1",
            "include_followed_by": "1",
            "include_want_retweets": "1",
            "include_mute_edge": "1",
            "include_can_dm": "1",
            "include_can_media_tag": "1",
            "skip_status": "1",
            **self.poll_params,
        }

    async def inbox_initial_state(self) -> InitialStateResponse:
        """
        Get the initial DM inbox state, including conversations, user profiles and some messages.

        This also gets the initial :attr:`poll_cursor` value.

        Returns:
            The response data from the server.
        """
        url = (self.dm_url / "inbox_initial_state.json").with_query({
            **self.full_state_params,
            "filter_low_quality": "true",
            "include_quality": "all",
        })
        async with self.http.get(url, headers=self.headers) as resp:
            data = await resp.json()
            response = InitialStateResponse.deserialize(data["inbox_initial_state"])
            self.poll_cursor = response.cursor
            return response

    async def _poll_once(self) -> PollResponse:
        if not self.poll_cursor:
            raise RuntimeError("Cursor must be set to poll")
        url = (self.dm_url / "user_updates.json").with_query({
            **self.poll_params,
            "cursor": self.poll_cursor,
            "filter_low_quality": "true",
            "include_quality": "all",
        })
        async with self.http.get(url, headers=self.headers) as resp:
            data = await resp.json()
            response = PollResponse.deserialize(data["user_events"])
            self.poll_cursor = response.cursor
            return response

    async def poll_forever(self, raise_exceptions: bool = True) -> None:
        """
        Poll for updates forever.

        Args:
            raise_exceptions: Whether or not fatal errors should be raised after logging.
                :class:`asyncio.CancelledError` will not be raised in any case.
        """
        try:
            await self._poll_forever()
        except asyncio.CancelledError:
            self.log.debug("Polling stopped")
        except Exception:
            self.log.exception("Fatal error while polling")
            if raise_exceptions:
                raise

    async def dispatch(self, event: T) -> None:
        """
        Dispatch an event to handlers registered with :meth:`add_handler`.

        Args:
            event: The event to dispatch.
        """
        for handler in self._handlers[type(event)]:
            await handler(event)

    async def _poll_forever(self) -> None:
        if not self.poll_cursor:
            self.log.debug("Poll cursor not set, calling initial state to get cursor")
            await self.inbox_initial_state()
        while True:
            try:
                resp = await self._poll_once()
            except Exception:
                self.log.warning("Error while polling", exc_info=True)
                await asyncio.sleep(self.poll_sleep * 5)
                continue
            for user in (resp.users or {}).values():
                await self.dispatch(user)
            for conversation in (resp.conversations or {}).values():
                await self.dispatch(conversation)
            for entry in resp.entries or []:
                if entry.trust_conversation:
                    await self.dispatch(entry.trust_conversation)
                if entry.message:
                    await self.dispatch(entry.message)
                if entry.reaction_delete:
                    await self.dispatch(entry.reaction_delete)
                if entry.reaction_create:
                    await self.dispatch(entry.reaction_create)
            await asyncio.sleep(self.poll_sleep)

    def start(self) -> asyncio.Task:
        """
        Start polling forever in the background. This calls :meth:`poll_forever` and puts it in an
        asyncio Task. The task is stored so it can be cancelled with :meth:`stop`.

        Returns:
            The created asyncio task.
        """
        self._poll_task = self.loop.create_task(self.poll_forever())
        return self._poll_task

    def stop(self) -> None:
        """Stop the ongoing poll task. Any ongoing handlers will also be cancelled."""
        if self._poll_task:
            self._poll_task.cancel()

    def add_handler(self, event_type: Type[T], handler: Handler) -> None:
        """
        Add an event handler.

        Args:
            event_type: The type of event to handle.
            handler: The handler function.
        """
        self._handlers[event_type].append(handler)

    def remove_handler(self, event_type: Type[T], handler: Handler) -> None:
        """
        Remove an event handler.

        Args:
            event_type: The type of event the handler was registered for.
            handler: The handler function to remove.
        """
        self._handlers[event_type].remove(handler)
