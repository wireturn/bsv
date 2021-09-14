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

	// Get a RSS feed
	var results string
	// results, _, err = client.CampaignsFeed(tonicpow.FeedTypeRSS)
	// results, _, err = client.CampaignsFeed(tonicpow.FeedTypeAtom)
	results, _, err = client.CampaignsFeed(tonicpow.FeedTypeJSON)
	if err != nil {
		log.Fatalf("error in CampaignsFeed: %s", err.Error())
	}

	log.Printf(results)
}
