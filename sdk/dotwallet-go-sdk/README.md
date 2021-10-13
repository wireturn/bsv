# DotWallet SDK
> A golang client for the [DotWallet](https://dotwallet.com) [API](https://developers.dotwallet.com/documents/en/#intro)

[![Release](https://img.shields.io/github/release-pre/dotwallet/dotwallet-go-sdk.svg?logo=github&style=flat&v=1)](https://github.com/dotwallet/dotwallet-go-sdk/releases)
[![Build Status](https://img.shields.io/github/workflow/status/dotwallet/dotwallet-go-sdk/run-go-tests?logo=github&v=1)](https://github.com/dotwallet/dotwallet-go-sdk/actions)
[![Report](https://goreportcard.com/badge/github.com/dotwallet/dotwallet-go-sdk?style=flat&v=1)](https://goreportcard.com/report/github.com/dotwallet/dotwallet-go-sdk)
[![codecov](https://codecov.io/gh/dotwallet/dotwallet-go-sdk/branch/master/graph/badge.svg?v=1)](https://codecov.io/gh/dotwallet/dotwallet-go-sdk)
[![Go](https://img.shields.io/github/go-mod/go-version/dotwallet/dotwallet-go-sdk?v=1)](https://golang.org/)

<br/>

## Table of Contents
- [Installation](#installation)
- [Documentation](#documentation)
- [Examples & Tests](#examples--tests)
- [Benchmarks](#benchmarks)
- [Code Standards](#code-standards)
- [Usage](#usage)
- [Maintainers & Contributors](#maintainers--contributors)
- [Contributing](#contributing)
- [License](#license)

<br/>

## Installation

**dotwallet-go-sdk** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/dotwallet/dotwallet-go-sdk
```

<br/>

## Documentation
View the generated [documentation](https://pkg.go.dev/github.com/dotwallet/dotwallet-go-sdk)

[![GoDoc](https://godoc.org/github.com/dotwallet/dotwallet-go-sdk?status.svg&style=flat)](https://pkg.go.dev/github.com/dotwallet/dotwallet-go-sdk)

### Features
- [Client](client.go) is completely configurable
- Using default [Resty http client](https://github.com/go-resty/resty) with exponential backoff & more
- Use your own custom HTTP resty client
- Autoload the `application_access_token` when required by `Authorization` in `client.Request()`
- Current coverage for the [DotWallet API](https://developers.dotwallet.com/documents/en/#intro)
  - [ ] [Authorization](https://developers.dotwallet.com/documents/en/#authorization)
    - [x] Application Authorization
    - [x] User Authorization
    - [ ] Signature Authorization
  - [ ] [User Info](https://developers.dotwallet.com/documents/en/#user-info)
    - [x] Get User Info
    - [x] Get User Receive Address
    - [ ] Get User Badge Balance
  - [ ] [Transactions](https://developers.dotwallet.com/documents/en/#transactions)
    - [ ] Single Payment Order
    - [ ] Query Order Status
  - [ ] [Automatic Payment](https://developers.dotwallet.com/documents/en/#automatic-payment)
    - [ ] Automatic Payment Flow
    - [ ] Insufficient Balance
    - [ ] Get Autopay Balance
  - [ ] [Badge](https://developers.dotwallet.com/documents/en/#badge)
    - [ ] Create Badge
    - [ ] Transfer Badge to User
    - [ ] Transfer Badge to Address
    - [ ] Query Badge Balance
    - [ ] Query Badge Transaction
  - [ ] [Merchant API](https://developers.dotwallet.com/documents/en/#merchant-api)
    - [ ] Query Miner Fees
    - [ ] Query Transaction Status
    - [ ] Send Transaction
  - [ ] [Blockchain Queries](https://developers.dotwallet.com/documents/en/#blockchain-queries)
    - [ ] Block Header
    - [ ] Block Info from Hash
    - [ ] Block Info from Height
    - [ ] Block Info Batch Query from Height
    - [ ] Block Transaction List
    - [ ] Blockchain Basic Info
    - [ ] Transaction Inquiry
    - [ ] Merkle Branch Inquiry
  - [ ] [DAPP API](https://developers.dotwallet.com/documents/en/#dapp-api)
    - [ ] List UTXOS
    - [ ] Sign Raw Transaction
    - [ ] Get Signature
    - [ ] Broadcast Transaction
    - [ ] Get Raw Change Address
    - [ ] Get Balance
    - [ ] Get Public Key

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
All unit tests and [examples](examples) run via [Github Actions](https://github.com/dotwallet/dotwallet-go-sdk/actions) and
uses [Go version(s) 1.13.x, 1.14.x and 1.15.x](https://golang.org/doc/go1.15). View the [configuration file](.github/workflows/run-tests.yml).

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
Read more about this Go project's [code standards](CODE_STANDARDS.md).

<br/>

## Usage
Checkout all the [examples](examples)!

<br/>

## Maintainers & Contributors
| <img src="https://i.imgur.com/sAc5hoe.png" height="50" alt="吴浩瑜" /> | <img src="https://i.imgur.com/sAc5hoe.png" height="50" alt="chenhao" /> | [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
|:---:|:---:|:---:|
| 吴浩瑜 | chenhao | [MrZ](https://github.com/mrz1836) |

<br/>

## Contributing

View the [contributing guidelines](CONTRIBUTING.md) and follow the [code of conduct](CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:!
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:.
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/dotwallet) :clap:!
<br/>

## License

![License](https://img.shields.io/github/license/dotwallet/dotwallet-go-sdk.svg?style=flat&v=1)