# v0.18 migration

> NOTE: This applies when upgrading to the latest `v0.18.x` release.

When migrating from Vela version [v0.17](../../releases/v0.17.md) to [v0.18](../../releases/v0.18.md) the Vela
administrator will want to ensure the following actions are being performed:

1. Updating tables in the database:
    * `ALTER TABLE users ALTER COLUMN token TYPE varchar(1000);`

## Utility

For your convenience, we've provided a `vela-migration` utility in this directory to help execute the database
operations.

This utility supports invoking the following actions when migrating to `v0.18.x`:

* `action.all` - run all supported actions (below) configured in the migration utility
* `alter.tables` - runs the required queries to alter the database tables

More information can be found in the [`DOCS.md` for the utility](DOCS.md).
