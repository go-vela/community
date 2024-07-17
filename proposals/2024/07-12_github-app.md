# Vela GitHub App

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value |
| :-----------: | :---: |
| **Author(s)** | Win San |
| **Reviewers** |       |
| **Date**      | July 12th, 2024 |
| **Status**    | Reviewed |

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

- Create a GitHub App to replace the Vela OAuth App

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

- [GitHub Apps are preferred to OAuth apps](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/differences-between-github-apps-and-oauth-apps) because they use fine-grained permissions, give more control over which repositories the app can access, and use short-lived tokens. 
- These properties can harden security by limiting the damage that could be done if the app's credentials were leaked.

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->
- Permissions and scope could be more granular.

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

- No

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->

- No

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

### Authorization Flow (First Time)
**Current flow:**
- User clicks "GitHub" button in the Vela UI.
- Users gets directed to GitHub Enterprise sign in page.
- User is presented with a [page](https://pages.git.target.com/vela/doc-site/getting-started/logging-in/index.html) outlining resources that the OAuth app wants access to:
    - Personal user data: Email addresses (read-only), profile information (read-only)
    - Organizations and teams: Read access to user's organization, team membership, and private project boards.
    - Repositories: Read and write access to ALL public and private repository data.
- User gets directed back to Vela UI.

**Proposed flow:**
- User clicks "GitHub" button in the Vela UI.
- Users gets directed to GitHub Enterprise sign in page.
- User is presented with a page where they allow the GitHub App to verify their GitHub identity.
- User is presented with a page outlining resources that the GitHub App wants access to and selects which repositories to install the app on:
    - Repositories: Read and write access to either ALL repositories or SELECT repositories.
- User gets directed back to Vela UI.

### Permissions and Scope
What [permissions](https://docs.github.com/en/enterprise-server@3.13/rest/authentication/permissions-required-for-github-apps?apiVersion=2022-11-28) should the GitHub App have?

**List of permission:**
- [Repository permissions for "Checks"](https://docs.github.com/en/enterprise-server@3.13/rest/authentication/permissions-required-for-github-apps?apiVersion=2022-11-28#repository-permissions-for-checks)
- [Repository permissions for "Commit statuses"](https://docs.github.com/en/enterprise-server@3.13/rest/authentication/permissions-required-for-github-apps?apiVersion=2022-11-28#repository-permissions-for-commit-statuses)
- [Repository permissions for "Pull requests"](https://docs.github.com/en/enterprise-server@3.13/rest/authentication/permissions-required-for-github-apps?apiVersion=2022-11-28#repository-permissions-for-pull-requests)
- [Repository permissions for "Webhooks"](https://docs.github.com/en/enterprise-server@3.13/rest/authentication/permissions-required-for-github-apps?apiVersion=2022-11-28#repository-permissions-for-webhooks)
- [User permissions for "Email addresses"](https://docs.github.com/en/enterprise-server@3.13/rest/authentication/permissions-required-for-github-apps?apiVersion=2022-11-28#user-permissions-for-email-addresses)
- [User permissions for "Profile"](https://docs.github.com/en/enterprise-server@3.13/rest/authentication/permissions-required-for-github-apps?apiVersion=2022-11-28#user-permissions-for-profile)

### Code Integration
How does the GitHub app integrate with Vela?

### Rollout
How do we fully switch from OAuth App to GitHub App?

**Phase 1:** Offering improved status checks
- Making a GitHub App with advanced status checks will incentive teams to install the app in their organizational repos. 
- This can be accomplished by using GitHub's [REST API to manage checks](https://docs.github.com/en/enterprise-server@3.13/rest/checks/runs?apiVersion=2022-11-28), which is exclusive to GitHub Apps.
- How do we track GH App installations in the Vela server?

**Phase 2:** Replacing NETRC

**Phase 3:** Moving off of OAuth entirely

## Implementation

<!--
This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

<!-- Answer here -->

Yes

2. What's the estimated time to completion?

<!-- Answer here -->

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

## Questions

**Please list any questions you may have:**

<!-- Answer here -->