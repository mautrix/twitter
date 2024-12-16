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
