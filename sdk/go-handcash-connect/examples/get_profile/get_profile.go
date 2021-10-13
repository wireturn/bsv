package main

import (
	"context"
	"log"
	"os"

	"github.com/tonicpow/go-handcash-connect"
)

func main() {

	// Create a new client (Beta ENV)
	client := handcash.NewClient(nil, nil, handcash.EnvironmentBeta)

	// Get the current user's profile (Auth token was from oAuth callback)
	profile, err := client.GetProfile(context.Background(), os.Getenv("AUTH_TOKEN"))
	if err != nil {
		log.Fatalln("error: ", err)
	}
	log.Println("found profile: ", profile)
}
