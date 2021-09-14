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

	// Query the transaction status
	var response *minercraft.QueryTransactionResponse
	if response, err = client.QueryTransaction(context.Background(), miner, "950a10beb1650e91621f748c408f7024f2082408a93c11cecc1ab4b5f440ac12"); err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}

	// Display the results
	log.Printf("miner: %s", response.Miner.Name)
	log.Printf("status: %s [%s]", response.Query.ReturnResult, response.Query.ResultDescription)
	log.Printf("payload validated: %v", response.Validated)
}
