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

	// Get the profile
	var profile *tonicpow.AdvertiserProfile
	if profile, _, err = client.GetAdvertiserProfile(23); err != nil {
		log.Fatalf("error in GetAdvertiserProfile: %s", err.Error())
	}

	log.Println("advertiser profile: " + profile.Name)
}
