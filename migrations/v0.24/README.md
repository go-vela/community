# v0.24 migration

> NOTE: This applies when upgrading to the latest `v0.24.x` release.

When migrating from Vela version [v0.23](../../releases/v0.23.md) to [v0.24](../../releases/v0.24.md) the Vela administrator will want to ensure the following actions are being performed. All queries are available in the [SQL migration scripts](./scripts/).

1. `v0.24.x` introduces a new `sender_scm_id` field to the `builds` table that will be set to the string identifier that maps to the `sender` within the SCM.

  - https://github.com/go-vela/server/pull/1133

2. For pre-existing builds, the `sender_scm_id` column should be prepopulated to `0` using the SQL query in the migration script.
