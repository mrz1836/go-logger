# go-logger
> Easy to use, extendable and superfast logging package for Go

[![Release](https://img.shields.io/github/release-pre/mrz1836/go-logger.svg?logo=github&style=flat)](https://github.com/mrz1836/go-logger/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/mrz1836/go-logger/run-tests.yml?branch=master&logo=github&v=3)](https://github.com/mrz1836/go-logger/actions)
[![Report](https://goreportcard.com/badge/github.com/mrz1836/go-logger?style=flat)](https://goreportcard.com/report/github.com/mrz1836/go-logger)
[![codecov](https://codecov.io/gh/mrz1836/go-logger/branch/master/graph/badge.svg)](https://codecov.io/gh/mrz1836/go-logger)
[![Go](https://img.shields.io/github/go-mod/go-version/mrz1836/go-logger)](https://golang.org/)
[![Sponsor](https://img.shields.io/badge/sponsor-MrZ-181717.svg?logo=github&style=flat&v=3)](https://github.com/sponsors/mrz1836)
[![Donate](https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat)](https://mrz1818.com/?tab=tips&af=go-logger)

<br/>

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

<br/>

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

<br/>

## Documentation
View the generated [documentation](https://pkg.go.dev/github.com/mrz1836/go-logger)

[![GoDoc](https://godoc.org/github.com/mrz1836/go-logger?status.svg&style=flat)](https://pkg.go.dev/github.com/mrz1836/go-logger)

### Features
- Native logging package (extends log package)
- Native support for [Log Entries (Rapid7)](https://www.rapid7.com/products/insightops/) with queueing
- Test coverage on all custom methods
- Supports different Rapid7 endpoints & ports

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

[goreleaser](https://github.com/goreleaser/goreleaser) for easy binary or library deployment to Github and can be installed via: `brew install goreleaser`.

The [.goreleaser.yml](.goreleaser.yml) file is used to configure [goreleaser](https://github.com/goreleaser/goreleaser).

Use `make release-snap` to create a snapshot version of the release, and finally `make release` to ship to production.
</details>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>
<br/>

View all `makefile` commands
```shell script
make help
```

List of all current commands:
```text
all                  Runs multiple commands
clean                Remove previous builds and any test cache data
clean-mods           Remove all the Go mod cache
coverage             Shows the test coverage
godocs               Sync the latest tag with GoDocs
help                 Show this help message
install              Install the application
install-go           Install the application (Using Native Go)
lint                 Run the golangci-lint application (install if not found)
release              Full production release (creates release in Github)
release              Runs common.release then runs godocs
release-snap         Test the full release (build binaries)
release-test         Full production test release (everything except deploy)
replace-version      Replaces the version in HTML/JS (pre-deploy)
run-examples         Runs all the examples
tag                  Generate a new tag and push (tag version=0.0.0)
tag-remove           Remove a tag if found (tag-remove version=0.0.0)
tag-update           Update an existing tag to current commit (tag-update version=0.0.0)
test                 Runs vet, lint and ALL tests
test-ci              Runs all tests via CI (exports coverage)
test-ci-no-race      Runs all tests via CI (no race) (exports coverage)
test-ci-short        Runs unit tests via CI (exports coverage)
test-short           Runs vet, lint and tests (excludes integration tests)
uninstall            Uninstall the application (and remove files)
update-linter        Update the golangci-lint package (macOS only)
vet                  Run the Go vet application
```
</details>

<br/>

## Examples & Tests
All unit tests and [examples](examples) run via [Github Actions](https://github.com/mrz1836/go-logger/actions) and
uses [Go version 1.16.x](https://golang.org/doc/go1.16). View the [configuration file](.github/workflows/run-tests.yml).

Run all tests (including integration tests)
```shell script
make test
```

Run tests (excluding integration tests)
```shell script
make test-short
```

<br/>

## Benchmarks
Run the Go [benchmarks](logger_test.go):
```shell script
make bench
```

<br/>

## Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

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

<br/>

## Maintainers
| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) | [<img src="https://github.com/kayleg.png" height="50" alt="kayleg" />](https://github.com/kayleg) |
|:------------------------------------------------------------------------------------------------:|:-------------------------------------------------------------------------------------------------:|
|                                [MrZ](https://github.com/mrz1836)                                 |                                [kayleg](https://github.com/kayleg)                                |

<br/>

## Contributing
View the [contributing guidelines](.github/CONTRIBUTING.md) and please follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:! 
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:. 
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/mrz1836) :clap: 
or by making a [**bitcoin donation**](https://mrz1818.com/?tab=tips&af=go-logger) to ensure this journey continues indefinitely! :rocket:

<br/>

## License

[![License](https://img.shields.io/github/license/mrz1836/go-logger.svg?style=flat&p=1)](LICENSE)
