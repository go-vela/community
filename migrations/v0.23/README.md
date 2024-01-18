# v0.23 migration

> NOTE: This applies when upgrading to the latest `v0.23.x` release.

When migrating from Vela version [v0.19](../../releases/v0.19.md) to [v0.23](../../releases/v0.23.md) the Vela
administrator will want to ensure the following actions are being performed:

1. `v0.23.x` introduces a new `build_executables` table that will be created automatically by default. However, if you are running the `server` component with the `VELA_DATABASE_SKIP_CREATION` set to `true`, you will need to manually create this table using the following:
  * Postgres: 
    ```sql
    CREATE TABLE
    IF NOT EXISTS
    build_executables (
      id               SERIAL PRIMARY KEY,
      build_id         INTEGER,
      data             BYTEA,
      UNIQUE(build_id)
    );
    ```
  * SQLite:
    ```sql
    CREATE TABLE
    IF NOT EXISTS
    build_executables (
      id               INTEGER PRIMARY KEY AUTOINCREMENT,
      build_id         INTEGER,
      data             BLOB,
      UNIQUE(build_id)
    );
    ```

1. `v0.23.x` introduces a new `branch` column to the repos table for targeting specific branches. In order to effectively use this enhancement, the platform administrators will need to run the following query:
  ```sql
  ALTER TABLE repos
      ADD COLUMN IF NOT EXISTS branch VARCHAR(250)
      ;
  ```


## Utility

For your convenience, we've provided a `vela-migration` utility in this directory to help execute the database operations.

This utility supports invoking the following actions when migrating to `v0.23.x`:

* `action.all` - run all supported actions (below) configured in the migration utility
* `alter.tables` - runs the required queries to alter the database tables

More information can be found in the [`DOCS.md` for the utility](DOCS.md).