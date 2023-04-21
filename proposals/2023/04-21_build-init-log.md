# Use `Log` directly for Build Init instead of using a `Step`'s `Log`

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value              |
| :-----------: | :----------------: |
| **Author(s)** | Jacob Floyd        |
| **Reviewers** |                    |
| **Date**      | April 21st, 2023   |
| **Status**    | Waiting for Review |

<!--
If you're already working with someone, please add them to the proper author/reviewer category.

If not, please leave the reviewer category empty and someone from the Vela team will assign it to themself.

Here is a brief explanation of the different proposal statuses:

1. Reviewed: The proposal is currently under review or has been reviewed.

2. Accepted: The proposal has been accepted and is ready for implementation.

3. In Progress: An accepted proposal is being implemented by actual work.

NOTE: The design is subject to change during this phase.

4. Cancelled: While or before implementation the proposal was cancelled.

NOTE: This can happen for a multitude of reasons.

5. Complete: This feature/change is implemented.
-->

## Background

<!--
This section is intended to describe the new feature, redesign or refactor.
-->

**Please provide a summary of the new feature, redesign or refactor:**

<!--
Provide your description here.
-->

Drop the magic "init" `Step` (and the "init" `Stage` for stages pipelines). Instead, associate a `Log` directly with a Build and use that to store "init" logs.

