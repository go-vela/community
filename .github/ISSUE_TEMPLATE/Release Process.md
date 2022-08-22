---
name: Release Process
about: Repeatable release process
title: ''
labels: feature
assignees: ''

---

for rc1:

## Release Story 1: cut v0.x.0-rc1
External

## Timing
Can be concurrent with stories 2 and 3

## Background
This cuts the go-vela external (pre-release) rc1 https://pages.git.target.com/rapid/doc-site/13-vela/upgrades/#external
https://github.com/go-vela/server/releases (similar for types, sdk-go, worker, cli, ui)
fyi, ui can be released separately, so is currently n+1 version (server 14, ui 15)



## Release Story 2: prep draft pr for release notes
External

## Timing
Can be concurrent with stories 1 and 3

## Background
pr sometimes combined with migration tool, but could still be separate story
https://github.com/go-vela/community/pull/591



## Release Story 3: prep draft pr for migration tool
External

## Timing
Can be concurrent with stories 1 and 2

## Background
pr sometimes combined with release notes, but could still be separate story
migration tool often necessary for database changes
simple, if there is no migration tool necessary for this release like https://github.com/go-vela/community/pull/439 (fyi, this one combined with release notes)
more involved, if there is a migration tool like https://github.com/go-vela/community/pull/592



## Release Story 4: release rc1 to dev / run dev vela-canary tests
Internal

## Timing
After Story 3 (external)

## Action
https://pages.git.target.com/rapid/doc-site/13-vela/upgrades/#internal and 
specifically https://pages.git.target.com/rapid/doc-site/13-vela/upgrades/#dev

https://vela.dev.target.com/vela-canary/all-repo restart a build => results in all the other canary repos kicking off a build near the same time (a good way to watch them all is https://vela.dev.target.com/vela-canary/builds)



## Release Story 5: test rc1 in dev
Internal 

## Timing
After Step 4 (internal)

## Action
manually test any new bug fixes / enhancements / features
look for bugs
add vela-canary modifications/additions to story 7



## Release Story 6: rc1 bugs
External

## Background
keeping 1 story for the sake of this summary exercise; could be 1 where we link all bugs external bugs, or separate stories 1 per bug



## Release Story 7: add/modify / validate canary tests
Internal

## Background
keeping 1 story for the sake of this summary exercise; could be 1 or several stories for a single release flow
either split by add/modify vs validate, or split by canary test, or keep 1 long-running story for entire release process, or by bug/enhancement/feature (this way is gross), or…hmm



for rc2+ (if needed):
equivalent stories 1, 4, 5, 6; review 2, 3, 7 for any pertinent changes



for external release of v0.x.0:
this cuts the go-vela external release https://pages.git.target.com/rapid/doc-site/13-vela/upgrades/#external (same link as for rc)



## Release Story 8: cut v0.x.0
External 

## Timing
Can be concurrent with stories 9 and 10

## Action
This cuts the go-vela external release https://pages.git.target.com/rapid/doc-site/13-vela/upgrades/#external
https://github.com/go-vela/server/releases (similar for types, sdk-go, worker, cli, ui)
fyi, ui can be released separately, so is currently n+1 version (server 14, ui 15)
update homebrew https://github.com/go-vela/homebrew-vela/blob/master/Formula/vela%400.13.rb



## Release Story 9: approve/merge pr for release notes
External

## Timing
Can be concurrent with stories 8 and 10

## Action
pr sometimes combined with migration tool, but could still be separate story
https://github.com/go-vela/community/pull/591



## Release Story 10 approve/merge pr for migration tool
External 

## Timing
Can be concurrent with stories 8 and 9

## Background
pr sometimes combined with release notes, but could still be separate story
migration tool often necessary for database changes
simple, if there is no migration tool necessary for this release like https://github.com/go-vela/community/pull/439 (fyi, this one combined with release notes)
more involved, if there is a migration tool like https://github.com/go-vela/community/pull/592



## Release Story 11: release to prod / run prod vela-canary tests
Internal 

## Timing
After Story 10

## Action
https://pages.git.target.com/rapid/doc-site/13-vela/upgrades/#internal and specifically https://pages.git.target.com/rapid/doc-site/13-vela/upgrades/#prod
https://vela.prod.target.com/vela-canary/all-repo restart a build => results in all the other canary repos kicking off a build near the same time (a good way to watch them all is https://vela.prod.target.com/vela-canary/builds)
