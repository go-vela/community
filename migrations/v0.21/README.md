# v0.21 migration

> NOTE: This applies when upgrading to the latest `v0.21.x` release.

When migrating from Vela version [v0.19](../../releases/v0.19.md) to [v0.21](../../releases/v0.21.md) the Vela
administrator will want to ensure the following actions are being performed:

1. `v0.21.x` introduces a new `build_executables` table that will be created automatically by default. However, if you are running the `server` component with the `VELA_DATABASE_SKIP_CREATION` set to `true`, you will need to manually create this table using the following:
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

1. `v0.21.x` also introduces a new `branch` column to the schedules table for targeting specific branches. In order to effectively use this enhancement, the platform administrators will need to run the following query:
  ```sql
  ALTER TABLE schedules
      ADD COLUMN IF NOT EXISTS branch VARCHAR(250)
      ;
  ```

1. It is also imperative to run the following query to adjust the existing branch for schedule records created prior to the migration.
  ```sql
  UPDATE schedules 
      SET branch = r.branch 
      FROM (SELECT id, branch FROM repos) r 
      WHERE schedules.repo_id = r.id
	    AND schedules.branch IS NULL
      ; 
  ```


## Utility

For your convenience, we've provided a `vela-migration` utility in this directory to help execute the database operations.

This utility supports invoking the following actions when migrating to `v0.21.x`:

* `action.all` - run all supported actions (below) configured in the migration utility
* `alter.tables` - runs the required queries to alter the database tables
* `update.schedules` - runs the required queries to update existing schedules to leverage the corresponding repo default branch

More information can be found in the [`DOCS.md` for the utility](DOCS.md).