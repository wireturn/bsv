package main

import (
	"log"

	"github.com/tonicpow/go-minercraft"
)

func main() {

	// Create a new client
	client, err := minercraft.NewClient(nil, nil, nil)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Add a custom miner!
	if err = client.AddMiner(minercraft.Miner{
		Name: "Custom",
		URL:  "https://mapi.customminer.com",
	}); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Show all miners loaded
	for _, miner := range client.Miners {
		log.Printf("miner: %s (%s)", miner.Name, miner.URL)
	}
}
