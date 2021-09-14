/*
Package main is an example of using the go-bsvrates package using custom providers
*/
package main

import (
	"context"
	"log"

	"github.com/tonicpow/go-bsvrates"
)

func main() {

	// Create a new client (custom providers)
	client := bsvrates.NewClient(nil, nil, bsvrates.ProviderWhatsOnChain, bsvrates.ProviderCoinPaprika)

	// Get rates
	rate, provider, _ := client.GetRate(context.Background(), bsvrates.CurrencyDollars)
	log.Printf("found rate: %v %s from provider: %s", rate, bsvrates.CurrencyToName(bsvrates.CurrencyDollars), provider.Name())
}
