# v0.19 migration

> NOTE: This applies when upgrading to the latest `v0.19.x` release.

When migrating from Vela version [v0.18](../../releases/v0.18.md) to [v0.19](../../releases/v0.19.md) the Vela
administrator will want to ensure the following actions are being performed:

1. `v0.19.x` introduces a new enhancement where repository topics will be injected as an environment variable and available to all builds. These topics will be stored in the `repos` table. In order to effectively use this enhancement, the platform administrators will need to run the following query:
  ```sql
  ALTER TABLE repos ADD COLUMN IF NOT EXISTS topics VARCHAR(1020);
  ```


## Utility

For your convenience, we've provided a `vela-migration` utility in this directory to help execute the database operations.

This utility supports invoking the following actions when migrating to `v0.19.x`:

* `action.all` - run all supported actions (below) configured in the migration utility
* `alter.tables` - runs the required queries to alter the database tables

More information can be found in the [`DOCS.md` for the utility](DOCS.md).