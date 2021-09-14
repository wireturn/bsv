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

	// Get the campaigns
	var results *tonicpow.CampaignResults
	if results, _, err = client.ListCampaignsByAdvertiserProfile(
		23, 1, 5, "", "",
	); err != nil {
		log.Fatalf("error in ListCampaignsByAdvertiserProfile: %s", err.Error())
	}

	log.Printf("results: %d", results.Results)
}
