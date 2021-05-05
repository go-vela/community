# Multi worker builds

<!-- Please leave this commented out section.

The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in YYYY-MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value            |
| :-----------: | :--------------: |
| **Author(s)** |  Neal.Coleman    |
| **Reviewers** |                  |
| **Date**      | April 27th, 2021 |
| **Status**    | Pending          |

<!-- Please leave this commented out section.

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

The goal of this feature is to enable builds to fan-in-fan-out builds across workers. Today, builds are always isolated to a single worker which means build throughput can be an issue when you're trying to multiplex a large amount of work.

Sometimes running all the work on a single worker isn't efficient because of the size limitations of the worker. The ability to design a pipeline that can spawn a set of pipelines across workers can solve this by taking advantage of either different pools of worker sizes or even different flavors of a worker to ensure the code is compatible with multiple operating system environments.

Some examples of this work in other CI platforms:

- [Drone](https://docs.drone.io/pipeline/configuration/#multiple-pipelines)
- [Actions](https://docs.github.com/en/actions/managing-workflow-runs/using-the-visualization-graph)
- [CircleCI](https://circleci.com/docs/2.0/parallelism-faster-jobs/)

_Note: GitHub Actions is hard to point to a single point in the docs. It's somewhat inherent to the design of the pipeline spec._

<!-- Please leave this commented out section.

This section is intended to describe the new feature, redesign or refactor.
-->

**Please provide a summary of the new feature, redesign or refactor:**

<!-- Please remove this commented out section.

Provide your description here.
-->

A user should be able to write a workflow that can expand the build out to execute on a set of Vela workers.


**Please briefly answer the following questions:**

1. Why is this required?

<!-- Please remove this commented out section.

Provide your answer here.
-->

This feature is required to allow builds to spawn a set of pipelines across workers. Some core use cases:

* build that requires high amounts of throughput
* builds that want independent failures of pipelines within a single build
* builds that require code to be tested on more than just one runtime or platform

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Please remove this commented out section.

Provide your answer here.
-->

Not a redesign or refactor. N/A

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Please remove this commented out section.

Provide your answer here.
-->

Resource constraint builds could be solved today by having very large workers that can handle the workload locally in parallel. However, there are known limitations with the number of stages that can execute in parallel.

4. Are there any related issues? Please provide them below if any exist.

<!-- Please remove this commented out section.

Provide your answer here.
-->

N/A

## Design

<!-- Please leave this commented out section.

This section is intended to explain the solution design for the proposal.

NOTE: If there are no current plans for a solution, please leave this section blank.
-->

**Please describe your solution to the proposal. This includes, but is not limited to:**

* new/updated endpoints or url paths

The goal will be to keep the feature isolated and available to all existing endpoints.

This will reduce/limit the number of changes to existing CLI/UI workflows.

* new/updated configuration variables (environment, flags, files, etc.)

New boolean compiler/server variable to hide compile process behind a feature flag: `VELA_COMPILER_FAN_IN_FAN_OUT`

* performance and user experience tradeoffs

The other assumption at the core of implementation detail is this model is designed to not let workers share a workspace between workflows. If that is the desired feature it should be an enhancement. Workflows that can not be highly parallelizable are not recommended for this design.

Another recommendation when writing this feature is to create the ruleset available as a pipeline resource. By allowing it to exist at the root of the YAML file would allow pipelines not to reduce the number of lines of code by setting the ruleset once and not in each step for evaluation.

* security concerns or assumptions

At this time the only concern is around cluster performance. Because the size of any given Vela cluster is unknown without any limits in place a single repo could spawn enough concurrent builds to take up all available workers within the cluster. This could create pending builds for any other scheduled build behind this build pushed to the server.

* examples or (pseudo) code snippets

The following are examples of how a user would leverage writing a pipeline with the feature:

One: _a `.vela.yml` at the base of a repo_

```yaml
---
version: 1

metadata:
	name: build on worker A

steps:
  - name: hello world
  	image: alpine:latest
    commands:
    	- echo "Hello, World"

---
version: 1

metadata:
	name: build on worker B

steps:
  - name: hello world
  	image: alpine:latest
    commands:
    	- echo "Hello, World"

```

Two: _a `.vela/` at the base of a repo_

```yaml
---
# File in .vela/hello1.yml
version: 1

metadata:
	name: build on worker A

steps:
  - name: hello world
  	image: alpine:latest
    commands:
    	- echo "Hello, World"
```

```yaml
---
# File in .vela/hello2.yml
version: 1

metadata:
	name: build on worker B

steps:
  - name: hello world
  	image: alpine:latest
    commands:
    	- echo "Hello, World"

```

<!-- Please remove this commented out section.

Provide your answer here.
-->

## Implementation

<!-- Please leave this commented out section.

This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

<!-- Please remove this commented out section.

Provide your answer here.
-->

TBD

2. What's the estimated time to completion?

<!-- Please remove this commented out section.

Provide your answer here.
-->

3-4 weeks

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Please remove this commented out section.

Provide your answer here.
-->

* add the a name metadata field to types
* add pipeline level ruleset support _stretch goal_
* update compiler to support sets of pipelines
* update server to put support new compile functionality
* update server to pull `.vela/` directories
* Spike: evaluate need to update UI build page to support different visualization of the build
* Spike: evaluate need to update CLI for fan-in-fan-out support
* add documentation to the vela site

## Questions

**Please list any questions you may have:**

<!-- Please remove this commented out section.

Provide your answer here.
-->
