# go-bsvrates
> Real-time exchange rates for BitcoinSV using multiple providers with failover support

[![Release](https://img.shields.io/github/release-pre/tonicpow/go-bsvrates.svg?logo=github&style=flat&v=3)](https://github.com/tonicpow/go-bsvrates/releases)
[![Build Status](https://img.shields.io/github/workflow/status/tonicpow/go-bsvrates/run-go-tests?logo=github&v=3)](https://github.com/tonicpow/go-bsvrates/actions)
[![Report](https://goreportcard.com/badge/github.com/tonicpow/go-bsvrates?style=flat&v=3)](https://goreportcard.com/report/github.com/tonicpow/go-bsvrates)
[![codecov](https://codecov.io/gh/tonicpow/go-bsvrates/branch/master/graph/badge.svg?v=3)](https://codecov.io/gh/tonicpow/go-bsvrates)
[![Go](https://img.shields.io/github/go-mod/go-version/tonicpow/go-bsvrates?v=3)](https://golang.org/)
[![Mergify Status](https://img.shields.io/endpoint.svg?url=https://gh.mergify.io/badges/tonicpow/go-bsvrates&style=flat&)](https://mergify.io)

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

**go-bsvrates** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/tonicpow/go-bsvrates
```

<br/>

## Documentation
View the generated [documentation](https://pkg.go.dev/github.com/tonicpow/go-bsvrates)

[![GoDoc](https://godoc.org/github.com/tonicpow/go-bsvrates?status.svg&style=flat)](https://pkg.go.dev/github.com/tonicpow/go-bsvrates)

### Features
- [BSV Rates Client](client.go) is completely configurable
- Using default [heimdall http client](https://github.com/gojek/heimdall) with exponential backoff & more
- Use your own [HTTP client](client.go)
- Helpful currency conversion and formatting methods:
    - [ConvertFloatToIntBSV()](currency.go)
    - [ConvertIntToFloatUSD()](currency.go)
    - [ConvertPriceToSatoshis()](currency.go)
    - [ConvertSatsToBSV()](currency.go)
    - [FormatCentsToDollars()](currency.go)
    - [GetCentsFromSatoshis()](currency.go)
    - [GetDollarsFromSatoshis()](currency.go)
    - [TransformCurrencyToInt()](currency.go)
    - [TransformIntToCurrency()](currency.go)
- Supported Fiat Currencies:
    - USD
- Supported Providers:
    - **[Coin Paprika](https://api.coinpaprika.com/)**
      - [GetBaseAmountAndCurrencyID()](coinpaprika.go)
      - [GetMarketPrice()](coinpaprika.go)
      - [GetPriceConversion()](coinpaprika.go)
      - [IsAcceptedCurrency()](coinpaprika.go)
      - [GetHistoricalTickers()](coinpaprika.go)
    - **[What's On Chain](https://developers.whatsonchain.com/)**
      - [GetExchangeRate()](https://github.com/mrz1836/go-whatsonchain)
    - **[Preev](https://preev.pro/api/)**
      - [GetPair()](https://github.com/mrz1836/go-preev)
      - [GetPairs()](https://github.com/mrz1836/go-preev)
      - [GetTicker()](https://github.com/mrz1836/go-preev)
      - [GetTickers()](https://github.com/mrz1836/go-preev)

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
all                  Runs lint, test and vet
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
run-examples         Runs the basic example
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
All unit tests and [examples](examples) run via [Github Actions](https://github.com/tonicpow/go-bsvrates/actions) and 
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

Basic exchange rate implementation:

```go
package main

import (
  "context"
  "log"

  "github.com/tonicpow/go-bsvrates"
)

func main() {

  // Create a new client (all default providers)
  client := bsvrates.NewClient(nil, nil)

  // Get rates
  rate, provider, _ := client.GetRate(context.Background(), bsvrates.CurrencyDollars)
  log.Printf("found rate: %v %s from provider: %s", rate, bsvrates.CurrencyToName(bsvrates.CurrencyDollars), provider.Name())
}
``` 

Basic price conversion implementation:
```go
package main

import (
    "context"
    "log"

	"github.com/tonicpow/go-bsvrates"
)

func main() {

	// Create a new client (all default providers)
	client := bsvrates.NewClient(nil, nil)
    
	// Get a conversion from $ to Sats
	satoshis, provider, _ := client.GetConversion(context.Background(),bsvrates.CurrencyDollars, 0.01)
	log.Printf("0.01 USD = satoshis: %d from provider: %s", satoshis, provider.Name())
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
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/tonicpow) :clap: 
or by making a [**bitcoin donation**](https://tonicpow.com/?utm_source=github&utm_medium=sponsor-link&utm_campaign=go-bsvrates&utm_term=go-bsvrates&utm_content=go-bsvrates) to ensure this journey continues indefinitely! :rocket:


[![Stars](https://img.shields.io/github/stars/tonicpow/go-bsvrates?label=Please%20like%20us&style=social)](https://github.com/tonicpow/go-bsvrates/stargazers)

<br/>

### Credits

[Coin Paprika](https://tncpw.co/7c2cae76), [What's On Chain](https://tncpw.co/638d8e8a) and [Preev](https://tncpw.co/d19f43a3) for their hard work on their public API

[Jad](https://github.com/jadwahab) for his contributions to the package!

<br/>

## License

![License](https://img.shields.io/github/license/tonicpow/go-bsvrates.svg?style=flat&v=3)