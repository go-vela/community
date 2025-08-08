#!/usr/bin/env bash

# cross repo release notes generation
# relies heavily on the github cli 'gh'.
# part of the reason is that it provides
# an easy method of retrieving GitHub
# user name for attribution.
#
# this is not meant to produce publishable
# release notes, but rather to serve as
# a starting point. a human touch is expected
# after the notes have been generated.
#
# needs sponge, use brew install coreutils
# if full external release has not been cut, change line 45 to LAST_TAG="HEAD"
# TODO:
# - allow passing in flag to generate release notes
#   from last non-rc tag to HEAD
# - reorder output to follow current
#   pre-established output
# - help output

set -Eeuo pipefail

# check for 'gh'
if ! command -v gh &>/dev/null; then
	echo "⚠️ the github cli 'gh' is required. get it here: https://github.com/cli/cli#installation"
fi

RELEASE_FILE="release-$(date +%s).md"
echo "📣 release notes will be generated at: $RELEASE_FILE"

# all the core vela repos
repos=(cli sdk-go server ui worker)

# adding top of file
echo "⚓ adding main header"
printf '# __TARGET_VERSION__ 🚀\n\nThis document contains all release notes pertaining to the `__TARGET_VERSION__.x` releases of Vela.\n\n/' >>"$RELEASE_FILE"

