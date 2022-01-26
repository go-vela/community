# v0.9 migration

> NOTE: This applies when upgrading to the latest `v0.9.x` release.

When migrating from Vela version [v0.8](../../releases/v0.8.md) to [v0.9](../../releases/v0.9.md) the Vela administrator will want to ensure the following actions are being performed:

1. Updating tables in the database:
   * `ALTER TABLE repos ADD COLUMN IF NOT EXISTS pipeline_type TEXT DEFAULT 'yaml';`

## Utility

For your convenience, we've provided a `vela-migration` utility in this directory to help execute the database operations.

This utility supports invoking the following actions when migrating to `v0.9.x`:

* `action.all` - run all supported actions (below) configured in the migration utility

More information can be found in the [`DOCS.md` for the utility](DOCS.md).