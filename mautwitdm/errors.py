# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Any
import json
import time

from aiohttp import ClientResponse, ContentTypeError
from multidict import CIMultiDictProxy


class TwitterError(Exception):
    code: int
    message: str

    def __init__(self, code: int, message: str) -> None:
        super().__init__(f"{code}: {message}")
        self.code = code
        self.message = message


class RateLimitError(TwitterError):
    limit: int
    remaining: int
    reset: int

    def __init__(self, code: int, message: str, headers: CIMultiDictProxy[str]) -> None:
        self.code = code
        self.message = message
        # TODO make sure this works
        print(headers)
        self.limit = int(headers["x-rate-limit-limit"])
        # self.remaining = int(headers["x-rate-limit-remaining"])
        self.reset = int(headers["x-rate-limit-reset"])
        time_till_reset = int(time.time() - self.reset)
        Exception.__init__(self, f"Rate limit exceeded. Will reset in {time_till_reset} seconds "
                                 f"(endpoint is limited to {self.limit} requests in 15 minutes)")


async def check_error(resp: ClientResponse) -> Any:
    try:
        resp_data = await resp.json()
    except ContentTypeError:
        resp.raise_for_status()
        return
    except json.JSONDecodeError:
        resp.raise_for_status()
        raise

    if not isinstance(resp_data, dict) or "errors" not in resp_data:
        resp.raise_for_status()
        return resp_data

    try:
        error = resp_data["errors"][0]
        code = error["code"]
        message = error["message"]
    except (KeyError, IndexError):
        resp.raise_for_status()
        raise
    if code == 88:
        raise RateLimitError(code, message, resp.headers)
    raise TwitterError(code, message)
