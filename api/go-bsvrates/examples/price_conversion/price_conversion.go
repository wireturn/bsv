/*
Package main is an example of using the go-bsvrates package for price conversions
*/
package main

import (
	"context"
	"log"

	"github.com/tonicpow/go-bsvrates"
)

func main() {

	// Create a new client (all default providers)
	client := bsvrates.NewClient(nil, nil)

	// Get a conversion from $ to Sats
	satoshis, provider, _ := client.GetConversion(context.Background(), bsvrates.CurrencyDollars, 0.01)
	log.Printf("0.01 USD = satoshis: %d from provider: %s", satoshis, provider.Name())
}
