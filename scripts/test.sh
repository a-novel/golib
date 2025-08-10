#!/bin/bash

# shellcheck disable=SC2046
go tool gotestsum --format pkgname -- -count=1 -cover $(go list ./... | grep -v /mocks)
