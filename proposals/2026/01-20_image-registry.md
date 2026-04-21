# Image Registry

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value              |
| :-----------: | :----------------: |
| **Author(s)** | Jordan Brockopp    |
| **Reviewers** | TBD                |
| **Date**      | January 20th, 2026 |
| **Status**    | Reviewed           |

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

The idea is one unique feature that would cover two distinct use-cases:

1. As an administrator (platform admin) of a Vela instance, I want to restrict (block) certain container images from being used in a pipeline that we know are not valid. As an example, imagine one (or more) images that were designed for a specific use-case (i.e. service) that no longer exists (i.e. decommissioned). Another example would be images that were managed by a team that no longer exists (i.e. disbanded). In either scenario, the images are likely old and out-of-date (potentially years) so they're riddled with vulnerabilities and should be considered a security risk to the company. I fully acknowledge that this could be jarring to end-users which is where the second use-case comes into play.
2. As an administrator of a Vela instance, I want to promote awareness for certain container images being used in a pipeline that we know are subject to cause issues in the future. Per the examples from the first use-case, if I know a service is being decommissioned sometime in the future, I'd like to highlight that to end-users leveraging specific images (i.e. plugins) designed to interact with that service. Another example would be highlighting that a specific image or version of an image is deprecated (i.e. bug, vulnerability etc.) so I should upgrade to the latest version.

> NOTE: This is considered one feature because it supports configuration at different levels (i.e. block vs warn).

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

* Please see the description and background for the use-cases.

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

* N/A

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

* Not directly through Vela

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->

* https://github.com/go-vela/community/issues/1066

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

The idea is have one or more images configured on the "warn" list for a period of time. Any images used in a pipeline from this list would surface a "warning" on the `Pipeline` tab for a build similar to the behavior added for the YAML migration:

https://github.com/go-vela/community/blob/main/proposals/2024/12-04_yaml-migration.md

Ideally, we'd also highlight this on the `Builds` page to indicate that a respective build in the list had warnings. Depending on the interest, we could also try to surface these warnings on the `Audit` tab for a respective build.

Then, when the Vela administrators feel sufficient time has been given, they could "promote" the container image from the "warn" list to the "block" list. The significance of this change would be that the compiler for Vela would reject (throw an error) if an image on the "block" list was used thus preventing a build from even being created. This would use the same `error` field in the `hooks` table that we use for other compilation failures:

https://go-vela.github.io/docs/usage/troubleshooting#why-is-my-build-not-displaying

### Option 1

This option involves updating the [go-vela/server](https://github.com/go-vela/server) codebase to support ingesting an optional YAML file at startup. If the YAML file is valid, and matches the expected syntax from Vela, the server ingests it and then uses it to make determinations based off the images used in the pipeline depending on the configured "block" and/or "warn" lists.

Here's a simple, rough draft to demonstrate a potential example of the syntax for the YAML file:

```yaml
---
block:
  - image: docker.example.com/disbanded-team/*
    reason: This image has been blocked by Vela administrators as the team that owns it has been disbanded:\n<link to docs>
  - image: docker.example.com/deprecated-image:*
    reason: This image has been blocked by Vela administrators as the image is no longer supported:\n<link to docs>
  - image: docker.example.com/valid-team/deprecated-version:<version>
    reason: The version of this image has been blocked by Vela administrators as it is no longer supported:\n<link to docs>
  - image: docker.example.com/valid-team/vulnerable-image:*
    reason: This image has been blocked by Vela administrators due to the vulnerabilities associated with it:\n<link to docs>

warn:
  - image: docker.example.com/disbanded-team/*
    reason: The team that owns these images has been disbanded so support for these will be blocked in the future:\n<link to docs>
  - image: docker.example.com/deprecated-image:*
    reason: This image has been deprecated so support for these will be blocked in the future:\n<link to docs>
  - image: docker.example.com/valid-team/deprecated-version:<version>
    reason: The version of this image has been deprecated so support for this will be blocked in the future:\n<link to docs>
  - image: docker.example.com/valid-team/vulnerable-image:*
    reason: This image has vulnerabilities that put the company at risk so support for this will be blocked in the future:\n<link to docs>
```

### Option 2

This option involves updating both the [go-vela/server](https://github.com/go-vela/server) and [go-vela/ui](https://github.com/go-vela/ui) codebases to support creating, updating and deleting individual values from the registry via the "site admin" page:

https://vela.example.com/admin/settings

The idea would be to have the server store configuration for both the "block" and "warn" lists directly in the database. This means the API (server) would have to offer new API endpoints to support CRUD operations for these lists. Also, the UI would have to be updated to support interacting with these new API endpoints. There is the potential for the UI page to support configuring these lists as a YAML file for configuration but it wouldn't be a requirement as we could adopt a similar approach to secrets where users can restrict the usage of secrets to specific container images.

## Implementation

<!--
This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

<!-- Answer here -->

* TBD based off the option we choose but yes, for at least the backend (API) portion

2. What's the estimated time to completion?

<!-- Answer here -->

* TBD based off the option we choose

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

* N/A

## Questions

**Please list any questions you may have:**

<!-- Answer here -->

* N/A
