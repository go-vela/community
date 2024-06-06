# v0.24 migration

> NOTE: This applies when upgrading to the latest `v0.24.x` release.

When migrating from Vela version [v0.23](../../releases/v0.23.md) to [v0.24](../../releases/v0.24.md) the Vela
administrator will want to ensure the following actions are being performed. All queries are available in the [SQL migration scripts](./scripts/).

1. `v0.24.x` introduces a new `dashboards` table that will be created automatically by default. However, if you are running the `server` component with the `VELA_DATABASE_SKIP_CREATION` set to `true`, you will need to manually create this table. 

  - https://github.com/go-vela/server/pull/1028

2. This `dashboards` table also impacts the `users` table, adding the column `dashboards` to attach dashboards to users.

3. `v0.24.x` introduces a new `settings` table that will be created automatically by default. However, if you are running the `server` component with the `VELA_DATABASE_SKIP_CREATION` set to `true`, you will need to manually create this table. 

  - https://github.com/go-vela/server/pull/1110


4. `v0.24.x` introduces a new `report_as` column to the steps table for SCM status reporting for specific steps.

  - https://github.com/go-vela/server/pull/1090


5. `v0.24.x` introduces a new `error` column to the schedules table for compile-time error reporting for scheduled builds.

  - https://github.com/go-vela/server/pull/1108

6. `v0.24.x` removes legacy columns from the database that referenced the old `allow_event` system. These columns should be removed from the database.

7. `v0.24.x` introduces a new `sender_scm_id` field to the `builds` table that will be set to the string identifier that maps to the `sender` within the SCM.

  - https://github.com/go-vela/server/pull/1133

For pre-existing builds, the `sender_scm_id` column should be prepopulated to `0` using the SQL query in the migration script.

## API Responses

There have been several server refactors included in `v0.24.x`. These refactors have changed the structure of API responses for the following resources:
    - `Worker`
    - `Repository`
    - `Build`
    - `Schedule`

These resources will now include any referenced resources in a nested response. For example, `Repository` will now include the `Owner` as a standard `User` resource instead of `user_id`:

```diff
{
  "id": 1,
-  "user_id": 1
+  "owner": {
+	"id": 1,
+	"name": "octocat",
+	"favorites": [],
+	"active": true,
+   "admin": false
+  },
  "org": "github",
  "counter": 10,
  "name": "octocat",
  "full_name": "github/octocat",
  "link": "https://github.com/github/octocat",
  "clone": "https://github.com/github/octocat",
  "branch": "main"
  ...
}
```
