# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Any
import json

from aiohttp import ClientResponse


class TwitterError(Exception):
    code: int
    message: str

    def __init__(self, code: int, message: str) -> None:
        super().__init__(f"{code}: {message}")
        self.code = code
        self.message = message


async def check_error(resp: ClientResponse) -> Any:
    try:
        resp_data = await resp.json()
    except json.JSONDecodeError:
        resp.raise_for_status()
        raise

    if not isinstance(resp_data, dict) or "errors" not in resp_data:
        resp.raise_for_status()
        return resp_data

    try:
        error = resp_data["errors"][0]
        raise TwitterError(error["code"], error["message"])
    except KeyError:
        resp.raise_for_status()
        raise
