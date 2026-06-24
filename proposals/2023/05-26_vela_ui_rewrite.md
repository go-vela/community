# Rewrite UI in a more common framework

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

|      Key      |                     Value                      |
| :-----------: | :--------------------------------------------: |
| **Author(s)** | Ryan Rampersad, Zac Skalko, James ⁠Christensen |
| **Reviewers** |                                                |
|   **Date**    |                   2023-06-30                   |
|  **Status**   |                                                |

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

Currently, the Vela UI is written in [Elm](https://elm-lang.org/). Elm is a domain-specific programming language for declaratively creating web browser-based graphical user interfaces. Elm is purely functional, and is developed with emphasis on usability, performance, and robustness.

Elm is a fairly obscure language with a small community, anecdotally in our contributor's experience, anecdotally by way of these recent Hacker News threads about Elm [(1)](https://news.ycombinator.com/item?id=34746161) [(2)](https://news.ycombinator.com/item?id=35495910) [(3)](https://news.ycombinator.com/item?id=36272692) and statistically according to the [2022 State of JS Survey](https://2022.stateofjs.com/en-US/other-tools/#javascript_flavors) and [2023 Stackoverflow Survey](https://survey.stackoverflow.co/2023/#technology).

Therefore we propose that a using more popular and widely used language such as TypeScript and a more popular and widely used framework such as React would significantly lower the barrier to entry to casual contributors to the Vela UI.

<!-- We understand that Elm was originally for this codebase because... _please contribute this information_. -->

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

Not required, however it should increase the number of contributors and contributions to the UI.

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

- It is difficult to contribute to the UI, it is unlikely people who are interested in contributing know anything about Elm
- In order to contribute it is likely the developer would need to learn a new language, which would put people off from contributing
- The main project files are 5,000 - 7,000 lines of code each, which makes it difficult to find / figure out how things are tied together
- ELM does not lend itself to code usability resulting in a copy / paste / rename / repurpose coding style.
- There hasn't been any active development in the ELM framework for nearly 4 years, there are [many pull requests that have been sitting for 3+ years](https://github.com/elm/elm-lang.org/pulls?q=is%3Apr+is%3Aopen+sort%3Acreated-asc)

3. Are there any other workarounds, and if so, what are the drawbacks?
   <!-- Answer here -->

   Alternatives to a redesign and rewrite could be...

   1. maintaining the existing Elm codebase
   2. creating and maintaining an extensive Elm-centric guide on maintenance activities, and new feature creation activities

   In general, while possible, this may not necessarily increase the overall casual contributions.

   However, with a redesign and rewrite...

   With a focus on React and other widely popular and community accessible tooling, we think we can increase casual contribution.

   This new rewrite will have:

   - common React patterns
   - common SPA patterns
   - reusable components
   - baseline guidance for creating new pages, widgets & components
   - baseline guidance for introducing new dependencies

4. Are there any related issues? Please provide them below if any exist.

  As part of our initial exploration, we noticed a few discrepancies with the Swagger/OpenAPI spec. We will contribute to help address those issues.


## Design

<!--
This section is intended to explain the solution design for the proposal.

NOTE: If there are no current plans for a solution, please leave this section blank.
-->

We have a suggested set of technologies to be considered for this redesign and rewrite, with justifications:

- [TypeScript](https://www.typescriptlang.org/)
  - [popular](https://survey.stackoverflow.co/2023/#programming-scripting-and-markup-languages), widely used, industry standard language with well known tooling
  - enables type safety for better reliability and better developer experience
- [React](https://react.dev/)
  - [popular](https://survey.stackoverflow.co/2023/#technology), widely used
  - enables common component patterns
  - can leverage popular react libraries
  - most _app_ developers have used react or a framework resembling react (vue, svelte)
- [TailwindCSS](https://tailwindcss.com/)
  - [popular](https://2022.stateofcss.com/en-US/css-frameworks/), utility first centric "framework" for styles
  - easily portable to another framework in the future if necessary
  - easy to isolate styles for particular components
  - easy to create ad-hoc or situational tweaks
  - co-locates component markup/structure and style, reducing cross referencing and hunting
- [Radix UI](https://www.radix-ui.com/)
  - Radix UI is a "headless" that provides accessibility features and enhanced experiences
- [React Query](https://tanstack.com/query/latest/docs/react/overview)
  - [popular](https://www.npmjs.com/package/@tanstack/react-query), a data-fetching and state management library for React applications that simplifies fetching, caching, and updating data, and increases user experience by streamlining loading ui and increases developer experience by providing a consistent api for loading data and maintaining transient state
- [Vite](https://vitejs.dev/)
  - [popular](https://www.npmjs.com/package/vite), [increasingly used](https://2022.stateofjs.com/en-US/libraries/), alternative to the now effectively deprecated [create react app](https://create-react-app.dev/)
  - enables an out of box single-page-app solution
  - enables better developer experience with fast refresh and fast build times
  - still extensible should special needs arise
- [Vitest](https://vitest.dev/)
  - [popular](https://www.npmjs.com/package/vitest), a testing library with strong integration with Vite, and familiar to those with Jest flavored unit testing experience
  - works well in conjunction with [react-testing-library](https://testing-library.com/docs/react-testing-library/intro/)
- [Cypress](https://www.cypress.io/)
  - re-using the existing cypress tests as a foundation to reach and maintain coverage
- [Prettier](https://prettier.io/)
  - `gofmt` alike, automatic formatting without debate
- [eslint](https://eslint.org/)
  - `golangci-lint` alike, automatic linting, recommending best coding practices

<details>
<summary><b>Why react?</b></summary>

React has been a staple frontend "framework" in 2023. It has wide adoption and has a massive practitioner base. In the [2023 Stackoverflow Survey](https://survey.stackoverflow.co/2023/#section-admired-and-desired-web-frameworks-and-technologies), it continues to earn top a top position in desirability and admiration.

We think that React, among others "frameworks" like Angular, Vue or Svelte fulfills our desire to increase casual contributors in this codebase. The principle reason is the sheer momentum that React has.

React, like most frontend tooling, will enter its halcyon days and the industry will move to new tooling. There is no such thing as future-proof. However, in an effort to be future-resistant we are specifically making choices that will be easy to change around build tooling. We are not entirely leveraging the React core team's ecosystem (Next, CRA, Jest), and instead using a broader community ecosystem (Vite).

</details>

<details>
<summary><b>Can we achieve the same and better results than Elm with this new approach?</b></summary>

Elm offered a bundled developer experience for webapps. The React-centric ecosystem is much more fragmented, but it lets you pick from a variety of options that are easier to plug in and to swap as needs change.

</details>

We have considered these functional changes:

- new/updated endpoints or url paths
  - re-evaluating the `/-/secrets` path
  - re-evaluating paths that shadow potentially valid Organization, Repository, or Secret names
- new/updated configuration variables (environment, flags, files, etc.)
  - environment variables for features (e.g. default pull request settings, auto-linkifying in logs, etc)

From an overall perspective, the goal is to:

- Replicate the current user interface (look'n'feel; design) in the technology stack listed above as closely as possible
- Replicate the current functionality in the technology stack listed above as closely as possible

Consider this redesign and rewrite a `1-to-1.3̅3̅3̅` port rather than a `1-to-1` port.

Further enhancements can and will be made after.

## Implementation

<!--
This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

- Since the UI is a different repo and there are no co-dependency's the new UI should be a drop-in replacement and in theory should be able to run in parallel to the existing version if necessary.

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

Yes. There is a small team working on this (see below; hackathon activity).

2. What's the estimated time to completion?

30-60 days

1. During the 2023 Summer Vela Hackathon, we explored this space and within roughly 62 hours we have ported a majority of the existing look'n'feel and functionality
2. While a majority of read-centric screens exist, there are write-centric screens that need implementation
3. There is still plenty of polishing, documentation and cleanup activities thereafter

**Gists, issues, pull requests, etc**

Please refer to the ongoing work in [`go-vela/ui-hackathon`](https://github.com/go-vela/ui-hackathon) repo for more details.

## Questions

**Please list any questions you may have:**

- Can you share some additional details about how Elm was picked as the initial language/framework option for this codebase? We'd like to include this history and the pick's original goals as part of the proposal background.
