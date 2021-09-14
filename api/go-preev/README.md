# go-preev
> The unofficial golang implementation for the [Preev.pro API](https://preev.pro/api/)

[![Release](https://img.shields.io/github/release-pre/mrz1836/go-preev.svg?logo=github&style=flat&v=4)](https://github.com/mrz1836/go-preev/releases)
[![Build Status](https://img.shields.io/github/workflow/status/mrz1836/go-preev/run-go-tests?logo=github&v=3)](https://github.com/mrz1836/go-preev/actions)
[![Report](https://goreportcard.com/badge/github.com/mrz1836/go-preev?style=flat&v=4)](https://goreportcard.com/report/github.com/mrz1836/go-preev)
[![codecov](https://codecov.io/gh/mrz1836/go-preev/branch/master/graph/badge.svg?v=4)](https://codecov.io/gh/mrz1836/go-preev)
[![Go](https://img.shields.io/github/go-mod/go-version/mrz1836/go-preev?v=4)](https://golang.org/)
[![Sponsor](https://img.shields.io/badge/sponsor-MrZ-181717.svg?logo=github&style=flat&v=4)](https://github.com/sponsors/mrz1836)
[![Donate](https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat)](https://mrz1818.com/?tab=tips&utm_source=github&utm_medium=sponsor-link&utm_campaign=go-preev&utm_term=go-preev&utm_content=go-preev)

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

**go-preev** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/mrz1836/go-preev
```

<br/>

## Documentation
View the generated [documentation](https://pkg.go.dev/github.com/mrz1836/go-preev)

[![GoDoc](https://godoc.org/github.com/mrz1836/go-preev?status.svg&style=flat)](https://pkg.go.dev/github.com/mrz1836/go-preev)

### Features
- [Client](client.go) is completely configurable
- Customize the network per request (`main`, `test` or `stn`)
- Using default [heimdall http client](https://github.com/gojek/heimdall) with exponential backoff & more
- Use your own HTTP client
- Current coverage for the [Preev.pro API](https://preev.pro/api/)
    - [x] Pairs
        - [x] Get All Pairs
        - [x] Get Pair By ID
    - [x] Tickers
        - [x] Get All Tickers
        - [x] Get Ticker By Pair ID
        - [x] Get Ticker By Pair ID (Historical)

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
all                  Runs lint, test-short and vet
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
All unit tests and [examples](examples) run via [Github Actions](https://github.com/mrz1836/go-preev/actions) and
uses [Go version 1.15.x](https://golang.org/doc/go1.15). View the [configuration file](.github/workflows/run-tests.yml).

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
Run the Go [benchmarks](preev_test.go):
```shell script
make bench
```

<br/>

## Code Standards
Read more about this Go project's [code standards](CODE_STANDARDS.md).

<br/>

## Usage
View the [preev examples](examples/examples.go)

Basic implementation:
```go
package main

import (
    "fmt"
    
    "github.com/mrz1836/go-preev"
)

func main() {

    // Create a client
    client := preev.NewClient(nil, nil) 
    
    // Get pairs
    pairs, _ := client.GetPairs()
	
	fmt.Println("Found Active Pair(s):", pairs.BsvUsd.Name)
}
```
 
<br/>

## Maintainers
| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
|:---:|
| [MrZ](https://github.com/mrz1836) |
              
<br/>

## Contributing
View the [contributing guidelines](CONTRIBUTING.md) and please follow the [code of conduct](CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:! 
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:. 
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/mrz1836) :clap: 
or by making a [**bitcoin donation**](https://mrz1818.com/?tab=tips&utm_source=github&utm_medium=sponsor-link&utm_campaign=go-preev&utm_term=go-preev&utm_content=go-preev) to ensure this journey continues indefinitely! :rocket:


### Credits

[Preev](https://tncpw.co/d19f43a3) for their hard work on the Preev API

<br/>

## License

![License](https://img.shields.io/github/license/mrz1836/go-preev.svg?style=flat&v=4)