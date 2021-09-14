# go-mattercloud
> The unofficial Go implementation for the [MatterCloud API](https://developers.mattercloud.net/)

[![Go](https://img.shields.io/github/go-mod/go-version/mrz1836/go-mattercloud?v=1)](https://golang.org/)
[![Build Status](https://img.shields.io/github/workflow/status/mrz1836/go-mattercloud/run-go-tests?logo=github&v=3)](https://github.com/mrz1836/go-mattercloud/actions)
[![Report](https://goreportcard.com/badge/github.com/mrz1836/go-mattercloud?style=flat&v=1)](https://goreportcard.com/report/github.com/mrz1836/go-mattercloud)
[![Release](https://img.shields.io/github/release-pre/mrz1836/go-mattercloud.svg?style=flat&v=1)](https://github.com/mrz1836/go-mattercloud/releases)
[![GoDoc](https://godoc.org/github.com/mrz1836/go-mattercloud?status.svg&style=flat)](https://pkg.go.dev/github.com/mrz1836/go-mattercloud)
[![Sponsor](https://img.shields.io/badge/sponsor-MrZ-181717.svg?logo=github&style=flat&v=4)](https://github.com/sponsors/mrz1836)
[![Donate](https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat)](https://mrz1818.com/?tab=tips&af=go-mattercloud)      

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

**go-mattercloud** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/mrz1836/go-mattercloud
```

## Documentation
You can view the generated [documentation here](https://pkg.go.dev/github.com/mrz1836/go-mattercloud).

You can also view the [MatterCloud API](https://developers.mattercloud.net/) documentation.

### Features
- [Client](client.go) is completely configurable
- Customize the network per request (`main`, `test` or `stn`)
- Using [heimdall http client](https://github.com/gojek/heimdall) with exponential backoff & more
- Current (V3) coverage for the [MatterCloud API](https://developers.mattercloud.net/) API
    - [x] Authentication
    - [x] Address
    - [x] Transaction

<details>
<summary><strong><code>Library Deployment</code></strong></summary>

[goreleaser](https://github.com/goreleaser/goreleaser) for easy binary or library deployment to Github and can be installed via: `brew install goreleaser`.

The [.goreleaser.yml](.goreleaser.yml) file is used to configure [goreleaser](https://github.com/goreleaser/goreleaser).

Use `make release-snap` to create a snapshot version of the release, and finally `make release` to ship to production.
</details>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>

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

## Examples & Tests
All unit tests and [examples](mattercloud_test.go) run via [Github Actions](https://github.com/mrz1836/go-mattercloud/actions) and
uses [Go version 1.15.x](https://golang.org/doc/go1.15). View the [configuration file](.github/workflows/run-tests.yml).

Run all tests (including integration tests)
```shell script
make test
```

Run tests (excluding integration tests)
```shell script
make test-short
```

## Benchmarks
Run the Go [benchmarks](mattercloud_test.go):
```shell script
make bench
```

## Code Standards
Read more about this Go project's [code standards](CODE_STANDARDS.md).

## Usage
- View the [mattercloud examples](#examples--tests) above

Basic implementation:
```go
package main

import (
	"log"

	"github.com/mrz1836/go-mattercloud"
)

func main() {

	// Create a new client
	client, _ := mattercloud.NewClient("your-secret-api-key", mattercloud.NetworkMain, nil)

	// Get balance for an address
	balance, _ := client.AddressBalance("16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")

	// What's the confirmed balance?
	log.Println("confirmed balance:", balance.Confirmed)
}
```

## Maintainers

| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
|:---:|
| [MrZ](https://github.com/mrz1836) |

## Contributing
View the [contributing guidelines](CONTRIBUTING.md) and please follow the [code of conduct](CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:!
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:.
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/mrz1836) :clap:
or by making a [**bitcoin donation**](https://mrz1818.com/?tab=tips&af=go-mattercloud) to ensure this journey continues indefinitely! :rocket:

#### Credits

[@Attila](https://github.com/attilaaf) & [MatterCloud](https://mattercloud.net/) for their hard work on the [MatterCloud API](https://developers.mattercloud.net/)

Looking for a Javascript version? Check out the [MatterCloud JS SDK](https://github.com/MatterCloud/mattercloudjs)

## License

![License](https://img.shields.io/github/license/mrz1836/go-mattercloud.svg?style=flat&v=2)