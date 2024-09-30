#!/bin/bash

COMPOSE_FILE="docker-compose.test.yaml"
TEST_TOOL_PKG="gotest.tools/gotestsum@latest"

# Ensure containers are properly shut down when the program exits abnormally.
int_handler()
{
    docker compose -f ${COMPOSE_FILE} down
}
trap int_handler INT

# Setup test containers.
docker compose -f ${COMPOSE_FILE} up -d

# Execute tests.
export CGO_ENABLED=1
go run ${TEST_TOOL_PKG} --format pkgname -- -count=1 -coverprofile=cover.out -p 1 $(go list ./... | grep -v /mocks)
go tool cover -html=cover.out -o cover.html

# Normal execution: containers are shut down.
docker compose -f ${COMPOSE_FILE} down
