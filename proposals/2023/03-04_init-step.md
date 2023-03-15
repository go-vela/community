# Add `InitStep` resource

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
| **Date**      | March 15th, 2023   |
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

Replace the injected "init" step (a pseudo-container) with `InitStep`, a new resource.

Each `Build` will have one or more `InitSteps` that logically group related parts of the log data.

An `InitStep` is a report by a "reporter" about a discrete part of the build setup. The "reporter" would be a logical part of the stack, like the "Pipeline Compiler" or the Runtime. For example, we could have these reporters be discrete initsteps (Reporter: name of the `InitStep`):
- Pipeline Compiler: report info/debug logs
- Docker Runtime: report network setup
- Docker Runtime: report volume setup
- Kubernetes Runtime: Pod YAML

Similar to a `Step`, where the log data is stored in a separate `Log` struct/table, the `InitStep` is also associated with a `Log` that stores the actual log/report.

In the future, each `InitStep`'s `Log` can be given an optional `Mimetype`, which is meant for eventual consumption by the UI to do syntax highlighting. `Mimetype` will be a field on either `InitStep` or `Log` (I'm not sure which yet).

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

2. If this is a redesign or refactor, what issues exist in the current implementation?

Currently, we're abusing the `Step` and `Container` models to allow reporting on build setup (eg when the Docker Runtime initializes the Network and Volumes). This means that we have to check for the special "init" stage/step/container in many places. 

We are also simulating a shell in that init step, printing simulated commands and output even though the runner does not actually run them, so the init step does not accurately represent which part of vela is doing what thing to prepare for running the user's pipeline.

Also, there is only one "init" step. So editing it's `Log` must be managed by one thing, the worker. That way nothing gets lost. If we want to allow different parts of the worker to report status asyncronously, or to capture user-visible logs during compile or other steps, the responsibility for managing that one log entry gets muddy.

<!-- Answer here -->

3. Are there any other workarounds, and if so, what are the drawbacks?

The "init" step is one big workaround. It is an excellent MVP, but we need a way to clean it up.

I looked at adding an `IsInit` bool flag to steps instead of relying on the magic `init` string, but it has to be serializable, and I don't want to add it to the pipeline where users can set it. It has to be serialized when the compiler sends it to the external modification endpoint and when added to the queue for the worker.

<!-- Answer here -->

4. Are there any related issues? Please provide them below if any exist.

In the worker, we frequently need to iterate over the containers for steps/services. But the "init" stage/step is not really a container, so we have to identify which container's are not actually containers so they can be skipped. So far, the worker relies on `Name="init"`, but that does not work in all cases. when the executor is checking trusted+privileged settings in `AssembleBuild`, it checks for `Image="#init"` instead because service containers can be named "init" by users.

This issue is even worse with the kubernetes runtime. There, the number of containers has to be counted and indexed. Given a particular step or service the Kubernetes runtime has to look up which container it needs to edit. So there are many places where that count/index has to be adjusted by one to account for the init step. Then with the injected clone step, figuring out when to add or subtract one or 2 to get the index can be confusing. Also, the kubernetes runtime breaks when running a pipeline with a service named "init" because the container setup is skipped in  one place but not another. That was uncovered by attempting to use it in the executor AssembleBuild test.

So, relying on a magic `"init"` string is surprising and problematic. Relying on `"#init"` as a magic string on step.Image is only marginally better.

<!-- Answer here -->

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

This is description is based on this draft PR: https://github.com/go-vela/types/pull/280

Similar to a `Step`, where the log data is stored in a separate `Log` struct/table,
the `InitStep` is also associated with a `Log` that stores the actual log/report.

Adds structs:

- `pipeline.InitStep`: a new struct meant to be used in communication between worker <=> server to report about "InitStep".
- `library.InitStep`: represents a discrete "InitStep" report (actually, just the metadata for it - the report is stored in a `Log`)
- `database.InitStep`: to persist the library.InitStep

