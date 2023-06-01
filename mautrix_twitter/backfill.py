# mautrix-twitter - A Matrix-Twitter DM puppeting bridge
# Copyright (C) 2022 Tulir Asokan, Max Sandholm
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.
from __future__ import annotations

import asyncio
import logging
import traceback

# from .config import Config
from .db import BackfillStatus as DBBackfillStatus
from .portal import Portal
from .user import User


class NoFirstMessageException(BaseException):
    pass


log = logging.getLogger("mau.backfill_loop")


class BackfillStatus(DBBackfillStatus):
    recheck_queues: set[asyncio.Queue] = set()

    @classmethod
    async def recheck(cls) -> None:
        for q in cls.recheck_queues:
            await q.put(True)

    @classmethod
    async def get_next_backfill_status(cls, recheck_queue: asyncio.Queue) -> BackfillStatus:
        while True:
            status = await cls.get_next_unfinished_status()
            if status != None:
                return status

            try:
                await asyncio.wait_for(recheck_queue.get(), 10)
            except asyncio.exceptions.TimeoutError:
                pass

    @classmethod
    async def backfill_loop(cls) -> None:
        recheck_queue = asyncio.Queue()
        cls.recheck_queues.add(recheck_queue)

        while True:
            await asyncio.sleep(2)
            state = await cls.get_next_backfill_status(recheck_queue)
            portal = await Portal.get_by_twid(twid=state.twid, receiver=state.receiver)
            source = await User.get_by_twid(state.backfill_user)
            try:
                state.dispatched = True
                await state.update()
                num_filled = await portal.backfill(source, is_initial=state.state == 0)

                state.message_count += num_filled
                if num_filled == 0:
                    state.state = 2
                elif state.state == 0:
                    state.state = 1
            except NoFirstMessageException:
                log.error(f"No first message found to do backfill for {state.twid}")
                state.state = 0
            except Exception:
                log.exception(f"Error handling backfill task for {state.twid}")
                state.state = 3
            finally:
                state.dispatched = False
                await state.update()
