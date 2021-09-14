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

	// List campaign
	var results *tonicpow.CampaignResults
	results, _, err = client.ListCampaigns(
		1, 25, "", "", "", 0, false,
	)
	if err != nil {
		log.Fatalf("error in ListCampaigns: %s", err.Error())
	}

	log.Printf("campaigns found: %d", results.Results)
}
