# go-logger
**go-logger** is an easy to use, super fast and extendable logging package for Go

| | | | | | | |
|-|-|-|-|-|-|-|
| ![License](https://img.shields.io/github/license/mrz1836/go-logger.svg?style=flat&p=1) | [![Report](https://goreportcard.com/badge/github.com/mrz1836/go-logger?style=flat&p=1)](https://goreportcard.com/report/github.com/mrz1836/go-logger)  | [![Codacy Badge](https://api.codacy.com/project/badge/Grade/01708ca3079e4933bafb3b39fe2aaa9d)](https://www.codacy.com/app/mrz1818/go-logger?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mrz1836/go-logger&amp;utm_campaign=Badge_Grade) |  [![Build Status](https://travis-ci.com/mrz1836/go-logger.svg?branch=master)](https://travis-ci.com/mrz1836/go-logger)   |  [![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat)](https://github.com/RichardLitt/standard-readme) | [![Release](https://img.shields.io/github/release-pre/mrz1836/go-logger.svg?style=flat)](https://github.com/mrz1836/go-logger/releases) | [![GoDoc](https://godoc.org/github.com/mrz1836/go-logger?status.svg&style=flat)](https://godoc.org/github.com/mrz1836/go-logger) |

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

**go-logger** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy) and [dep](https://github.com/golang/dep).
```bash
$ go get -u github.com/mrz1836/go-logger
```

Updating dependencies in **go-logger**:
```bash
$ cd ../go-logger
$ dep ensure -update -v
```

## Documentation
You can view the generated [documentation here](https://godoc.org/github.com/mrz1836/go-logger).

### Features
- todo: @mrz

## Examples & Tests
All unit tests and [examples](logger_test.go) run via [Travis CI](https://travis-ci.com/mrz1836/go-logger) and uses [Go version 1.12.x](https://golang.org/doc/go1.12). View the [deployment configuration file](.travis.yml).

- [examples & tests](logger_test.go)

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
- View the [examples & benchmarks](logger_test.go)

Basic implementation:
- todo: @mrz

## Maintainers

[@MrZ1836](https://github.com/mrz1836) | [@kayleg](https://github.com/kayleg)

## Contributing

View the [contributing guidelines](CONTRIBUTING.md) and follow the [code of conduct](CODE_OF_CONDUCT.md).

Support the development of this project 🙏

[![Donate](https://img.shields.io/badge/donate-bitcoin-brightgreen.svg)](https://mrz1818.com/?tab=tips&af=go-logger)

## License

![License](https://img.shields.io/github/license/mrz1836/go-logger.svg?style=flat&p=1)
