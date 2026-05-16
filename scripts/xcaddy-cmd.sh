#!/bin/sh
set -e

BUNDLE="${1:-homelab}"
MODULES_JSON="$(cd "$(dirname "$0")/.." && pwd)/modules.json"

RETAG=$(jq -r --arg b "$BUNDLE" '.bundles[] | select(.name == $b) | .retag // ""' "$MODULES_JSON")
if [ -n "$RETAG" ]; then
    echo "Bundle '$BUNDLE' is a retag of $RETAG; no xcaddy build needed."
    exit 0
fi

MODULES=$(jq -r --arg b "$BUNDLE" \
    '.bundles[] | select(.name == $b) | .modules[] | "--with \(.import)@\(.version)"' \
    "$MODULES_JSON" | tr '\n' ' ')

echo "xcaddy build $MODULES"
