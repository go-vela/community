# External Secrets

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in YYYY-MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value                                                                                |
| :-----------: | :----------------------------------------------------------------------------------: |
| **Author(s)** | Neal.Coleman, Jordan.Brockopp                                                        |
| **Reviewers** | Neal.Coleman, David.May, Emmanuel.Meinen, Kelly.Merrick, David.Vader, Jordan.Sussman |
| **Date**      | May 22nd, 2020                                                                       |
| **Status**    | Complete                                                                             |

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

This redesign is to accommodate the integration of non-native secret providers in a less intrusive way to Vela.

The goal of this proposal is to simplify the current implementation of secrets and allow the users to have more control over secrets without server admins gating the addition of new secret providers.

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

We need to provide better support for other secret providers.

e.g. Vault, Password Vault, cloud provider secret stores (Secret Manager)

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

Users are not able to bring a secret store, or add new ones without the Vela admins involvement.

We have reviewed feedback of the current implementation and heard it's cumbersome for customers adding secrets in their pipelines, maintaining them via CLI, etc.

It also is not friendly to use cases around teams having a shared secret store for engineers on board and does not natively integrate with their existing secret stores.

Other struggles are complex API, UI, CLI designs that don't easily integrate or future proof integrations with non-native secret engines.

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

Things this will add:

* Bring you own secret store
* Pattern for external secret providers
* New Vault secret plugin

A big difference in this design compared to today's is the system administrator will no longer be required to enable new secret engines.

We will establish a pattern for reading secrets with plugins. As long as the plugin follows the secret interface guidelines a plugin can read external secrets into Vela.

This allows Vela to no longer be a source of truth for secrets but uses an encapsulation pattern to store secrets to external secret stores.

The pattern is used by CI systems like GitHub Actions and Jenkins to allow reading from external secret stores.

> **NOTE:**
>
> * Since we are no longer in control of all secret providers, failures pulling secrets are dependent on the upstream secret service
> * Vela team will no longer provide a catch all Vault or other secret stores. Teams will be on their own for onboarding to secret stores and must provide credentials to access their secrets.

### Example

This example will focus on "Vault" but note the secret provider can be any secret store

```diff
steps:
  - name: build
    image: target/vela-docker:latest
    parameters:
      registry: index.docker.io
      repo: index.docker.io/target/test-docker
    secrets: [ docker_username, docker_password ]

secrets:
  # Implicit secret definition. This definition is only supported for native secrets of repository type.
  - name: vault_username
  - name: vault_password

  # Declarative secret definition.
  - name: foo1
    key: <org>/<repo>/<secret>
    engine: native
    type: repo
  - name: foo2
    key: <org>/<secret>
    engine: native
    type: org
  - name: foo3
    key: <org>/<team>/<secret>
    engine: native
    type: shared

# New plugins will enable reading secrets from a store and sourcing into secrets map
+  - name: docker_username
# Alternatives to "look_up" could be: "from", "with", "origin"
+   look_up:
+     image: target/vela-vault:latest
+     secrets: [ vault_username, vault_password ]
# Alternatives to "parameters" could be: "spec", "detail"
+     parameters:
+      addr: http://vault.company.com
+      auth_method: ldap
+      path: path/to/nuid_username
```

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

2 weeks

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

* Add new secret look up design into types
* Update complier to parse new secret syntax
* Add new environment variable on server to enable secret plugins whitelist
* Update executor to run new secret plugins
* Update dependencies in worker
* Write Vault plugin secret provider compatible with enterprise Vault

After some discussion amongst the team, we've decided to swap the `look_up` key for `origin`.

It was also decided to add a variable to the server to enable administrators to control which plugins can be executed

## Questions

**Please list any questions you may have:**

<!-- Answer here -->

N/A
