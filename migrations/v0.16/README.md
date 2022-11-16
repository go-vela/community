# v0.16 migration

> NOTE: This applies when upgrading to the latest `v0.16.x` release.

When migrating from Vela version [v0.15](../../releases/v0.15.md) to [v0.16](../../releases/v0.16.md) the Vela
administrator will want to ensure the following actions are being performed:

1. `v0.16.x` introduces a new security enhancement where privileged docker images will
only be executed when the repo is `trusted`. This functionality will be enabled by default but can be ignored setting the `VELA_EXECUTOR_ENFORCE_TRUSTED_REPOS` worker flag to `false`. In order to effectively use this enhancement, the platform administrators will need to run the following query to start:
  ```sql
  UPDATE repos SET trusted = 'false';
  ```
1. Further, if you would like to grant `trusted` to repos that have already been using privileged images during a certain time frame, you can execute the below query:
  ```sql
  UPDATE repos
  SET trusted = 'true'
  WHERE id IN (
    SELECT id
    FROM repos r 
    INNER JOIN (
        SELECT 
            repo_id
        FROM steps 
        WHERE image LIKE '%<your_image>%' AND 
        created > (
            SELECT EXTRACT(EPOCH FROM (NOW() - INTERVAL '<no. of days>' DAY))
        )   
        GROUP BY repo_id
    ) t 
    ON r.id = t.repo_id 
  WHERE active = 't'
  );
  ```

## Utility

For your convenience, we've provided a `vela-migration` utility in this directory to help execute the database operations.

This utility supports invoking the following actions when migrating to `v0.16.x`:

* `action.all` - run all supported actions (below) configured in the migration utility
* `action.untrusted` - runs the required queries to set all repos `trusted` to false
* `action.update-trusted` - runs the required queries to give already privileged repos `trusted` to true.

Options to supply:
* `trusted-update.privileged-images` - string slice of privileged images (default ['target/vela-docker'])
* `trusted-update.allow-personal-orgs` - bool to allow personal orgs to be trusted (default `true`)
* `trusted-update.days-back` - string of how many days back to track repo usage of privileged images (default `"90"`)
