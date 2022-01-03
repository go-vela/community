# Security Measures for Pipeline Secret Protection

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
| **Date**      | December 27th, 2021                                                                      |
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

  This new feature would allow users of Vela to implement a protective system to their pipeline. By doing one of the two designs proposed below, any edits to the `.vela.yml` would be verified before executing. 
  
  The verification process would be executable only by admins of the repository. By doing this, any pipeline that has been edited would fail to run unless an admin approves.

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

* This addresses a very key security issue: pipelines that execute on pull requests can be manipulated to expose secrets -- even on forked repositories. By adding verification, malicious edits to `.vela.yml` would never get the chance to run.
* The added level of security to our pipelines will be more attractive to users with highly sensitive information. As Vela expands, the likelihood that it will be utilized in an organization that practices the repo fork workflow will greatly increase.  

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

As mentioned, the primary issue with the current implementation is the possibility for a code contributor to expose secrets through a Vela pipeline.

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

One workaround, which is currently being worked on, is masking any secret exposed in the logs of a pipeline. While this may alleviate some consequences of careless secret usage in pipelines, it does not eliminate the possibility of a more persistent attacker from exposing secrets in other ways, such as through interacting with a foreign server.

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->

https://github.com/go-vela/community/issues/62

## Design Option 1 - Pipeline Signature

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

### Ideal User Workflow 

1. Repository admin determines it would be beneficial to enable a pipeline signature for added security.
2. Repository admin executes some command that generates a signature (some long string hash) that is based on the contents of the `.vela.yml` file. 
3. Admin commits and merges that signature to the repository.

**Case 1**
- Another repository admin makes changes to the `.vela.yml` pipeline. Perhaps they are adding Slack plugin functionality.
- Upon submitting a pull request, they realize that their build did not execute in Vela. Instead, they see in the audit that their pipeline did not match the signature.
- The repository admin then executes a command to re-sign the repository. They are able to do this because they are an admin. 
- The pipeline executes, and the PR is able to be reviewed.

**Case 2**
- A code contributor who is not a repository admin makes changes to the `.vela.yml` file. They are not malicious.
- Upon submitting a pull request, their code does not execute for the same reasons in Case 1.
- The code contributor executes a command to re-sign the pipeline. The command fails with a message telling the user that they can only sign the pipeline if they are an admin of the repository.
- The code contributor reaches out to an admin, who reviews the changes, pulls down the branch edits to his own machine, and re-signs the pipeline. 
- The pipeline executes, and the PR is able to be reviewed.

**Case 3**
- A hacker notices that you can trigger a build to run on Pull Requests within Vela. They decide to fork a repository and make PRs to his fork that include changes to the `.vela.yml` file. These changes include a curl request using the Vela native database secret associated with the repository.
- The hacker realizes that his pipeline will not run because his changes do not match the generated signature.
- Hacker proceeds to attempt to exploit Jenkins pipelines instead

### How It Works On Our End

**Changes to Server**
- Create a POST API endpoint that takes in the contents of the `.vela.yml` file (see Decision 3 for more details on this) 
- Endpoint Functionality
  - Check if repository has signature processing enabled (see Decision 2 for more details).
  - Combine the hash field captured from the repository with the representation of the pipeline file to create a signature.
  - Return that signature

**Changes to CLI**
- Create a new pipeline command named `sign`
- Command Functionality
  - Flags include org, repo, file, and path
  - File would search in whichever directory is explicated in the path (defaulting to the cwd).
  - Call the sdk to execute the sign endpoint implemented above
  - Take the signature and add it to the repository. (see Decision 1 for more details).

### Decision 1 - Include signature in .vela.yml or in a separate file (such as .vela.sig)

I can see a case being made for both of these options. On one hand, having a separate file adds a bit of confusion for our users and could potentially run into issues with source control, though I figure both of these complications would be largely infrequent. 

