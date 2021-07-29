# Log Streaming

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value                                                                                |
| :-----------: | :----------------------------------------------------------------------------------: |
| **Author(s)** | Jordan.Brockopp                                                                      |
| **Reviewers** | Neal.Coleman, David.May, Emmanuel.Meinen, Kelly.Merrick, David.Vader, Matthew.Fevold |
| **Date**      | July 22nd, 2021                                                                      |
| **Status**    | Reviewed                                                                             |

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

This enhancement will enable viewing logs for a service/step in near real-time.

Currently, when watching logs for a service/step they aren't provided in real-time.

Instead, we use a buffer mechanism for controlling how logs are published:

1. service/step starts running on a worker producing logs
2. logs that are produced from the service/step are pushed to a buffer
3. if the buffer exceeds `1000` bytes
   * publish the logs from the buffer via API call to the server
   * flush the buffer so we can push more logs to it
4. circle back to number 2 until the service/step is complete
5. if the service/step is complete publish logs from the buffer and flush it

The end-behavior produced by this method is the logs appear in delayed chunks.

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

* provide compatible functionality with existing CI solutions
* improve user experience when viewing logs for a service/step
* improve ability to troubleshoot pipelines by seeing what parts take longer

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

The current implementation leaves a lot to be desired for user experience.

With the logs appearing in delayed chunks, it's difficult to troubleshoot pipelines.

This behavior can make processes appear to be "stuck" or "hung" when inspecting the logs.

This can make it almost impossible to determine if something is running or how long it takes to run.

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

N/A

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->

https://github.com/go-vela/community/issues/156

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

### Pipeline

To resolve the undesired behavior, we need to create options that can fix it.

However, to craft those options, we need a pipeline that can show the behavior.

We'll use the below pipeline to demonstrate this behavior:

```yaml
version: "1"

steps:
  - name: logs
    image: alpine:latest
    pull: not_present
    commands:
      - sleep 1
      - echo "hello one"
      - sleep 2
      - echo "hello two"
      - sleep 3
      - echo "hello three"
      - sleep 4
      - echo "hello four"
      - sleep 5
      - echo "hello five"
```

### Option 1

This option involves updating the [go-vela/pkg-executor](https://github.com/go-vela/pkg-executor) codebase to upload logs on a regular time interval.

A brief explanation of how the code works:

1. service/step starts running on a worker producing logs
2. create a channel to signal to stop processing logs
3. logs that are produced from the service/step are pushed to a buffer
4. spawn a go routine to start polling the buffer
   * spawn an "infinte" `for` loop that will upload logs
     * sleep for `1s`
     * if the channel is closed, terminate the go routine
     * if the channel is not closed
       * publish the logs from the buffer via API call to the server
       * flush the buffer so we can push more logs to it
5. once the service/step is complete, close the channel to terminate the go routine

The code changes can be found below:

* [go-vela/pkg-executor](https://github.com/go-vela/pkg-executor/compare/feature/log_streaming/opt_one?expand=1)
* [go-vela/worker](https://github.com/go-vela/worker/compare/feature/log_streaming/opt_one?expand=1)

The time interval I chose to use in the above code is `1s`.

However, we could choose any time interval we deem fit for this use-case.

Also, we'd likely make this time interval configurable to provide more flexibility.

### Option 2

This option involves updating the [go-vela/pkg-executor](https://github.com/go-vela/pkg-executor) codebase to stream logs directly to the [go-vela/server](https://github.com/go-vela/server).

To accomplish this, new endpoints were added to the server that can accept streaming connections and upload logs to the database on a regular time interval.

A brief explanation of how the code works:

1. service/step starts running on a worker producing logs
2. worker begins streaming logs via API call to the server
3. server accepts the streaming logs from the worker
4. create a channel to signal to stop processing the streaming logs
5. streamed logs that are are pushed to a buffer
6. spawn a go routine to start polling the buffer
   * spawn an "infinte" `for` loop that will upload logs
     * sleep for `1s`
     * if the channel is closed, terminate the go routine
     * if the channel is not closed
       * publish the streamed logs from the buffer to the database
       * flush the buffer so we can push more logs to it
7. once the streaming is complete, close the channel to terminate the go routine

The code changes can be found below:

* [go-vela/pkg-executor](https://github.com/go-vela/pkg-executor/compare/feature/log_streaming/opt_two?expand=1)
* [go-vela/worker](https://github.com/go-vela/worker/compare/feature/log_streaming/opt_two?expand=1)
* [go-vela/server](https://github.com/go-vela/server/compare/feature/log_streaming/opt_two?expand=1)

The time interval I chose to use in the above code is `1s`.

However, we could choose any time interval we deem fit for this use-case.

Also, we'd likely make this time interval configurable to provide more flexibility.

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

1 month

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

## Questions

**Please list any questions you may have:**

<!-- Answer here -->

N/A
