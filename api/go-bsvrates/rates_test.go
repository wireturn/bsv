package bsvrates

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/mrz1836/go-preev"
	"github.com/mrz1836/go-whatsonchain"
	"github.com/stretchr/testify/assert"
)

// mockWOCValid for mocking requests
type mockWOCValid struct{}

// GetExchangeRate is a mock response
func (m *mockWOCValid) GetExchangeRate() (rate *whatsonchain.ExchangeRate, err error) {

	rate = &whatsonchain.ExchangeRate{
		Rate:     "159.01",
		Currency: CurrencyToName(CurrencyDollars),
	}

	return
}

// mockWOCFailed for mocking requests
type mockWOCFailed struct{}

// GetExchangeRate is a mock response
func (m *mockWOCFailed) GetExchangeRate() (rate *whatsonchain.ExchangeRate, err error) {

	return
}

// mockPaprikaValid for mocking requests
type mockPaprikaValid struct{}

// GetMarketPrice is a mock response
func (m *mockPaprikaValid) GetMarketPrice(ctx context.Context, coinID string) (response *TickerResponse, err error) {

	response = &TickerResponse{
		BetaValue:         1.39789,
		CirculatingSupply: 18448838,
		ID:                coinID,
		LastRequest: &lastRequest{
			Method:     http.MethodGet,
			StatusCode: http.StatusOK,
		},
		LastUpdated: "2020-07-01T18:36:56Z",
		MaxSupply:   21000000,
		Name:        "Bitcoin SV",
		Quotes: &currency{USD: &quote{
			Price:              158.49415248,
			Volume24h:          719426754.25105,
			Volume24hChange24h: -4.48,
			MarketCap:          2924031833,
		}},
		Rank:        6,
		Symbol:      "BSV",
		TotalSupply: 18448838,
	}

	return
}

// GetBaseAmountAndCurrencyID is a mock response
func (m *mockPaprikaValid) GetBaseAmountAndCurrencyID(currency string, amount float64) (string, float64) {

	// This is just a mock request

	return currency, 0.01
}

// GetPriceConversion is a mock response
func (m *mockPaprikaValid) GetPriceConversion(ctx context.Context, baseCurrencyID, quoteCurrencyID string, amount float64) (response *PriceConversionResponse, err error) {

	response = &PriceConversionResponse{
		Amount:                amount,
		BaseCurrencyID:        baseCurrencyID,
		BaseCurrencyName:      "US Dollars",
		BasePriceLastUpdated:  "2020-07-01T22:03:14Z",
		Price:                 0.006331560350007446,
		QuoteCurrencyID:       quoteCurrencyID,
		QuoteCurrencyName:     "Bitcoin SV",
		QuotePriceLastUpdated: "2020-07-01T22:03:14Z",
	}

	return
}

// GetHistoricalTickers is a mock response
func (m *mockPaprikaValid) GetHistoricalTickers(ctx context.Context, coinID string, start, end time.Time, limit int,
	quote tickerQuote, interval tickerInterval) (response *HistoricalResponse, err error) {

	// This is just a mock response

	return
}

// IsAcceptedCurrency is a mock response
func (m *mockPaprikaValid) IsAcceptedCurrency(currency string) bool {

	// This is just a mock response

	return true
}

// mockPreevValid for mocking requests
type mockPreevValid struct{}

// GetTicker is a mock response
func (m *mockPreevValid) GetTicker(ctx context.Context, pairID string) (ticker *preev.Ticker, err error) {

	ticker = &preev.Ticker{
		ID:        pairID,
		Timestamp: 1593628860,
		Tx: &preev.Transaction{
			Hash:      "175d87a3656a5d745af9fe9cee6afc0297a83fb317255962c40085eb31f06a4b",
			Timestamp: 1593628871,
		},
		Prices: &preev.PriceSource{
			Ppi: &preev.Price{
				LastPrice: 159.17,
				Volume:    935279,
			},
		},
	}

	return
}

// GetPair is a mock response
func (m *mockPreevValid) GetPair(ctx context.Context, pairID string) (pair *preev.Pair, err error) {

	return
}

// GetPairs is a mock response
func (m *mockPreevValid) GetPairs(ctx context.Context) (pairList *preev.PairList, err error) {

	return
}

// GetTickers is a mock response
func (m *mockPreevValid) GetTickers(ctx context.Context) (tickerList *preev.TickerList, err error) {

	return
}

// mockPaprikaFailed for mocking requests
type mockPaprikaFailed struct{}

// GetMarketPrice is a mock response
func (m *mockPaprikaFailed) GetMarketPrice(ctx context.Context, coinID string) (response *TickerResponse, err error) {
	err = fmt.Errorf("request to paprika fails... 502")
	return
}

// GetBaseAmountAndCurrencyID is a mock response
func (m *mockPaprikaFailed) GetBaseAmountAndCurrencyID(currency string, amount float64) (string, float64) {

	return "", 0
}

// GetPriceConversion is a mock response
func (m *mockPaprikaFailed) GetPriceConversion(ctx context.Context, baseCurrencyID, quoteCurrencyID string, amount float64) (response *PriceConversionResponse, err error) {

	return nil, fmt.Errorf("some error occurred")
}

// GetHistoricalTickers is a mock response
func (m *mockPaprikaFailed) GetHistoricalTickers(ctx context.Context, coinID string, start, end time.Time, limit int,
	quote tickerQuote, interval tickerInterval) (response *HistoricalResponse, err error) {

	// This is just a mock response

	return nil, fmt.Errorf("some error occurred")
}

// IsAcceptedCurrency is a mock response
func (m *mockPaprikaFailed) IsAcceptedCurrency(currency string) bool {

	return false
}

