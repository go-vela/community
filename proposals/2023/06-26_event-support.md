# Optimize Allowed Events for Repos

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
| **Date**      | June 26th, 2023                                                                    |
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

Some of the groundwork has already been laid with the implementation of [enhanced rulesets](https://github.com/go-vela/community/blob/main/proposals/2022/02-17_enhanced_event_support.md) in 2022. However, the resulting design left a lot to be desired in terms of scalability as well as the user experience. For one, if we wanted to truly account for new events / event actions, we would have to add a DB field to the `repos` table in some form like `Allow<Event><Action>`. Not only that, but we would have to fission a couple of our existing allow fields (`AllowPull` and `AllowComment`). 

This proposal is for a scalable storage system for repository allowed events, which frees up development for handling various events and their sub-actions.

Presently, Vela has the following fields in the `repos` table:

* AllowPush
* AllowPull (covers both `opened` and `synchronize` actions)
* AllowTag
* AllowDeployment
* AllowComment (covers both `created` and `edited` actions)

I believe all this data can be represented as a single field in the database, and, in conjucture with the [nested API objects design](https://github.com/go-vela/community/pull/639) can be represented as nested JSON in the library.

In the `Design` section, I will outline the two proposed solutions for this.


**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

* All of the good intentions of the original [enhanced rulesets proposal](https://github.com/go-vela/community/blob/main/proposals/2022/02-17_enhanced_event_support.md) simply cannot be actualized without proper handling of allowed events repo-side.
* Future developers will not have to weigh the importance of adding support for an event/action over the addition of another column in the database.

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

Too little data taking up valuable space and limiting freedom of development.

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

No

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->

https://github.com/go-vela/community/issues/159 is our most upvoted issue in the backlog

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

## Library Structure

Before presenting the choice of database representation, here is what I envision the library representation of these allowed events will be:

```go
// PRActions is the library representation of the various actions associated
// with the pull_request event webhook from the SCM.
type PRActions struct {
	Opened        *bool `json:"opened"`
	Edited        *bool `json:"edited"`
	Synchronize   *bool `json:"synchronize"`
	Labeled       *bool `json:"labeled"`
	ReviewRequest *bool `json:"review_request"`
}

// CommentActions is the library representation of the various actions associated
// with the comment event webhook from the SCM.
type CommentActions struct {
	Created *bool `json:"created"`
	Edited  *bool `json:"edited"`
}

// ReviewActions is the library representation of the various actions associated
// with the pull_request_review event webhook from the SCM.
type ReviewActions struct {
	Submitted *bool `json:"submitted"`
	Edited    *bool `json:"edited"`
}

// Events is the library representation of the various events that generate a
// webhook from the SCM.
type Events struct {
	Push        *bool           `json:"push"`
	PullRequest *PRActions      `json:"pull_request"`
	Tag         *bool           `json:"tag"`
	Deployment  *bool           `json:"deployment"`
	Comment     *CommentActions `json:"comment"`
	Schedule    *bool           `json:"schedule"`
	PullReview  *ReviewActions  `json:"pull_review"`
}

type Repo struct {
	ID           *int64    `json:"id,omitempty"`
	// ...
	Active       *bool     `json:"active,omitempty"`
	AllowEvents  *Events   `json:"allow_events,omitempty"`
	PipelineType *string   `json:"pipeline_type,omitempty"`
    // ...
}
```

Which gives the JSON response of:
```json
	"allow_events": {
		"push": true,
		"pull_request": {
			"opened": true,
			"edited": false,
			"synchronize": false,
			"labeled": true,
			"review_request": false
		},
		"tag": true,
		"deployment": false,
		"comment": {
			"created": false,
			"edited": true
		},
		"schedule": false,
		"pull_review": {
			"submitted": false,
			"edited": false
		}
	},
```

I think this is a clean representation of the data that is easy for the UI and CLI to interpret. However, if a list (like we do with Secret.Events) is preferred, that can work too.

## DB Design 1 - Bit masking iota integer representation

```go
// Allowed repo events.
const (
	AllowPush   = 1 << iota // 00000001 = 1
	AllowPROpen             // 00000010 = 2
	AllowPREdit             // 00000100 = 4
	AllowPRSync             // ...
	AllowPRLabel
	AllowPRReviewRequest
	AllowTag
	AllowDeploy
	AllowCommentCreate
	AllowCommentEdit
	AllowReviewSubmit
	AllowReviewEdit
	AllowSchedule
)

type Repo struct {
	ID           sql.NullInt64  `sql:"id"`
	// ...
	Active       sql.NullBool   `sql:"active"`
	AllowEvents  sql.NullInt64  `sql:"allow_events"`
	PipelineType sql.NullString `sql:"pipeline_type"`
	// ...
}
``` 

Pros (as I see them): 
* Scalable (adding more events is simple adding another bit)
* Very fast calculations to convert into library representation (nested structs)

Cons:
* Complicated
* Not readable from a DB Select perspective. `SELECT allow_events FROM repos LIMIT 1` would return an integer.

Note: we would probably implement `Secret.Events` in a simlar fashion if this was the chosen design.

## DB Design 2 - String Array (like Secret.Events today)
```go
type Repo struct {
	ID           sql.NullInt64  `sql:"id"`
	// ...
	Active       sql.NullBool   `sql:"active"`
	AllowEvents  pq.StringArray `sql:"events" gorm:"type:varchar(1000)"`
	PipelineType sql.NullString `sql:"pipeline_type"`
	// ...
}    
```

Pros:
* Very easy to understand
* Familiar structure to developers
* Human readable from a DB Select perspective

Cons:
* Have to perform slice logic to generate library representation (e.g. "is pull_request:opened in the slice?")
* Potential column size adjustments if we add more events 

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

2 weeks, once design is decided

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->


## Questions
