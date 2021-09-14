# go-whatsonchain
> The unofficial golang implementation for the [whatsonchain.com API](https://developers.whatsonchain.com/)

[![Release](https://img.shields.io/github/release-pre/mrz1836/go-whatsonchain.svg?logo=github&style=flat&v=1)](https://github.com/mrz1836/go-whatsonchain/releases)
[![Build Status](https://img.shields.io/github/workflow/status/mrz1836/go-whatsonchain/run-go-tests?logo=github&v=3)](https://github.com/mrz1836/go-whatsonchain/actions)
[![Report](https://goreportcard.com/badge/github.com/mrz1836/go-whatsonchain?style=flat&v=2)](https://goreportcard.com/report/github.com/mrz1836/go-whatsonchain)
[![codecov](https://codecov.io/gh/mrz1836/go-whatsonchain/branch/master/graph/badge.svg)](https://codecov.io/gh/mrz1836/go-whatsonchain)
[![Go](https://img.shields.io/github/go-mod/go-version/mrz1836/go-whatsonchain)](https://golang.org/)
[![Sponsor](https://img.shields.io/badge/sponsor-MrZ-181717.svg?logo=github&style=flat&v=3)](https://github.com/sponsors/mrz1836)
[![Donate](https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat)](https://mrz1818.com/?tab=tips&utm_source=github&utm_medium=sponsor-link&utm_campaign=go-whatsonchain&utm_term=go-whatsonchain&utm_content=go-whatsonchain)

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

**go-whatsonchain** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/mrz1836/go-whatsonchain
```

<br/>

## Documentation
View the generated [documentation](https://pkg.go.dev/github.com/mrz1836/go-whatsonchain)

[![GoDoc](https://godoc.org/github.com/mrz1836/go-whatsonchain?status.svg&style=flat)](https://pkg.go.dev/github.com/mrz1836/go-whatsonchain)

### Features
- [Client](client.go) is completely configurable
- Customize the network per request (`main`, `test` or `stn`)
- Using default [heimdall http client](https://github.com/gojek/heimdall) with exponential backoff & more
- Use your own custom HTTP client
- Current coverage for the [whatsonchain.com API](https://developers.whatsonchain.com/)
    - [x] Health
        - [x] Get API Status
    - [x] Chain Info
        - [x] Get Blockchain Info
        - [x] Get Circulating Supply
    - [x] Block
        - [x] Get by Hash
        - [x] Get by Height
        - [x] Get Block Pages
        - [x] Get Header by Hash
        - [x] Get Headers
    - [x] Transaction
        - [x] Get by TX Hash
        - [x] Broadcast Transaction
        - [x] Bulk Broadcast
        - [x] Bulk Transaction Details
        - [x] Decode Transaction
        - [x] Download Receipt
        - [x] Get Raw Transaction Data
        - [x] Get Raw Transaction Output
        - [x] Get Merkle Proof
        - [x] Get Bulk Transaction Details (Custom)
    - [x] Mempool
        - [x] Get Mempool Info
        - [x] Get Mempool Transactions
    - [x] Address
        - [x] Get Address Info
        - [x] Get Balance
        - [x] Get History
        - [x] Get Unspent Transactions
        - [x] Get Unspent Transaction Details (Custom)
        - [x] Bulk Balance
        - [x] Bulk UTXOs
        - [x] Download Statement
    - [x] Script
        - [x] Get Script History
        - [x] Get Script Unspent Transactions
        - [x] Bulk UTXOs
    - [x] Exchange Rate
        - [x] Get Exchange Rate
    - [x] Search
        - [x] Get Explorer Links
    - [x] Merchant API
        - [x] Fee Quotes
        - [x] Submit Transaction
        - [x] Transaction Status

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
All unit tests and [examples](examples) run via [Github Actions](https://github.com/mrz1836/go-whatsonchain/actions) and
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
Run the Go [benchmarks](whatsonchain_test.go):
```shell script
make bench
```

<br/>

## Code Standards
Read more about this Go project's [code standards](CODE_STANDARDS.md).

<br/>

## Usage
View the [whatsonchain examples](examples)

Basic implementation:
```go
package main

import (
    "fmt"
    
    "github.com/mrz1836/go-whatsonchain"
)

func main() {

    // Create a client
    client := whatsonchain.NewClient(whatsonchain.NetworkMain, nil, nil)

    // Get a balance for an address
    balance, _ := client.AddressBalance("16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
    fmt.Println("confirmed balance", balance.Confirmed)
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
or by making a [**bitcoin donation**](https://mrz1818.com/?tab=tips&utm_source=github&utm_medium=sponsor-link&utm_campaign=go-whatsonchain&utm_term=go-whatsonchain&utm_content=go-whatsonchain) to ensure this journey continues indefinitely! :rocket:


### Credits

[WhatsOnChain](https://tncpw.co/65733e42) for their hard work on the Whatsonchain API

[AustEcon's Python Version](https://github.com/AustEcon/whatsonchain)

<br/>

## License

![License](https://img.shields.io/github/license/mrz1836/go-whatsonchain.svg?style=flat)