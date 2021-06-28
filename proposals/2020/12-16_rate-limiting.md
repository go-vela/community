# Rate limiting builds

<!-- Please leave this commented out section.

The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value |
| :-----------: | :---: |
| **Author(s)** |  Matthew Fevold     |
| **Reviewers** |       |
| **Date**      |  12/16/2020     |
| **Status**    |  In Review    |

<!-- Please leave this commented out section.

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

<!-- Please leave this commented out section.

This section is intended to describe the new feature, redesign or refactor.
-->

**Please provide a summary of the new feature, redesign or refactor:**

This feature would attempt to limit the amount of running Vela builds at one time by a given party for system stability purposes.

How is it today:
- Vela customers have no way to achieve this.
- Drone customers can set a `concurrency` value in their `.drone.yml` to limit deployments (_notably little use however_)
- Jenkins customers can set a ton of options here, but are not competing for resources.


**Please briefly answer the following questions:**

1. Why is this required?

- Protect system stability.
  - "Noisy neighbors" - those who use a lot of Vela's resources with many commits/builds preventing other users from running builds.
  - A `Restart build` storm - when a developer notices an issue with their build, often times the easiest troubleshooting step is to restart the build. Now imagine 100+ people doing this, but the issue was low system resources.
- Allows customers to protect their deployments from being triggered twice. (a prod deploy for example)
- Opens up future patterns that could be potentially harmful without these limits (scheduled builds for example)



2. If this is a redesign or refactor, what issues exist in the current implementation?
3. Are there any other workarounds, and if so, what are the drawbacks?
4. Are there any related issues? Please provide them below if any exist.

## Design

**Please describe your solution to the proposal. This includes, but is not limited to:**

* new/updated endpoints or url paths
* new/updated configuration variables (environment, flags, files, etc.)
* performance and user experience tradeoffs
* security concerns or assumptions
* examples or (pseudo) code snippets

A few things to consider with this design (with options indented):
- At what level should we limit?
  - Org, total number for an org
  - Repo, total number for a repo
  - Pipeline, total number of pipelines with same org+repo+branch name (protecting deployments)
- How could it be configurable?
  - Admin set
  - User set in .vela.yml
  - User set in settings UI
- What is the UX when the limit is reached?
  - fail builds
  - queue builds (can be broken down further for queue strategies)


Given all above, we're left to pick and choose which parts of the above we like.


#### Option MVP:
```
- System set max builds at the repo level.
- Allow an override for particular repos (give larger max for larger team is the use case)
- User set repo max in the setting page. (possibly set a default for repos at the org level)
- When max builds is reached:
  - queue the build (FIFO queue)
  - builds triggered by [a configurable set of] events ignore queue

pro:
- protects system stability
- gives users some control (which we'd be sure to get feedback for how they really want it to be)
- expandable later down the line
- allows prod deploys to be prioritized

con:
- doesn't prevent a queue from getting large for a team
- lack of control around how builds are prioritized
- a large team could run into issues with
```

#### Option A:
```
- System set max builds at the repo level.
- Users can configure max builds in:
  - settings page (upto to the system max)
  - .vela.yml at the pipeline level (to enforce a single deployment or to avoid going to the settings page)
  - if configured in both, .vela.yml takes priority
- When max builds is reached:
  - queue the build
  - .vela.yml or settings page user set queue strategy [FIFO, priority queue, auto-cancel policy] default FIFO.

pro:
- protects system stability
- maximum user configuration with a default to protect stability
- default follows current user expectations
- easy to enforce 1 deployment
- can expand on the queue idea later

con:
- lots of different paths to test and implement.
- not very MVP
```


#### Option B:
```
- System set max builds at the repo level.
- Users can set a max builds in their settings repo settings page (up to the system max)
- When max builds is reached, fail.

Pro:
- simple.
- Easy to understand from a user perspective.

Con:
- to enforce 1 deployment, the user has to set to 1 before tagging
- failing builds might make for bad UX
  - likely lots of questions for build failures...
```

#### Option C:
```
same as B, but queue builds on failure.
```

#### Other options

I could go on, but it would just be iterating on the different possibilities listed above.
Do you see an option you'd prefer or tweak some?


## Implementation

<!-- Please leave this commented out section.

This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?
2. What's the estimated time to completion?


**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**


## Questions

**Please list any questions you may have:**

- I omitted time based limitation (e.g. N builds per hour) from my proposal as I didn't see a meaningful thing they provided that a concurrency based limitation didn't provide. But I'd be open to revisit this if there were ideas around it.
- How could we utilize Redis's priority?
- Would we have to setup a new queue?  
