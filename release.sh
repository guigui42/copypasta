#!/bin/bash
set -euo pipefail

# Usage: ./release.sh [--patch | --minor | --major]
# Defaults to --minor if no argument is given.

BUMP="minor"
while [[ $# -gt 0 ]]; do
  case "$1" in
    --patch) BUMP="patch"; shift ;;
    --minor) BUMP="minor"; shift ;;
    --major) BUMP="major"; shift ;;
    *) echo "Usage: $0 [--patch | --minor | --major]"; exit 1 ;;
  esac
done

# Get the latest semver tag (strip leading 'v')
LATEST_TAG=$(git tag --sort=-v:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | head -1)
if [[ -z "$LATEST_TAG" ]]; then
  echo "No existing semver tag found, starting at v0.1.0"
  LATEST_TAG="v0.0.0"
fi

VERSION="${LATEST_TAG#v}"
IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION"

case "$BUMP" in
  major) MAJOR=$((MAJOR + 1)); MINOR=0; PATCH=0 ;;
  minor) MINOR=$((MINOR + 1)); PATCH=0 ;;
  patch) PATCH=$((PATCH + 1)) ;;
esac

NEW_TAG="v${MAJOR}.${MINOR}.${PATCH}"

echo "Current version: $LATEST_TAG"
echo "Bump type:       $BUMP"
echo "New version:     $NEW_TAG"
echo ""

# Safety check: ensure working tree is clean
if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "Error: Working tree has uncommitted changes. Commit or stash them first."
  exit 1
fi

# Confirm
read -r -p "Tag and push $NEW_TAG? [y/N] " CONFIRM
if [[ "$CONFIRM" != [yY] ]]; then
  echo "Aborted."
  exit 0
fi

git tag "$NEW_TAG"
git push origin "$NEW_TAG"

echo ""
echo "✓ Tagged and pushed $NEW_TAG — GitHub Actions will handle the release."