On the other hand, having the signature embedded in the `.vela.yml` file lends its own drawbacks. Perhaps the most clear one to me is the potential for yaml errors with special characters, spaces, and what not. Furthermore, the length and content of a signature would be very incongruent to any pipeline that is otherwise clean and simple.

### Decision 2 - Utilize unused "Trusted" field in Repo or create new field

A relatively small decision: there exists a boolean field, `Trusted`, within the Repo type. Would it be prudent to use this field as a means to determine whether or not a repository has enabled signatures? I.e. if `repo.Trusted == true` then it requires a signature to run builds? This would save us from having to add a new column to the repo table in the database. 

### Decision 3 - Represent .vela.yml file as its full contents or a SHA checksum for encryption

In order to create a signature, we will need to utilize the contents of the `.vela.yml` file to generate the hash. 

When thinking of how to represent these contents, the idea of using a checksum comes to mind. This would be a fixed length representation of the entire contents of the `.vela.yml` file. Even very small changes to the file would produce very different looking checksums. The potential drawback to this is the small possibility of collision. Due to the nature of this feature, I believe the likelihood of collision is incredibly small. However, it is certainly worth mentioning and getting feedback about.

The obvious plus to implementing a checksum is the fact that we would not be sending gigantic pipeline files across the web and generating hash based on that. 

As a note, when Drone implemented signatures, they used checksums for generation.

## Design Option 2 -  Pulling Secrets Validation

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

### Ideal User Workflow 

Ideally, this solution would be the default behavior of Vela. Perhaps we could provide the option to opt out of the validation if a user really wants to, but I don't see why this would be asked for. Therefore, no initial set up would be required for the user like there is for the pipeline solution.

**Case 1**
- A contributor with write access to a repository makes changes to the `.vela.yml` pipeline. Perhaps they are adding Slack plugin functionality.
- Upon submitting a pull request, Vela works just as it currently does.
- The pipeline executes, and the PR is able to be reviewed.

**Case 2**
- A contributor without write access wants to make changes to the `.vela.yml` file. They are not malicious.
- They fork the repository in order to submit a PR
- Upon submitting a pull request, their code does not execute. Instead, an error that looks something like this appears:
```
Error: cannot pull <organization, repository> secrets without write access
```
- An admin (or anyone with write access) assesses the changes to the `.vela.yml` file and comments something that would cause the build to run, or something of that nature.
- The pipeline executes, and the PR is able to be reviewed.

**Case 3**
- A hacker notices that you can trigger a build to run on Pull Requests within Vela. They decide to fork a repository and make a PR that includes changes to the `.vela.yml` file. These changes include a curl request using the Vela native database secret associated with the repository.
- The hacker realizes that his pipeline will not run because they are attempting to pull secrets that belong to an org or a repository of which they do not have write access.
- Hacker proceeds to attempt to exploit Jenkins pipelines instead

### How It Works On Our End

**Changes to Server**
- Upon receiving webhook from SCM, capture user permissions
- After compiling the pipeline, check for the presence of secrets in the `pipeline.Build` type generated by compiler.
- Cross reference user permissions and presence of secrets to determine whether to execute the build. This could be limited to only PR events as well.

### Decision 1 - Opt In, Opt Out, or just be part of Vela

Should we give the users the option to customize this security feature, or should it just be how Vela works? And if it should be customizable, should it be an opt-in feature, or an opt-out feature (default behavior)?

### Decision 2 - Admin or Write Access Approval Necessary

Should the ability to manually check the contents of the `.vela.yml` file and approve it for running be reserved for admins only, or anyone with write access? Should it be customizable?

### Decision 3 - How to enable the build to be triggered by an admin/trusted contributor?

The idea of having a prescribed PR comment event that would release the restrictions on the PR comes to mind. I'm not sure if that would be too clunky. I am open to and eager to hear any ideas on this front. 

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

2-3 weeks

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

## Questions