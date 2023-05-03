# Build Tokens

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           |                               Value                               |
| :-----------: |:-----------------------------------------------------------------:|
| **Author(s)** | Erik Pearson, James Christensen, Jordan Sussman & Jordan Brockopp |
| **Reviewers** |                                                                   |
| **Date**      |                         March 16th, 2023                          |
| **Status**    |                            In Progress                            |

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

Currently, Vela has no mechanism to enable periodically triggering builds for a repo.

The idea is to add support for running builds for a repo on a configurable cadence to Vela.

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

There are a wide variety of use-cases that exist for this today so we'll name a few:

* automate unit, integration and/or functional tests for a product on a regular cadence
* automate building a product that requires pulling dependencies from an upstream repo on a regular cadence
* automate deploying a product that requires renewing a token which expires on a regular cadence

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

N/A

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

Use another mechanism (`cron`, CI/CD system etc.) to automatically trigger builds for a repo

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->

https://github.com/go-vela/community/issues/538

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

### Database

Add a new `schedules` table to the database and each row will be associated with a repo via `repo_id` column.

The plan is to support multiple `schedules` per repo and the `name` column will be the unique identifier.

This follows a similar approach as `builds` with the `number` column being the unique identifier per repo.

Also, the `name` column will be used to control behavior in the pipeline via `ruleset`.

The `entry` column will hold the cadence for how often to trigger a build for the repo.

This will follow syntax for [cron](https://en.wikipedia.org/wiki/Cron) but won't support every cron capability.

```sql
CREATE TABLE
IF NOT EXISTS
schedules (
	id           SERIAL PRIMARY KEY,
	repo_id      INTEGER,
	name         VARCHAR(100),
	entry        VARCHAR(100),
	created_at   INTEGER,
	created_by   VARCHAR(250),
	updated_at   INTEGER,
	updated_by   VARCHAR(250),
	scheduled_at INTEGER,
	UNIQUE(repo_id, name)
);
```

### UI

A tab (a.k.a. page) would be created for the new `schedules` functionality.

This will follow a similar approach as seen on the `secrets` tab for a repo:

https://vela.example.com/-/secrets/native/repo/MyOrg/MyRepo

| Name                 | Entry     | Created At     | Created By     | Updated At     | Updated By    | Active |
| -------------------- | --------- | -------------- | -------------- | -------------- | ------------- | ------ |
| [unit_test]()        | 0 2 * * * | 5 minutes ago  | JordanBrockopp | 1 minute ago   | JordanSussman | false  |
| [integration_test]() | 0 4 * * * | 10 minutes ago | JordanBrockopp | 5 minutes ago  | JordanSussman | true   |
| [nightly]()          | 0 0 * * * | 15 minutes ago | JordanBrockopp | 10 minutes ago | JordanSussman | true   |

When you click on one of those `schedules`, you'll navigate to a "View/Edit" page for the schedule like `secrets`:

https://vela.example.com/-/secrets/native/repo/MyOrg/MyRepo/MySecret

On this page, you should be able to perform the following interactions with that schedule:

* update the `name` field
* update the `entry` field
* update the `active` field (enable/disable)
* remove the schedule

The following URIs would be created for the new `schedules` functionality:

```
// NEW   /:org/:repo/add-schedule        | create a new schedule for a repo
// LIST  /:org/:repo/schedules           | list schedules for a repo
// EDIT  /:org/:repo/schedules/:schedule | view/update/remove an existing schedule
```

### API

The following endpoints would be created for the new `schedules` functionality:

```
// POST   /api/v1/schedules/:org/:repo           | create a new schedule for a repo
// GET    /api/v1/schedules/:org/:repo           | list schedules for a repo
// GET    /api/v1/schedules/:org/:repo/:schedule | view an existing schedule for a repo
// PUT    /api/v1/schedules/:org/:repo/:schedule | update an existing schedule for a repo
// DELETE /api/v1/schedules/:org/:repo/:schedule | remove an existing schedule for a repo
```

### CLI

The following commands would be created for the new `schedules` functionality:

```
vela add schedule --org MyOrg --repo MyRepo --name nightly --entry '0 0 * * *'    | create a new schedule for a repo
vela get schedule --org MyOrg --repo MyRepo                                       | list schedules for a repo
vela view schedule --org MyOrg --repo MyRepo --name nightly                       | view an existing schedule for a repo
vela update schedule --org MyOrg --repo MyRepo --name nightly --entry '0 0 * * *' | update an existing schedule for a repo
vela remove schedule --org MyOrg --repo MyRepo --name nightly                     | remove an existing schedule for a repo
```

### YAML

A `schedule` event type will be added under `ruleset` for the new `schedules` functionality.

This will instruct Vela to execute the step if the `event` for a build is `schedule`.

The `target` field under a `ruleset` can be used to control what schedule a step will run for.

This will enable further customization when multiple schedules are setup for a repo.

```yaml
steps:
  - name: unit-test
    ruleset:
      event: [ deployment, schedule ]
      target: [ unit_test, integration_test, nightly ]
    commands:
      - echo "I run when a schedule with the name 'unit_test', 'integration_test' or 'nightly' is executed"

  - name: integration-test
    ruleset:
      event: [ deployment, schedule ]
      target: [ integration_test, nightly ]
    commands:
      - echo "I run when a schedule with the name 'integration_test' or 'nightly' is executed"
      
  - name: publish
    ruleset:
      event: [ deployment, schedule ]
      target: [ nightly ]
    commands:
      - echo "I run when a schedule with the name 'nightly' is executed"
      
  - name: notification
    ruleset:
      event: [ deployment, schedule ]
    commands:
      - echo "I run when any schedule is executed"
```

> NOTE: This leverages the existing `target` field but we could add a new field i.e. `schedule_target`.
>
> However, we elected to use `target` to reduce the complexity and code required to deliver on this feature.
>
> Also, this would offer the ability to trigger a build via a `schedule` or a `deployment` for various use cases.

## Implementation

<!--
This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

Users will be able to manage schedules for a repo via the UI, API and CLI.

Only repos within the allowlist (`VELA_SCHEDULE_ALLOWLIST`) will be able to create schedules.

Only users with `admin` access to a repo will be able to manage schedules for that repo.

Vela will support having multiple schedules setup for a repo.

By default, when a repo is enabled in Vela it will not create a `schedule`.

When a `schedule` is created for a repo, a new row will be added to a table for it.

By default, Vela will support a minimum of `1hr` as the cadence for the schedule.

However, similar to build limits and timeouts, Vela will support configuring that minimum via environment variable.

The build rate limit per repo will protect the system from a DDOS effect for frequent schedules.

Each build created from a schedule will leverage the configured `branch` of a repo for the pipeline (`.vela.yml`).

The server will create a thread on startup to scan the table containing all schedules and spawn builds accordingly.

The thread will be setup to scan the table every half of the minimum supported cadence i.e. `1hr` / 2 = `30m`

The server will use [jitter](https://en.wikipedia.org/wiki/Jitter) to minimize the amount of builds created at the same time for all repos.

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

<!-- Answer here -->

Yes, Cargill will add the code required to deliver on this feature

2. What's the estimated time to completion?

<!-- Answer here -->

TBD - we assume this will require at least a couple weeks to write the code and test the functionality

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

* https://github.com/go-vela/community/issues/538
* https://github.com/go-vela/server/pull/833
* https://github.com/go-vela/server/pull/834
* https://github.com/go-vela/server/pull/836

## Questions

1. Should we create a new field (i.e. `schedule`/`schedule_target`) or reuse existing `target` field in the `ruleset`?

2. Should we prevent running multiple builds for the same schedule for a repo?