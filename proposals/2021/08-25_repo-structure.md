# Repo Structure

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value                                                                  |
| :-----------: | :--------------------------------------------------------------------: |
| **Author(s)** | Jordan.Brockopp                                                        |
| **Reviewers** | David.May, Emmanuel.Meinen, Kelly.Merrick, David.Vader, Matthew.Fevold |
| **Date**      | July 22nd, 2021                                                        |
| **Status**    | Reviewed                                                               |

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

Currently, Vela has an inconsistent structure of how the "core" repos are setup within the [go-vela](https://github.com/go-vela) org.

In this context, a "core" repo includes any repos that require code or dependency changes in order to release the Vela product.

Today, we have the following "core" repos:

* [cli](https://github.com/go-vela/cli)
* [compiler](https://github.com/go-vela/compiler)
* [mock](https://github.com/go-vela/mock)
* [pkg-executor](https://github.com/go-vela/pkg-executor)
* [pkg-queue](https://github.com/go-vela/pkg-queue)
* [pkg-runtime](https://github.com/go-vela/pkg-runtime)
* [sdk-go](https://github.com/go-vela/sdk-go)
* [server](https://github.com/go-vela/server)
* [types](https://github.com/go-vela/types)
* [ui](https://github.com/go-vela/ui)
* [worker](https://github.com/go-vela/worker)

## Context

Before attempting to propose any changes to this structure, it's important to provide context and clarity on why this structure exists.

In the early beginnings of Vela, we thought that each individual capability/service/tool that Vela integrates with would get a unique repo.

The naming convention for these repos would have the syntax of `pkg-<capability>` to denote the code housed in that repo.

To provide an example with the current structure, we'll reference the [pkg-runtime](https://github.com/go-vela/pkg-runtime) repo.

Vela has to integrate with different runtime services or tools in order to execute workloads.

So, the [pkg-runtime](https://github.com/go-vela/pkg-runtime) repo contains all the code necessary to integrate with the different supported runtimes:

* [docker](https://github.com/go-vela/pkg-runtime/tree/master/runtime/docker)
* [kubernetes](https://github.com/go-vela/pkg-runtime/tree/master/runtime/kubernetes)

Likewise for the [pkg-queue](https://github.com/go-vela/pkg-queue) repo, Vela has to integrate with different queue services or tools in order to schedule workloads.

And the [pkg-queue](https://github.com/go-vela/pkg-queue) repo contains all the code necessary to integrate with the different supported queues:

* [redis](https://github.com/go-vela/pkg-queue/tree/master/queue/redis)

However, this structure has some problems, which are detailed below, that leave some room for improvement.

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

This is not required so there is a lot of room for discussion on this.

The intention or goal of this proposal is to accomplish the following:

* condense and simplify the repo structure
* improve consistency among the repos
* ease the burden necessary to contribute functionality to the product
* reduce the level of overhead when attempting to publish new releases

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

## Problem 1

The first problem with this structure is ensuring consistency, both currently and going forward, among the repos in the [go-vela](https://github.com/go-vela) org.

Currently, not all "core" repos follow the `pkg-<capability>` naming convention. i.e. [compiler](https://github.com/go-vela/compiler)

Also, some "core" repos have the packages for different capabilities in the repo itself. i.e. [server](https://github.com/go-vela/server)

If we plan to keep the `pkg-<capability>` structure, those should be pulled out into unique repos following that structure:

* [server/database](https://github.com/go-vela/server/tree/master/database) -> `pkg-database`
* [server/secret](https://github.com/go-vela/server/tree/master/secret) -> `pkg-secret`
* [server/source](https://github.com/go-vela/server/tree/master/source) -> `pkg-source`

## Problem 2

The second problem with this structure is the increased burden, especially for new users, trying to contribute functionality.

Some, or most, of the "core" repos have dependencies on one or more repos.

This means that in order to make certain changes to the product, you have to introduce code changes to multiple repos.

To provide an example with the current structure, we'll reference a recent feature that allows you to customize the pipeline type for a repo.

This functionality was added as a part of [go-vela/community#326](https://github.com/go-vela/community/issues/326).

To achieve this, a new `pipeline_type` field was added to the database which required code changes for multiple repos:

* [types](https://github.com/go-vela/types/pull/188)
* [compiler](https://github.com/go-vela/compiler/pull/199)
* [server](https://github.com/go-vela/server/pull/444)
* [ui](https://github.com/go-vela/ui/pull/421)

This also increases the complexity when attempting to functionally test code changes for new logic.

This is due to needing to modify each repo locally with the code changes for the functionality.

Once that is complete, the application must be rebuilt from the local code changes for each repo.

## Problem 3

The third problem with this structure is level of overhead when attempting to publish new releases for the product.

The way the repos are setup today, there is a dependency tree that must be followed in order to properly release the product.

The way the repos must be released is in groups that use the following order:

1. [types](https://github.com/go-vela/types), [ui](https://github.com/go-vela/ui)
2. [compiler](https://github.com/go-vela/compiler), [mock](https://github.com/go-vela/mock)
3. [pkg-runtime](https://github.com/go-vela/pkg-runtime), [sdk-go](https://github.com/go-vela/sdk-go)
4. [pkg-executor](https://github.com/go-vela/pkg-executor), [pkg-queue](https://github.com/go-vela/pkg-queue)
5. [cli](https://github.com/go-vela/cli), [server](https://github.com/go-vela/server), [worker](https://github.com/go-vela/worker)

After a release is completed for the first group, you must update the code dependencies for the second group before proceeding.

This pattern repeats, for each group, until you get to the last group which has dependencies on most of the preceding groups.

With this structure, we simply have to create more releases (one per repo) than we would if we were to condense the repos.

It also adds more complexity when releasing because you have to remember the order and ensure each repo has the latest dependencies.

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

N/A

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->

N/A

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

## Implementation

<!--
This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

<!-- Answer here -->

2. What's the estimated time to completion?

<!-- Answer here -->

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

## Questions

**Please list any questions you may have:**

<!-- Answer here -->
