/*
Package bsvrates brings multiple providers into one place to obtain the current BSV exchange rate
*/
package bsvrates

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mrz1836/go-preev"
	"github.com/mrz1836/go-whatsonchain"
)

// GetRate will get a BSV->Currency rate from the list of providers.
// The first provider that succeeds is the rate that is returned
func (c *Client) GetRate(ctx context.Context, currency Currency) (rate float64, providerUsed Provider, err error) {

	// Check if currency is accepted across all providers
	if !currency.IsAccepted() {
		err = fmt.Errorf("currency [%s] is not accepted by all providers at this time", currency.Name())
		return
	}

	// Loop providers and get a rate
	for _, provider := range c.Providers {
		providerUsed = provider
		switch provider {
		case ProviderCoinPaprika:
			var response *TickerResponse
			if response, err = c.CoinPaprika.GetMarketPrice(ctx, CoinPaprikaQuoteID); err == nil && response != nil {
				rate = response.Quotes.USD.Price
			}
		case ProviderWhatsOnChain:
			var response *whatsonchain.ExchangeRate
			if response, err = c.WhatsOnChain.GetExchangeRate(); err == nil && response != nil {
				rate, err = strconv.ParseFloat(response.Rate, 64)
			}
		case ProviderPreev:
			var response *preev.Ticker
			if response, err = c.Preev.GetTicker(ctx, PreevTickerID); err == nil && response != nil {
				rate = response.Prices.Ppi.LastPrice
			}
		case providerLast:
			err = fmt.Errorf("provider unknown")
			return
		}

		// todo: log the error for sanity in case the user want's to see the failure?

		// Did we get a rate? Otherwise, keep looping
		if rate > 0 {
			return
		}
	}

	return
}

// todo: create a new method to get all three and then average the results
