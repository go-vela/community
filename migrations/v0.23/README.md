# v0.23 migration

> NOTE: This applies when upgrading to the latest `v0.23.x` release.

When migrating from Vela version [v0.22](../../releases/v0.22.md) to [v0.23](../../releases/v0.23.md) the Vela
administrator will want to ensure the following actions are being performed. All queries are available in the [SQL migration scripts](./scripts/).

1. `v0.23.x` introduces a new `deployments` table that will be created automatically by default. However, if you are running the `server` component with the `VELA_DATABASE_SKIP_CREATION` set to `true`, you will need to manually create this table. 

  - https://github.com/go-vela/server/pull/1031


2. `v0.23.x` introduces a new `approve_build` column to the repos table for approve build policy setting. The change also impacts the builds table, with the addition of the `approved_at` and `approved_by` columns.

  - https://github.com/go-vela/server/pull/1016

3. The `approve_build` column should be prepopulated using the SQL query in the migration script to ensure a smooth transition. Feel free to adjust the default setting for already-enabled repos.


4. `v0.23.x` also introduces a new `allow_events` column to the repos table, which is an integer mask that represents all possible event subscription configurations. The secrets table will also utilize the new `allow_events` column. 

  - https://github.com/go-vela/server/pull/1023
  - https://github.com/go-vela/server/pull/1033

5. The `allow_events` column should be prepopulated using the SQL query in the migration script to ensure a smooth transition to the new system.

6. `v0.23.x` (starting in v0.23.2) also adds an `allow_substitution` column to the secrets table to give more control on secret usage.

## Vault Secret Engine

Rather than migrate all Vault internal secrets to leverage `allow_events` and `allow_substition` via scripting, the Vela code has been updated to include [translation per request](https://github.com/go-vela/server/pull/1086) for Vault secrets that have yet to be updated since an upgrade to `v0.23.3`. These values will be rectified over time whenever users update the secret.

For existing native secrets, please follow the provided DML query in the migration script.

## Recommended

For increased security we recommend to set `allow_command` and `allow_substitution` to `false` for shared secrets in your secrets table. You can use the following SQL commands to do so:

```sql
UPDATE secrets SET allow_command = false WHERE type = 'shared';
```

```sql
UPDATE secrets SET allow_substitution = false WHERE type = 'shared';
```
