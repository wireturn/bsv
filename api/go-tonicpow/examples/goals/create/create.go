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

	// Start goal
	goal := &tonicpow.Goal{
		CampaignID:     23,
		Description:    "Example goal description",
		MaxPerPromoter: 1,
		Name:           "example_goal",
		PayoutRate:     0.01,
		PayoutType:     "flat",
		Title:          "Example Goal",
	}

	// Create a goal
	_, err = client.CreateGoal(goal)
	if err != nil {
		log.Fatalf("error in CreateGoal: %s", err.Error())
	}

	log.Printf("created goal: %d", goal.ID)
}
