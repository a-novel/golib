# Go lib

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/a-novel/golib/main.yaml)
[![codecov](https://codecov.io/gh/a-novel/golib/graph/badge.svg?token=LQMRBETC8K)](https://codecov.io/gh/a-novel/golib)

![GitHub repo file or directory count](https://img.shields.io/github/directory-file-count/a-novel/golib)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/a-novel/golib)

![Coverage graph](https://codecov.io/gh/a-novel/golib/graphs/sunburst.svg?token=LQMRBETC8K)

Shared libraries for Go projects.

## Use in a project

```bash
go get github.com/a-novel/golib
```

## Run the project locally

### Prerequisites

- [Go](https://go.dev/doc/install)
- [Mockery](https://vektra.github.io/mockery/latest/installation/)
- Make
    - macOS:
      ```bash
      brew install make
      ```
    - Ubuntu:
      ```bash
      sudo apt-get install make
      ```
    - Windows: Install [chocolatey](https://chocolatey.org/install) (from a PowerShell with admin privileges), then run:
      ```bash
      choco install make
      ```

Install the project dependencies.

```bash
go get ./... && go mod tidy
```

## Work on the project

Make sure the project files are properly formatted.

```bash
make format
```

Run tests.

```bash
make test
```

Make sure your code is compliant with the linter.

```bash
make lint
```

If you create / update interfaces signatures, make sure to update the mocks.

```bash
mockery
```
