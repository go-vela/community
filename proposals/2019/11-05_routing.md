# Routing

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value                                                                   |
| :-----------: | :---------------------------------------------------------------------: |
| **Author(s)** | Neal.Coleman                                                            |
| **Reviewers** | Jordan.Brockopp, David.May, Emmanuel.Meinen, Kelly.Merrick, David.Vader |
| **Date**      | November 5th, 2019                                                      |
| **Status**    | Complete                                                                |

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

This feature will enable users to route a build to a specific worker.

This functionality does not exist in other CI systems but we need it to enable more complex workflows.

This functionality will act as a seed to support the following use cases:

* Supporting different worker OS (i.e. Mac, Windows, Linux) types.
* Running builds on specific runtimes (i.e. Docker, Kubernetes, TAP).
* Running builds on workers with extra security measures for compliance (i.e. PCI, SOX).

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

To support users running builds on specific runtimes or OS.

We need to route them to select group of workers

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

N/A

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

The Server will connect to a set of queues. Those queues will be labeled with one to many types( i.e. PCI, linux).

If a build needs a specific worker, the user can specify one or many keys in their YAML that allows for that build to get routed to a specific worker.

Queues will always be dynamic.

The admin of the Vela server can input whatever number of required queues via an environment flag (`VELA_QUEUE_KEYS=PCI,linux`) on the server.

This will allow the admin to specify which keys their deployed instance supports.

Workers will support the same environment flag and only execute builds that are pushed into its queue.

### Keys

Possible YAML keys:

* `label`
* `tag`
* `run`
* `transport`
* `route`
* `use`

### Option One

```yaml
version: "1"

metadata:
  label: [ pci, linux ]
  
steps:
  - name: echo
    image: alpine:latest
    commands:
      - echo "I'm batman!"
```

### Option Two

```yaml
version: "1"

worker:
  label: [ pci, linux ]

steps:
  - name: echo
    image: alpine:latest
    commands:
      - echo "I'm batman!"
```

### Option Three

```yaml
version: "1"

label: [ pci, linux ]
  
steps:
  - name: echo
    image: alpine:latest
    commands:
      - echo "I'm batman!"
```

### Option Four

```yaml
version: "1"

ruleset:
  label: [ pci, linux ]
  
steps:
  - name: echo
    image: alpine:latest
    commands:
      - echo "I'm batman!"
```

### Chosen Design

```yaml
worker:
  name: 
  flavor:
  runtime:
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

Chosen Design:

```yaml
worker:
  name: 
  flavor:
  runtime:
```

## Questions

**Please list any questions you may have:**

<!-- Answer here -->

N/A
