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

	// Get a goal
	var goal *tonicpow.Goal
	goal, _, err = client.GetGoal(13)
	if err != nil {
		log.Fatalf("error in GetGoal: %s", err.Error())
	}

	log.Printf("goal: %s", goal.Name)
}