// mockPreevFailed for mocking requests
type mockPreevFailed struct{}

// GetPair is a mock response
func (m *mockPreevFailed) GetPair(ctx context.Context, pairID string) (pair *preev.Pair, err error) {

	return nil, fmt.Errorf("some error occurred")
}

// GetTicker is a mock response
func (m *mockPreevFailed) GetTicker(ctx context.Context, pairID string) (ticker *preev.Ticker, err error) {

	return nil, fmt.Errorf("some error occurred")
}

// GetTickers is a mock response
func (m *mockPreevFailed) GetTickers(ctx context.Context) (tickerList *preev.TickerList, err error) {

	return nil, fmt.Errorf("some error occurred")
}

// GetPairs is a mock response
func (m *mockPreevFailed) GetPairs(ctx context.Context) (pairList *preev.PairList, err error) {

	return nil, fmt.Errorf("some error occurred")
}

// newMockClient returns a client for mocking
func newMockClient(wocClient whatsOnChainInterface, paprikaClient coinPaprikaInterface, preevClient preevInterface, providers ...Provider) *Client {
	client := NewClient(nil, nil, providers...)
	client.WhatsOnChain = wocClient
	client.CoinPaprika = paprikaClient
	client.Preev = preevClient
	return client
}

// TestClient_GetRate will test the method GetRate()
func TestClient_GetRate(t *testing.T) {
	// t.Parallel()

	t.Run("valid get rate - default", func(t *testing.T) {
		client := newMockClient(&mockWOCValid{}, &mockPaprikaValid{}, &mockPreevValid{})
		assert.NotNil(t, client)

		rate, provider, err := client.GetRate(context.Background(), CurrencyDollars)
		assert.NoError(t, err)
		assert.Equal(t, 158.49415248, rate)
		assert.Equal(t, true, provider.IsValid())
		assert.Equal(t, "CoinPaprika", provider.Name())
	})

	t.Run("valid get rate - preev", func(t *testing.T) {
		client := newMockClient(&mockWOCValid{}, &mockPaprikaValid{}, &mockPreevValid{}, ProviderPreev)
		assert.NotNil(t, client)

		rate, provider, err := client.GetRate(context.Background(), CurrencyDollars)
		assert.NoError(t, err)
		assert.Equal(t, 159.17, rate)
		assert.Equal(t, true, provider.IsValid())
		assert.Equal(t, "Preev", provider.Name())
	})

	t.Run("valid get rate - whats on chain", func(t *testing.T) {
		client := newMockClient(&mockWOCValid{}, &mockPaprikaValid{}, &mockPreevValid{}, ProviderWhatsOnChain)
		assert.NotNil(t, client)

		rate, provider, err := client.GetRate(context.Background(), CurrencyDollars)
		assert.NoError(t, err)
		assert.Equal(t, 159.01, rate)
		assert.Equal(t, true, provider.IsValid())
		assert.Equal(t, "WhatsOnChain", provider.Name())
	})

	t.Run("valid get rate - custom providers", func(t *testing.T) {
		client := newMockClient(&mockWOCValid{}, &mockPaprikaValid{}, &mockPreevValid{}, ProviderPreev, ProviderWhatsOnChain)
		assert.NotNil(t, client)

		rate, provider, err := client.GetRate(context.Background(), CurrencyDollars)
		assert.NoError(t, err)
		assert.Equal(t, 159.17, rate)
		assert.Equal(t, true, provider.IsValid())
		assert.Equal(t, "Preev", provider.Name())
	})

	t.Run("non accepted currency", func(t *testing.T) {
		client := newMockClient(&mockWOCValid{}, &mockPaprikaFailed{}, &mockPreevValid{})
		assert.NotNil(t, client)

		_, _, rateErr := client.GetRate(context.Background(), 123)
		assert.Error(t, rateErr)
	})

	t.Run("failed rate - default", func(t *testing.T) {
		client := newMockClient(&mockWOCValid{}, &mockPaprikaFailed{}, &mockPreevValid{})
		assert.NotNil(t, client)

		rate, provider, err := client.GetRate(context.Background(), CurrencyDollars)
		assert.NoError(t, err)
		assert.Equal(t, 159.01, rate)
		assert.Equal(t, true, provider.IsValid())
		assert.Equal(t, "WhatsOnChain", provider.Name())
	})

	t.Run("failed rate - preev", func(t *testing.T) {
		client := newMockClient(&mockWOCValid{}, &mockPaprikaValid{}, &mockPreevFailed{}, ProviderPreev&ProviderCoinPaprika)
		assert.NotNil(t, client)

		rate, provider, err := client.GetRate(context.Background(), CurrencyDollars)
		assert.NoError(t, err)
		assert.Equal(t, 158.49415248, rate)
		assert.Equal(t, true, provider.IsValid())
		assert.Equal(t, "CoinPaprika", provider.Name())
	})

	t.Run("failed rate - whats on chain", func(t *testing.T) {
		client := newMockClient(&mockWOCFailed{}, &mockPaprikaValid{}, &mockPreevFailed{})
		assert.NotNil(t, client)

		rate, provider, err := client.GetRate(context.Background(), CurrencyDollars)
		assert.NoError(t, err)
		assert.Equal(t, 158.49415248, rate)
		assert.Equal(t, true, provider.IsValid())
		assert.Equal(t, "CoinPaprika", provider.Name())
	})

	t.Run("failed rate - all providers", func(t *testing.T) {
		client := newMockClient(&mockWOCFailed{}, &mockPaprikaFailed{}, &mockPreevFailed{})
		assert.NotNil(t, client)

		rate, _, err := client.GetRate(context.Background(), CurrencyDollars)
		assert.Error(t, err)
		assert.Equal(t, float64(0), rate)
	})
}
