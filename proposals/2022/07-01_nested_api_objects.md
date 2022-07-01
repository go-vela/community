# Nested API Objects

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value |
| :-----------: | :---: |
| **Author(s)** | Jordan.Brockopp, Easton.Crupper |
| **Reviewers** |       |
| **Date**      |       |
| **Status**    |       |

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

This change would be considered a redesign/refactor to modify the behavior of what information is returned by the API.

The idea is to no longer return the ID fields for resources that have a relationship with one another.

For example, today when you query a repo (`GET /api/v1/repo/:org/:repo`), a `user_id` field is returned in the response.

The `user_id` field contains the primary key for a row in the `users` table that represents the "owner" of the repo.

Unfortunately, that `user_id` field isn't providing value for administrators or consumers (end-users).

At this time, the below table contains a list of all resources that have a ID field nested under them:

| Resource           | Fields |
| :-----------: | :---: |
| [library.Build](https://github.com/go-vela/types/blob/master/library/build.go#L16-L51) | `repo_id`, `pipeline_id` |
| [library.Deployment](https://github.com/go-vela/types/blob/master/library/deployment.go#L13-L28) | `repo_id`       |
| [library.Hook](https://github.com/go-vela/types/blob/master/library/hook.go#L11-L29) | `repo_id`, `build_id` |
| [library.Log](https://github.com/go-vela/types/blob/master/library/log.go#L14-L25) | `repo_id`, `build_id`, `service_id`, `step_id` |
| [library.Pipeline](https://github.com/go-vela/types/blob/master/library/pipeline.go#L11-L31) | `repo_id` |
| [library.Repo](https://github.com/go-vela/types/blob/master/library/repo.go#L11-L38) | `user_id` |
| [library.Service](https://github.com/go-vela/types/blob/master/library/service.go#L16-L35) | `repo_id`, `build_id` |
| [library.Step](https://github.com/go-vela/types/blob/master/library/step.go#L16-L36) | `repo_id`, `build_id` |

<!--
Provide your description here.
-->

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

4. Are there any related issues? Please provide them below if any exist.

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

### Option 1 - use existing v1

### Option 2 - create new v2

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
