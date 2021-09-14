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

	// Start campaign
	campaign := &tonicpow.Campaign{
		AdvertiserProfileID: 23,
		Description:         "example campaign",
		PayPerClickRate:     1,
		TargetURL:           "https://tonicpow.com",
		Title:               "Example Campaign",
	}

	// Create a campaign
	_, err = client.CreateCampaign(campaign)
	if err != nil {
		log.Fatalf("error in CreateCampaign: %s", err.Error())
	}

	log.Printf("created campaign: %d", campaign.ID)
}
