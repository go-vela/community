# Platform Settings

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value |
| :-----------: | :-: |
| **Author(s)** | David.Vader |
| **Reviewers** |  |
| **Date**      | March 7th, 2024 |
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

* Build a system that gives platform admins the ability to update platform settings on the fly, without requiring a platform redeploy.
* Potentially support multiple configs and an ability to switch the "active" configuration.

**Please briefly answer the following questions:**

1. Why is this required?

* This allows admins to respond quicker to configuration change requirements and configuration rollbacks.
* It may set the groundwork for feature flags, meaning faster feature rollout and easier feature rollback.

2. If this is a redesign or refactor, what issues exist in the current implementation?

* Platform settings are currently configured using environment variables and are basically treated as insensitive secrets. This is fine and should still be supported moving forward (at least for the next release) but it would be ideal to set these in a provider that can be updated.

3. Are there any other workarounds, and if so, what are the drawbacks?

* n/a

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

## Solution 1: Use A New Database Table

### Database 

* A new table called `settings` with fields such as:
  * `clone-image`
  * `log.level`
  * `log.formatter`
  * `queue.routes`
  * `repo-allowlist`
  * `schedule-allowlist`
  * etc, anything non-sensitive.

### Server

* Logic for creating the initial platform settings on startup (using environment variables?).
* A new API handler `GET /admin/settings` that returns a JSON for platform settings.
* A new API handler `PUT /admin/settings` that receives a JSON payload for updating platform settings.
* Alerting/logging specifically for tracking changes made.
* A source of truth definition for all configurable platform settings.
* A system for mapping platform settings database type to urfave flags.

### CLI

* Commands for managing platform settings.

### Pros

* It is a simple and straight-forward solution that solves the problem.
* It doesn't require any additional libraries or fancy tooling.
* It would allow admins to configure the platform directly in the UI.

### Cons

* It does not scale well without some boiler plate code. It might get annoying to modify database tables every time you add a server configuration.
* It would require full-stack development every time a new field is added that we should be able to update.
* It leans heavily on the `admin` user flag which does not have the best control experience.
* More code (hidden pro?).

## Solution 2: Use a Feature Flag Provider

### Database 

* n/a

### Server

* An authentication process for fetching flags from the provider.
* A pattern for defining feature flags that can be fetched.
* A pattern for actually fetching a feature flag and executing custom code.

### Pros

* It would scale well.
* It would allow us to implement actual feature flags.
* Feature flag providers are getting pretty slick, see [Unleash](https://github.com/Unleash/unleash) and [PostHog](https://posthog.com/).

### Cons

* It would potentially require us to self-host another solution.
* It would make our codebase slightly harder to test.
* It would require another layer of communication.
* It would require admins to visit a different source of truth to apply and set configurations.

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

