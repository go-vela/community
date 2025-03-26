# Dynamic ruleset

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

|      Key      |        Value        |
|:-------------:|:-------------------:|
| **Author(s)** |   Jordan Sussman    |
| **Reviewers** |                     |
|   **Date**    | December 27th, 2024 |
|  **Status**   |     In Progress     |

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

Currently, Vela currently offers several methods to filter which steps are executed:

1. [`ruleset` key](https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key).
2. Conditional logic with [templating directly in the Vela file](https://go-vela.github.io/docs/templates/#templating-directly-in-velayml).
3. Conditional logic with [external templates](https://go-vela.github.io/docs/templates/).

This proposal aims to introduce more dynamic conditionals, similar to those available with templates, but extend this functionality to YAML pipelines that do not use any templates.

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

There are a wide variety of use-cases that exist for this today so we'll name a few:

* filter steps based on conditions that the `ruleset doesn't support`.
* Remove the necessity of potentially complex templates merely to circumvent the limitations of `ruleset`.

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

N/A

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

As mentioned previously, you can use [templates](https://go-vela.github.io/docs/templates/).

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->

- https://github.com/go-vela/community/issues/826
- https://github.com/go-vela/community/issues/524
- https://github.com/go-vela/community/issues/306

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

### Database

N/A

### UI

N/A

### API

N/A

### CLI

N/A

### YAML

A `eval` key will be added under `ruleset` for the new dynamic functionality.

This will instruct Vela to execute the step if the `eval` expression returns `true`.

Example pipeline:

```yaml
steps:
  - name: simple-match
    ruleset:
      eval: "VELA_BUILD_AUTHOR == 'JordanSussman'"
    commands:
      - echo "this will only run if the build author name is JordanSussman"

  - name: function-match
    ruleset:
      eval: "hasPrefix(VELA_BUILD_AUTHOR, 'Jordan')"
    commands:
      - echo "this will only run if the build author name has the prefix of Jordan"
      
  - name: new-and-old
    ruleset:
      eval: "VELA_BUILD_AUTHOR == 'JordanSussman'"
      event: [tag]
    commands:
      - echo "this will only run if JordanSussman triggered the tag event"
```

## Implementation

<!--
This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

Users will be able to use [expr](https://expr-lang.org/) syntax within the `eval` key to filter steps.

Expr describes itself as:

> Expr is a Go-centric expression language designed to deliver dynamic configurations with unparalleled accuracy, safety, and speed. Expr combines simple syntax with powerful features for ease of use.

All of the [built-in Vela environment variables](https://go-vela.github.io/docs/reference/environment/variables/) will be available within the `eval` key to perform logic with.

For those interested in understanding how `expr` works, the project [README](https://github.com/expr-lang/expr/tree/fb6792b2486778dd8a3eb5ab2e7550f5b1dad150?tab=readme-ov-file#examples) has a nice example of managing the underlying code for injecting environment variables. Additionally, I've created a [Go playground](https://go.dev/play/p/6dEEUPTzK8r) to demonstrate a more Vela-specific implementation.

The server code changes will involve updating the existing ruleset match function to consider both the `ruleset` and `eval` keys when determining if the step should execute. Both `ruleset` and `eval` must return true for the step to proceed.

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

<!-- Answer here -->

Yes, Cargill will add the code required to deliver on this feature

2. What's the estimated time to completion?

<!-- Answer here -->

The code changes required to incorporate this should be relatively minor, so the work involved shouldn't take too long. However, the timeline for when we will start is still to be determined.

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

## Questions

1. Should we name the new key `eval` or choose a different name?
2. Should the new key be nested under `ruleset` or be at the same level?
3. Should it be possible to use both the existing `ruleset` and the new `eval` feature simultaneously? If so, should both need to return true for the step to execute?