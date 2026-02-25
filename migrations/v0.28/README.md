# v0.28 migration

> NOTE: This applies when upgrading to the latest `v0.28.x` release.

When migrating from Vela version [v0.27](../../releases/v0.27.md) to [v0.28](../../releases/v0.28.md) the Vela
administrator will want to ensure the following actions are being performed. All database queries are available in the [SQL migration scripts](./scripts/).

## Database Migrations

1. The `pipelines` table adds a `artifacts` (BOOLEAN) column.

2. The `repos` table adds columns `merge_queue_events` (VARCHAR-500) and `hook_counter` (BIGINT).

3. The `settings` table adds columns `enable_shared_secrets` (BOOLEAN), `enable_org_secrets` (BOOLEAN), and `enable_repo_secrets` (BOOLEAN).

## Cache Service (Required for GitHub App Integration)

`v0.28` introduces a new Redis-driven cache that stores installation token metadata and check run information. Administrators can decide whether the cache is a part of the same Redis cluster as the queue or separate.

Below are the required environment variables for setting up the cache:

- `VELA_CACHE_DRIVER`: cache driver. Only available driver is `redis`.
- `VELA_CACHE_ADDR`: connection address to the cache service.
- `VELA_CACHE_CLUSTER`: enables connection to a cache cluster.
- `VELA_CACHE_INSTALL_TOKEN_KEY`: key used for hex-signing install tokens as cache keys. Can be generated with `openssl rand -hex 32`.

## Artifact Storage Service

`v0.28` introduces a new MinIO-driven storage service for build artifacts. Users will be able to upload these artifacts to an S3 bucket of the administrators choosing. These artifacts will show up as available to be downloaded in the Vela UI.

Below are the required environment variables for setting up the storage service:

- `VELA_STORAGE_ENABLE`: on/off flag for storage service.
- `VELA_STORAGE_DRIVER`: storage driver. Only availabe driver is `minio`
- `VELA_STORAGE_ADDRESS`: address of storage service.
- `VELA_STORAGE_ACCESS_KEY`: key identification for the storage service.
- `VELA_STORAGE_SECRET_KEY`: key credential for the storage service.
- `VELA_STORAGE_BUCKET`: name of configured bucket.
