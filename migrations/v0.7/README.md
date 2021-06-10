# v0.7 migration

> NOTE: This applies when upgrading to the latest `v0.7.x` release.

When migrating from Vela version [v0.6.0](../../releases/v0.6.0.md) to [v0.7](../../releases/v0.7.md) the Vela administrator will want to ensure the following actions are being performed:

1. Updating tables in the database:
   * `ALTER TABLE builds ADD COLUMN IF NOT EXISTS deploy_payload VARCHAR(2000);`
   * `ALTER TABLE users ADD COLUMN IF NOT EXISTS refresh_token VARCHAR(500);`
   * `ALTER TABLE workers ADD COLUMN IF NOT EXISTS build_limit INTEGER;`

1. Dropping unused indexes from the database:
   * `DROP INDEX IF EXISTS builds_repo_id_number;`
   * `DROP INDEX IF EXISTS hooks_repo_id_number;`
   * `DROP INDEX IF EXISTS logs_step_id;`
   * `DROP INDEX IF EXISTS logs_service_id;`
   * `DROP INDEX IF EXISTS repos_full_name;`
   * `DROP INDEX IF EXISTS secrets_type;`
   * `DROP INDEX IF EXISTS services_build_id_number;`
   * `DROP INDEX IF EXISTS steps_build_id_number;`
   * `DROP INDEX IF EXISTS users_name;`

1. Compressing the `data` field for all rows in the `logs` table in the database:
   * Use the `vela-migration` utility following [the documentation](DOCS.md).

1. Encrypting the `value` field for all rows in the `secrets` table in the database:
   * Use the `vela-migration` utility following [the documentation](DOCS.md).

1. Configure the `Authorization callback URL` for the Vela OAuth App (in GitHub) to point to `<vela-server-address>/authenticate`.
   * DO NOT use the UI address for the callback url.

1. If configured, ensure the `VELA_PORT` environment variable does not include a leading `:`.
   * Removing the environment variable altogether is acceptable if the `8080` default setting works.

1. Configure the `VELA_DATABASE_ENCRYPTION_KEY` environment variable for the [go-vela/server](https://github.com/go-vela/server) with a random string.
   * The encryption key SHOULD NOT contain any special characters and must be 32 characters in length.

## Utility

For your convenience, we've provided a `vela-migration` utility in this directory to help execute the database operations.

This utility supports invoking the following actions when migrating to `v0.7.x`:

* `all` - run all supported actions (below) configured in the migration utility
* `alter.tables` - runs the required queries to alter the database tables
* `compress.logs` - runs the required queries to compress the `data` field for all rows in the `logs` table
* `drop.indexes` - runs the required queries to drop the unused indexes
* `encrypt.secrets` - runs the required queries to encrypt the `value` field for all rows in the `secrets` table

More information can be found in the [`DOCS.md` for the utility](DOCS.md).
