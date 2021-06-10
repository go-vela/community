# v0.8 migration

> NOTE: This applies when upgrading to the latest `v0.8.x` release.

When migrating from Vela version [v0.7.4](../../releases/v0.7.4.md) to [v0.8.1](../../releases/v0.8.1.md) the Vela administrator will want to ensure the following actions are being performed:

1. Updating tables in the database:
   * `ALTER TABLE repos ADD COLUMN IF NOT EXISTS counter INTEGER;`

1. Encrypting the `hash`, `refresh_token` and `token` fields for all rows in the `users` table in the database:
   * Use the `vela-migration` utility following [the documentation](DOCS.md).

1. Syncing the build number `value` field for all rows in the `repos` table in the database:
   * Use the `vela-migration` utility following [the documentation](DOCS.md).

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) and [go-vela/worker](https://github.com/go-vela/worker) configuration for queue address:
   * use `queue.addr` to set the via CLI flag
   * use `QUEUE_ADDR` or `VELA_QUEUE_ADDR` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) and [go-vela/worker](https://github.com/go-vela/worker) configuration for the queue routes:
   * use `queue.routes` to set the via CLI flag
   * use `QUEUE_ROUTES` or `VELA_QUEUE_ROUTES` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) and [go-vela/worker](https://github.com/go-vela/worker) configuration for the queue timeout:
   * use `queue.pop.timeout` to set via the CLI flag
   * use `QUEUE_POP_TIMEOUT` or `VELA_QUEUE_POP_TIMEOUT` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the database address:
   * use `database.addr` to set via CLI flag
   * use `DATABASE_ADDR` or `VELA_DATABASE_ADDR` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the secret Vault driver:
   * use `secret.vault.driver` to set via CLI flag
   * use `SECRET_VAULT` or `VELA_SECRET_VAULT` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the secret Vault address:
   * use `secret.vault.addr` to set via CLI flag
   * use `SECRET_VAULT_ADDR` or `VELA_SECRET_VAULT_ADDR` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the secret Vault authentication method:
   * use `secret.vault.auth-method` to set via CLI flag
   * use `SECRET_VAULT_AUTH_METHOD` or `VELA_SECRET_VAULT_AUTH_METHOD` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the secret Vault prefix:
   * use `secret.vault.prefix` to set via CLI flag
   * use `SECRET_VAULT_PREFIX` or `VELA_SECRET_VAULT_PREFIX` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the secret Vault token renewal duration:
   * use `secret.vault.renewal` to set via CLI flag
   * use `SECRET_VAULT_RENEWAL` or `VELA_SECRET_VAULT_RENEWAL` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the secret Vault token:
   * use `secret.vault.token` to set via CLI flag
   * use `SECRET_VAULT_TOKEN` or `VELA_SECRET_VAULT_TOKEN` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the secret Vault version:
   * use `secret.vault.version` to set via CLI flag
   * use `SECRET_VAULT_VERSION` or `VELA_SECRET_VAULT_VERSION` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the source driver:
   * use `source.driver` to set via CLI flag
   * use `SOURCE_DRIVER` or `VELA_SOURCE_DRIVER` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the source address:
   * use `source.addr` to set via CLI flag
   * use `SOURCE_ADDR` or `VELA_SOURCE_ADDR` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the source OAuth client id:
   * use `source.client` to set via CLI flag
   * use `SOURCE_CLIENT` or `VELA_SOURCE_CLIENT` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the source OAuth client secret:
   * use `source.secret` to set via CLI flag
   * use `SOURCE_SECRET` or `VELA_SOURCE_SECRET` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the source status context:
   * use `source.context` to set via CLI flag
   * use `SOURCE_CONTEXT` or `VELA_SOURCE_CONTEXT` to set via environment variable

1. If needed, update the [go-vela/worker](https://github.com/go-vela/worker) configuration for privileged images:
   * use `runtime.privileged-images` to set via CLI flag
   * use `RUNTIME_PRIVILEGED_IMAGES` or `VELA_RUNTIME_PRIVILEGED_IMAGES` to set via environment variable

## Utility

For your convenience, we've provided a `vela-migration` utility in this directory to help execute the database operations.

This utility supports invoking the following actions when migrating to `v0.8.x`:

* `action.all` - run all supported actions (below) configured in the migration utility
* `alter.tables` - runs the required queries to alter the database tables
* `encrypt.users` - runs the required queries to encrypt the the `hash`, `refresh_token` and `token` fields for all rows in the `users` table
* `sync.counter` - runs the required queries to grab the build number `value` field and update the repo counter for all rows in the `repos` table

More information can be found in the [`DOCS.md` for the utility](DOCS.md).
