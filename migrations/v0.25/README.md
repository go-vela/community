# v0.25 migration

> NOTE: This applies when upgrading to the latest `v0.25.x` release.

**No database changes are required** when migrating from Vela version [v0.24](../../releases/v0.24.md) to [v0.25](../../releases/v0.25.md).

## API Responses

A continuation of server refactors is included in `v0.25.x`. These refactors have changed the structure of API responses for the following resources:
    - `Hook`

These resources will now include any referenced resources in a nested response. For example, `RepoID` has been replaced with the `Repo` object:

```diff
{
  "id": 1,
-  "repo_id": 1
+  "repo": {
+	"id": 1,
+ ...
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