Unlike the [`InitStep`](https://github.com/go-vela/community/pull/771) proposal, this proposal only allows for one "init" `Log` per `Build`.

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

- To drop use of magic `"init"` and `"#init"` strings in the worker and remove all of the associated TODOs throughtout the codebase.
- So that the worker does not have special-casing logic to avoid handling the pseudo-`Step` that build init logs are associated with.
    - The Kubernetes runtime has to access containers in a `Pod` using indexes, which is confusing thanks to the "init" and "clone" Steps, especially since "init" is never a container. And if also switching between 0-based and 1-based indexes, then the off-by-one bugs become very confusing and require extensive comments (subtract 2 from the index here, subtract 1 from the index there, subtract 3 over yonder...).

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

Currently, we're abusing the `Step` and `Container` models to allow reporting on build setup (eg when the Docker Runtime initializes the Network and Volumes). This means that we have to check for the special "init" stage/step/container in many places.

The more special-casing around "init" and "#init" strings is an infestation that spreads. It encourages duplicate checks for those strings throughout the code base, and we keep finding more edge cases where that special casing is missing.

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

The "init" step is one big workaround. It is an excellent MVP, but we need a way to clean it up.

I looked at adding an `IsInit` bool flag to steps instead of relying on the magic `init` string, but it has to be serializable, and I don't want to add it to the pipeline where users can set it. It has to be serialized when the compiler sends it to the external modification endpoint and when added to the queue for the worker.

I also looked at introducing an `InitStep` resource to allow multiple `InitStep`s (each with their own log) per `Build`, but that had probable DB performance implications and so was closed: https://github.com/go-vela/community/pull/771

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->

In the worker, we frequently need to iterate over the containers for steps/services. But the "init" stage/step is not really a container, so we have to identify which container's are not actually containers so they can be skipped. So far, the worker relies on `Name="init"`, but that does not work in all cases. when the executor is checking trusted+privileged settings in `AssembleBuild`, it checks for `Image="#init"` instead because service containers can be named "init" by users.

This issue is even worse with the kubernetes runtime. There, the number of containers has to be counted and indexed. Given a particular step or service the Kubernetes runtime has to look up which container it needs to edit. So there are many places where that count/index has to be adjusted by one to account for the init step. Then with the injected clone step, figuring out when to add or subtract one or 2 to get the index can be confusing. Also, the kubernetes runtime breaks when running a pipeline with a service named "init" because the container setup is skipped in  one place but not another. That was uncovered by attempting to use it in the executor AssembleBuild test.

So, relying on a magic `"init"` string is surprising and problematic. Relying on `"#init"` as a magic string on step.Image is only marginally better.

## Design

<!--
This section is intended to explain the solution design for the proposal.

NOTE: If there are no current plans for a solution, please leave this section blank.
-->

**Please describe your solution to the proposal. This includes, but is not limited to:**

* new/updated endpoints or url paths
* new/updated configuration variables (environment, flags, files, etc.)
* performance and user experience tradeoffs
* security concerns or assumptions
* examples or (pseudo) code snippets

<!-- Answer here -->

### Types

No changes required for `database.Log` or `library.Log` structs since they are already associated with a `BuildID`.
Validation, however, will need to allow for both `StepID` and `ServiceID` to be `NULL`.

Drop this check in `Validate()` in `database/log.go` since no `StepID` and no `ServiceID` means we have a `Build` `Log`:
<!-- https://github.com/go-vela/types/blob/main/database/log.go#L129-L132 -->
```
	// verify the has StepID or ServiceID field populated
	if l.StepID.Int64 <= 0 && l.ServiceID.Int64 <= 0 {
		return ErrEmptyLogStepOrServiceID
	}
```
<!-- https://github.com/go-vela/types/blob/main/database/log_test.go#L298-L299
And this needs to be failure: false
```
		{ // no service_id or step_id set for log
			failure: true,
```
-->

No changes are needed in the `pipeline` or the `yaml` layers because there is not a `pipeline.Log` nor a `yaml.Log`, so the relationship between `Build` and `Log` does not need to be represented.

### Server

#### Server - Database

Adjust `CreateLog`, `UpdateLog`, and `DeleteLog` to handle the build `Log` (adding a `logger.Tracef` entry and an `Error` specific to the Build Log).

Add `database/log/get_build.go` with:
```
// GetLogForBuild gets a log by build ID from the database.
func (e *engine) GetLogForBuild(b *library.Build) (*library.Log, error) {
```

See if we can make the build log sort before the step logs in `ListLogsForBuild`:
<!-- https://github.com/go-vela/server/blob/main/database/log/list_build.go#L37-L40 -->
```
	err = e.client.
		Table(constants.TableLog).
		Where("build_id = ?", b.GetID()).
		Order("step_id ASC").
```

Sqlite sorts with `NULLS FIRST` by default, and Postgres sorts with `NULLS LAST` (see: [How Are NULLs Sorted by Default?](https://learnsql.com/blog/how-to-order-rows-with-nulls/)).
So this query is inconsistent between databases. And, the service logs returned by this query are not sorted. It would be nice to have this method return the logs sorted in a consistent manner.
I would like to see the Build Log, then Step Logs sorted by step_id, then Service Logs sorted by service_id. Any suggestions on how to do that with gorm?

We might need some kind of index or constraint on the `Log` table so that rows with NULL step_id and NULL service_id must have a unique build_id. But, I'm not sure how to create such an index/constraint.

#### Server - API

We will need new endpoints to create/get/update/delete the build init logs (using `/logs` to be consistent with the other endpoints that use `/logs` to get the `Log` associated with a `Step` or `Service`).

```
- CreateBuildInitLog:  POST   /api/v1/repos/:org/:repo/builds/:build/init/logs
- GetBuildInitLog:     GET    /api/v1/repos/:org/:repo/builds/:build/init/logs
- UpdateBuildInitLog:  PUT    /api/v1/repos/:org/:repo/builds/:build/init/logs
- DeleteBuildInitLog:  DELETE /api/v1/repos/:org/:repo/builds/:build/init/logs
```

The API path ends up being the only place that needs this "init" string. This avoids ambiguity between retrieve all logs (including step and service logs) and the one `Log` associated only with the `Build` (not a `Step` or `Service`).

These endpoints need to be registered in the router<!-- `LogBuildInitHandlers` in `router/log.go` -->, but no new middleware is required.

Also include the endpoints in the mock server.

#### Server - Compiler

The compiler needs to deprecate and stop adding the `InitStage` and `InitStep`. That includes:
- remove `compiler/native/initialize{,_test}.go`,
- remove `Init*` from the `Engine` interface in `compiler/engine.go`,
- remove the `Init*()` calls from `compiler/native/compile.go`, and
- remove the magic `"init"` string special-casing in `compiler/native/validate.go`.

Backwards compatibility is a concern with the compiler. Referencing old builds/pipelines should work just fine as it includes the "init" step and stage, thus preserving references to the underlying log. But, recompiling the pipeline will create a slightly different pipeline--one without the injected "init" stage and/or step. The worker will no longer handle the old pipeline with the injected "init" step, so any re-runs MUST re-compile the pipeline.

Does re-running a build ALWAYS re-compile the pipeline?

#### Server - Queue

Once the worker is upgraded to stop special-casing the `"init"` and `"#init"` strings, any queued builds that have that `"init"` pseudo-`Step` will fail.

Also, when upgrading, we need to ensure that:
- the queue is empty, or
- all queued builds get re-compiled (or more particularly `item.Pipeline` which is a `types.pipeline.Build` and includes the injected init step) before execution in the worker starts.

To make the worker backwards compatible, could we trigger the rebuild in `server.queue.redis.Pop()` after the item is created?
Or perhaps we'll just need to include queued build migrations in the migration script.

### Worker

Save init logs to the Build `Log` instead of the magic init step.
Nothing in the worker should check for these magic strings any more: `"init"`, `"#init"`

We might need something that rejects an old build when popped off the queue if it is an older version that includes the injected "init" `Step`.

We might need to spread this change over a couple of versions.
- In one version, we start logging to the Build `Log` and deprecate support for the magic "init" step.
- In the next (or a future) version, we remove support for the magic string checking.

In my previous [`InitStep` proposal](https://github.com/go-vela/community/pull/771), I suggested logging to both places for at least one version while waiting for the UI to catch up.
However, the primary objection to that proposal was increased database storage costs. So, this proposal recommends a hard break the worker will only create the init logs in one place.
The backwards compatibility is achieved by continuing to ignore Stages/Steps with the magic "init" or "#init" strings. Even if those strings are present, the worker will still send
the build init logs to the server via the new build init logs endpoints.

### SDK

Needs support for the new Server endpoints.

### CLI

Any `vela get log` requests for a build should automatically include the
initstep logs without change by virue of re-using the `Log` types for this.

We also need a new command to retrieve the "init" logs for a Build.

One way to provide this would be an `--init` flag to specify we only want the build's init log, not all of the logs for the build:
```
  3. View init logs for a build.
    $ vela view log --org MyOrg --repo MyRepo --build 1 --init
```

Another alternitve is adding a separate "initlog" subcommand:
```
  3. View init logs for a build.
    $ vela view initlog --org MyOrg --repo MyRepo --build 1
```

I have not used the vela CLI much, so I'm not sure which option would be more ergonomic, or if we need something else entirely.

### UI

The UI needs to stream the build init log just like it does for steps and services. This can be shown as a step that comes before any stages or steps.

The UI must handle older builds gracefully as they will not have this build init log.

The UI must be upgraded in tandem with the worker so that when the worker starts logging build init logs with the new endpoint, those logs can be displayed in the UI.

## Implementation

<!--
This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

<!-- Answer here -->

Yes for the go code in Types, Server, Worker, sdk-go, and CLI.

The UI, however, is beyond me at this point. I need someone familiar eith elm to handle the UI.

2. What's the estimated time to completion?

<!-- Answer here -->

TBD

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

## Questions

**Please list any questions you may have:**

<!-- Answer here -->

This is a summary of the questions listed above:
- How can we add an index or constraint to the database so that each build can only have one init Log? (see "Server" > "Server - Database" above)
- Does re-running a build ALWAYS re-compile the pipeline? (see  "Server" > "Server - Compiler" above)
- Do we need to trigger a re-compile for old builds on the queue when the worker Pops them off the queue? (see  "Server" > "Server - Queue" above)
- Do we need to migrate / re-compile any builds on the queue in an upgrade migration script? (see  "Server" > "Server - Queue" above)
- Does anyone have bandwidth to help with the UI piece? (see "UI" above)
