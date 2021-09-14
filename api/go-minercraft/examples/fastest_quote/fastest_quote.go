package main

import (
	"context"
	"log"
	"time"

	"github.com/tonicpow/go-minercraft"
)

func main() {

	// Create a new client
	client, err := minercraft.NewClient(nil, nil, nil)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	log.Printf("querying %d miners for the fastest response...", len(client.Miners))

	// Fetch the fastest quote from all miners
	var response *minercraft.FeeQuoteResponse
	response, err = client.FastestQuote(context.Background(), 10*time.Second)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	log.Printf("found quote: %s", response.Miner.Name)
}
