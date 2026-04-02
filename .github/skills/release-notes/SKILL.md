---
name: release-notes
description: "Generates Vela release notes for a given target version by fetching commits from the five core go-vela repos (cli, sdk-go, server, ui, worker), filtering and categorizing them by conventional commit type, writing an AI-authored Highlights section, and producing a finished draft at releases/<version>.md. Supports two modes: tagged (compare previous x.(y-1).z to target .0, with a forward-bridge if a repo has already released past target) and unreleased (compare latest non-prerelease tag to HEAD on each repo)."
argument-hint: "<target-version> [--unreleased] (e.g. v0.28 or v0.28 --unreleased)"
user-invokable: true
---

# Release Notes Skill

## Overview

This skill produces a structured release notes draft for a Vela minor release,
following the exact format defined in `template.md`. It is designed to replace
the manual work in `releases/scripts/relgen.sh` with AI-driven categorization
and narrative writing.

## Prerequisites

- The `gh` CLI must be authenticated (`gh auth status`).
- Run from within the `go-vela/community` repository.

## Invocation

```
/release-notes v0.28           # tagged mode  — compares previous x.(y-1).z tag to v0.28.0
/release-notes v0.28 --unreleased  # unreleased mode — compares latest tag to HEAD
```

---

## Constraints

> **All operations performed by this skill must be purely local and read-only
> with respect to remote systems.**

- Do **not** run any `git push`, `git push --tags`, or any variant.
- Do **not** create, publish, edit, or delete any GitHub release (`gh release create`, etc.).
- Do **not** create, merge, or close any pull request or issue.
- Do **not** commit or stage any files (`git commit`, `git add`).
- `gh api` calls are permitted **only** for read operations (GET requests):
  fetching releases, tags, and commit comparisons.
- Do **not** generate helper scripts under `releases/scripts/` (or anywhere in
  the repo) as part of this workflow.
- The only write actions permitted are creating or overwriting
  `releases/<TARGET_VERSION>.md` and `releases/<TARGET_VERSION>-commits.txt`.
  If temporary automation is needed, use ephemeral locations outside the repo
  (for example `/tmp`) and do not leave generated files behind in the workspace.

---

## Step-by-Step Instructions

### 1. Parse Arguments

- Extract `TARGET_VERSION` from the first argument (e.g. `v0.28`).
- Set `UNRELEASED=true` if `--unreleased` is present, otherwise `UNRELEASED=false`.
- Derive `MINOR_SERIES` = major.minor portion (e.g. `0.28`).

### 2. Determine Commit Range Per Repo

The five core repos are: `cli`, `sdk-go`, `server`, `ui`, `worker`.

For each repo, use the `gh` CLI to fetch stable release tags, then sort them:

```bash
gh api "repos/go-vela/<repo>/releases" \
  --jq '[.[] | select(.prerelease != true)] | map(.tag_name)'
```

> ⚠️ The GitHub releases API does **not** return tags in version order — it
> returns them in the order releases were created. Always sort the resulting
> list by semver (descending) before using it. For example, `v0.28.0` may
> appear after `v0.27.3` in the raw API response even though it is the newer
> release. Sort by splitting on `.`, comparing each numeric component
> (major, minor, patch) as integers, newest first.

#### Tagged mode (`UNRELEASED=false`)

Goal: compare the `x.y.0` tag of the target version to the `x.(y-1).z` tag of
the previous minor series, unless the repo has already cut a newer stable
release in which case compare from `x.y.0` forward to that newer tag.

1. Find `TARGET_TAG`: the tag matching `v<MINOR_SERIES>.0` exactly (e.g. `v0.28.0`).
  If no such tag exists for a given repo, skip that repo and note it in the output.
2. Find `LATEST_STABLE_TAG`: the newest semver-sorted stable tag for that repo.
3. Determine range mode:
  - **Standard mode** (default): if `LATEST_STABLE_TAG == TARGET_TAG`, set
    `LAST_TAG = TARGET_TAG`, then iterate the semver-sorted tag list starting
    after `LAST_TAG` and pick the first tag whose major.minor differs from
    `MINOR_SERIES` as `PREVIOUS_TAG`.
  - **Forward-bridge mode**: if `LATEST_STABLE_TAG > TARGET_TAG`, set
    `PREVIOUS_TAG = TARGET_TAG` and `LAST_TAG = LATEST_STABLE_TAG`.
    Example: for target `v0.28`, if UI has `v0.29.0`, use
    `v0.28.0...v0.29.0` for UI.
4. If Standard mode cannot find a prior minor-series tag, skip that repo and
  note it in the output.
5. Do **not** derive `PREVIOUS_TAG` from the raw API order. Always derive both
   `LATEST_STABLE_TAG` and `PREVIOUS_TAG` from the semver-sorted list.
6. Add a per-repo range validation check before fetching commits:
   - `TARGET_TAG` must exist for tagged mode, otherwise skip.
   - In forward-bridge mode, ensure `PREVIOUS_TAG == TARGET_TAG` and
     `LAST_TAG` is the highest semver tag.
   - In standard mode, ensure `PREVIOUS_TAG` and `LAST_TAG` are different.
   If a check fails, skip that repo and note the reason in the output.

#### Unreleased mode (`UNRELEASED=true`)

Goal: capture all work since the last non-prerelease tag on each repo, regardless
of version series.

1. `LAST_TAG` = `HEAD`
2. `PREVIOUS_TAG` = the first (newest) stable tag for that repo — i.e. the
  most recent non-prerelease tag on that repo. Repos may be on
  different series; treat each independently.

### 3. Fetch Commits

For each repo, run:

