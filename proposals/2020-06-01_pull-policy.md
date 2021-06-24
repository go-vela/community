# Pull Policy

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in YYYY-MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value                                                                |
| :-----------: | :------------------------------------------------------------------: |
| **Author(s)** | Jordan.Brockopp                                                      |
| **Reviewers** | Neal.Coleman, David.May, Emmanuel.Meinen, Kelly.Merrick, David.Vader |
| **Date**      | June 1st, 2020                                                       |
| **Status**    | Completed                                                            |

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

This enhancement will enable the ability to control when an image is retrieved to be used in a pipeline.

By expanding upon the existing `pull` YAML declaration, we can add more flexibility in how the worker will capture the images to be executed in a pipeline.

It will also give customers the flexibility to use images that are built during the execution of the pipeline in a later, subsequent step.

<!--
Provide your description here.
-->

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

* provide compatible functionality with existing CI solutions
* provide greater control of how and when images will be retrieved for a pipeline
* enable customers to use an image that is built during the execution of the pipeline

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

Currently, users are unable to use an image that is built directly in the pipeline.

In the below design, the pipeline would fail due to the worker attempting to pull all images before executing the pipeline.

Even if the image does exist, and a customer wants to overwrite the existing image with the one built in their pipeline, they would still end up running the pipeline with the old image.

```yaml
version: "1"

steps:
  - name: publish
    image: target/vela-docker:v0.2.1
    pull: true
    parameters:
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
      tags: [ myNewTag ]

  - name: test
    image: index.docker.io/octocat/hello-world:myNewTag
    pull: true
    commands:
      - echo 'this step will fail'
```

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

Expand the existing `pull` YAML declaration to allow other configurations beyond "pull or don't pull".

5 policies will exist:

* `always` (equal to the old `true` policy)
* `not_present` (default)
* `on_start`
* `never`
* `true` (for backwards compatibility)

### Backwards Compatibility

* `pull` will remain **an optional configuration**

All pipelines that exist in Vela today will not need to be modified for this enhancement.

This is because the new policies that can be applied will be compatible with the old policies.

* providing `pull: true` == `pull: always`
* not providing `pull: true` == `pull: not_present`
* `pull: not_present` will be the default so you can omit providing it

```yaml
version: "1"

steps:
  - name: old-pull-policy
    image: alpine:latest
    pull: true
    commands:
      - echo "example of old pull true setup"

  - name: new-pull-policy
    image: alpine:latest
    pull: always
    commands:
      - echo "example of new pull true setup"
```

### Always

The `always` policy will be the functional equivalent of the old, deprecated `true` policy.

Using this policy will instruct the worker to pull the image, even if it already exists locally.

```yaml
version: "1"

steps:
  - name: pull-always
    image: alpine:latest
    # pull policy declaration
    pull: always
    commands:
      - echo "this will always pull the latest tag of the alpine image"
```

### Not Present

The `not_present` policy will be the functional equivalent of omitting the `pull` YAML declaration.

Using this policy will instruct the worker to pull the image only if it doesn't already exist locally.

> **NOTE:** This will be the new default policy.

```yaml
version: "1"

steps:
  - name: pull-not_present
    image: alpine:latest
    # pull policy declaration
    pull: not_present
    commands:
      - echo "this will pull the latest tag of the alpine image if it doesn't exist locally"
```

### Never

The `never` policy will be a new concept introduced with this enhancement.

The intention of this policy is to enforce the ability of only using an image from the existing worker cache.

This can allow Vela admins to guarantee a specific version of an image being executed across all pipelines.

Using this policy will instruct the worker to never pull the image

> **NOTE:** If the image doesn't already exist locally, this will cause a pipeline failure.

```yaml
version: "1"

steps:
  - name: pull-never
    image: alpine:latest
    # pull policy declaration
    pull: never
    commands:
      - echo "this will never pull the latest tag of the alpine image"
```

### On Start

The `on_start` policy will be a new concept introduced with this enhancement.

The intention of this policy is to enable using an image that was built during the execution of the pipeline.

Using this policy will instruct the worker to pull the image before starting the container.

```yaml
version: "1"

steps:
  - name: pull-on_start
    image: alpine:latest
    # pull policy declaration
    pull: on_start
    commands:
      - echo "this will pull the latest tag of the alpine image at container startup"
```

### True

The `true` policy will be the old, deprecated functional equivalent of the `always` policy.

Using this policy will instruct the worker to pull the image, even if it already exists locally.

This policy will only be here for backwards compatibility purpose.

> **NOTE:** This policy should be planned for removal in a future release.

```yaml
version: "1"

steps:
  - name: pull-true
    image: alpine:latest
    # pull policy declaration
    pull: true
    commands:
      - echo "this is the old deprecated form of pull always"
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

* Implement `always`, `not_present`, `never`, `on_start`, `true` constants in [go-vela/types/constants](https://github.com/go-vela/types/tree/master/constants)
* Implement a `PullPolicy` type in [go-vela/types/pipeline](https://github.com/go-vela/types/tree/master/pipeline)
* Implement a `PullPolicy` type in [go-vela/types/yaml](https://github.com/go-vela/types/tree/master/yaml)
* The YAML `pull` declaration will be used to control the configuration for the `PullPolicy`.
* The `not_present` policy will be the new default for all steps.

## Questions

**Please list any questions you may have:**

<!-- Answer here -->

N/A
