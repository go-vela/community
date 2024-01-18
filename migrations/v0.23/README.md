# v0.23 migration

> NOTE: This applies when upgrading to the latest `v0.23.x` release.

When migrating from Vela version [v0.22](../../releases/v0.22.md) to [v0.23](../../releases/v0.23.md) the Vela
administrator will want to ensure the following actions are being performed:

1. `v0.23.x` introduces a new `approve_build` column to the repos table for approve build policy setting. In order to effectively use this enhancement, the platform administrators will need to run the following query:
  ```sql
  ALTER TABLE repos
      ADD COLUMN IF NOT EXISTS approve_build VARCHAR(250)
      ;
  ```

1. The above change also impacts the builds table, with the addition of the `approved_at` and `approved_by` columns. In order to effectively use this enhancement, the platform administrators will need to run the following queries:
  ```sql
	ALTER TABLE builds
	  ADD COLUMN IF NOT EXISTS approved_at INTEGER
	  ;

	ALTER TABLE builds
	  ADD COLUMN IF NOT EXISTS approved_by INTEGER
	  ;
  ```

1. `v0.23.x` also introduces a new `allow_events` column to the repos table, which is an integer mask that represents all possible event subscription configurations. In order to effectively use this enhancement, the platform administrators will need to run the following query:
  ```sql
  ALTER TABLE repos
      ADD COLUMN IF NOT EXISTS allow_events INTEGER
      ;
  ```

1. The secrets table will also utilize the new `allow_events` column. In order to effectively use this enhancement, the platform administrators will need to run the following query:
  ```sql
	ALTER TABLE secrets
	  ADD COLUMN IF NOT EXISTS allow_events INTEGER
	  ;
  ```

1. With the allow_events and approve_build changes, there are a few initialization queries that make the transition to this new system much easier. These are the following queries:

  ```sql
    UPDATE repos 
      SET approve_build = 'fork-always'
	  ;

    UPDATE repos
    SET allow_events =
        (CASE WHEN allow_push = true THEN 1 ELSE 0 END) |
        (CASE WHEN allow_pull = true THEN 4 | 16 | 1024 ELSE 0 END) |
        (CASE WHEN allow_tag = true THEN 2 ELSE 0 END) |
        (CASE WHEN allow_deploy = true THEN 8192 ELSE 0 END) |
        (CASE WHEN allow_comment = true THEN 16384 | 32768 ELSE 0 END)
    ;

    UPDATE secrets
	  SET allow_events =
	    (CASE WHEN events LIKE '%push%' THEN 1 ELSE 0 END) |
        (CASE WHEN events LIKE '%pull_request%' THEN 4 | 16 | 1024 ELSE 0 END) |
        (CASE WHEN events LIKE '%tag%' THEN 2 ELSE 0 END) |
        (CASE WHEN events LIKE '%deployment%' THEN 8192 ELSE 0 END) |
        (CASE WHEN events LIKE '%comment%' THEN 16384 | 32768 ELSE 0 END) |
		    (CASE WHEN events LIKE '%schedule%' THEN 65536 ELSE 0 END)
    ;
  ```





## Utility

For your convenience, we've provided a `vela-migration` utility in this directory to help execute the database operations.

This utility supports invoking the following actions when migrating to `v0.23.x`:

* `action.all` - run all supported actions (below) configured in the migration utility
* `alter.tables` - runs the required queries to alter the database tables
* `update.repos` - runs the related updates for approve build policy and allow events on the repos table
* `update.secrets` - runs the related updates for the allow events on the secrets table

More information can be found in the [`DOCS.md` for the utility](DOCS.md).