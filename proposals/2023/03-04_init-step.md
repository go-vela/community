# Add `InitStep` resource

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value       |
| :-----------: | :---------: |
| **Author(s)** | Jacob Floyd |
| **Reviewers** |             |
| **Date**      |             |
| **Status**    | WIP         |

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



### Server

### Worker

### SDK

### CLI

### UI

## Implementation

<!--
This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

<!-- Answer here -->

yes for server and worker. And maybe CLI.

I need someone familiar eith elm to handle the UI.

2. What's the estimated time to completion?

<!-- Answer here -->

I have a PoC branch for types and server. So, if this gets accepted as is, I don't think wrapping it up will take long.

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

## Questions

**Please list any questions you may have:**

<!-- Answer here -->
