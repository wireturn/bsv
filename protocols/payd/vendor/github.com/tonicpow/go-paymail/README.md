# go-paymail
> Paymail client & server library for Golang

[![Release](https://img.shields.io/github/release-pre/tonicpow/go-paymail.svg?logo=github&style=flat&v=3)](https://github.com/tonicpow/go-paymail/releases)
[![Build Status](https://img.shields.io/github/workflow/status/tonicpow/go-paymail/run-go-tests?logo=github&v=3)](https://github.com/tonicpow/go-paymail/actions)
[![Report](https://goreportcard.com/badge/github.com/tonicpow/go-paymail?style=flat&v=3)](https://goreportcard.com/report/github.com/tonicpow/go-paymail)
[![codecov](https://codecov.io/gh/tonicpow/go-paymail/branch/master/graph/badge.svg?v=3)](https://codecov.io/gh/tonicpow/go-paymail)
[![Go](https://img.shields.io/github/go-mod/go-version/tonicpow/go-paymail?v=3)](https://golang.org/)

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

**go-paymail** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/tonicpow/go-paymail
```

<br/>

## Documentation
View the generated [documentation](https://pkg.go.dev/github.com/tonicpow/go-paymail)

[![GoDoc](https://godoc.org/github.com/tonicpow/go-paymail?status.svg&style=flat)](https://pkg.go.dev/github.com/tonicpow/go-paymail)

### Features
- [Paymail Client](client.go) (outgoing requests to other providers)
    - Use a custom [Resty HTTP client](https://github.com/go-resty/resty)
    - Use custom [client options](client.go)
    - Use a custom [net.Resolver](resolver_test.go)
    - [Get & Validate SRV records](srv.go)
    - [Check SSL Certificates](ssl.go)
    - [Check & Validate DNNSEC](dns_sec.go)
    - [Generate, Validate & Load Additional BRFC Specifications](brfc.go)
    - [Fetch, Get and Has Capabilities](capabilities.go)
    - [Get Public Key Information - PKI](pki.go)
    - [Basic Address Resolution](resolve_address.go)
    - [Verify PubKey & Handle](verify_pubkey.go)
    - [Get Public Profile](public_profile.go)
    - [P2P Payment Destination](p2p_payment_destination.go)
    - [P2P Send Transaction](p2p_send_transaction.go)
- [Paymail Server](server) (basic example for hosting your own paymail server)
    - [Example Showing Capabilities](server/capabilities.go) 
    - [Example Showing PKI](server/pki.go)
    - [Example Verifying a PubKey](server/verify.go)
    - [Example Address Resolution](server/resolve_address.go)
    - [Example Getting a P2P Payment Destination](server/p2p_payment_destination.go)
    - [Example Receiving a P2P Transaction](server/p2p_receive_transaction.go)
- [Paymail Utilities](utilities.go) (handy methods)
    - [Sanitize & Validate Paymail Addresses](utilities.go)
    - [Sign & Verify Sender Request](sender_request.go)
    
<details>
<summary><strong><code>Package Dependencies</code></strong></summary>
<br/>

Client Packages:
- [BitcoinSchema/go-bitcoin](https://github.com/BitcoinSchema/go-bitcoin)
- [bitcoinsv/bsvd](https://github.com/bitcoinsv/bsvd)
- [bitcoinsv/bsvutil](https://github.com/bitcoinsv/bsvutil)
- [go-resty/resty](https://github.com/go-resty/resty/v2)
- [jarcoal/httpmock](https://github.com/jarcoal/httpmock)
- [miekg/dns](https://github.com/miekg/dns)
- [mrz1836/go-sanitize](https://github.com/mrz1836/go-sanitize)
- [mrz1836/go-validate](https://github.com/mrz1836/go-validate)

Server Packages:
- [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)
- [mrz1836/go-api-router](https://github.com/mrz1836/go-api-router)
- [mrz1836/go-logger](https://github.com/mrz1836/go-logger)
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
All unit tests and [examples](examples) run via [Github Actions](https://github.com/tonicpow/go-paymail/actions) and
uses [Go version(s) 1.14.x and 1.15.x](https://golang.org/doc/go1.15). View the [configuration file](.github/workflows/run-tests.yml).

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
Checkout all the [client examples](examples/client) or [server examples](examples/server)!

<br/>

## Maintainers
| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
|:---:|
| [MrZ](https://github.com/mrz1836) |

<br/>

## Contributing

View the [contributing guidelines](CONTRIBUTING.md) and follow the [code of conduct](CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:! 
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:. 
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/tonicpow) :clap: 
or by making a [**bitcoin donation**](https://tonicpow.com/?utm_source=github&utm_medium=sponsor-link&utm_campaign=go-paymail&utm_term=go-paymail&utm_content=go-paymail) to ensure this journey continues indefinitely! :rocket:

<br/>

## License

![License](https://img.shields.io/github/license/tonicpow/go-paymail.svg?style=flat&v=3)
