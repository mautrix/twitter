# v0.4.1 (2025-05-16)

* Added support for voice messages in both directions.
* Fixed certain reactions not being bridged to Twitter.
* Fixed private chats having explicit names/avatars even if the bridge wasn't
  configured to set them.
* Fixed mention bridging not including the mentioned user displayname.
* Fixed handling member join events from Twitter.

# v0.4.0 (2025-04-16)

* Added support for own read status bridging.
* Added support for sending mentions and intentional mentions in incoming messages.
* Fixed newlines in incoming formatted messages.
* Stopped creating portals for message requests automatically.

# v0.3.0 (2025-03-16)

* Added support for tweet attachments.
* Added support for unshortening URLs in incoming messages.
* Fixed sending media in unencrypted Matrix rooms.
* Fixed chats not being bridged in some cases even after accepting the message
  request on Twitter.

# v0.2.1 (2025-01-16)

* Fixed various bugs.

# v0.2.0 (2024-12-16)

* Rewrote bridge in Go using bridgev2 architecture.
  * To migrate the bridge, simply upgrade in-place. The database and config
    will be migrated automatically, although some parts of the config aren't
    migrated (e.g. log config).
  * It is recommended to check the config file after upgrading. If you have
    prevented the bridge from writing to the config, you should update it
    manually.

# v0.1.8 (2024-07-16)

No changelog available.

# v0.1.7 (2023-09-19)

No changelog available.

# v0.1.6 (2023-05-22)

No changelog available.

# v0.1.5 (2022-08-23)

No changelog available.

# v0.1.4 (2022-03-30)

No changelog available.

# v0.1.3 (2022-01-15)

No changelog available.

# v0.1.1 (2020-12-11)

No changelog available.

# v0.1.1 (2020-11-10)

No changelog available.

# v0.1.0 (2020-09-04)

Initial release.
