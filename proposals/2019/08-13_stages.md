# Stages (Concurrency)

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value                                                   |
| :-----------: | :-----------------------------------------------------: |
| **Author(s)** | Jordan.Brockopp                                         |
| **Reviewers** | Neal.Coleman, David.May, Emmanuel.Meinen, Kelly.Merrick |
| **Date**      | August 13th, 2019                                       |
| **Status**    | Complete                                                |

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

This feature will enable the ability to execute multiple tasks, concurrently, for users.

By having some type of YAML attribute that can be provided in a pipeline, this will enable users to utilize functionality that already exists in other CI solutions.

It will also give users the flexibility to improve their build performance and throughput by running any and all tasks concurrently that can support it.

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

* provide compatible functionality with other CI solutions
* enable users to optimize performance and throughput of builds
* provide a fan-in, fan-out mechanism in Vela

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

N/A

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

N/A

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->

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

* `stages` will **always run in parallel**
* Each `stage` is required to have a unique name
* A single `stage` can contain one or many `steps`
* `steps` will **always execute sequentially**
* You can use the `needs` label to specify one stage's dependency on another

### Backwards Compatibility

* `stages` declaration **is optional**
* `needs` declaration **is optional**

All pipelines that exist in Vela today will not need to be modified for this new feature.

The plan is to keep the initial steps only syntax to allow simple use cases that don't require concurrent execution of processes.

To be more forward, using the `stages` and `needs` declarations **should only be necessary** when trying to **execute processes in parallel.**

### Why

It is understood that this might appear as a barrier to entry for some, but the need to run processes in parallel is viewed as an advanced use-case.

So the barrier is necessary to mitigate the misuse or misunderstanding of the purpose of `stages` and `needs`.

There is also a hope that by explicitly defining the `stages`, this will help both the consumers and the producers of pipelines gain an elevated understanding for how the pipeline will be executed.

This should lead to an increased ability for all parties involved to discover, create and troubleshoot advanced use-cases of pipelines.

### Sample

```yaml
version: "1"

stages:

  backend:
    steps:
      - name: install
        image: golang:latest
        commands:
          - go get ./...
      - name: test
        image: golang:latest
        commands:
          - go test ./...
      - name: build
        image: golang:latest
        commands:
          - go build

  frontend:
    steps:
      - name: install
        image: npm:latest
        commands:
          - npm install
      - name: test
        image: npm:latest
        commands:
          - npm run test
      - name: build
        image: npm:latest
        commands:
          - npm run build

  publish_backend:
    needs: [ backend ]
    steps:
      - name: publish
        image: target/vela-docker:latest
        parameters:
          dockerfile: Dockerfile.backend
          repo: octocat/hello-world
          tags: [ backend ]

  publish_frontend:
    needs: [ frontend ]
    steps:
      - name: publish
        image: target/vela-docker:latest
        parameters:
          dockerfile: Dockerfile.frontend
          repo: octocat/hello-world
          tags: [ frontend ]

  notify:
    needs: [ publish_backend, publish_frontend ]
    steps:
      - name: notify
        image: target/vela-slack:latest
        parameters:
          webhook: https://api.slack.com
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
