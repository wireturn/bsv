# go-minercraft
> Interact with Bitcoin Miner APIs (unofficial Go library of [Minercraft](https://github.com/interplanaria/minercraft))

[![Release](https://img.shields.io/github/release-pre/tonicpow/go-minercraft.svg?logo=github&style=flat&v=3)](https://github.com/tonicpow/go-minercraft/releases)
[![Build Status](https://img.shields.io/github/workflow/status/tonicpow/go-minercraft/run-go-tests?logo=github&v=3)](https://github.com/tonicpow/go-minercraft/actions)
[![Report](https://goreportcard.com/badge/github.com/tonicpow/go-minercraft?style=flat&v=3)](https://goreportcard.com/report/github.com/tonicpow/go-minercraft)
[![codecov](https://codecov.io/gh/tonicpow/go-minercraft/branch/master/graph/badge.svg?v=3)](https://codecov.io/gh/tonicpow/go-minercraft)
[![Go](https://img.shields.io/github/go-mod/go-version/tonicpow/go-minercraft?v=3)](https://golang.org/)
[![Mergify Status](https://img.shields.io/endpoint.svg?url=https://gh.mergify.io/badges/tonicpow/go-minercraft&style=flat&v=3)](https://mergify.io)  

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

**go-minercraft** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/tonicpow/go-minercraft
```

<br/>

## Documentation
View the generated [documentation](https://pkg.go.dev/github.com/tonicpow/go-minercraft)

[![GoDoc](https://godoc.org/github.com/tonicpow/go-minercraft?status.svg&style=flat&v=3)](https://pkg.go.dev/github.com/tonicpow/go-minercraft)
          
This package interacts with BSV miners using the [Merchant API](https://github.com/bitcoin-sv-specs/brfc-merchantapi) specification.

View documentation on hosting your own [mAPI server](https://github.com/bitcoin-sv/merchantapi-reference).

### Features
- Merchant API Support:
  - [x] [Fee Quote](https://github.com/bitcoin-sv-specs/brfc-merchantapi#get-fee-quote)
  - [x] [Query Transaction Status](https://github.com/bitcoin-sv-specs/brfc-merchantapi#Query-transaction-status)
  - [x] [Submit Transaction](https://github.com/bitcoin-sv-specs/brfc-merchantapi#Submit-transaction)
  - [ ] [Submit Multiple Transactions](https://github.com/bitcoin-sv-specs/brfc-merchantapi#Submit-multiple-transactions) `(Miners have not implemented as of 10/15/20)`
- Custom Features:
  - [Client](client.go) is completely configurable
  - Using default [heimdall http client](https://github.com/gojektech/heimdall) with exponential backoff & more
  - Use your own [HTTP client](client.go)
  - Use your own [miner configuration](client.go)
  - Uses common type: [`bt.Fee`](https://github.com/libsv/go-bt/blob/master/fees.go) for easy integration across projects 
  - Current miner information located at `response.Miner.name` and [defaults](config.go)
  - Automatic Signature Validation `response.Validated=true/false`
  - `AddMiner()` for adding your own customer miner configuration
  - `RemoveMiner()` for removing any miner configuration
  - `FastestQuote()` asks all miners and returns the fastest quote response
  - `BestQuote()` gets all quotes from miners and return the best rate/quote
  - `CalculateFee()` returns the fee for a given transaction

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
generate             Runs the go generate command in the base of the repo
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
test-unit            Runs tests and outputs coverage
uninstall            Uninstall the application (and remove files)
update-linter        Update the golangci-lint package (macOS only)
vet                  Run the Go vet application
```
</details>

<br/>

## Examples & Tests
All unit tests and [examples](examples) run via [Github Actions](https://github.com/tonicpow/go-minercraft/actions) and
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
Run the Go [benchmarks](client.go):
```shell script
make bench
```

<br/>

## Code Standards
Read more about this Go project's [code standards](CODE_STANDARDS.md).

<br/>

## Usage
View the [examples](examples)

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
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/tonicpow) :clap:
or by making a [**bitcoin donation**](https://tonicpow.com/?utm_source=github&utm_medium=sponsor-link&utm_campaign=go-minercraft&utm_term=go-minercraft&utm_content=go-minercraft) to ensure this journey continues indefinitely! :rocket:

[![Stars](https://img.shields.io/github/stars/tonicpow/go-minercraft?label=Please%20like%20us&style=social)](https://github.com/tonicpow/go-minercraft/stargazers)

### Credits

[Unwriter & Interplanaria](https://github.com/interplanaria) for their original contribution: [Minercraft](https://github.com/interplanaria/minercraft) which was the inspiration for this library.
      
nChain & team for developing the [brfc-merchant-api](https://github.com/bitcoin-sv-specs/brfc-merchantapi) specifications.

<br/>

## License

![License](https://img.shields.io/github/license/tonicpow/go-minercraft.svg?style=flat&v=3)