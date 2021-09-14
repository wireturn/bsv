/*
Package main is an example of using the go-bsvrates package
*/
package main

import (
	"context"
	"log"
	"time"

	"github.com/tonicpow/go-bsvrates"
)

func main() {

	// Create a new client (all default providers)
	client := bsvrates.NewClient(nil, nil)

	// Get historical tickers
	response, err := client.CoinPaprika.GetHistoricalTickers(
		context.Background(),
		bsvrates.CoinPaprikaQuoteID,
		time.Now().UTC().Add(-1*24*time.Hour),
		time.Now().UTC(),
		0,
		bsvrates.TickerQuoteUSD,
		bsvrates.TickerInterval5m,
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("found %d historical rates", len(response.Results))
}
