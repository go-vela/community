# Templates

<!-- Please leave this commented out section.

The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value                          |
| :-----------: | :----------------------------: |
| **Author(s)** | Neal.Coleman, Jordan.Brockopp  |
| **Reviewers** | Emmanuel.Meinen, Kelly.Merrick |
| **Date**      | July 1st, 2019                 |
| **Status**    | Complete                       |

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

This feature will allow customers to create pipelines across multiple repos from a single blueprint.

This blueprint (a.k.a. template) can be a customized version written by the customer, or an "officially supported" template written by the Vela admins themselves.

The template will support the ability to inject variables from the pipeline sourcing the template to enable flexible utilization and account for atypical use-cases.

This has the added benefit of enabling features like "matrix pipelines" that allow repeating the same set of tasks with different variables.

Similar to the pipeline configuration itself, templates will be written in YAML and enable specific attributes, or even entire sections of the template to be controlled through defined variables.

**Please briefly answer the following questions:**

1. Why is this required?

* customers spend a large amount of time copying CI pipelines across repos and updating few fields.
* bridge the gap of supporting mono repos - Mono repos tend to have identical code in multiple places
* support "matrix builds" that existing CI solutions can perform

2. If this is a redesign or refactor, what issues exist in the current implementation?

N/A

3. Are there any other workarounds, and if so, what are the drawbacks?

The only workaround that exists is to manually copy and paste the same YAML pipeline across multiple repos. Several drawbacks exist for this method:

* miscalculating the part(s) of the YAML pipeline that are to be copied
* amount of time and effort required to copy/paste between repos
* keeping all repos in sync (up to date) that attempt to use the same YAML

4. Are there any related issues? Please provide them below if any exist.

N/A

## Design

<!-- Please leave this commented out section.

This section is intended to explain the solution design for the proposal.

NOTE: If there are no current plans for a solution, please leave this section blank.
-->

**Please describe your solution to the proposal. This includes, but is not limited to:**

* new/updated endpoints or url paths
* new/updated configuration variables (environment, flags, files, etc.)
* performance and user experience tradeoffs
* security concerns or assumptions
* examples or (pseudo) code snippets

### Option 1

```yaml
version: "1"
steps:
  - name: tmpl
    template: true
    source: github
      local: git.example.com/template-repo/stable/spring
    variables:
      image_name: hello_world
      application: hello_world

  - name: slack
    image: target/vela-slack:latest
    format:
      message: You ran a build
```

### Option 2

```yaml
version: "1"

metadata:
  template:
    name: git.example.com/template-repo/stable/spring
    vars:
      image_name: hello_world
      application: hello_world

steps:
  - name: tmpl
    template: true

  - name: slack
    image: target/vela-slack:latest
    format:
      message: You ran a build
```

### Option 3

```yaml
version: "1"

templates:
  - name: spring
    source: git.example.com/template-repo/stable/spring
    vars:
      image_name: hello_world
      application: hello_world

steps:
  - template: spring

  - name: slack
    image: target/vela-slack:latest
    format:
      message: You ran a build
```

### Option 4

```yaml
version: "1"

steps:
  - template:
      name: spring
      source: git.example.com/template-repo/stable/spring
      vars:
        image_name: hello_world
        application: hello_world

  - name: slack
    image: target/vela-slack:latest
    format:
      message: You ran a build
```

### Option 5

```yaml
version: "1"

templates:
  - name: spring
    source: git.example.com/template-repo/stable/spring

steps:
  - name: tmpl # When using templates name is optional
    template:
      name: spring
      vars:
        image_name: hello_world
        application: hello_world


  - name: slack
    image: target/vela-slack:latest
    format:
      message: You ran a build
```

## Implementation

<!-- Please leave this commented out section.

This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

Yes

2. What's the estimated time to completion?

2 weeks

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

After some trials and discussions, [Option 5](07-01_templates.md#option-5) was chosen!

## Questions

**Please list any questions you may have:**

N/A