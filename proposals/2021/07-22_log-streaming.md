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
   * flush the buffer so we can push more logs to it from the service/step
4. circle back to number 2 until the service/step is complete
5. once the service/step is complete, publish remaining logs from the buffer

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

We could explore decreasing the limit we impose on the log buffer we use.

Currently, we're using `1000` bytes as that limit but we could set that to something smaller (`100`, `500`, etc.).

This would mean that the worker uploads logs to the server more frequently which could improve the experience.

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

> DISCLAIMER:
>
> The options below **only** cover updating code for the **backend** components for log streaming.
>
> This means the proposal **does not** cover any updates to the [go-vela/ui](https://github.com/go-vela/ui) codebase.
>
> The reason for this is no UI changes are required to produce a streaming effect for logs.
>
> This is due to the UI already polling ~ every `5s` in its current state.
>
> However, this also means we leave room for improvement in regards to the user experience.
>
> This experience could be augmented by making changes to the UI (and server?) along with the below options.

### Option 1

This option involves updating the [go-vela/pkg-executor](https://github.com/go-vela/pkg-executor) codebase to upload logs on a regular time interval to simulate a streaming effect.

**A brief explanation of how the code works:**

1. service/step starts running on a worker producing logs
2. worker creates a channel to signal when to stop processing logs
3. logs that are produced from the service/step are pushed to a buffer
4. worker spawns a `go` routine to start polling the buffer
   * (inside the `go` routine) spawn an "infinite" `for` loop
     * (inside the `for` loop) sleep for `1s`
     * if the channel is closed, terminate the `go` routine
     * if the channel is not closed
       * publish the logs from the buffer via API call to the server
       * flush the buffer so we can push more logs to it from the service/step
5. once the service/step is complete, worker closes the channel to terminate the `go` routine

The code changes can be found below:

* [go-vela/pkg-executor](https://github.com/go-vela/pkg-executor/compare/feature/log_streaming/opt_one?expand=1)
* [go-vela/worker](https://github.com/go-vela/worker/compare/feature/log_streaming/opt_one?expand=1)

> NOTE:
>
> The time interval I chose to use in the above code is `1s`.
>
> However, we could choose any time interval we deem fit for this use-case.
>
> Also, we'd likely make this time interval configurable to provide more flexibility.

### Option 2

This option involves updating the [go-vela/pkg-executor](https://github.com/go-vela/pkg-executor) codebase to stream logs to the [go-vela/server](https://github.com/go-vela/server) via HTTP.

To accomplish this, new endpoints were added to the server that can accept streaming connections.

Once a streaming connection is open, the server will capture and upload logs to the database on a regular time interval.

**A brief explanation of how the code works:**

1. service/step starts running on a worker producing logs
2. worker begins streaming logs via HTTP call to the server
3. server accepts the streaming logs from the worker
4. server creates a channel to signal when to stop processing streamed logs
5. streamed logs are pushed to a buffer by server
6. server spawns a `go` routine to start polling the buffer
   * (inside the `go` routine) spawn an "infinite" `for` loop
     * (inside the `for` loop) sleep for `1s`
     * if the channel is closed, terminate the `go` routine
     * if the channel is not closed
       * publish the streamed logs from the buffer to the database
       * flush the buffer so we can push more logs to it from the service/step
7. once the service/step is complete, worker terminates the HTTP call
8. once the streaming is complete, server closes the channel to terminate the `go` routine

The code changes can be found below:

* [go-vela/pkg-executor](https://github.com/go-vela/pkg-executor/compare/feature/log_streaming/opt_two?expand=1)
* [go-vela/worker](https://github.com/go-vela/worker/compare/feature/log_streaming/opt_two?expand=1)
* [go-vela/server](https://github.com/go-vela/server/compare/feature/log_streaming/opt_two?expand=1)

> NOTE:
>
> The time interval I chose to use in the above code is `1s`.
>
> However, we could choose any time interval we deem fit for this use-case.
>
> Also, we'd likely make this time interval configurable to provide more flexibility.

### Option 3

This option involves updating the [go-vela/pkg-executor](https://github.com/go-vela/pkg-executor) codebase to stream logs to the [go-vela/server](https://github.com/go-vela/server) via WebSocket.

To accomplish this, new endpoints were added to the server that can accept websocket connections.

Once a websocket connection is open, the server will capture and upload logs to the database on a regular time interval.

**A brief explanation of how the code works:**

1. service/step starts running on a worker producing logs
2. worker opens a websocket connection to the server
3. server accepts the websocket connection for streaming logs from the worker
4. worker begins streaming logs via websocket to the server
5. server creates a channel to signal to stop processing the streaming logs
6. streamed logs from the websocket connection are pushed to a buffer by server
7. server spawns a `go` routine to start polling the buffer
   * (inside the `go` routine) spawn an "infinite" `for` loop
     * (inside the `for` loop) sleep for `1s`
     * if the channel is closed, terminate the `go` routine
     * if the channel is not closed
       * publish the streamed logs from the buffer to the database
       * flush the buffer so we can push more logs to it from the service/step
8. once the service/step is complete, worker closes the websocket connection
9. once the streaming is complete, server closes the channel to terminate the `go` routine

The code changes can be found below:

* [go-vela/pkg-executor](https://github.com/go-vela/pkg-executor/compare/feature/log_streaming/opt_three?expand=1)
* [go-vela/worker](https://github.com/go-vela/worker/compare/feature/log_streaming/opt_three?expand=1)
* [go-vela/server](https://github.com/go-vela/server/compare/feature/log_streaming/opt_three?expand=1)

> NOTE:
>
> The time interval I chose to use in the above code is `1s`.
>
> However, we could choose any time interval we deem fit for this use-case.
>
> Also, we'd likely make this time interval configurable to provide more flexibility.

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

After some discussion amongst the team, we've decided to progress forward with Option 2.

This decision was driven via a vote that was done and results provided [here](https://github.com/go-vela/community/pull/339#issuecomment-903837207).

A concern that was brought up among those discussions was how much resources (CPU/RAM) were required for each option.

As we look to actually implement the functionality for Option 2, we should evaluate what changed in resource consumption (if any).

This will likely involve looking at how much CPU/RAM is consumed by both the [server](https://github.com/go-vela/server) and [worker](https://github.com/go-vela/worker) when streaming logs with Option 2.

* https://github.com/go-vela/community/issues/366
* https://github.com/go-vela/community/issues/367
* https://github.com/go-vela/community/issues/368
* https://github.com/go-vela/community/issues/369
* https://github.com/go-vela/community/issues/370

## Questions

**Please list any questions you may have:**

<!-- Answer here -->

N/A
