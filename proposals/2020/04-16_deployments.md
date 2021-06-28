# Deployments

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value                                                                |
| :-----------: | :------------------------------------------------------------------: |
| **Author(s)** | Jordan Brockopp                                                      |
| **Reviewers** | Neal Coleman, David May, Emmanuel Meinen, Kelly Merrick, David Vader |
| **Date**      | May 1st, 2020                                                        |
| **Status**    | Complete                                                             |

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

This feature will enable customers to trigger deployments in their pipelines.

This should include supporting the `deployment` event for repos along with the targeted environment for the deployment as conditions for executing a task in the YAML pipeline.

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

* provide compatible functionality with existing CI solutions
* better support the ability to perform continuous delivery

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

N/A

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

N/A

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->

* https://github.com/go-vela/server/issues/89
* https://github.com/go-vela/cli/issues/49
* https://github.com/go-vela/sdk-go/issues/32

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

### Option 1

This option would be to give customers the greatest control over spawning a new build. The customer would be able to provide every single field necessary for creating a build.

At a minimum, we'd need the following information:

* `org`
* `repo`
* `base-ref`
* `branch`
* `commit`
* `event`
* `ref`
* [Variable for Deployment Environment](04-16_deployments.md#questions)

```sh
$ vela add build --help

NAME:
	vela add build - Add a build

USAGE:
	vela add build [command options] [arguments...]

DESCRIPTION:
	Use this command to add a build.

OPTIONS:
	--org value                                    Provide the organization for the repository [$BUILD_ORG]
	--repo value                                   Provide the repository contained within the organization [$BUILD_REPO]
	--base-ref                                     Provide the base commit reference for the new build
	--branch                                       Provide the branch for the new build (default: master)
	--commit                                       Provide the commit for the new build
	--deploy                                       Provide the target environment for the new build
	--event                                        Provide the event for the new build (default: push)
	--ref                                          Provide the commit reference for the new build (default: refs/heads/master)
```

In order to create a `deployment`, you'd run this command:

```sh
$ vela add build \
 --org github \
 --repo octocat \
 --base-ref refs/heads/master \
 --branch master \
 --commit 7fd1a60b01f91b314f59955a4e4d4e80d8edf11d \
 --deploy production \
 --event deployment \
 --ref refs/heads/master
```

When executing this command, it would send a request directly to the Vela API to create this new build:

```sh
POST   /api/v1/repos/:org/:repo/builds
```

> **NOTE:**
>
> This method carries some risk to customers. They could forget to provide some of the input or provide invalid input.

### Option 2

This option would be to adopt a pattern from existing CI solutions.

The customer provides the build they are targeting for deployment and then we create a new deployment event from that existing build.

At a minimum, we'd need the following information:

* `org`
* `repo`
* `build`
* [Variable for Deployment Environment](04-16_deployments.md#questions)

```sh
$ vela deploy --help

NAME:
	vela deploy - Add a build as a deployment

USAGE:
	vela deploy [command options] [arguments...]

DESCRIPTION:
	Use this command to add a build as a deployment.

OPTIONS:
	--org value                                    Provide the organization for the repository [$BUILD_ORG]
	--repo value                                   Provide the repository contained within the organization [$BUILD_REPO]
	--build-number value, --build value, -b value  Provide the build number (default: 0) [$BUILD_NUMBER]
	--deploy                                       Provide the target environment for the new build
```

In order to create a `deployment`, you'd run this command:

```sh
$ vela deploy \
 --org github \
 --repo octocat \
 --build 5 \
 --deploy production
```

When executing this command, it would send a request directly to the Vela API to create this new build:

```sh
POST   /api/v1/repos/:org/:repo/builds
```

The idea is that you provide the specific build you're looking to trigger a deployment from.

This is very powerful because customers could trigger a deployment from any kind of existing build in our system:

* `comment`
* `pull_request`
* `push`
* `tag`

> **NOTE:**
>
> This method reduces the risk from [Option 1](04-16_deployments.md#option-1) of customers forgetting to provide some of the input or providing invalid input.

### Option 3

This option would be to create an actual deployment on the source repository via API call to ensure the repo itself tracks the information.

This also makes deployments a native resource to Vela requiring new API/CLI endpoints.

Several source providers support deployments so there isn't much of a risk for different source providers not supporting them:

* [GitHub](https://developer.github.com/v3/repos/deployments/)
* [GitLab](https://docs.gitlab.com/ee/api/deployments.html)
* [Bitbucket](https://confluence.atlassian.com/bitbucket/bitbucket-deployments-guidelines-941599590.html)

At a minimum, we'd need the following information:

* `org`
* `repo`
* `commit`
* `ref`
* [Variable for Deployment Environment](04-16_deployments.md#questions)

```sh
$ vela add deployment --help

NAME:
	vela add deployment - Add a deployment

USAGE:
	vela add deployment [command options] [arguments...]

DESCRIPTION:
	Use this command to add a deployment.

OPTIONS:
	--org value                                    Provide the organization for the repository [$BUILD_ORG]
	--repo value                                   Provide the repository contained within the organization [$BUILD_REPO]
	--commit                                       Provide the commit for the new build
	--deploy                                       Provide the target environment for the new build
	--ref                                          Provide the commit reference for the new build (default: refs/heads/master)
```

In order to create a `deployment`, you'd run this command:

```sh
$ vela add deployment \
 --org github \
 --repo octocat \
 --commit 7fd1a60b01f91b314f59955a4e4d4e80d8edf11d \
 --deploy production \
 --event deployment \
 --ref refs/heads/master
```

When executing this command, it would send a request directly to the Vela API to create the new deployment:

```sh
POST   /api/v1/deployments/:org/:repo/:build
```

This new endpoint would parse the information and prepare to send a request directly to the source provider.

An example for GitHub:

https://developer.github.com/v3/repos/deployments/#create-a-deployment

```
POST /api/v3/repos/github/octocat/deployments

{
  "ref": "refs/heads/master",
  "task": "deploy",
  "environment": "production",
}
```

Once GitHub received that deployment, it would then trigger our configured webhook for the repository, thus treating a `deployment` for a repository the same as a commit, pull request, tag etc.

The reasoning why we'd create the deployment directly on the repository is so that a "system of record" is created for the transaction, but it also gives us the added benefit of not requiring us to create our own `deployments` table in our database backend.

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

After some discussion amongst the team, we've decided to progress forward with Option 3.

It was also decided to move forward with the `--target` from the question posed below.

A concern that was brought up among those discussions was the increased dependence on the source system and Vela.

This is because we plan to keep the source repo as the official record for deployments rather then our database.

To address this concern, our plan is to go back and later implement a `deployments` table in our database.

## Questions

**Please list any questions you may have:**

<!-- Answer here -->

Which of the following options do you prefer to specify the target environment we are setting for the deployment?

* `--deploy`
* `--deploy-target`
* `--environment`
* `--target`

It's worth mentioning that the choice of flag will likely define our YAML specification under the `ruleset` attribute:

```yaml
steps:
  - name: deploy
    image: alpine
    commands:
      - echo deploy
    ruleset:
      event: deployment
      deploy: dev

  - name: deploy_target
    image: alpine
    commands:
      - echo deploy_target
    ruleset:
      event: deployment
      deploy_target: test

  - name: environment
    image: alpine
    commands:
      - echo environment
    ruleset:
      event: deployment
      environment: stage

  - name: target
    image: alpine
    commands:
      - echo target
    ruleset:
      event: deployment
      target: production
```
