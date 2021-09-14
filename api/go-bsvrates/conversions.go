package bsvrates

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mrz1836/go-preev"
	"github.com/mrz1836/go-whatsonchain"
)

// GetConversion will get the satoshi amount for the given currency + amount provided.
// The first provider that succeeds is the conversion that is returned
func (c *Client) GetConversion(ctx context.Context, currency Currency, amount float64) (satoshis int64, providerUsed Provider, err error) {

	// Check if currency is accepted across all providers
	if !currency.IsAccepted() {
		err = fmt.Errorf("currency [%s] is not accepted by all providers at this time", currency.Name())
		return
	}

	// Loop providers and get a conversion value
	for _, provider := range c.Providers {
		providerUsed = provider
		switch provider {
		case ProviderCoinPaprika:
			var response *PriceConversionResponse
			if response, err = c.CoinPaprika.GetPriceConversion(
				ctx, USDCurrencyID, CoinPaprikaQuoteID, amount,
			); err == nil && response != nil {
				satoshis, err = response.GetSatoshi()
			}
		case ProviderWhatsOnChain:
			var response *whatsonchain.ExchangeRate
			if response, err = c.WhatsOnChain.GetExchangeRate(); err == nil && response != nil {
				var rate float64
				if rate, err = strconv.ParseFloat(response.Rate, 64); err == nil {
					satoshis, err = ConvertPriceToSatoshis(rate, amount)
				}
			}
		case ProviderPreev:
			var response *preev.Ticker
			if response, err = c.Preev.GetTicker(
				ctx, PreevTickerID,
			); err == nil && response != nil {
				satoshis, err = ConvertPriceToSatoshis(response.Prices.Ppi.LastPrice, amount)
			}
		case providerLast:
			err = fmt.Errorf("provider unknown")
			return
		}

		// todo: log the error for sanity in case the user want's to see the failure?

		// Did we get a satoshi value? Otherwise, keep looping
		if satoshis > 0 {
			return
		}
	}

	return
}

// todo: create a new method to get all three and then average the results
