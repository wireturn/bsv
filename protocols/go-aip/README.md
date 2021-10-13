# go-aip
> Library for working with [Author Identity Protocol](https://github.com/BitcoinFiles/AUTHOR_IDENTITY_PROTOCOL) in Go

[![Release](https://img.shields.io/github/release-pre/BitcoinSchema/go-aip.svg?logo=github&style=flat&v=3)](https://github.com/BitcoinSchema/go-aip/releases)
[![Build Status](https://img.shields.io/github/workflow/status/BitcoinSchema/go-aip/run-go-tests?logo=github&v=3)](https://github.com/BitcoinSchema/go-aip/actions)
[![Report](https://goreportcard.com/badge/github.com/BitcoinSchema/go-aip?style=flat&v=3)](https://goreportcard.com/report/github.com/BitcoinSchema/go-aip)
[![codecov](https://codecov.io/gh/BitcoinSchema/go-aip/branch/master/graph/badge.svg?v=3)](https://codecov.io/gh/BitcoinSchema/go-aip)
[![Go](https://img.shields.io/github/go-mod/go-version/BitcoinSchema/go-aip?v=3)](https://golang.org/)
<br>

[![Mergify Status](https://img.shields.io/endpoint.svg?url=https://gh.mergify.io/badges/BitcoinSchema/go-aip&style=flat&v=3)](https://mergify.io)
[![Sponsor](https://img.shields.io/badge/sponsor-BitcoinSchema-181717.svg?logo=github&style=flat&v=3)](https://github.com/sponsors/BitcoinSchema)
[![Donate](https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat&v=3)](https://gobitcoinsv.com/#sponsor?utm_source=github&utm_medium=sponsor-link&utm_campaign=go-aip&utm_term=go-aip&utm_content=go-aip)

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

**go-aip** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/bitcoinschema/go-aip
```

<br/>

## Documentation
View the generated [documentation](https://pkg.go.dev/github.com/bitcoinschema/go-aip)

[![GoDoc](https://godoc.org/github.com/bitcoinschema/go-aip?status.svg&style=flat)](https://pkg.go.dev/github.com/bitcoinschema/go-aip)

### Features
- [Sign Message (using ECDSA, BitcoinSignedMessage)](aip.go)
- [Sign OpReturn](aip.go)
- [Validate Signatures (ECDSA & Paymail)](aip.go)
- [Parse from BOB](bob.go)
- [Validate BOB Tape](bob.go)

<details>
<summary><strong><code>Package Dependencies</code></strong></summary>
<br/>

- [bitcoinschema/go-bitcoin](https://github.com/bitcoinschema/go-bitcoin)
- [bitcoinschema/go-bob](https://github.com/bitcoinschema/go-bob)
- [libsv/go-bt](https://github.com/libsv/go-bt)
</details>

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
All unit tests and [examples](examples) run via [Github Actions](https://github.com/BitcoinSchema/go-aip/actions) and
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
Run the Go benchmarks:
```shell script
make bench
```

<br/>

## Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## Usage
Checkout all the [examples](examples)!

<br/>

## Maintainers
| [<img src="https://github.com/rohenaz.png" height="50" alt="MrZ" />](https://github.com/rohenaz) | [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
|:---:|:---:|
| [Satchmo](https://github.com/rohenaz) | [MrZ](https://github.com/mrz1836) |

<br/>

## Contributing

View the [contributing guidelines](.github/CONTRIBUTING.md) and follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:!
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:.
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/BitcoinSchema) :clap:
or by making a [**bitcoin donation**](https://gobitcoinsv.com/#sponsor?utm_source=github&utm_medium=sponsor-link&utm_campaign=go-aip&utm_term=go-aip&utm_content=go-aip) to ensure this journey continues indefinitely! :rocket:

<br/>

## License

![License](https://img.shields.io/github/license/BitcoinSchema/go-aip.svg?style=flat&v=3)