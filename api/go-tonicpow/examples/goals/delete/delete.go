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

	// Delete a goal
	var deleted bool
	deleted, _, err = client.DeleteGoal(63)
	if err != nil {
		log.Fatalf("error in DeleteGoal: %s", err.Error())
	}

	log.Printf("goal deleted: %t", deleted)
}
