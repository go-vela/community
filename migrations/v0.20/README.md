# v0.20 migration

> NOTE: This applies when upgrading to the latest `v0.20.x` release.

When migrating from Vela version [v0.19](../../releases/v0.19.md) to [v0.20](../../releases/v0.20.md) the Vela
administrator will want to ensure the following actions are being performed:

1. `v0.20.x` introduces multiple new columns in `workers` table to enhance worker's visibility. In order to effectively use this enhancement, the platform administrators will need to run the following query:
  ```sql
  ALTER TABLE workers 
      ADD COLUMN IF NOT EXISTS status VARCHAR(50),
      ADD COLUMN IF NOT EXISTS last_status_update_at INTEGER,
      ADD COLUMN IF NOT EXISTS running_build_ids VARCHAR(500),
      ADD COLUMN IF NOT EXISTS last_build_started_at INTEGER,
      ADD COLUMN IF NOT EXISTS last_build_finished_at INTEGER
      ;
  ```


## Utility

For your convenience, we've provided a `vela-migration` utility in this directory to help execute the database operations.

This utility supports invoking the following actions when migrating to `v0.20.x`:

* `action.all` - run all supported actions (below) configured in the migration utility
* `alter.tables` - runs the required queries to alter the database tables

More information can be found in the [`DOCS.md` for the utility](DOCS.md).