# go-logger
**go-logger** is an easy to use, extendable and super fast logging package for Go

[![Go](https://img.shields.io/github/go-mod/go-version/mrz1836/go-logger)](https://golang.org/)
[![Build Status](https://travis-ci.org/mrz1836/go-logger.svg?branch=master)](https://travis-ci.org/mrz1836/go-logger)
[![Report](https://goreportcard.com/badge/github.com/mrz1836/go-logger?style=flat)](https://goreportcard.com/report/github.com/mrz1836/go-logger)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/de9d8cd1e21445e9823b005e4f7dcf20)](https://www.codacy.com/app/mrz1818/go-logger?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mrz1836/go-logger&amp;utm_campaign=Badge_Grade)
[![Release](https://img.shields.io/github/release-pre/mrz1836/go-logger.svg?style=flat)](https://github.com/mrz1836/go-logger/releases)
[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat)](https://github.com/RichardLitt/standard-readme)
[![GoDoc](https://godoc.org/github.com/mrz1836/go-logger?status.svg&style=flat)](https://godoc.org/github.com/mrz1836/go-logger)

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
```bash
$ go get -u github.com/mrz1836/go-logger
```

For use with [Log Entries](https://logentries.com/), change the environment variables:
```bash
export LOG_ENTRIES_TOKEN=your-token-here
```

_(Optional)_ Set custom endpoint or port parameters
```bash
export LOG_ENTRIES_ENDPOINT=us.data.logs.insight.rapid7.com
export LOG_ENTRIES_PORT=514
```

## Documentation
You can view the generated [documentation here](https://godoc.org/github.com/mrz1836/go-logger).

### Features
- Native logging package (extends log package)
- Native support for [Log Entries](https://logentries.com/) with queueing
- Test coverage on all custom methods
- Supports different Rapid7 endpoints & ports

## Examples & Tests
All unit tests and [examples](example/example.go) run via [Travis CI](https://travis-ci.org/mrz1836/go-logger) and uses [Go version 1.14.x](https://golang.org/doc/go1.14). View the [deployment configuration file](.travis.yml).

- [examples](example/example.go)
- [tests](logger_test.go)

Run all tests (including integration tests)
```bash
$ cd ../go-logger
$ go test ./... -v
```

Run tests (excluding integration tests)
```bash
$ cd ../go-logger
$ go test ./... -v -test.short
```

## Benchmarks
Run the Go [benchmarks](logger_test.go):
```bash
$ cd ../go-logger
$ go test -bench . -benchmem
```

## Code Standards
Read more about this Go project's [code standards](CODE_STANDARDS.md).

## Usage
- View the [examples](example/example.go)
- View the [tests](logger_test.go)

Basic implementation:
```golang
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
