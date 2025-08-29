# v0.27 migration

> NOTE: This applies when upgrading to the latest `v0.27.x` release.

When migrating from Vela version [v0.26](../../releases/v0.26.md) to [v0.27](../../releases/v0.27.md) the Vela
administrator will want to ensure the following actions are being performed. All queries are available in the [SQL migration scripts](./scripts/).

1. A new table, `secret_repo_allowlists`, is introduced to support allowlisting repositories for secrets. This table includes `id` (BIGSERIAL, primary key), `secret_id` (BIGINT), and `repo` (VARCHAR-500) with a unique constraint on `(secret_id, repo)`. An index on `secret_id` (`secret_repo_allowlists_secret_id`) is also created.

2. The `builds` table adds a `route` (VARCHAR-250) column.

3. The `logs` table adds a `created_at` (BIGINT) column to track when logs were created. An additional index is created to improve query performance: `logs_created_at` (on `created_at`).

4. The `repos` table adds `custom_props` (JSON, nullable) for storing custom repository properties.

5. The `settings` table adds `scm` (JSON, nullable) for source control management configuration regarding role mappings.

6. The `settings` table adds `max_dashboard_repos` (INTEGER).

7. The `settings` table adds `queue_restart_limit` (INTEGER).

8. Backfill `logs.created_at` for existing rows. The migration attempts to populate values in the following order for accuracy:
	- From `steps.created` where `logs.step_id = steps.id`.
	- From `services.created` where `logs.service_id = services.id`.
	- From `builds.created` via joins through steps or services if still NULL.
	- Final fallback sets remaining NULLs to the current epoch time.
	
   If accuracy is less critical and speed is preferred, you can simplify to only the final fallback to current time or leave as-is.

   For large `logs` tables, consider running the backfill in smaller batches to avoid long locks or timeouts.