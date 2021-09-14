/*
Package main is an example of using the go-bsvrates package
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

	// Get rates
	rate, provider, _ := client.GetRate(context.Background(), bsvrates.CurrencyDollars)
	log.Printf("found rate: %v %s from provider: %s", rate, bsvrates.CurrencyToName(bsvrates.CurrencyDollars), provider.Name())
}
