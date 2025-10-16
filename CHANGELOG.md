# v25.10

* Switched to calendar versioning.
* Removed legacy provisioning API and database legacy migration.
  Upgrading directly from versions prior to v0.2.0 is not supported.
  * If you've been using the bridge since before v0.2.0 and have prevented the
    bridge from writing to the config, you must either update the config
    manually or allow the bridge to update it for you **before** upgrading to
    this release (i.e. run v0.5.0 once with config writing allowed).
* Fixed parsing data from Twitter index page.
* Fixed handling auth errors while polling for new messages.

# v0.5.0 (2025-08-16)

* Deprecated legacy provisioning API. The `/_matrix/provision/v1` endpoints will
  be deleted in the next release.
* Bumped minimum Go version to 1.24.
* Added option to reconnect faster after restart by caching connection state.
* Added support for bridging group chat participant leaves from Twitter.
* Fixed ogg voice messages converted from Twitter including a video stream.

# v0.4.3 (2025-07-16)

* Fixed forward backfill not fetching more messages if there's a large gap.

# v0.4.2 (2025-06-16)

* Added notice about missed calls.
* Added basic support for direct media.
* Updated Docker image to Alpine 3.22.
* Fixed sending attachments with no caption.
* Fixed management room notices not being sent if polling Twitter failed.
* Fixed portals not being created after accepting message request that Twitter
  had marked as "low quality".

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
