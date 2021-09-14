package bsvrates

import "strings"

const (

	// version is the current package version
	version = "v0.1.12"

	// defaultUserAgent is the default user agent for all requests
	defaultUserAgent string = "go-bsvrates: " + version

	// CoinPaprikaQuoteID is the id for CoinPaprika (BSV)
	CoinPaprikaQuoteID = "bsv-bitcoin-sv"

	// PreevTickerID is the id for Preev (BSV)
	PreevTickerID = "12eLTxv1vyUeJtp5zqWbqpdWvfLdZ7dGf8"
)

var (
	// defaultProviders (if no provider slice is set, use this as the default set)
	defaultProviders = []Provider{ProviderCoinPaprika, ProviderWhatsOnChain, ProviderPreev}
)

// Provider is a provider for rates or prices
type Provider uint8

// Provider constants for the different available rate providers.
// Leave the start and last constants in place
const (
	_ Provider = iota // 0

	ProviderWhatsOnChain // 1
	ProviderCoinPaprika  // 2
	ProviderPreev        // 3
	providerLast         // 4
)

// IsValid tests if the provider is valid or not
func (p Provider) IsValid() bool {
	return p >= ProviderWhatsOnChain && p < providerLast
}

// Name will return the display name for the given provider
func (p Provider) Name() string {
	switch p {
	case ProviderWhatsOnChain:
		return "WhatsOnChain"
	case ProviderCoinPaprika:
		return "CoinPaprika"
	case ProviderPreev:
		return "Preev"
	case providerLast:
		return ""
	default:
		return ""
	}
}

// ProviderToName helper function to convert the provider value to it's associated name
func ProviderToName(provider Provider) string {
	return provider.Name()
}

// Currency is a valid currency for rates or prices
type Currency uint8

// Currency constants for the different available currencies.
// Leave the start and last constants in place
const (
	_               Currency = iota
	CurrencyDollars          = 1
	CurrencyBitcoin          = 2

	currencyLast = iota
)

// IsValid tests if the provider is valid or not
func (c Currency) IsValid() bool {
	return c >= CurrencyDollars && c < currencyLast
}

// IsAccepted tests if the currency is accepted by all providers
func (c Currency) IsAccepted() bool {
	return c == CurrencyDollars
}

// Name will return the display name for the given currency
func (c Currency) Name() string {
	switch c {
	case CurrencyDollars:
		return usd
	case CurrencyBitcoin:
		return "bsv"
	default:
		return ""
	}
}

// CurrencyToName helper function to convert the currency value to it's associated name
func CurrencyToName(currency Currency) string {
	return currency.Name()
}

// CurrencyFromName helper function to convert the name into it's Currency type
func CurrencyFromName(name string) Currency {
	switch strings.ToLower(name) {
	case usd:
		return CurrencyDollars
	case "bsv":
		return CurrencyBitcoin
	default:
		return CurrencyDollars
	}
}
