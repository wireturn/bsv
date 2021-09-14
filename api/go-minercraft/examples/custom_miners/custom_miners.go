package main

import (
	"log"

	"github.com/tonicpow/go-minercraft"
)

func main() {

	// Create a new client using a custom set of miners
	client, err := minercraft.NewClient(
		nil,
		nil,
		[]*minercraft.Miner{{
			MinerID: "example-miner-id",
			Name:    "example-miner",
			Token:   "",
			URL:     "https://example.miner.endpoint.com",
		}},
	)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Show all miners loaded
	for _, miner := range client.Miners {
		log.Printf("miner: %s (%s)", miner.Name, miner.URL)
	}
}
