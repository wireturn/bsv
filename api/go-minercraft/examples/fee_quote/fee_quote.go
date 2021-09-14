package main

import (
	"context"
	"log"

	"github.com/tonicpow/go-minercraft"
)

func main() {

	// Create a new client
	client, err := minercraft.NewClient(nil, nil, nil)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Select the miner
	miner := client.MinerByName(minercraft.MinerTaal)

	// Get a fee quote from a miner
	var response *minercraft.FeeQuoteResponse
	if response, err = client.FeeQuote(context.Background(), miner); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Example Tx Size (computed from the rawTx.ToBytes())
	txSizeInBytes := uint64(300) // The Tx is 300 bytes in size, let's get the fee for that Tx

	// Get the fee for a specific tx size (for mining and for data)
	var fee uint64
	if fee, err = response.Quote.CalculateFee(minercraft.FeeCategoryMining, minercraft.FeeTypeData, txSizeInBytes); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Display the results
	log.Printf("miner: %s", response.Miner.Name)
	log.Printf("fee quote expires at: %s", response.Quote.ExpirationTime)
	log.Printf("tx size in bytes: %d and mining fee: %d", txSizeInBytes, fee)
	log.Printf("payload validated: %v", response.Validated)
}
