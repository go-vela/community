#!/usr/bin/env bash

# cross repo release notes generation
# relies heavily on the github cli 'gh'
# and `gsed` (for Mac usage)
#
# TODO:
# - allow passing in flag to generate release notes
#   from last non-rc tag to HEAD
# - reorder output to follow current
#   pre-established output
# - help output
# - sed/gsed detection

set -Eeuo pipefail

# check for 'gh'
if ! command -v gh &>/dev/null; then
	echo "⚠️ the github cli 'gh' is required. get it here: https://github.com/cli/cli#installation"
fi

RELEASE_FILE="release-$(date +%s).md"
echo "📣 release notes will be generated at: $RELEASE_FILE"

# all the core vela repos
repos=(cli sdk-go server types ui worker)

# main loop to iterate over repos
echo "📣 generating release notes for core vela repos"
for repo in "${repos[@]}"; do
    # get the tags
    LAST_TWO_TAGS="$(gh api repos/go-vela/$repo/releases --jq '[.[] | select(.prerelease != true) | .tag_name] | join(" ")')"
    LAST_TAG="$(echo $LAST_TWO_TAGS | awk '{print $1}')"
    PREVIOUS_TAG="$(echo $LAST_TWO_TAGS | awk '{print $2}')"

	LAST_NON_RC_TAG="$(gh release view --repo "go-vela/$repo" --json tagName --jq '.tagName')"
	printf "📝 fetching entries for for %s (from %s to %s)\n" "$repo" "$PREVIOUS_TAG" "$LAST_TAG"

	# fetch the changes since last tag and append to file
	# awk prepends each line with "- (<repo>) " to maintain a reference
	gh api "repos/go-vela/$repo/compare/$PREVIOUS_TAG...$LAST_TAG" \
		--jq '.commits.[] | "\(.commit.message|split(" \\(#[0-9]{3,}\\)";"")[0]) [\(.commit.message|capture("\\((?<pr>#[0-9]{3,})"; "")|.pr)](\(.html_url)) - thanks [@\(.author.login)](\(.author.html_url))!"' |
		awk -v repo="$repo" '{print "- (" repo ") " $0}' \
			>>"$RELEASE_FILE"
done

# getting unique contributors
echo "📣 creating contributor list"
CONTRIBUTORS="$(perl -ne 'if(/\[(@[a-z0-9\[\]_-]+)\]\(/) { print "- $1\n";}' "$RELEASE_FILE" | sort -u)"

# filter commits
# - only keep conventional commit formatted commits
# - ignore dependency updates, reverts, and release commits
cat "$RELEASE_FILE" |
grep --ignore-case --extended-regexp "^-\s+\([a-z\-]+\)\s+[a-z]+(\([a-z\-_ ]+\))?!?:\s.+" |
grep --invert-match --ignore-case --extended-regexp "^-\s+\([a-z\-]+\)\s+(chore|fix)\(deps\)" |
grep --invert-match --ignore-case --extended-regexp "^-\s+\([a-z\-]+\)\s+revert(\([a-z\-_ ]+\))?:" |
grep --invert-match --ignore-case --extended-regexp "^-\s+\([a-z\-]+\)\s+chore.*release" |
sponge "$RELEASE_FILE"

# sort releases by type (fixes, features, etc) and and then by component (server, cli, etc)
echo "🔃 sorting release notes"
sort --ignore-case --ignore-leading-blanks --key=3,3.3 --key=2,2 "$RELEASE_FILE" | sponge "$RELEASE_FILE"

# move breaking changes to the beginning of the file
# store matching and non-matching lines in array and return them
perl -ne 'push @{/!:/ ? \@A : \@B}, $_}{print @A, @B' "$RELEASE_FILE" | sponge "$RELEASE_FILE"

echo "⚓ adding main header"
gsed -i '1s/^/## Release Notes 🚀\n\n/' "$RELEASE_FILE"

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
printf "\n## 🔗 %s\n\n" "Release Links" >> "$RELEASE_FILE"
for repo in "${repos[@]}"; do
	echo "- https://github.com/go-vela/${repo}/releases" >> "$RELEASE_FILE"
done

echo "🥹 adding contributors"
printf "\n## 💟 Thank you to all the contributors in this release!\n\n%s\n" "$CONTRIBUTORS" >> "$RELEASE_FILE"

echo "🪩 all done! see '$RELEASE_FILE' for your release notes."