# Run tests.
test:
	bash -c "set -m; bash '$(CURDIR)/scripts/test.sh'"

# Check code quality.
lint:
	go tool golangci-lint run
	go tool buf lint
	pnpm lint

# Generate Go code.
generate-go:
	go generate ./...

generate: generate-go

# Reformat code so it passes the code style lint checks.
format:
	go mod tidy
	go tool golangci-lint run --fix
	go tool buf format -w
	go tool buf dep update
	pnpm format
