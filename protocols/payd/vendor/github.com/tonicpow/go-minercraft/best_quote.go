package minercraft

import (
	"context"
	"sync"
)

// BestQuote will check all known miners and compare rates, returning the best rate/quote
//
// Note: this might return different results each time if miners have the same rates as
// it's a race condition on which results come back first
func (c *Client) BestQuote(ctx context.Context, feeCategory, feeType string) (*FeeQuoteResponse, error) {

	// Best rate & quote
	var bestRate uint64
	var bestQuote FeeQuoteResponse

	// The channel for the internal results
	resultsChannel := make(chan *internalResult, len(c.Miners))

	// Loop each miner (break into a Go routine for each quote request)
	var wg sync.WaitGroup
	for _, miner := range c.Miners {
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, client *Client,
			miner *Miner, resultsChannel chan *internalResult) {
			defer wg.Done()
			resultsChannel <- getQuote(ctx, client, miner)
		}(ctx, &wg, c, miner, resultsChannel)
	}

	// Waiting for all requests to finish
	wg.Wait()
	close(resultsChannel)

	// Loop the results of the channel
	var testRate uint64
	var quoteFound bool
	var lastErr error
	for result := range resultsChannel {

		// Check for error?
		if result.Response.Error != nil {
			lastErr = result.Response.Error
			continue
		}

		// Parse the response
		var quote FeeQuoteResponse
		if quote, lastErr = result.parseQuote(); lastErr != nil {
			continue
		}

		// Get a test rate
		if testRate, lastErr = quote.Quote.CalculateFee(feeCategory, feeType, 1000); lastErr != nil {
			continue
		}

		// (Never set) || (or better than previous rate)
		quoteFound = true
		if bestRate == 0 || testRate < bestRate {
			bestRate = testRate
			bestQuote = quote
		}
	}

	// No quotes?
	if !quoteFound && lastErr != nil {
		return nil, lastErr
	}

	// Return the best quote found
	return &bestQuote, nil
}
