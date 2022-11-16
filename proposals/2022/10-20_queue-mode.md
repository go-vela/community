# Server + Worker Queue Mode

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value                                                                                                |
| :-----------: | :--------------------------------------------------------------------------------------------------: |
| **Author(s)** | David.Vader                                                                                          |
| **Reviewers** | Neal.Coleman, David.May, Jordan.Brockopp, Kelly.Merrick, Easton.Crupper, Jordan.Sussman |
| **Date**      | October 20th, 2022                                                                                      |
| **Status**    | Under Review                                                                                         |

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

The feature proposed is a new build queue/secrets consumption model where the queue listener is moved to the server, build secrets are fetched and packaged by the server, then the payload is pushed from the server to the worker. Along with the build and secrets, the server would also generate and package a BUILD_TOKEN to send to the worker. That token would have access to update the resources related to that build and would expire when the build is complete or the duration has passed. Ideally that BUILD_TOKEN would also ONLY work when used by the worker that it was provisioned for. The server would need the ability to track worker state and verify that it is pushing builds to the correct host, therefore the `workers` table in the database will need to be extended and other efforts such as a worker onboarding process may end up being a requirement.

This feature would add new queue modes for the server and worker. If enabled on the server, the server instance would have a queue listener attached to it that is finding available workers running in `queueless` mode. If enabled on the worker, the worker instance would not pop items from the queue, but rather wait for build allocations from the server. 

The main advantage to this approach would be the replacement of the  `WORKER_SECRET` with the `BUILD_TOKEN` and removing the worker's ability to fetch secrets from the server.


**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

* provides a significant improvement to security (removes the need for a static WORKER_SECRET that can access secrets)
* takes a greater "dumb worker" approach by making the execution less contingent on access to the server

1. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

The worker code does not allow for skipping secret pulling.
The worker code does not allow API based build execution.
The server has no awareness of worker activity and availability (for running builds).
Worker async token management will need to be implemented in the server.


1. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

We could accomplish the same security improvements mentioned above using ephemeral VMs and a short-lived WORKER_SECRET, but ephemeral VMs are not being pursued at this time.

1. Are there any related issues? Please provide them below if any exist.

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

From a high level we need to implement the following:
1. `server`: pops items from the queue
    - duplicate listening code from worker to server
    - enabled by feature flag for server queue mode
1. `server`: process queue item, validate, gather secrets / other resources
    - server code
    - enabled by feature flag for server queue mode
1. `server`: identify and validate available worker from pool
    - issue challenge to worker
        - worker would need to provide registered `worker_token` :warning:
1. `server`: push build + secrets + `build_token` to available worker
    - server generates `build_token` that CANNOT access secrets
        - worker would need to exchange `worker_token` for a `build_token` :warning:
        - that `worker_token` would become **locked** until the build is complete
        - what happens if a worker crashes during the build, how do you **unlock** it?
        - what happens if the `WORKER_BUILD_LIMIT` is > 1?
    - potentially encrypt the secrets portion of the payload
    - enabled by feature flag for server queue mode
1. `worker`: does not pop items from the queue
    - enabled by feature flag for worker queue mode
1. `worker`: receives build + build_token and executes it
    - validate server request using some form of token :warning:
    - implement as worker api endpoint 
    - requires worker availability
    - execute the build with the provided secrets
    - enabled by feature flag for worker queue mode

New functionality that is absolutely required:
#### Worker
- POST `api/v1/execute_build_package`
  - would take a payload containing a packaged build with secrets/token and execute it
- `api/v1/availability`
  - would respond with a server challenge and information regarding availability to run a build
- middleware in the worker API to expose the executor capabilities to the API context
- feature flag to allow swapping between queue runtime modes

#### Server
- `goroutine` queue listener ripped from the worker codebase
- `type Package` that is a build containing secrets, with a build_token
- feature flag to allow swapping between queue runtime modes

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
