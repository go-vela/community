---
name: Proposal
about: Submit a potential feature or modification
title: ''
labels: proposal
assignees: ''

---

# Title

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           |       Value        |
| :-----------: |:------------------:|
| **Author(s)** |    Aaron.Hooper    |
| **Reviewers** |                    |
| **Date**      | December 6th, 2023 |
| **Status**    |                    |

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

- Add ability to properly separate application builds in a monorepo.
    - `.vela.yml` files can be found at the root of each application in the repository.
    - Each application will get its own build queue.
    - Previously successful build information will be injected as ENV variables to each build
      - Includes build commit and branch (useful for docker promotions of previously build docker images)
    - Monorepo option will be in repository settings.
    - Choosing which application build queue to view will happen under the `Builds` tab

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->
- Building the entire monorepo on each change is slow
- Build results aren't specific to an application
- Releasing a single application from a monorepo is not possible

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->
- n/a

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->
- The current workaround is an external application that receives the webhooks and triggers deployments in Vela with a payload
  - Drawbacks
    - Deployments are limited compared to builds (i.e. rulesets)
    - Another application needs to be maintained and kept up to date with Vela API changes
    - The link between source control and Vela is compromised

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->
- n/a

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
- Yes

2. What's the estimated time to completion?

<!-- Answer here -->
- 1-2 months

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

## Questions

**Please list any questions you may have:**

<!-- Answer here -->
