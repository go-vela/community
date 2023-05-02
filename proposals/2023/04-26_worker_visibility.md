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
| **Date**      | April 25th, 2023 |
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

Extend the `worker` table with additional fields that will provide more visibility into a worker's status.

This is a step towards a push vs pull queue model, as well as set the groundwork for admin views/dashboard.

**Please briefly answer the following questions:**

1. Why is this required?

* These additional fields will expand visibility into the current status of each worker.

2. If this is a redesign or refactor, what issues exist in the current implementation?

n/a

3. Are there any other workarounds, and if so, what are the drawbacks?

The current workaround would be frequent additional database calls to the `builds` and `workers` tables to gather various pieces of data to extrapolate `available` or `busy` statuses and other added fields. This could impact database performance on the already frequently updated `builds` table.

There is not currently a workaround for the other proposed statuses.

4. Are there any related issues? Please provide them below if any exist.

n/a

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

Fields to add to the `workers` table:
* `status`
  * `unregistered`: potential default value
  * `available`: at least 1 executor is available
  * `busy`: all executors are busy
  * `maintenance`: unavailable due to maintenance activities
  * `error`
* `last_status_update_at`
* `running_build_ids`
  * An array of currently running `build_id`s
* `last_build_finished_at`

Endpoints:

New: 
* PUT /api/v1/admin/worker admin AdminUpdateWorker

Expand functionality/usage:
* PUT /api/v1/workers/{worker} UpdateWorker

## Implementation

<!--
This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

<!-- Answer here -->
Yes, with help
2. What's the estimated time to completion?

<!-- Answer here -->
1-3 weeks, pending discussion
**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->
Server:
  - https://github.com/go-vela/server/pull/773
  - https://github.com/go-vela/server/pull/772
Types:
  - https://github.com/go-vela/types/pull/277
