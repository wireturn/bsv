package minercraft

import (
	"context"
	"errors"
	"sync"
	"time"
)

// FastestQuote will check all known miners and return the fastest quote response
//
// Note: this might return different results each time if miners have the same rates as
// it's a race condition on which results come back first
func (c *Client) FastestQuote(ctx context.Context, timeout time.Duration) (*FeeQuoteResponse, error) {

	// No timeout (use the default)
	if timeout.Seconds() == 0 {
		timeout = defaultFastQuoteTimeout
	}

	// Get the fastest quote
	result := c.fetchFastestQuote(ctx, timeout)
	if result == nil {
		return nil, errors.New("no quotes found")
	}

	// Check for error?
	if result.Response.Error != nil {
		return nil, result.Response.Error
	}

	// Parse the response
	quote, err := result.parseQuote()
	if err != nil {
		return nil, err
	}

	// Return the quote
	return &quote, nil
}

// fetchFastestQuote will return a quote that is the quickest to resolve
func (c *Client) fetchFastestQuote(ctx context.Context, timeout time.Duration) *internalResult {

	// The channel for the internal results
	resultsChannel := make(chan *internalResult, len(c.Miners))

	// Create a context (to cancel or timeout)
	ctxWithCancel, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Loop each miner (break into a Go routine for each quote request)
	var wg sync.WaitGroup
	for _, miner := range c.Miners {
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, client *Client, miner *Miner) {
			defer wg.Done()
			res := getQuote(ctx, client, miner)
			if res.Response.Error == nil {
				resultsChannel <- res
			}
		}(ctxWithCancel, &wg, c, miner)
	}

	// Waiting for all requests to finish
	go func() {
		wg.Wait()
		close(resultsChannel)
	}()

	return <-resultsChannel
}
