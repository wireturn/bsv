package bsvrates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestProvider_IsValid will test the method IsValid()
func TestProvider_IsValid(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase      string
		provider      Provider
		expectedValid bool
	}{
		{"provider 0", 0, false},
		{"provider 1", 1, true},
		{"provider 2", 2, true},
		{"provider 3", 3, true},
		{"provider 4", 4, false},
		{"ProviderWhatsOnChain", ProviderWhatsOnChain, true},
		{"ProviderCoinPaprika", ProviderCoinPaprika, true},
		{"ProviderPreev", ProviderPreev, true},
		{"providerLast", providerLast, false},
	}
	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			isValid := test.provider.IsValid()
			assert.Equal(t, test.expectedValid, isValid)
		})
	}
}

// TestProvider_Name will test the method Name()
func TestProvider_Name(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase     string
		provider     Provider
		expectedName string
	}{
		{"provider 0", 0, ""},
		{"provider 1", 1, "WhatsOnChain"},
		{"provider 2", 2, "CoinPaprika"},
		{"provider 3", 3, "Preev"},
		{"provider 4", 4, ""},
		{"ProviderWhatsOnChain", ProviderWhatsOnChain, "WhatsOnChain"},
		{"ProviderCoinPaprika", ProviderCoinPaprika, "CoinPaprika"},
		{"ProviderPreev", ProviderPreev, "Preev"},
		{"providerLast", providerLast, ""},
	}
	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			name := test.provider.Name()
			assert.Equal(t, test.expectedName, name)
		})
	}
}

// TestProviderToName will test the method ProviderToName()
func TestProviderToName(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase     string
		provider     Provider
		expectedName string
	}{
		{"provider 0", 0, ""},
		{"provider 1", 1, "WhatsOnChain"},
		{"provider 2", 2, "CoinPaprika"},
		{"provider 3", 3, "Preev"},
		{"provider 4", 4, ""},
		{"ProviderWhatsOnChain", ProviderWhatsOnChain, "WhatsOnChain"},
		{"ProviderCoinPaprika", ProviderCoinPaprika, "CoinPaprika"},
		{"ProviderPreev", ProviderPreev, "Preev"},
		{"providerLast", providerLast, ""},
	}
	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			name := ProviderToName(test.provider)
			assert.Equal(t, test.expectedName, name)
		})
	}
}

// TestCurrency_IsValid will test the method IsValid()
func TestCurrency_IsValid(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		testCase      string
		currency      Currency
		expectedValid bool
	}{
		{"currency 0", 0, false},
		{"currency 1", 1, true},
		{"currency 2", 2, true},
		{"currency 3", 3, false},
		{"CurrencyDollars", CurrencyDollars, true},
		{"CurrencyBitcoin", CurrencyBitcoin, true},
		{"currencyLast", currencyLast, false},
	}
	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			isValid := test.currency.IsValid()
			assert.Equal(t, test.expectedValid, isValid)
		})
	}
}

// TestCurrency_Name will test the method Name()
func TestCurrency_Name(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		testCase     string
		currency     Currency
		expectedName string
	}{
		{"currency 0", 0, ""},
		{"currency 1", 1, usd},
		{"currency 2", 2, "bsv"},
		{"currency 3", 3, ""},
		{"CurrencyDollars", CurrencyDollars, usd},
		{"CurrencyBitcoin", CurrencyBitcoin, "bsv"},
		{"currencyLast", currencyLast, ""},
	}
	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			name := test.currency.Name()
			assert.Equal(t, test.expectedName, name)
		})
	}
}

// TestCurrencyToName will test the method CurrencyToName()
func TestCurrencyToName(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		testCase     string
		currency     Currency
		expectedName string
	}{
		{"currency 0", 0, ""},
		{"currency 1", 1, usd},
		{"currency 2", 2, "bsv"},
		{"currency 3", 3, ""},
		{"CurrencyDollars", CurrencyDollars, usd},
		{"CurrencyBitcoin", CurrencyBitcoin, "bsv"},
		{"currencyLast", currencyLast, ""},
	}
	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			name := CurrencyToName(test.currency)
			assert.Equal(t, test.expectedName, name)
		})
	}
}

// TestCurrencyFromName will test the method CurrencyFromName()
func TestCurrencyFromName(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase         string
		currency         string
		expectedCurrency Currency
	}{
		{"", "", CurrencyDollars},
		{usd, usd, CurrencyDollars},
		{"bsv", "bsv", CurrencyBitcoin},
		{"bogus", "bogus", CurrencyDollars},
	}
	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			currency := CurrencyFromName(test.currency)
			assert.Equal(t, test.expectedCurrency, currency)
		})
	}
}

// TestCurrency_IsAccepted will test the method IsAccepted()
func TestCurrency_IsAccepted(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase      string
		currency      Currency
		expectedValid bool
	}{
		{"currency 0", 0, false},
		{"currency 1", 1, true},
		{"currency 2", 2, false},
		{"currency 3", 3, false},
		{"CurrencyDollars", CurrencyDollars, true},
		{"CurrencyBitcoin", CurrencyBitcoin, false},
		{"currencyLast", currencyLast, false},
	}
	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			isValid := test.currency.IsAccepted()
			assert.Equal(t, test.expectedValid, isValid)
		})
	}
}
