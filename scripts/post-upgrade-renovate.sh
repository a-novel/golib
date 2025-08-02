#!/bin/bash

set -e

# ======================================================================================================================
# Install node on the renovate image.
FNM_DIR="$HOME/.fnm"
FNM="$FNM_DIR/fnm"

curl -o- https://fnm.vercel.app/install | bash -s -- --install-dir "$FNM_DIR"
eval "$("$FNM" env --log-level --fnm-dir "$FNM_DIR")"
"$FNM" install --latest --log-level error

which node || echo "node install failed" && exit 1
which npm || echo "npm install failed" && exit 1
which npx || echo "npx install failed" && exit 1
# ======================================================================================================================

npx -y prettier . --write
