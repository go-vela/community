# v0.8 migration

> NOTE: This applies when upgrading to the latest `v0.8.x` release.

When migrating from Vela version [v0.7.4](../../releases/v0.7.4.md) to [v0.8.0](../../releases/v0.8.0.md) the Vela administrator will want to ensure the following actions are being performed:

1. Updating tables in the database:
   * `ALTER TABLE repos ADD COLUMN IF NOT EXISTS counter INTEGER;`

1. Encrypting the `hash`, `refresh_token` and `token` fields for all rows in the `users` table in the database:
   * Use the `vela-migration` utility following [the documentation](DOCS.md).

1. Syncing the build number `value` field for all rows in the `repos` table in the database:
   * Use the `vela-migration` utility following [the documentation](DOCS.md).

## Utility

For your convenience, we've provided a `vela-migration` utility in this directory to help execute the database operations.

This utility supports invoking the following actions when migrating to `v0.8.x`:

* `action.all` - run all supported actions (below) configured in the migration utility
* `alter.tables` - runs the required queries to alter the database tables
* `encrypt.users` - runs the required queries to encrypt the the `hash`, `refresh_token` and `token` fields for all rows in the `users` table
* `sync.counter` - runs the required queries to grab the build number `value` field and update the repo counter for all rows in the `repos` table

More information can be found in the [`DOCS.md` for the utility](DOCS.md).
