#!/bin/bash

TEST_TOOL_PKG="gotest.tools/gotestsum@latest"

# shellcheck disable=SC2046
go run ${TEST_TOOL_PKG} --format pkgname -- -count=1 -cover $(go list ./... | grep -v /mocks)
