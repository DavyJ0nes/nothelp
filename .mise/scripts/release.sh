#!/usr/bin/env bash
# Cut a release: choose the version, tag it, and publish binaries + release notes
# with GoReleaser.
#
# Version selection:
#   - Override: if RELEASE_VERSION (or the first argument) is set, that exact
#     version is used; a leading "v" is optional. Use this to skip the automatic
#     bump — e.g. to step over a version number GitHub has wedged.
#   - Auto: otherwise the next version is computed from Conventional Commits since
#     the latest tag, matching the old allow-initial-development-versions cadence:
#       feat -> minor, fix/perf -> patch, breaking -> minor while 0.x else major.
#     A push containing only chore/docs/ci/etc. commits cuts no release.
#
# RELEASE_DRY_RUN=1 prints the chosen version and exits before tagging/publishing.
# The publish step needs GITHUB_TOKEN.
set -euo pipefail

override="${RELEASE_VERSION:-${1:-}}"
override="${override#v}"

if [ -n "${override}" ]; then
	if ! printf '%s' "${override}" | grep -qE '^[0-9]+\.[0-9]+\.[0-9]+$'; then
		echo "release: invalid RELEASE_VERSION '${override}' (expected MAJOR.MINOR.PATCH)" >&2
		exit 1
	fi
	next="${override}"
	echo "release: using override version ${next}"
else
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
	echo "release: ${current} -> ${next} (${bump} bump)"
fi

tag="v${next}"

if [ "${RELEASE_DRY_RUN:-}" = "1" ]; then
	echo "release: dry run, stopping before tag/publish."
	exit 0
fi

# Only override the committer identity in CI; locally use the developer's own.
if [ -n "${GITHUB_ACTIONS:-}" ]; then
	git config user.name "github-actions[bot]"
	git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
fi

# Create the tag locally if it isn't there yet.
if ! git rev-parse -q --verify "refs/tags/${tag}" >/dev/null; then
	git tag -a "${tag}" -m "${tag}"
fi

# Push the tag only if the remote doesn't already have it.
if ! git ls-remote --exit-code --tags origin "refs/tags/${tag}" >/dev/null 2>&1; then
	git push origin "${tag}"
fi

goreleaser release --clean
