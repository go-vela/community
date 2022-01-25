# v0.12 migration

> NOTE: This applies when upgrading to the latest `v0.12.x` release.

When migrating from Vela version [v0.11](../../releases/v0.11.md) to [v0.12](../../releases/v0.12.md) the Vela
administrator will want to ensure the following actions are being performed:

1. Updating tables in the database:
    * `ALTER TABLE repos ADD COLUMN IF NOT EXISTS build_limit INTEGER DEFAULT 10;`
    * `ALTER TABLE repos ADD COLUMN IF NOT EXISTS previous_name VARCHAR(100);`
    * `ALTER TABLE secrets ADD COLUMN IF NOT EXISTS created_at INTEGER;`
    * `ALTER TABLE secrets ADD COLUMN IF NOT EXISTS created_by VARCHAR(250);`
    * `ALTER TABLE secrets ADD COLUMN IF NOT EXISTS updated_at INTEGER;`
    * `ALTER TABLE secrets ADD COLUMN IF NOT EXISTS updated_by VARCHAR(250);`

1. Adding new indexes to the database:
    * `CREATE INDEX CONCURRENTLY IF NOT EXISTS builds_created ON builds (created);`

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the SCM driver:
    * use `scm.driver` to set via CLI flag
    * use `SCM_DRIVER` or `VELA_SCM_DRIVER` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the SCM address:
    * use `scm.addr` to set via CLI flag
    * use `SCM_ADDR` or `VELA_SCM_ADDR` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the SCM OAuth client id:
    * use `scm.client` to set via CLI flag
    * use `SCM_CLIENT` or `VELA_SCM_CLIENT` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the SCM OAuth client
   secret:
    * use `scm.secret` to set via CLI flag
    * use `SCM_SECRET` or `VELA_SCM_SECRET` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the SCM status context:
    * use `scm.context` to set via CLI flag
    * use `SCM_CONTEXT` or `VELA_SCM_CONTEXT` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the SCM OAuth scopes:
    * use `scm.scopes` to set via CLI flag
    * use `SCM_SCOPES` or `VELA_SCM_SCOPES` to set via environment variable

1. If needed, update the [go-vela/server](https://github.com/go-vela/server) configuration for the SCM webhook address:
    * use `scm.webhook.addr` to set via CLI flag
    * use `SCM_WEBHOOK_ADDR` or `VELA_SCM_WEBHOOK_ADDR` to set via environment variable

## Utility

For your convenience, we've provided a `vela-migration` utility in this directory to help execute the database
operations.

This utility supports invoking the following actions when migrating to `v0.12.x`:

* `action.all` - run all supported actions (below) configured in the migration utility
* `alter.tables` - runs the required queries to alter the database tables
* `create.indexees` - runs the required queries to create the database indexes

More information can be found in the [`DOCS.md` for the utility](DOCS.md).