# mautrix-twitter - A Matrix-Twitter DM puppeting bridge
# Copyright (C) 2020 Tulir Asokan
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
from . import command_handler, CommandEvent, SECTION_AUTH


@command_handler(needs_auth=False, management_only=True,
                 help_section=SECTION_AUTH, help_text="Log in to Twitter manually")
async def login_cookie(evt: CommandEvent) -> None:
    if evt.sender.client:
        await evt.reply("You're already logged in")
        return
    evt.sender.command_status = {
        "action": "Login",
        "room_id": evt.room_id,
        "next": enter_login_cookies,
        "auth_token": None,
    }
    await evt.reply(
        "1. Log in to [Twitter](https://www.twitter.com/) in a private/incognito window.\n"
        "2. Press `F12` to open developer tools.\n"
        "3. Select the \"Application\" (Chrome) or \"Storage\" (Firefox) tab.\n"
        "4. In the sidebar, expand \"Cookies\" and select `https://twitter.com`.\n"
        "5. In the cookie list, find the `auth_token` row and double click on the value"
        r", then copy the value and send it here.")


async def enter_login_cookies(evt: CommandEvent) -> None:
    if not evt.sender.command_status["auth_token"]:
        if len(evt.args) == 0:
            await evt.reply("Please enter the value of the `auth_token` cookie, or use "
                            "the `cancel` command to cancel.")
            return
        evt.sender.command_status["auth_token"] = evt.args[0]
        await evt.reply("Now do the last step again, but find the value of the `ct0` row instead. "
                        "Before you send the value, close the private window.")
        return
    if len(evt.args) == 0:
        await evt.reply("Please enter the value of the `ct0` cookie, or use "
                        "the `cancel` command to cancel.")
        return

    try:
        await evt.sender.connect(auth_token=evt.sender.command_status["auth_token"],
                                 csrf_token=evt.args[0])
    except Exception as e:
        evt.sender.command_status = None
        await evt.reply(f"Failed to log in: {e}")
        evt.log.exception("Failed to log in")
        return

    await evt.reply(f"Successfully logged in as @{evt.sender.username}")
    evt.sender.command_status = None


@command_handler(needs_auth=True, help_section=SECTION_AUTH, help_text="Disconnect the bridge from"
                                                                       "your Twitter account")
async def logout(evt: CommandEvent) -> None:
    await evt.sender.logout()
    await evt.reply("Successfully logged out")
