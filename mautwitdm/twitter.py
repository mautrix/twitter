# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Optional, Dict, Any
from uuid import UUID, uuid1, getnode
from collections import defaultdict
from http.cookies import SimpleCookie
import logging
import asyncio

from aiohttp import ClientSession
from yarl import URL

from .conversation import Conversation
from .uploader import TwitterUploader
from .poller import TwitterPoller


class TwitterAPI(TwitterUploader, TwitterPoller):
    """The main entrypoint for using the internal Twitter DM API."""
    base_url: URL = URL("https://api.twitter.com/1.1")
    dm_url: URL = base_url / "dm"

    loop: asyncio.AbstractEventLoop
    http: ClientSession
    log: logging.Logger

    node_id: int
    active: bool

    _csrf_token: str

    def __init__(self, http: Optional[ClientSession] = None, log: Optional[logging.Logger] = None,
                 loop: Optional[asyncio.AbstractEventLoop] = None, node_id: Optional[int] = None
                 ) -> None:
        self.loop = loop or asyncio.get_event_loop()
        self.http = http or ClientSession(loop=self.loop)
        self.log = log or logging.getLogger("mautwitdm")
        self.node_id = node_id or getnode()
        self.poll_cursor = None
        self._handlers = defaultdict(lambda: [])
        self.active = True

    def set_tokens(self, auth_token: str, csrf_token: str) -> None:
        """
        Set the authentication tokens. After this, use :meth:`get_user_identifier` to check if the
        auth is working correctly.

        Args:
            auth_token: The auth_token cookie value.
            csrf_token: The ct0 cookie/x-csrf-token header value.
        """
        cookie = SimpleCookie()
        cookie["auth_token"] = auth_token
        cookie["auth_token"].update({"domain": "twitter.com", "path": "/"})
        cookie["ct0"] = csrf_token
        cookie["ct0"].update({"domain": "twitter.com", "path": "/"})
        self._csrf_token = csrf_token
        self.http.cookie_jar.update_cookies(cookie, URL("https://twitter.com/"))

    @property
    def headers(self) -> Dict[str, str]:
        """
        Get the headers to use with every request to Twitter.

        Returns:
            A key-value HTTP header list.
        """
        return {
            # Hardcoded authorization header from the web app
            "authorization": "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs"
                             "%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA",
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) "
                          "Gecko/20100101 Firefox/78.0",
            "Accept": "*/*",
            "Accept-Language": "en-US,en;q=0.5",
            "DNT": "1",
            "Origin": "https://twitter.com",
            "Referer": "https://twitter.com/messages",
            "x-twitter-auth-type": "OAuth2Session",
            "x-twitter-client-language": "en",
            "x-twitter-active-user": "yes",
            "x-csrf-token": self._csrf_token,
        }

    def new_request_id(self) -> UUID:
        """
        Create a new request ID for DM send requests.

        Returns:
            A v1 UUID.
        """
        return uuid1(self.node_id)

    def conversation(self, id: str) -> Conversation:
        return Conversation(self, id)

    async def update_last_seen_event_id(self, last_seen_event_id: str) -> None:
        await self.http.post(self.dm_url / "update_last_seen_event_id.json",
                             data={"last_seen_event_id": last_seen_event_id,
                                   "trusted_last_seen_event_id": last_seen_event_id},
                             headers=self.headers)

    async def get_user_identifier(self) -> Optional[str]:
        async with self.http.post(self.base_url / "branch" / "init.json", json={},
                                  headers=self.headers) as resp:
            resp.raise_for_status()
            return (await resp.json()).get("user_identifier", None)

    async def get_settings(self) -> Dict[str, Any]:
        """Get the account settings of the currently logged in account."""
        async with self.http.get(self.base_url / "account" / "settings.json",
                                 headers=self.headers) as resp:
            resp.raise_for_status()
            return await resp.json()