We explicitly do NOT want to add a `yaml.InitStep` struct, as this is only for server/worker to report on build init. It shouldn't be exposed in the user's pipeline.

Also add these fields:

- `pipeline.Build.InitSteps`: an `InitStepSlice` / a slice of `pipeline.InitStep` structs
- `library.Log.InitStepID`: to associate a `library.Log` with a `library.InitStep` instead of a step or a service.
- `database.Log.InitStepID`: to associate a `database.Log` with a `database.InitStep` instead of a step or a service.

Future enhancement, also add `Mimetype` to the `Log` structs.

### Server

This is description is based on this draft PR: https://github.com/go-vela/server/pull/779

#### Server API

Add `api/initstep` package with endpoints for `InitStep` that mirror the `Step` endpoints:

```
- UpdateInitStep:     PUT    /api/v1/admin/initstep
- CreateInitStep:     POST   /api/v1/repos/:org/:repo/builds/:build/initsteps
- ListInitSteps:      GET    /api/v1/repos/:org/:repo/builds/:build/initsteps
- GetInitStep:        GET    /api/v1/repos/:org/:repo/builds/:build/initsteps/:initstep
- UpdateInitStep:     PUT    /api/v1/repos/:org/:repo/builds/:build/initsteps/:initstep
- DeleteInitStep:     DELETE /api/v1/repos/:org/:repo/builds/:build/initsteps/:initstep
- CreateInitStepLog:  POST   /api/v1/repos/:org/:repo/builds/:build/initsteps/:initstep/logs
- GetInitStepLog:     GET    /api/v1/repos/:org/:repo/builds/:build/initsteps/:initstep/logs
- UpdateInitStepLog:  PUT    /api/v1/repos/:org/:repo/builds/:build/initsteps/:initstep/logs
- DeleteInitStepLog:  DELETE /api/v1/repos/:org/:repo/builds/:build/initsteps/:initstep/logs
- PostInitStepStream: POST   /api/v1/repos/:org/:repo/builds/:build/initsteps/:initstep/stream
```

Add `middleware/initstep` package to inject `InitStep` in the `*/initsteps/*` APIs.

Also include the endpoints in the mock server.

#### Server Database

Add `database/initstep` package for storing the new `InitStep` object:

- `CountInitSteps`: gets the count of all InitSteps
- `CountInitStepsForBuild`: gets the count of InitSteps by build ID
- `CreateInitStep`: creates a new InitStep.
- `DeleteInitStep`: deletes an existing InitStep.
- `GetInitStep`: gets a InitStep by ID.
- `GetInitStepForBuildgets`: gets an InitStep by build ID and InitStep number.
- `ListInitSteps`: gets a list of all InitSteps.
- `ListInitStepsForBuild`: gets a list of InitSteps by build ID.
- `UpdateInitStep`: updates an existing InitStep.

And methods for the InitStep Logs:
- GetLogForInitStep: gets a log by init step ID from the database.

#### End-User visible logging in the server

Today, only the worker can safely add details to the init step log. With this change,
the server can also add end-user visible logging, especially in places like the compiler
where it would be helpful to highlight compile errors in the UI/CLI.

In the Build and Webhook API endpoints, we can create the `InitStep` and `Log` as soon as we
have the `RepoID` and `Build.Number`. Failures before this point can't be associated
with a particular build, so they cannot bubble up to the end-user. The endpiont will
then pass in a `library.Log` to the compiler that it creates specifically for the compiler.
It can also record relevant its own log messages when handling requests like:

- `CreateBuild`
- `RestartBuild`
- `PostWebHook`

In the compiler, we can fill in an `InitStep` `Log` for end-user visible logging.
The compiler does not (and probably should not) have access to the database.
So, add `WithLog(*library.Log) Engine` to the `compiler.Engine` interface. To use it
the endpoint methods that create the compiler will pass in a `library.Log` type that the
compiler can use to report any end-user visibile log messages. Once the compiler finishes,
the API endpoint will handle saving it to the database.

