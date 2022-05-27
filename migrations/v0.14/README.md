# v0.14 migration

> NOTE: This applies when upgrading to the latest `v0.14.x` release.

When migrating from Vela version [v0.13](../../releases/v0.13.md) to [v0.14](../../releases/v0.14.md) the Vela
administrator will want to ensure the following actions are being performed:

1. `v0.14.x` introduces a new `pipelines` table that will be created automatically by default. However, if you are running the `server` component with the `VELA_DATABASE_SKIP_CREATION` set to `true`, you will need to manually create this table using the following:
  * Postgres: 
    ```sql
    CREATE TABLE
    IF NOT EXISTS
    pipelines (
      id               SERIAL PRIMARY KEY,
      repo_id          INTEGER,
      commit           VARCHAR(500),
      flavor           VARCHAR(100),
      platform         VARCHAR(100),
      ref              VARCHAR(500),
      type             VARCHAR(100),
      version          VARCHAR(50),
      external_secrets BOOLEAN,
      internal_secrets BOOLEAN,
      services         BOOLEAN,
      stages           BOOLEAN,
      steps            BOOLEAN,
      templates        BOOLEAN,
      data             BYTEA,
      UNIQUE(repo_id, commit)
    );
    ```
  * SQLite:
    ```sql
    CREATE TABLE
    IF NOT EXISTS
    pipelines (
      id               INTEGER PRIMARY KEY AUTOINCREMENT,
      repo_id          INTEGER,
      'commit'         TEXT,
      flavor           TEXT,
      platform         TEXT,
      ref              TEXT,
      type             TEXT,
      version          TEXT,
      external_secrets BOOLEAN,
      internal_secrets BOOLEAN,
      services         BOOLEAN,
      stages           BOOLEAN,
      steps            BOOLEAN,
      templates        BOOLEAN,
      data             BLOB,
      UNIQUE(repo_id, 'commit')
    );
    ```
1. Updating tables in the database:
  * `ALTER TABLE builds ADD COLUMN IF NOT EXISTS event_action VARCHAR(250);`
  * `ALTER TABLE builds ADD COLUMN IF NOT EXISTS pipeline_id INTEGER;` 
  * `ALTER TABLE hooks ADD COLUMN IF NOT EXISTS event_action VARCHAR(250);`
  * `ALTER TABLE hooks ADD COLUMN IF NOT EXISTS webhook_id INTEGER;`
1. Although not required for the release, we recommend adding an additional index on the `source` column in the `builds table`.
  * `CREATE INDEX CONCURRENTLY IF NOT EXISTS builds_source ON builds (source);`

## Utility

For your convenience, we've provided a `vela-migration` utility in this directory to help execute the database
operations.

This utility supports invoking the following actions when migrating to `v0.14.x`:

* `action.all` - run all supported actions (below) configured in the migration utility
* `alter.tables` - runs the required queries to alter the database tables
* `create.indexees` - runs the required queries to create the database indexes

More information can be found in the [`DOCS.md` for the utility](DOCS.md).
