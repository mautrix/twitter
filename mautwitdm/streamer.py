# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import AsyncGenerator, Set, Dict
import logging
import asyncio
import json
import io

from aiohttp import ClientSession
from yarl import URL

from .types import StreamEvent
from .dispatcher import TwitterDispatcher
from .errors import check_error


class TwitterStreamer(TwitterDispatcher):
    """
    This class handles receiving live events like typing notifications
    via ``/live_pipeline/events``.
    """
    pipeline_url = URL("https://api.twitter.com/live_pipeline/events")
    pipeline_update_url = URL("https://api.twitter.com/1.1/live_pipeline/update_subscriptions")

    log: logging.Logger
    loop: asyncio.AbstractEventLoop
    http: ClientSession
    headers: Dict[str, str]
    user_agent: str

    topics: Set[str]
    _stream_task: asyncio.Task

    async def _stream(self) -> AsyncGenerator[StreamEvent, None]:
        url = self.pipeline_url.with_query(topics=",".join(self.topics))
        headers = {
            "User-Agent": self.user_agent,
            "Accept": "text/event-stream",
            "Accept-Language": "en-US,en;q=0.5",
            "DNT": "1",
            "Origin": "https://twitter.com",
            "Referer": "https://twitter.com/messages",
            "Pragma": "no-cache",
        }
        empty_bytes = b""
        data_prefix = b"data: "
        chunk_separator = b"\n\n"
        empty_chunk = b":\n\n"
        async with self.http.get(url, headers=headers) as resp:
            buffer = io.BytesIO()
            partial = False
            while True:
                chunk, end = await resp.content.readchunk()
                if not end and chunk == empty_bytes:
                    break
                if not chunk.endswith(chunk_separator):
                    buffer.write(chunk)
                    partial = True
                    continue
                if partial:
                    chunk = buffer.getvalue() + chunk
                    buffer = io.BytesIO()
                    partial = False
                if chunk == empty_chunk:
                    continue
                elif chunk.startswith(data_prefix):
                    data = json.loads(chunk[len(data_prefix):-len(chunk_separator)])
                    yield StreamEvent.deserialize(data["payload"])

    async def update_topics(self, subscribe: Set[str], unsubscribe: Set[str]) -> None:
        """
        Update the topics the client is subscribed to.

        Args:
            subscribe: The list of topics to subscribe to.
            unsubscribe: The list of topics to unsubscribe from.
        """
        subscribe = subscribe - self.topics
        unsubscribe = unsubscribe & self.topics
        self.topics = self.topics - unsubscribe | subscribe
        if self._stream_task is None or self._stream_task.done():
            self.log.debug("Not sending update_subscriptions request: no ongoing stream task")
            return
        url = self.pipeline_update_url
        req = {"sub_topics": ",".join(subscribe), "unsub_topics": ",".join(unsubscribe)}
        async with self.http.post(url, data=req, headers=self.headers) as resp:
            resp_data = await check_error(resp)
        event = StreamEvent.deserialize(resp_data)
        for event_type in event.all_types:
            await self.dispatch(event_type)

    async def stream_forever(self, raise_exceptions: bool = True) -> None:
        """
        Stream for events from the server forever. Events are dispatched using :meth:`dispatch`,
        which means handlers can be added with :meth:`add_handler`.

        You can set which updates to listen to by setting :attribute:`topics` before calling this,
        or by calling :meth:`update_topics` while streaming.

        Args:
            raise_exceptions: Whether or not errors should be raised after logging. If set to
                ``False``, errors will be logged, and the stream will retry in 10 seconds.
                :class:`asyncio.CancelledError` will not be raised and will simply return.
        """
        while True:
            try:
                async for event in self._stream():
                    for event_type in event.all_types:
                        await self.dispatch(event_type)
            except asyncio.CancelledError:
                self.log.debug("Streaming stopped")
                break
            except Exception:
                self.log.exception("Error while streaming events")
                if raise_exceptions:
                    raise
                await asyncio.sleep(10)

    def start_streaming(self) -> asyncio.Task:
        """
        Start polling forever in the background. This calls :meth:`stream_forever` and puts it in
        an asyncio Task. The task is stored so it can be cancelled with :meth:`stop_streaming`.

        Returns:
            The created asyncio task.
        """
        self.log.debug("Starting streaming task")
        self._stream_task = self.loop.create_task(self.stream_forever())
        return self._stream_task

    def stop_streaming(self) -> None:
        """Stop the ongoing stream task. Any ongoing handlers will also be cancelled."""
        if self._stream_task:
            self.log.debug("Cancelling ongoing stream task")
            self._stream_task.cancel()
