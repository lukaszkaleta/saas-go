#!/usr/bin/env bash
set -e

# --- CONFIGURATION ---
VERSION_FILE="VERSION"   # file storing the version string (e.g., v0.2.190)

# --- READ CURRENT VERSION ---
if [[ ! -f "$VERSION_FILE" ]]; then
  echo "VERSION file not found. Create one with initial version like: v0.2.190"
  exit 1
fi

CURRENT_VERSION=$(cat "$VERSION_FILE")
echo "Current version: $CURRENT_VERSION"

# --- PARSE VERSION vMAJOR.MINOR.PATCH ---
if [[ ! $CURRENT_VERSION =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
  echo "Invalid version format in VERSION file. Expected: v0.2.190"
  exit 1
fi

MAJOR="${BASH_REMATCH[1]}"
MINOR="${BASH_REMATCH[2]}"
PATCH="${BASH_REMATCH[3]}"

# --- INCREMENT PATCH ---
NEW_PATCH=$((PATCH + 1))
NEW_VERSION="v${MAJOR}.${MINOR}.${NEW_PATCH}"

echo "Bumping version to: $NEW_VERSION"

# --- REPLACE OCCURRENCES IN ALL FILES ---
# macOS: use -i '' ; Linux: use -i
if sed --version >/dev/null 2>&1; then
    # GNU sed (Linux)
    sed -i "s/${CURRENT_VERSION}/${NEW_VERSION}/g" -r $(grep -rl "$CURRENT_VERSION" .)
else
    # BSD sed (macOS)
    sed -i '' "s/${CURRENT_VERSION}/${NEW_VERSION}/g" $(grep -rl "$CURRENT_VERSION" .)
fi

# --- UPDATE VERSION FILE ---
echo "$NEW_VERSION" > "$VERSION_FILE"

echo "Done!"
