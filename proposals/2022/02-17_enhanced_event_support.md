# Enhanced Event Support

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value                                                                                |
| :-----------: | :----------------------------------------------------------------------------------: |
| **Author(s)** | Easton.Crupper                                                                     |
| **Reviewers** |  |
| **Date**      | February 17th, 2022                                                                      |
| **Status**    | In Progress                                                                             |

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

Currently, users are able to configure their repositories to run their pipeline on the following events: `push`, `pull_request`, `deployment`, `tag`, and `comment`. Further, they are allowed to leverage these events within the `ruleset` tag in their pipeline. 

This proposal will largely focus on the `pull_request` event, but will have potentially widespread implementations in the future. 

There are several `actions` that trigger a `pull_request` event (see [GitHub's Docs](https://docs.github.com/en/enterprise-server@3.1/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#pull_request)). Vela currently processes two of these actions: `opened` and `synchronize`. 

There have been many requests from users to accommodate other actions, such as `edited`, or `labeled`. Instead of continuing to add more actions that would launch more builds for users who may not want them, a potential path forward is to allow users to configure which actions they want processed by Vela. 

Implementation of this configurability will be discussed below.

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

* There have been many users who have brought use cases where they would like to process an action that is not `opened` or `synchronize`. 
* Being able to customize actions is not present in Travis CI, Drone, or Jenkins. The only CI tool that manages to do this is [GH Actions](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#using-activity-types). This makes Vela a more desirable tool.

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

The `pull_request` event has two major actions, for which we currently account. Every additional action we handle would be incredibly cumbersome/annoying for any team that does not wish to process builds on these actions. The _only_ way forward is to make them optional.

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

No

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->

Plenty of discussion has already happened related to this change [here](https://github.com/go-vela/community/issues/159). 

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
## YAML Changes

### Design 1 - A more verbose event array

```yaml
steps:
    - name: current behavior of PR event but with this design
      ruleset:
          events: 
              - type: push
              - type: pull_request
                actions: [open, commit]

    - name: example customized event set
      ruleset:
          events:
              - type: push
              - type: pull_request
                actions: [edit, label]
                label: bug
``` 

Pros (as I see them): 
* Very easy to understand
* Follows a similar structure that Vela has established in the past
* Keeps the door open for handling additional actions for other events, such as `push`, `pull_request_review`, etc. 

Cons:
* Would probably need to be a `version 2` feature. Users would not be able to simply put `event: [push, pull_request]` and keep the original behavior. That is a level of breaking that would have to be optional.

### Design 2 - Add an actions tag 
```yaml
steps:
  - name: current behavior of PR event but with this design
    ruleset:
      event: [ pull_request ]
      action: [ open, commit ]

  - name: example customized event set
    ruleset:
      event: [ pull_request ]
      action: [ open, edit, label ]
      label: bug
```

Pros:
* Very easy to understand
* Familiar structure to users
* Would not need a `version 2` requirement since it would be very easy to keep default behavior is pipelines were unchanged.

Cons:
* reserves `action` field for the `pull_request` event. Unless we named it `pr_action` or something along those lines, it would not leave the door open for implementation for other events. 


### Design 3 - Scoping Events
```yaml
steps:
    - name: current behavior of PR event but with this design
      ruleset:
          event: [push, pull_request:open, pull_request:commit]

    - name: example customized event set
      ruleset:
          event: [push, pull_request:edit, pull_request:label]
          label: bug
```

Pros:
* Follows a recognizable pattern seen in GitHub Oauth scoping
* No ambiguity, and no need for an additional actions tag

Cons:
* Would need `version 2` requirement unless we accepted simply `pull_request` and had a default configuration. This adds complications such as: what if a user has both `pull_request` and `pull_request:open`. 
* A tinge messy

## Database / Repo Object Changes

We will need to keep track of what events are allowed for a repository, much like we do now with the `AllowPush`, `AllowPull`, etc fields. With the addition of actions, the list is getting quite expansive. 

### Design 1 - Make a singular events array, which has scoped values

```go
events := _repo.GetEvents()
fmt.Println(events)
  // ["push", "pull_request:open", "pull_request:commit", "tag"]
```

### Design 2 - Keep events as is and add actions array
```go
allowPull := _repo.GetAllowPull()
actions := _repo.GetAllowedActions()

fmt.Println(allowPull) // true
fmt.Println(actions)
// ["open", "edit", "commit"]
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

2 weeks, pending discussion

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

After discussion with the vela admins and vela committers, it was decided that design 3 would be the
choice for handling specific PR actions. 

Further, while a version 2 was discussed, it was determined that we can utilize the scoping structure
so long as we preserve the legacy method of simply putting `pull_request` in the events array.

Labeling and PR review events will be handled separately. In the meantime, `pull_request:opened`, 
`pull_request:edited`, and `pull_request:synchronized` will be implemented.

PRs:

* https://github.com/go-vela/types/pull/236 (adding action constants and handling legacy method)

## Questions