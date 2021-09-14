/*
Package main is an example of using the go-bsvrates package using custom options
*/
package main

import (
	"context"
	"log"
	"time"

	"github.com/tonicpow/go-bsvrates"
)

func main() {

	// Set your own custom options
	options := bsvrates.DefaultClientOptions()
	options.UserAgent = "custom-user-agent"
	options.TransportIdleTimeout = 30 * time.Second

	// Create a new client (custom options & providers)
	client := bsvrates.NewClient(options, nil, bsvrates.ProviderWhatsOnChain, bsvrates.ProviderCoinPaprika)

	// Get rates
	rate, provider, _ := client.GetRate(context.Background(), bsvrates.CurrencyDollars)
	log.Printf("found rate: %v %s from provider: %s", rate, bsvrates.CurrencyToName(bsvrates.CurrencyDollars), provider.Name())
}
