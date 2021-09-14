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

	// Create conversion
	var conversion *tonicpow.Conversion
	conversion, _, err = client.CreateConversion(
		tonicpow.WithGoalName("example_goal"),
		tonicpow.WithUserID(43),
	)
	if err != nil {
		log.Fatalf("error in CreateConversion: %s", err.Error())
	}

	log.Printf("created conversion: %d", conversion.ID)
}
