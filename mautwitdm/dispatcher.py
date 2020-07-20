# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.'
from typing import Dict, Type, TypeVar, List, Callable, Awaitable, Any
import logging

T = TypeVar('T')
Handler = Callable[[T], Awaitable[Any]]
HandlerMap = Dict[Type[T], List[Handler]]


class TwitterDispatcher:
    """
    This class is used to dispatch events that :class:`TwitterPoller` and :class:`TwitterStreamer`
    get from Twitter.
    """
    log: logging.Logger
    _handlers: HandlerMap

    async def dispatch(self, event: T) -> None:
        """
        Dispatch an event to handlers registered with :meth:`add_handler`.

        Args:
            event: The event to dispatch.
        """
        for handler in self._handlers[type(event)]:
            try:
                await handler(event)
            except Exception:
                self.log.exception(f"Error while handling event of type {type(event)}")

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
