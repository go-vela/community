# <VELA_VERSION> migration

> NOTE: This applies when upgrading to the latest `<VELA_VERSION>.x` release.

When migrating from Vela version [<VELA_PREV_VERSION>](../../releases/<VELA_PREV_VERSION>.md) to [<VELA_VERSION>](../../releases/<VELA_VERSION>.md) the Vela
administrator will want to ensure the following actions are being performed:

1. `<VELA_VERSION>.x` introduces a new enhancement where repository topics will be injected as an environment variable and available to all builds. These topics will be stored in the `repos` table. In order to effectively use this enhancement, the platform administrators will need to run the following query:
  ```sql
  ALTER TABLE repos ADD COLUMN IF NOT EXISTS topics VARCHAR(1020);
  ```


## Utility

For your convenience, we've provided a `vela-migration` utility in this directory to help execute the database operations.

This utility supports invoking the following actions when migrating to `<VELA_VERSION>.x`:

* `action.all` - run all supported actions (below) configured in the migration utility

More information can be found in the [`DOCS.md` for the utility](DOCS.md).