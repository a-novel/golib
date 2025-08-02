#!/bin/bash

set -e


curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | bash
\. "$HOME/.nvm/nvm.sh"
#nvm install node

ls -al ./scripts && exit 1

npx -y prettier . --write
