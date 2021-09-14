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

	// Start with an existing profile (from the GET request)
	profile := &tonicpow.AdvertiserProfile{
		HomepageURL: "https://tonicpow.com",
		IconURL:     "https://i.imgur.com/HvVmeWI.png",
		PublicGUID:  "a4503e16b25c29b9cf58eee3ad353410",
		Name:        "TonicPow",
		ID:          23,
	}

	// Update the profile
	profile.Name = "TonicPow Test"
	if _, err = client.UpdateAdvertiserProfile(profile); err != nil {
		log.Fatalf("error in UpdateAdvertiserProfile: %s", err.Error())
	}

	log.Println("profile updated!")
}
