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

	// Get a conversion
	var conversion *tonicpow.Conversion
	if conversion, _, err = client.GetConversion(99); err != nil {
		log.Fatalf("error in GetConversion: %s", err.Error())
	}

	// Cancel a conversion
	if conversion, _, err = client.CancelConversion(
		conversion.ID, "my custom reason",
	); err != nil {
		log.Fatalf("error in CancelConversion: %s", err.Error())
	}

	log.Printf("conversion: %d:%s", conversion.ID, conversion.Status)
}
