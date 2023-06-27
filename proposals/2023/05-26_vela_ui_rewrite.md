# Rewrite UI in a more common framework

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value |
| :-----------: | :---: |
| **Author(s)** |       |
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

<!--
Provide your description here.
-->

Currently, the Vela UI is written in [ELM](https://elm-lang.org/), which is a fairly obscure language with a small community.
Using a more common framework, such as React or Angular would significantly lower the barrier to entry to contribute to the UI.

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->
Not required, however it should increase the number of contributors and contributions to the UI.

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->
- It is difficult to contribute to the UI, it is unlikely people who are interested in contributing know anything about ELM.
- In order to contribute it is likely the developer would need to learn a new language, which would put people off from contributing.
- The main project files are 5,000 - 7,000 lines of code each, which makes it difficult to find / figure out how things are tied together.
- ELM does not lend itself to code usability resulting in a copy / paste / rename / repurpose coding style.
- There hasn't been any active development in the ELM framework for nearly 4 years, there are many pull requests that have been sitting for 3+ years [https://github.com/elm/elm-lang.org](https://github.com/elm/elm-lang.org)

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->
- A major refactor of the current UI
   - Clear and consistent patterns
   - Reusable components
   - A clear guide around the expected patterns

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->

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

- Replicate current UI in another more common framework. Functionality should match or exceed current functionality.
- Design should align with current style and visual styles.

## Implementation

<!--
This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

- Since the UI is a different repo and there are no co-dependency's the new UI should be a drop-in replacement.

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

<!-- Answer here -->
I plan to contribute along with a few others.

2. What's the estimated time to completion?

<!-- Answer here -->
30-60 days

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

## Questions

**Please list any questions you may have:**

<!-- Answer here -->