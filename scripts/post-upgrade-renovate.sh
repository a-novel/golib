#!/bin/bash

set -e

# Install node on the renovate image.
if ! which node; then
   curl -o- https://fnm.vercel.app/install | bash
   fnm install 24

   node -v || echo "node install failed" && exit 1
   npm -v || echo "npm install failed" && exit 1
fi

npx -y prettier . --write
