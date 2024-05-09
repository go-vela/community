# Nested API Objects

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value |
| :-----------: | :---: |
| **Author(s)** | Jordan.Brockopp, Easton.Crupper |
| **Reviewers** |       |
| **Date**      |       |
| **Status**    |       |

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

This change would be considered a redesign/refactor to modify the behavior of what information is returned by the API.

The idea is to no longer return the ID fields for resources that have a relationship with one another.

For example, today when you query a repo (`GET /api/v1/repo/:org/:repo`), a `user_id` field is returned in the response.

The `user_id` field contains the primary key for a row in the `users` table that represents the "owner" of the repo.

Unfortunately, that `user_id` field isn't providing value for administrators or consumers (end-users).

At this time, the below table contains a list of all resources that have a ID field nested under them:

| Resource           | Fields |
| :-----------: | :---: |
| [library.Build](https://github.com/go-vela/types/blob/master/library/build.go#L16-L51) | `repo_id`, `pipeline_id` |
| [library.Deployment](https://github.com/go-vela/types/blob/master/library/deployment.go#L13-L28) | `repo_id`       |
| [library.Hook](https://github.com/go-vela/types/blob/master/library/hook.go#L11-L29) | `repo_id`, `build_id` |
| [library.Log](https://github.com/go-vela/types/blob/master/library/log.go#L14-L25) | `repo_id`, `build_id`, `service_id`, `step_id` |
| [library.Pipeline](https://github.com/go-vela/types/blob/master/library/pipeline.go#L11-L31) | `repo_id` |
| [library.Repo](https://github.com/go-vela/types/blob/master/library/repo.go#L11-L38) | `user_id` |
| [library.Service](https://github.com/go-vela/types/blob/master/library/service.go#L16-L35) | `repo_id`, `build_id` |
| [library.Step](https://github.com/go-vela/types/blob/master/library/step.go#L16-L36) | `repo_id`, `build_id` |

> NOTE: To provide credit where it's due, this pattern is being adopted from an existing, prominent API (GitHub).
>
> To see an example of how this may look in a real world scenario, you can reference their [docs on getting a repo](https://docs.github.com/en/rest/repos/repos#get-a-repository).
>
> The `owner` field from the GitHub docs displays the behavior being referred to.

<!--
Provide your description here.
-->

**Please briefly answer the following questions:**

**1. Why is this required?**

<!-- Answer here -->

This functionality should not be considered required.

Instead, this proposal aims to enhance/augment the existing API implementation.

**2. If this is a redesign or refactor, what issues exist in the current implementation?**

<!-- Answer here -->

Today, the API returns pieces of information that aren't valuable without directly querying the database.

This functionality would improve the experience for both administrators and end-users.

i.e. an end-user is able to answer the question "Who is the owner of a repo in Vela?"

**3. Are there any other workarounds, and if so, what are the drawbacks?**

<!-- Answer here -->

N/A

**4. Are there any related issues? Please provide them below if any exist.**

<!-- Answer here -->

https://github.com/go-vela/community/issues/69

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

As described, the idea is to no longer return an ID field for related objects from the API.

Using the above endpoint (`GET /api/v1/repo/:org/:repo`), the result would replace the `user_id` field with an `owner` field:

```diff
{
  "id": 1,
- "user_id": 1,
+ "owner": {
+   "id": 1,
+   "name": "OctoKitty",
+   "favorites": ["github/octocat"],
+   "active": true,
+   "admin": false
+ },
  "org": "github",
  "name": "octocat",
  "full_name": "github/octocat",
  "link": "https://github.com/github/octocat",
  "clone": "https://github.com/github/octocat",
  "branch": "master",
  "build_limit": 10,
  "timeout": 60,
  "visibility": "public",
  "private": false,
  "trusted": true,
  "active": true,
  "allow_pr": false,
  "allow_push": true,
  "allow_deploy": false,
  "allow_tag": false
}
```

However, the implementation of this functionality has different options available which are described below.

### Option 1

This option would explore redesigning the existing endpoints under the `/api/v1` collection to return the nested objects.

The below contains a concise list of pros and cons for this option:

#### Pros:

* less code than option 2 (changes existing structs & API handlers)
* less effort than option 2 (less code translates to less effort)
* not having to support another API collection of endpoints

#### Cons:

* breaking change for any user workflow utilizing the ID fields
* existing code dependent on the ID fields needs to be updated

### Option 2

This option would explore adding a new `/api/v2` collection that would be supported in-tandem with the `/api/v1` collection.

> NOTE: Looking to the future, we would eventually remove the `/api/v1` collection at an undetermined date.

The below contains a concise list of pros and cons for this option:

#### Pros:

* backwards-compatible (not a breaking change with new `/api/v2` collection)
* no existing code needs to be updated (new `/api/v2` collection)

#### Cons:

* more code than option 1 (duplicate code i.e. structs (`library.Repo`/`library.V2Repo`), API endpoints/handlers/routes etc.)
* more effort than option 1 (more code translates to more effort)
* requires supporting multiple API collections of endpoints
* likely requires another proposal for determining what the `/api/v2` collection supports

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

N/A since it depends on the option

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

TODO since it depends on the option

## Questions

**Please list any questions you may have:**

<!-- Answer here -->

Which option do you prefer?

A. [Option 1](#option-1)

B. [Option 2](#option-2)
