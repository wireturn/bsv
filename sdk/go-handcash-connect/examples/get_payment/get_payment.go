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

	// Get the payment information (given AuthToken and TxID)
	payment, err := client.GetPayment(
		context.Background(),
		os.Getenv("AUTH_TOKEN"),
		os.Getenv("TX_ID"),
	)
	if err != nil {
		log.Fatalln("error: ", err)
	}
	log.Println("payment: ", payment)
}
