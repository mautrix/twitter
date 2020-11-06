# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Dict, Optional, Type, TypeVar, List, Callable, Awaitable, Any, Union
import logging
import asyncio

from attr import dataclass
from aiohttp import ClientSession
from yarl import URL

from .types import PollResponse, InitialStateResponse
from .errors import check_error
from .conversation import Conversation
from .dispatcher import TwitterDispatcher

T = TypeVar('T')
Handler = Callable[[T], Awaitable[Any]]
HandlerMap = Dict[Type[T], List[Handler]]


class PollingStarted:
    pass


class PollingStopped:
    pass


@dataclass
class PollingErrored:
    error: Exception
    fatal: bool
    count: int


class PollingErrorResolved:
    pass


class TwitterPoller(TwitterDispatcher):
    """This class handles polling for new messages using ``/dm/user_updates.json``."""
    dm_url: URL

    log: logging.Logger
    loop: asyncio.AbstractEventLoop
    http: ClientSession
    headers: Dict[str, str]
    skip_poll_wait: asyncio.Event

    poll_sleep: int = 3
    error_sleep: int = 5
    max_poll_errors: int = 12
    poll_cursor: Optional[str]
    dispatch_initial_resp: bool
    _poll_task: Optional[asyncio.Task]
    _typing_in: Optional[Conversation]

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
            "supports_reactions": "true",
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

    async def inbox_initial_state(self, set_poll_cursor: bool = True) -> InitialStateResponse:
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
            data = await check_error(resp)
            response = InitialStateResponse.deserialize(data["inbox_initial_state"])
            if set_poll_cursor:
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
            data = await check_error(resp)
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
            await self.dispatch(PollingStopped())
        except Exception as e:
            await self.dispatch(PollingErrored(e, fatal=True, count=0))
            self.log.exception("Fatal error while polling")
            if raise_exceptions:
                raise

    async def dispatch_all(self, resp: Union[PollResponse, InitialStateResponse]) -> None:
        for user in (resp.users or {}).values():
            await self.dispatch(user)
        for conversation in (resp.conversations or {}).values():
            await self.dispatch(conversation)
        for entry in resp.entries or []:
            if not entry:
                continue
            for entry_type in entry.all_types:
                entry_type.conversation = resp.conversations[entry_type.conversation_id]
                await self.dispatch(entry_type)

    async def _poll_forever(self) -> None:
        if not self.poll_cursor:
            self.log.debug("Poll cursor not set, calling initial state to get cursor")
            resp = await self.inbox_initial_state()
            if self.dispatch_initial_resp:
                await self.dispatch_all(resp)
        await self.dispatch(PollingStarted())
        errors = 0
        while True:
            try:
                resp = await self._poll_once()
            except Exception as e:
                if errors > self.max_poll_errors > 0:
                    self.log.debug(f"Error count ({errors}) exceeded maximum, "
                                   f"raising error as fatal")
                    raise
                errors += 1
                sleep = min(15 * 60, self.error_sleep * errors)
                self.log.warning(f"Error while polling, retrying in {sleep}s", exc_info=True)
                await self.dispatch(PollingErrored(e, fatal=False, count=errors))
                await asyncio.sleep(sleep)
                continue
            if errors > 0:
                errors = 0
                await self.dispatch(PollingErrorResolved())
            await self.dispatch_all(resp)
            try:
                await asyncio.wait_for(self.skip_poll_wait.wait(), self.poll_sleep)
            except asyncio.TimeoutError:
                pass
            if self._typing_in:
                await self._typing_in.mark_typing()

    def is_polling(self) -> bool:
        return self._poll_task and not self._poll_task.done()

    def start_polling(self) -> asyncio.Task:
        """
        Start polling forever in the background. This calls :meth:`poll_forever` and puts it in an
        asyncio Task. The task is stored so it can be cancelled with :meth:`stop_polling`.

        Returns:
            The created asyncio task.
        """
        self.stop_polling()
        self.log.debug("Starting poll task")
        self._poll_task = self.loop.create_task(self.poll_forever())
        return self._poll_task

    def stop_polling(self) -> None:
        """Stop the ongoing poll task. Any ongoing handlers will also be cancelled."""
        if self._poll_task:
            self.log.debug("Cancelling ongoing poll task")
            self._poll_task.cancel()
            self._poll_task = None
