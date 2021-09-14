package main

import (
	"log"
	"os"

	"github.com/tonicpow/go-tonicpow"
)

func main() {

	// Load the api client
	client, err := tonicpow.NewClient(
		tonicpow.WithAPIKey(os.Getenv("TONICPOW_API_KEY")),
		tonicpow.WithEnvironmentString(os.Getenv("TONICPOW_ENVIRONMENT")),
	)
	if err != nil {
		log.Fatalf("error in NewClient: %s", err.Error())
	}

	// Get current rate
	var rate *tonicpow.Rate
	rate, _, err = client.GetCurrentRate("usd", 1.00)
	if err != nil {
		log.Fatalf("error in GetCurrentRate: %s", err.Error())
	}

	log.Printf("rate: %s %f is %d sats", rate.Currency, rate.CurrencyAmount, rate.PriceInSatoshis)
}
