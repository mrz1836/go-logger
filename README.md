# go-logger
> Easy to use, extendable and super fast logging package for Go

[![Go](https://img.shields.io/github/go-mod/go-version/mrz1836/go-logger)](https://golang.org/)
[![Build Status](https://travis-ci.com/mrz1836/go-logger.svg?branch=master)](https://travis-ci.com/mrz1836/go-logger)
[![Report](https://goreportcard.com/badge/github.com/mrz1836/go-logger?style=flat)](https://goreportcard.com/report/github.com/mrz1836/go-logger)
[![codecov](https://codecov.io/gh/mrz1836/go-logger/branch/master/graph/badge.svg)](https://codecov.io/gh/mrz1836/go-logger)
[![Release](https://img.shields.io/github/release-pre/mrz1836/go-logger.svg?style=flat)](https://github.com/mrz1836/go-logger/releases)
[![GoDoc](https://godoc.org/github.com/mrz1836/go-logger?status.svg&style=flat)](https://pkg.go.dev/github.com/mrz1836/go-logger)

## Table of Contents
- [Installation](#installation)
- [Documentation](#documentation)
- [Examples & Tests](#examples--tests)
- [Benchmarks](#benchmarks)
- [Code Standards](#code-standards)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

## Installation

**go-logger** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/mrz1836/go-logger
```

For use with [Log Entries (Rapid7)](https://www.rapid7.com/products/insightops/), change the environment variables:
```shell script
export LOG_ENTRIES_TOKEN=your-token-here
```

_(Optional)_ Set custom endpoint or port parameters
```shell script
export LOG_ENTRIES_ENDPOINT=us.data.logs.insight.rapid7.com
export LOG_ENTRIES_PORT=514
``` 

## Documentation
You can view the generated [documentation here](https://pkg.go.dev/github.com/mrz1836/go-logger).

### Features
- Native logging package (extends log package)
- Native support for [Log Entries (Rapid7)](https://www.rapid7.com/products/insightops/) with queueing
- Test coverage on all custom methods
- Supports different Rapid7 endpoints & ports

<details>
<summary><strong><code>Library Deployment</code></strong></summary>

[goreleaser](https://github.com/goreleaser/goreleaser) for easy binary or library deployment to Github and can be installed via: `brew install goreleaser`.

The [.goreleaser.yml](.goreleaser.yml) file is used to configure [goreleaser](https://github.com/goreleaser/goreleaser).

Use `make release-snap` to create a snapshot version of the release, and finally `make release` to ship to production.
</details>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>

View all `makefile` commands
```shell script
make help
```

List of all current commands:
```text
all                            Runs lint, test-short and vet
bench                          Run all benchmarks in the Go application
clean                          Remove previous builds and any test cache data
clean-mods                     Remove all the Go mod cache
coverage                       Shows the test coverage
godocs                         Sync the latest tag with GoDocs
help                           Show all make commands available
lint                           Run the Go lint application
release                        Full production release (creates release in Github)
release-test                   Full production test release (everything except deploy)
release-snap                   Test the full release (build binaries)
run-examples                   Runs all the examples
tag                            Generate a new tag and push (IE: tag version=0.0.0)
tag-remove                     Remove a tag if found (IE: tag-remove version=0.0.0)
tag-update                     Update an existing tag to current commit (IE: tag-update version=0.0.0)
test                           Runs vet, lint and ALL tests
test-short                     Runs vet, lint and tests (excludes integration tests)
test-travis                    Runs tests via Travis (also exports coverage)
update                         Update all project dependencies
update-releaser                Update the goreleaser application
vet                            Run the Go vet application
```
</details>

## Examples & Tests
All unit tests and [examples](examples/examples.go) run via [Travis CI](https://travis-ci.org/mrz1836/go-logger) and uses [Go version 1.14.x](https://golang.org/doc/go1.14). View the [deployment configuration file](.travis.yml).

Run all tests (including integration tests)
```shell script
make test
```

Run tests (excluding integration tests)
```shell script
make test-short
```

## Benchmarks
Run the Go [benchmarks](logger_test.go):
```shell script
make bench
```

## Code Standards
Read more about this Go project's [code standards](CODE_STANDARDS.md).

## Usage
- View the [examples](examples/examples.go)
- View the [tests](logger_test.go)

Basic implementation:
```go
package main

import "github.com/mrz1836/go-logger"

func main() {
	logger.Data(2, logger.DEBUG, "testing the go-logger package")
	// Output: type="debug" file="example/example.go" method="main.main" line="12" message="testing the go-logger package"
}
```

## Maintainers

| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) | [<img src="https://github.com/kayleg.png" height="50" alt="kayleg" />](https://github.com/kayleg) |
|:---:|:---:|
| [MrZ](https://github.com/mrz1836) | [kayleg](https://github.com/kayleg) |

## Contributing

View the [contributing guidelines](CONTRIBUTING.md) and follow the [code of conduct](CODE_OF_CONDUCT.md).

Support the development of this project 🙏

[![Donate](https://img.shields.io/badge/donate-bitcoin-brightgreen.svg)](https://mrz1818.com/?tab=tips&af=go-logger)

## License

![License](https://img.shields.io/github/license/mrz1836/go-logger.svg?style=flat&p=1)
