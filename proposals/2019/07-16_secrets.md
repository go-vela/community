# Secrets

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value                                     |
| :-----------: | :---------------------------------------: |
| **Author(s)** | Neal.Coleman, Jordan.Brockopp             |
| **Reviewers** | David.May, Emmanuel.Meinen, Kelly.Merrick |
| **Date**      | July 17th, 2019                           |
| **Status**    | Complete                                  |

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

This feature will enable the ability to provide multiple secret engines (backends) for customers.

[Vault](https://www.vaultproject.io/) is a tool for securely storing and accessing secrets.

Adding another secret backend, like Vault, would provide a greater sophistication to our secret configuration options.

In order to accomplish adding more secret backends, we also must involve refactoring the existing secrets implementation.

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

* provide compatible functionality with existing CI solutions
* provide extra layers of redundancy to secrets management
* enables adding other secret backends (like Kubernetes)
* greater sophistication to secret configurations

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

The current API/CLI setup all assume we're interacting with a single secret backend.

An example of the current structure:

```sh
// POST   /api/v1/repositories/:owner/:repository/secrets
```

Unless we were to use query parameters, there isn't a way for us to address which backend we should be creating the new secret in.

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

The below represents the modifications necessary regarding the 3 secret types:

* `repo`
* `org`
* `shared`

### API

The below endpoints would be created for the new secrets functionality.

The idea is using the `:engine` and `:type` parameters in the path, will enable support for future secret engines while not having to actually add new endpoints for them.

```sh
// POST   /api/v1/secrets/:engine/:type/:owner/:name/
// GET    /api/v1/secrets/:engine/:type/:owner/:name/
// GET    /api/v1/secrets/:engine/:type/:owner/:name/
// PUT    /api/v1/secrets/:engine/:type/:owner/:name/
// DELETE /api/v1/secrets/:engine/:type/:owner/:name/
```

**NOTE: The below endpoints will be deprecated and supported only in 0.4.x versions.**

```sh
// POST   /api/v1/repositories/:owner/:repository/secrets
// GET    /api/v1/repositories/:owner/:repository/secrets
// GET    /api/v1/repositories/:owner/:repository/secrets/:secret
// PUT    /api/v1/repositories/:owner/:repository/secrets/:secret
// DELETE /api/v1/repositories/:owner/:repository/secrets/:secret
```

### CLI

#### Repo

```sh
# Example Vault path for repo secret storage
secret/repo/{org}/{repo}/foo

# Example Vela CLI command for repository secret storage
vela create secret --engine vault --type repo --org github --repo octocat --name foo --value bar

# Example Vault CLI return for repository secret storage
vault read secret/repo/github/octocat/foo
Key                 Value
---                 -----
refresh_interval    10h
repo                github/octocat
value               bar
```

#### Org

```sh
# Example Vault path for org secret storage
secret/org/{org}/foo

# Example Vela CLI command for repository secret storage
vela create secret --engine vault --type org --org github --repo * --name foo --value bar

# Example Vault CLI return for org secret storage
vault read secret/org/github/foo
Key                 Value
---                 -----
refresh_interval    10h
repo                github/*
value               bar
```

#### Shared

```sh
# Example Vault path for shared secret storage
secret/shared/{org}/{team}/foo

# Example Vela CLI command for shared secret storage
vela create secret --engine vault --type shared --org github --team octokitties --name foo --value bar

# Example Vault CLI return for shared secret storage
vault read secret/shared/github/octokitties/foo
Key                 Value
---                 -----
refresh_interval    10h
team                github/octokitties
value               bar
```

### Pipeline

Current `secrets` YAML definition:

```yaml
secrets:
  # native secrets implicit
  - name: docker_username

  # native secrets explicit
  - name: docker_password
    type: native
    key:  docker_username

  # external secrets
  - name: docker_token
    type: vault
    key:  secrets/org/docker/token
```

Proposed `secrets` YAML definition:

```yaml
secrets:
  # native secrets implicit
  - name: docker_username

  # native secrets explicit
  - name: docker_password
    key:  /{org}/{repo}/docker_password
    engine: native
    type: repo

  # vault secrets explicit
  - name: docker_token
    key:  /{org}/docker/token
    engine: vault
    type: org

  # vault secrets explicit
  - name: docker_url
    key:  /{org}/{team}/docker/url
    engine: vault
    type: shared
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

## Questions

**Please list any questions you may have:**

<!-- Answer here -->

N/A
