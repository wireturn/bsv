package main

import (
	"fmt"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {

	// Create a client
	client := whatsonchain.NewClient(whatsonchain.NetworkMain, nil, nil)

	// Get the balance for multiple addresses
	balances, _ := client.BulkBalance(
		&whatsonchain.AddressList{Addresses: []string{"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", "1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob"}},
	)

	for _, record := range balances {
		fmt.Printf(
			"address: %s confirmed: %d unconfirmed: %d \n",
			record.Address,
			record.Balance.Confirmed,
			record.Balance.Unconfirmed,
		)
	}
}
