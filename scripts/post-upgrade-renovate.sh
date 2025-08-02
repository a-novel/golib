#!/bin/bash

set -e

# Install node on the renovate image.
curl -o- https://fnm.vercel.app/install | bash
eval "$(fnm env)"
fnm install 24

node -v || echo "node install failed" && exit 1
npm -v || echo "npm install failed" && exit 1

npx -y prettier . --write
