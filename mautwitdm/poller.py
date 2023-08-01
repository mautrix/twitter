# Copyright (c) 2022 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from __future__ import annotations

from typing import Dict
import asyncio
import logging
import time

from aiohttp import ClientSession
from attr import dataclass
from yarl import URL

from . import conversation as c
from .dispatcher import TwitterDispatcher
from .errors import RateLimitError, TwitterAuthError, check_error
from .types import (
    Conversation,
    InboxTimeline,
    InitialStateResponse,
    PollResponse,
    TimelineStatus,
    User,
)


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
    headers: dict[str, str]
    skip_poll_wait: asyncio.Event

    poll_sleep: int = 3
    error_sleep: int = 5
    max_poll_errors: int = 12
    max_poll_auth_errors: int = 4
    poll_cursor: str | None
    dispatch_initial_resp: bool
    _poll_task: asyncio.Task | None
    _typing_in: c.Conversation | None

    @property
    def poll_query_params(self) -> dict[str, str]:
        return {
            "include_ext_alt_text": "true",
            "include_reply_count": "1",
            "tweet_mode": "extended",
            "include_groups": "true",
            "include_inbox_timelines": "true",
            "include_ext_media_color": "true",
            "supports_reactions": "true",
            "include_ext_limited_action_results": "false",
            "ext": ",".join(
                [
                    "mediaColor",
                    "altText",
                    "mediaStats",
                    "highlightedLabel",
                    "hasNftAvatar",
                    "voiceInfo",
                    "birdwatchPivot",
                    "enrichments",
                    "superFollowMetadata",
                    "unmentionInfo",
                    "editControl",
                    "vibe",
                ]
            ),
        }

    @property
    def poll_params(self) -> dict[str, str]:
        return {
            "cards_platform": "Web-12",
            "include_cards": "1",
            "include_quote_count": "true",
            "dm_users": "false",
            "include_ext_edit_control": "true",
            "nsfw_filtering_enabled": "false",
            "dm_secret_conversations_enabled": "false",
            "include_ext_views": "true",
            **self.poll_query_params,
        }

    @property
    def full_state_params(self) -> dict[str, str]:
        return {
            "include_profile_interstitial_type": "1",
            "include_blocking": "1",
            "include_blocked_by": "1",
            "include_followed_by": "1",
            "include_want_retweets": "1",
            "include_mute_edge": "1",
            "include_can_dm": "1",
            "include_can_media_tag": "1",
            "include_ext_alt_text": "1",
            "include_ext_has_nft_avatar": "1",
            "skip_status": "1",
            **self.poll_params,
        }

    async def all_trusted_conversations(
        self,
    ) -> tuple[Dict[str, Conversation], Dict[str, User]]:
        """
        Get all trusted conversations from the inbox (using pagination).

        Returns:
            Dictionary containing all conversations, keyed by their ID.
        """
        initial_state = await self.inbox_initial_state(set_poll_cursor=False)
        conversations = initial_state.conversations
        users = initial_state.users

        if initial_state.inbox_timelines.trusted.status == TimelineStatus.AT_END:
            return conversations, users

        min_entry_id = initial_state.inbox_timelines.trusted.min_entry_id

        while True:
            self.log.debug("Not at end, fetching more conversations with max_id %s", min_entry_id)
            inbox_timeline = await self.inbox_timeline("trusted", min_entry_id)

            if inbox_timeline.conversations is not None:
                conversations = {**conversations, **inbox_timeline.conversations}

            if inbox_timeline.users is not None:
                users = {**users, **inbox_timeline.users}

            if inbox_timeline.status == TimelineStatus.AT_END:
                return conversations, users

            if inbox_timeline.min_entry_id is not None:
                min_entry_id = inbox_timeline.min_entry_id
            else:
                return conversations, users

    async def inbox_timeline(self, inbox: str, max_id: str) -> InboxTimeline:
        """
        Gets a page of DMs from one of the possible inbox timelines.

        Returns:
            The response data from the server.
        """
        url = (self.dm_url / "inbox_timeline" / f"{inbox}.json").with_query(
            {**self.full_state_params, "max_id": max_id}
        )
        async with self.http.get(url, headers=self.headers) as resp:
            data = await check_error(resp)
            response = InboxTimeline.deserialize(data["inbox_timeline"])
            return response

    async def inbox_initial_state(self, set_poll_cursor: bool = True) -> InitialStateResponse:
        """
        Get the initial DM inbox state, including conversations, user profiles and some messages.

        This also gets the initial :attr:`poll_cursor` value.

        Returns:
            The response data from the server.
        """
        url = (self.dm_url / "inbox_initial_state.json").with_query(
            {
                **self.full_state_params,
                "filter_low_quality": "true",
                "include_quality": "all",
            }
        )
        async with self.http.get(url, headers=self.headers) as resp:
            data = await check_error(resp)
            response = InitialStateResponse.deserialize(data["inbox_initial_state"])
            if set_poll_cursor:
                self.poll_cursor = response.cursor
            return response

    async def _poll_once(self) -> PollResponse:
        if not self.poll_cursor:
            raise RuntimeError("Cursor must be set to poll")
        url = (self.dm_url / "user_updates.json").with_query(
            {
                **self.poll_params,
                "cursor": self.poll_cursor,
                "filter_low_quality": "true",
                "include_quality": "all",
            }
        )
        async with self.http.get(url, headers=self.headers) as resp:
            data = await check_error(resp)
            try:
                user_events = data["user_events"]
            except KeyError:
                try:
                    inbox_initial_state = data["inbox_initial_state"]
                except KeyError:
                    self.log.warning("Got unknown poll response: %s", data)
                    raise
                else:
                    response = InitialStateResponse.deserialize(inbox_initial_state)
            else:
                response = PollResponse.deserialize(user_events)
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
        except TwitterAuthError as e:
            self.log.exception("Auth error while polling")
            await self.dispatch(PollingErrored(e, fatal=True, count=0))
        except Exception as e:
            await self.dispatch(PollingErrored(e, fatal=True, count=0))
            self.log.exception("Fatal error while polling")
            if raise_exceptions:
                raise

    async def dispatch_all(self, resp: PollResponse | InitialStateResponse) -> None:
        for user in (resp.users or {}).values():
            await self.dispatch(user)
        for conversation in (resp.conversations or {}).values():
            await self.dispatch(conversation)
        for entry in resp.entries or []:
            if not entry:
                continue
            for entry_type in entry.all_types:
                try:
                    entry_type.conversation = resp.conversations[entry_type.conversation_id]
                except KeyError:
                    msg = (
                        "Poll response didn't contain conversation info "
                        f"for {entry_type.conversation_id} "
                        f"(entry type {type(entry_type)}, ID {entry_type.id}"
                    )
                    if entry_type.conversation:
                        self.log.debug(f"{msg}, but the entry had its own conversation info")
                        await self.dispatch(entry_type)
                    else:
                        # TODO we could still dispatch the event and guess the conversation type
                        #      in event handlers if necessary.
                        self.log.warning(msg)
                else:
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
            except RateLimitError as e:
                sleep = e.reset - int(time.time())
                self.log.warning(
                    f"Got rate limit until {e.reset} while polling, waiting {sleep} seconds"
                )
                if self.poll_sleep < 8:
                    self.poll_sleep += 1
                    self.log.debug(f"Increased poll sleep to {self.poll_sleep}")
                await asyncio.sleep(sleep)
                continue
            except Exception as e:
                max_errors = (
                    self.max_poll_auth_errors
                    if isinstance(e, TwitterAuthError)
                    else self.max_poll_errors
                )
                if errors > max_errors > 0:
                    self.log.debug(
                        f"Error count ({errors}) exceeded maximum, raising error as fatal"
                    )
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
        self._poll_task = asyncio.create_task(self.poll_forever())
        return self._poll_task

    def stop_polling(self) -> None:
        """Stop the ongoing poll task. Any ongoing handlers will also be cancelled."""
        if self._poll_task:
            self.log.debug("Cancelling ongoing poll task")
            self._poll_task.cancel()
            self._poll_task = None
