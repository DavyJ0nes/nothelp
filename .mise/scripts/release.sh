#!/usr/bin/env bash
# Cut a release: compute the next version from Conventional Commits, tag it, and
# publish binaries + release notes with GoReleaser.
#
# Versioning preserves the project's established cadence (and the previous
# go-semantic-release `allow-initial-development-versions` behaviour):
#   feat            -> minor
#   fix / perf      -> patch
#   breaking change -> minor while 0.x, major once >= 1.0
# A push containing only chore/docs/ci/etc. commits cuts no release.
#
# Set RELEASE_DRY_RUN=1 to print the computed version and exit before tagging or
# publishing. The publish step needs GITHUB_TOKEN.
set -euo pipefail

current="$(convco version)"
range="v${current}..HEAD"

subjects="$(git log --no-merges --format='%s' "${range}" || true)"
bodies="$(git log --no-merges --format='%B' "${range}" || true)"

bump=""
if printf '%s\n' "${subjects}" | grep -qE '^[a-zA-Z]+(\([^)]*\))?!:' ||
	printf '%s\n' "${bodies}" | grep -qE '^BREAKING[ -]CHANGE:'; then
	bump="major"
elif printf '%s\n' "${subjects}" | grep -qE '^feat(\([^)]*\))?:'; then
	bump="minor"
elif printf '%s\n' "${subjects}" | grep -qE '^(fix|perf)(\([^)]*\))?:'; then
	bump="patch"
fi

if [ -z "${bump}" ]; then
	echo "release: no feat/fix/perf/breaking commits since v${current}; nothing to release."
	exit 0
fi

IFS='.' read -r major minor patch <<<"${current}"

case "${bump}" in
major)
	if [ "${major}" -eq 0 ]; then
		# Stay in 0.x, matching the old allow-initial-development-versions setting.
		minor=$((minor + 1))
		patch=0
	else
		major=$((major + 1))
		minor=0
		patch=0
	fi
	;;
minor)
	minor=$((minor + 1))
	patch=0
	;;
patch)
	patch=$((patch + 1))
	;;
esac

next="${major}.${minor}.${patch}"
tag="v${next}"
echo "release: ${current} -> ${next} (${bump} bump)"

if [ "${RELEASE_DRY_RUN:-}" = "1" ]; then
	echo "release: dry run, stopping before tag/publish."
	exit 0
fi

if git rev-parse -q --verify "refs/tags/${tag}" >/dev/null; then
	echo "release: tag ${tag} already exists; skipping."
	exit 0
fi

git config user.name "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
git tag -a "${tag}" -m "${tag}"
git push origin "${tag}"

goreleaser release --clean
