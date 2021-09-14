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

	// Get campaign first
	var campaign *tonicpow.Campaign
	campaign, _, err = client.GetCampaign(23)
	if err != nil {
		log.Fatalf("error in GetCampaign: %s", err.Error())
	}

	// Update a campaign
	campaign.Title = "Updated Title Name"
	_, err = client.UpdateCampaign(campaign)
	if err != nil {
		log.Fatalf("error in UpdateCampaign: %s", err.Error())
	}

	log.Printf("updated campaign: %s", campaign.Title)
}