# main loop to iterate over repos
echo "📣 generating release notes for core vela repos"
for repo in "${repos[@]}"; do
	# get stable release tags (newest first)
    mapfile -t TAGS < <(gh api "repos/go-vela/$repo/releases" --jq '.[] | select(.prerelease != true) | .tag_name')

	# helpers to parse semver-like tags (e.g. v0.27.1)
    mm() { sed -E 's/^[^0-9]*//' <<<"$1" | awk -F. '{print $1 "." $2}'; }   # major.minor
    patch() { sed -E 's/^[^0-9]*//' <<<"$1" | awk -F. '{print ($3+0)}'; }   # patch as int
	
	# default to newest two tags
	LAST_TAG="${TAGS[0]:-}"
	PREVIOUS_TAG="${TAGS[1]:-}"

	# If newest is a patch (x.y.z, z>0), anchor to most recent x.y.0
    if [[ -n "$LAST_TAG" && "$(patch "$LAST_TAG")" -ne 0 ]]; then
        series="$(mm "$LAST_TAG")"
        for t in "${TAGS[@]}"; do
            if [[ "$(mm "$t")" == "$series" && "$(patch "$t")" -eq 0 ]]; then
                LAST_TAG="$t"
                break
            fi
        done
    fi

    # Pick previous tag from a different series than LAST_TAG
    if [[ -n "$LAST_TAG" ]]; then
        last_series="$(mm "$LAST_TAG")"
        last_idx=-1
        for i in "${!TAGS[@]}"; do
            if [[ "${TAGS[$i]}" == "$LAST_TAG" ]]; then
                last_idx="$i"
                break
            fi
        done
        for ((j=last_idx+1; j<${#TAGS[@]}; j++)); do
            t="${TAGS[$j]}"
            if [[ "$(mm "$t")" != "$last_series" ]]; then
                PREVIOUS_TAG="$t"
                break
            fi
        done
    fi

    # fallback if no prior series found
    : "${PREVIOUS_TAG:=${TAGS[1]:-}}"

	printf "📝 fetching entries for for %s (from %s to %s)\n" "$repo" "$PREVIOUS_TAG" "$LAST_TAG"

	# fetch the changes since between given tags
	# awk prepends each line with "- (<repo>) " to maintain a reference
	gh api "repos/go-vela/$repo/compare/$PREVIOUS_TAG...$LAST_TAG" \
		--jq '.commits.[]
			| select(((.author.login // "") | endswith("[bot]")) | not)
			| "\(.commit.message|split(" \\(#[0-9]{3,}\\)";"")[0]) [\(.commit.message|capture("\\((?<pr>#[0-9]{3,})"; "")|.pr)](\(.html_url)) - thanks [@\(.author.login)](\(.author.html_url))!"' |
		awk -v repo="$repo" '{print "- (" repo ") " $0}' \
			>>"$RELEASE_FILE"
done

# getting unique contributors
echo "📣 creating contributor list"
CONTRIBUTORS="$(perl -ne 'if(/\[(@[a-z0-9\[\]_-]+)\]\(/i) { print "- $1\n";}' "$RELEASE_FILE" | grep -vi '\[bot\]' | sort --ignore-case --unique)"

# filter commits
# - only keep conventional commit formatted commits
# - ignore dependency updates, reverts, and release commits
cat "$RELEASE_FILE" |
	grep --ignore-case --extended-regexp '^-\s+\([a-z\-]+\)\s+[a-z]+(\([a-z_\/[:space:]\-]+\))?!?:\s.+' |
	grep --invert-match --ignore-case --extended-regexp '^-\s+\([a-z\-]+\)\s+(chore|fix)\(deps\)' |
	grep --invert-match --ignore-case --extended-regexp '^-\s+\([a-z\-]+\)\s+revert(\([a-z_\/[:space:]\-]+\))?:' |
	grep --invert-match --ignore-case --extended-regexp '^-\s+\([a-z\-]+\)\s+[a-z]+\(ci(\/[a-z]+)?\)(!?)?:' |
	grep --invert-match --ignore-case --extended-regexp '^-\s+\([a-z\-]+\)\s+chore.*release' |
	sponge "$RELEASE_FILE"

# sort releases by type (fixes, features, etc) and and then by component (server, cli, etc)
echo "🔃 sorting release notes"
sort --ignore-case --ignore-leading-blanks --key=3,3.3 --key=2,2 "$RELEASE_FILE" | sponge "$RELEASE_FILE"

# move breaking changes to the beginning of the file
# store matching and non-matching lines in array and return them
perl -ne 'push @{/!:/ ? \@A : \@B}, $_}{print @A, @B' "$RELEASE_FILE" | sponge "$RELEASE_FILE"

echo "📑 demarking sections"
# each line here, find the first occurence of the given pattern
# in a line and inserts a heading above it.
awk '/!:/ && !x {print "\n### 💥 Breaking Changes\n"; x=1} 1' "$RELEASE_FILE" | sponge "$RELEASE_FILE"
awk '/(fix)(\([a-z\-_ ]+\))?:/ && !x {print "\n### 🐛 Bug Fixes\n"; x=1} 1' "$RELEASE_FILE" | sponge "$RELEASE_FILE"
awk '/(refactor)(\([a-z\-_ ]+\))?:/ && !x {print "\n### ♻️ Refactors\n"; x=1} 1' "$RELEASE_FILE" | sponge "$RELEASE_FILE"
awk '/(feat)(\([a-z\-_ ]+\))?:/ && !x {print "\n### ✨ Features\n"; x=1} 1' "$RELEASE_FILE" | sponge "$RELEASE_FILE"
awk '/(enhance)(\([a-z\-_ ]+\))?:/ && !x {print "\n### 🚸 Enhancements\n"; x=1} 1' "$RELEASE_FILE" | sponge "$RELEASE_FILE"
awk '/(chore)(\([a-z\-_ ]+\))?:/ && !x {print "\n### 🔧 Miscellaneous\n"; x=1} 1' "$RELEASE_FILE" | sponge "$RELEASE_FILE"

echo "🐙 adding repo release links"
printf "\n## 🔗 %s\n\n" "Release Links" >>"$RELEASE_FILE"
for repo in "${repos[@]}"; do
	echo "- https://github.com/go-vela/${repo}/releases" >>"$RELEASE_FILE"
done

echo "🥹 adding contributors"
printf "\n## 💟 Thank you to all the contributors in this release!\n\n%s\n" "$CONTRIBUTORS" >>"$RELEASE_FILE"

echo "🪩 all done! see '$RELEASE_FILE' for your release notes."