```bash
# Tagged mode
gh api "repos/go-vela/<repo>/compare/<PREVIOUS_TAG>...<LAST_TAG>" \
  --jq '.commits[] | ...'

# Unreleased mode (LAST_TAG is HEAD)
gh api "repos/go-vela/<repo>/compare/<PREVIOUS_TAG>...HEAD" \
  --jq '.commits[] | ...'
```

Extract per commit:
- `commit.message` (first line only — the subject)
- `html_url` (link to the commit on GitHub)
- `author.login` (GitHub username)
- `author.html_url`

Format each entry as:

```
- (<repo>) <subject with ([#NNN](https://github.com/go-vela/<repo>/pull/NNN)) inline when present> - thanks [@login](https://github.com/login)!
```

Where `#NNN` is the PR number extracted from the commit subject if present
(pattern: `(#NNN)` at the end of the subject line).
When present, replace the trailing `(#NNN)` in the subject with `([#NNN](...))`
instead of appending a second PR link at the end of the line.
If no PR number is present, append a link to the commit URL using the short SHA as the label.

As entries are collected from each repo, append them to an in-memory **raw
commit list** (all entries, unfiltered). This list is written to disk in Step 4.

### 4. Filter Commits

For each entry in the raw commit list, check it against the discard rules below.
Instead of silently dropping filtered entries, **append a `[FILTERED: <reason>]`
tag to the end of the line**. Kept entries remain unchanged.

Discard rules — tag any entry matching **any** of the following:

| Rule | Tag to append |
|---|---|
| Author login ends with `[bot]` (case-insensitive) | `[FILTERED: bot]` |
| Subject does not match conventional commits format | `[FILTERED: not conventional commit]` |
| Type is `chore` or `fix` AND scope contains `deps` | `[FILTERED: deps]` |
| Type is `revert` | `[FILTERED: revert]` |
| Scope starts with `ci` | `[FILTERED: ci scope]` |
| Type is `chore` AND subject contains the word `release` | `[FILTERED: release chore]` |

After tagging, write the **complete annotated list** (kept and filtered entries
together, grouped by repo, in the original fetch order) to:

```
releases/<TARGET_VERSION>-commits.txt
```

For each repo section header in this file, include the resolved compare range
and mode for traceability, for example:

```text
## ui (v0.28.0...v0.29.0, mode=forward-bridge)
## server (v0.27.5...v0.28.0, mode=standard)
```

This file is for human review only — it is excluded from version control via
`releases/.gitignore`. Do not commit it.

Proceed to Step 5 using only the entries that have **no** `[FILTERED: ...]` tag.

### 5. Categorize Commits

Assign each remaining commit to exactly one category based on its type prefix.
Apply the following rules in order:

| Condition in subject | Category |
|---|---|
| contains `!:` | 💥 Breaking Changes |
| type is `feat` | ✨ Features |
| type is `fix` | 🐛 Bug Fixes |
| type is `enhance` | 🚸 Enhancements |
| type is `refactor` | ♻️ Refactors |
| type is `chore` or `docs` or anything else | 🔧 Miscellaneous |

A commit with `!:` goes into Breaking Changes **only**, even if it is also a `feat`.

Within each category, sort entries first by repo name alphabetically, then by
commit subject alphabetically (case-insensitive).

### 6. Write the Highlights Section

Read the categorized commits and write a `### 📣 Highlights` section in plain
language. Guidelines:

- Select the 4–8 most meaningful changes across all repos — prioritize breaking
  changes, user-facing features, and significant enhancements.
- Write one short bullet per highlight (1–2 sentences). Do not just repeat the
  commit subject verbatim; paraphrase into human-readable language.
- Do not include routine chores, test improvements, or minor dependency bumps.
- For unreleased mode, use present-tense / future framing ("Adds...", "Improves...")
  since the release is not yet cut.

### 7. Write the Breaking Changes Admin/User Split

If any breaking changes exist, split them into two sub-sections under
`### 💥 Breaking Changes`:

- `#### Admins` — changes to server configuration, environment variables,
  deployment behavior, database schema.
- `#### Users` — changes to pipeline YAML, CLI behavior, or API responses
  that affect end users.

Each item should include 1–2 sentences of plain-language explanation, then the
raw formatted commit entry on a new line beneath it.

### 8. Assemble the Contributor List

Collect all unique `@login` values from the filtered commit entries. Exclude
any login ending in `[bot]`. Sort case-insensitively. Format as:

```
- @login
```

### 9. Output the Document

Using the structure defined in `template.md` as the exact model, assemble the
complete release notes document and **write it to
`releases/<TARGET_VERSION>.md`** in the repository.

- Replace all `__TARGET_VERSION__` placeholders with the actual version string.
- Populate each section with the categorized commits from step 5.
- Insert the Highlights section from step 6.
- Insert the split Breaking Changes from step 7.
- Use the release links footer and contributor list from steps 3 and 8.
- For unreleased mode, add a callout at the top of the document:
  ```
  > ⚠️ **Draft** — these notes cover unreleased changes as of `HEAD`.
  > Patch release sections will be added manually before publishing.
  ```
- Note at the bottom of the document that patch release sections (e.g. v0.28.1,
  v0.28.2) are not included and should be added manually following the same
  section pattern used in prior releases.

Two files are written in total:

| File | Purpose | Commit to repo? |
|---|---|---|
| `releases/<TARGET_VERSION>.md` | Finished release notes draft | Yes, after human review |
| `releases/<TARGET_VERSION>-commits.txt` | Annotated raw commit audit log | No — excluded via `releases/.gitignore` |

No additional generated helper scripts should be created or kept in the
repository during this process.
