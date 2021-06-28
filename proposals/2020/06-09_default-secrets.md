# Configurable Default Secrets

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value                                                                                |
| :-----------: | :----------------------------------------------------------------------------------: |
| **Author(s)** | Neal.Coleman, Jordan.Brockopp                                                        |
| **Reviewers** | Neal.Coleman, David.May, Emmanuel.Meinen, Kelly.Merrick, David.Vader, Jordan.Sussman |
| **Date**      | June 9th, 2020                                                                       |
| **Status**    | Accepted                                                                             |

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

There is a desire to have a more elaborate secret store that can handle secret policies, secret rotation, time-sensitive secrets, etc.

Today, we only allow the database secrets to be the default secret store within Vela. Additionally, we give administrators the option to add additional Vault stores.

The [external secrets](05-22_external-secrets.md) proposal solved the problem of integrating with many different types of secret stores when the user is the owner of that store.

Now we still have a problem with the default store not handling all the required uses-cases for modern secrets. I am proposing:

 1. We allow an administrator to select the Vault engine to be the secret store.
 2. We go to only allowing a single secret store specified on the server.

With Vault as a first-class citizen, we can begin pursuing embedding additional Vault specific features natively into Vela. This will enhance our security posture and give us the full power of Vault to help customers solve more advanced secret problems.

Things that will not change:

* unreadable secrets from the server

> **NOTE:**
>
> Aspects of this idea are coming from [GitLab's direction](https://about.gitlab.com/direction/release/secrets_management/) on how to approach storing secrets for enterprise installations.

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

Not required but an enhancement that can help reduce Vela's attack surface and improve on overall security by getting free enhanced security from Vault.

This change could be thought of as baby steps in the direction of Zero Trust Policies:

* [The Zero Trust Security Playbook](https://www.forrester.com/playbook/The+Zero+Trust+Security+Playbook+For+2020/-/E-PLA300#)
* [Google/Alphabet Security white papers](https://cloud.google.com/beyondcorp#researchPapers)

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

* Customers unable to use/want dynamic secrets
* Detailed audit logging around
* Time-based secrets
* Linking up to cloud provider native secrets
* Secret encryption for secrets at rest

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

No workarounds, the introduction of [external secrets](05-22_external-secrets.md) will help mitigate the need for supporting a larger list of secret engines embedded within Vela.

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->

* [Additional secret engines for Vault](https://github.com/go-vela/community/issues/20)

**GitLab:**

* [Manage Vault secrets using GitLab UI](https://gitlab.com/gitlab-org/gitlab/-/issues/20306)
* [Vault integration for key/value secrets MVC](https://gitlab.com/gitlab-org/gitlab-foss/-/issues/61053)
* [Vault integration for CI/CD proof-of-concept](https://gitlab.com/gitlab-org/gitlab/-/issues/9981)

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

No - open to contributors

2. What's the estimated time to completion?

<!-- Answer here -->

2 - 3 weeks

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

Must be defined when user picks up proposal for implementation

## Questions

**Please list any questions you may have:**

<!-- Answer here -->

* What does the difference will need to exist to support additional Vault engines?
* Can we leverage HashiCorps [Vault Integration Program](https://www.vaultproject.io/docs/partnerships) with our enterprise contract?
