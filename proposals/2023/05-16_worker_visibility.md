# Extend Worker Visibility

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value |
| :-----------: | :-: |
| **Author(s)** | Kelly.Merrick, Tim.Huynh, Easton.Crupper, David.Vader, David.May |
| **Reviewers** |  |
| **Date**      | May 16th, 2023 |
| **Status**    | In Progress |

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

* Extend the `worker` database table with additional fields that will provide more visibility into a worker's status, associated timestamps, and currently running builds.

* Give worker the ability to update it's own status:

  | Status | Definition | Build Limit |
  | - | - | - |
  | idle | no builds are running | =>1 |
  | available | at least one build is running and build limit is not met | >1 |
  | busy | build limit is met | =>1 |
  | error | an error has occurred | =>1 |

* Give worker the ability to update the `running_build_ids` array, `last_status_update_at`, `last_build_started_at`, and `last_build_finished_at` timestamps.

**Please briefly answer the following questions:**

1. Why is this required?

* These additional fields will expand visibility into the current status of each worker

* This sets the groundwork for future admin views/dashboards

2. If this is a redesign or refactor, what issues exist in the current implementation?

* n/a

3. Are there any other workarounds, and if so, what are the drawbacks?

* The current workaround would be frequent additional database calls to the `builds` and `workers` tables to gather data to extrapolate `idle`, `available`, or `busy` statuses, associated timestamps, and currently running builds. This could impact database performance on the already frequently updated `builds` table.

4. Are there any related issues? Please provide them below if any exist.

* n/a

## Design

<!--
This section is intended to explain the solution design for the proposal.

NOTE: If there are no current plans for a solution, please leave this section blank.
-->

**Please describe your solution to the proposal. This includes, but is not limited to:**

* new/updated configuration variables (environment, flags, files, etc.)
* performance and user experience tradeoffs
* security concerns or assumptions
* examples or (pseudo) code snippets

### Database 

* Fields to add to the `workers` table:
  * `status`:
    * `idle`
    * `available`
    * `busy`
    * `error`
  * `last_status_update_at`
  * `running_build_ids`
  * `last_build_started_at`
  * `last_build_finished_at`

### Worker

* Update status to `idle` after successful registration, and when no builds are running
* Update status to `available` when at least 1 build is running and build limit is not met
* Update status to `busy` when build limit is met
* Update status to `error` when applicable
* Update current `running_build_ids` when worker picks up a build or finishes a build

## Implementation

<!--
This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

<!-- Answer here -->
* Yes

2. What's the estimated time to completion?

<!-- Answer here -->
* 2-3 weeks

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

PRs:
* https://github.com/go-vela/types/pull/277
* https://github.com/go-vela/server/pull/772
* https://github.com/go-vela/worker/pull/482
* https://github.com/go-vela/community/pull/822