For things that use CompileLite, passing in a `Log` doesn't make sense because it is
not part of a build. So, using `WithLog` should be optional.

### Worker

The executor and runtime can create `InitStep` + `Log` entries wherever it makes sense
without worrying about retrieving prior logged steps. The worker could also add more
InitSteps not just before a Build, but also when preparing to run individual steps
or services.

Nothing in the worker should check for these magic strings any more: `"init"`, `"#init"`

### SDK

Needs support for the new Server endpoints.

### CLI

Add commands for InitStep similar to Step:

```
  1. Get initsteps for a repository.
    $ vela get initsteps --org MyOrg --repo MyRepo --build 1
  2. Get initsteps for a repository with wide view output.
    $ vela get initsteps --org MyOrg --repo MyRepo --build 1 --output wide
  3. Get initsteps for a repository with yaml output.
    $ vela get initsteps --org MyOrg --repo MyRepo --build 1 --output yaml
  4. Get initsteps for a repository with json output.
    $ vela get initsteps --org MyOrg --repo MyRepo --build 1 --output json
  5. Get initsteps for a build when config or environment variables are set.
    $ vela get initsteps --build 1

  1. View initstep details for a repository.
    $ vela view initstep --org MyOrg --repo MyRepo --build 1 --initstep 1
  2. View step details for a repository with json output.
    $ vela view initstep --org MyOrg --repo MyRepo --build 1 --initstep 1 --output json
  3. View step details for a repository config or environment variables are set.
    $ vela view initstep --build 1 --initstep 1

  3. View logs for an initstep.
    $ vela view log --org MyOrg --repo MyRepo --build 1 --initstep 1

```

Any `vela get log` requests for a build should automatically include the
initstep logs without change by virue of re-using the `Log` types for this.

### UI

The UI needs to stream initsteps just like it does for steps and services.
This can be added after all of the backend work for server, worker, sdk, etc.

The init steps should be presented separately (somehow) from the normal steps
because they represent something that Vela is doing, not the output of the actual
pipeline's steps. Many other CI services do not make much of a distinction here.
So, at first we could just include the list of initsteps in the same view as
the other steps.

Eventually, I would love to see syntax highlighting when an InitStep includes
a known Mimetype like YAML (For example, when the Kubernetes runtime logs the Pod
just before creating it, or perhaps when the pod gets updated). Doing that is
beyond the scope of this proposal, but that is part of what I would like to
facilitate eventually by separating the InitStep logs from standard Step logs now.

### Deprecating the legacy "init" step as a Step

The current init step should be deprecated for at least one release before removing it.
Any logging in the worker that currently gets put on that step should continue to
be logged there for at least one release. The worker should also create an InitStep
and log these things there. Then, 

If it takes an extra release for the UI to get these features, that's ok because the
old init step will still be getting all of the info that was logged before. Once the
UI work is complete, additional logs will become available from the compiler and
anything else that uses the InitStep to send logs to the end-user.

Once we're satisfied that log data in InitSteps sufficiently covers everything
that is currently logged via the pseudo-step named "init", then we can stop injecting
that init step (and the init stage) and all of the special-case handling of the "init"
container/stage/step.

## Implementation

<!--
This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

<!-- Answer here -->

yes for the go code in Types, Server, Worker, and sdk-go.
I have not used the CLI so far, but copying the step actions/commands to initstep seems
simple enough.

The UI, however, is beyond me at this point. I need someone familiar eith elm to handle the UI
after all the other components have been merged.

2. What's the estimated time to completion?

<!-- Answer here -->

I have a PoC branch for types and server. So, if this gets accepted as is, I don't think wrapping it up will take long.

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

## Questions

**Please list any questions you may have:**

<!-- Answer here -->
