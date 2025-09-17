# v0.26 migration

> NOTE: This applies when upgrading to the latest `v0.26.x` release.

When migrating from Vela version [v0.25](../../releases/v0.25.md) to [v0.26](../../releases/v0.26.md) the Vela
administrator will want to ensure the following actions are being performed. All queries are available in the [SQL migration scripts](./scripts/).

1. `v0.26.x` introduces a new opt-in GitHub App integration. This integration requires the addition of the `install_id` (INTEGER) column to the `repos` table.

  - https://github.com/go-vela/server/pull/1217

2. Repos will also be able to set approval timeouts for builds that are pending approval. This requires `approval_timeout` (INTEGER) to be added to the `repos` table.
  - https://github.com/go-vela/server/pull/1227

3. `v0.26.x` will use `go-yaml/v3` as the YAML parser for pipelines and templates but will support backwards comptability with `buildkite/yaml` until a later release. This change led to the creation of pipeline warnings, which necessitate the addition of `warnings` (VARCHAR-5000) to the `pipelines` table.

  - https://github.com/go-vela/server/pull/1232


4. An additional OIDC claim will be added to identity tokens in `v0.26.x`. This will be the `fork` field, which should be added to `builds` table as a BOOLEAN.

  - https://github.com/go-vela/server/pull/1221


5. `v0.26.x` introduces a new index (`builds_event`) to the default configuration. This index will be created automatically. However, if you are running the `server` component with the `VELA_DATABASE_SKIP_CREATION` set to `true`, you will need to manually create this index. This change accompanies a deletion of an unused index (`builds_source`). You will need to manually delete this index if it exists, regardless of the skip creation setting.

  - https://github.com/go-vela/server/pull/1228

6. Specifically `v0.26.5` adds the `route` column to the `builds` table as well as the `queue_restart_limit` column to the `settings` table.