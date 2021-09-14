package main

import (
	"fmt"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {

	// Create a client
	client := whatsonchain.NewClient(whatsonchain.NetworkMain, nil, nil)

	// Get UTXOs for an address
	history, err := client.AddressUnspentTransactionDetails("16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", 0)
	if err != nil {
		fmt.Printf("error getting utxos: %s", err.Error())
	} else if len(history) == 0 {
		fmt.Println("no utxos found")
	} else {
		for index, utxo := range history {
			fmt.Printf("(%d) %s | Sats: %d \n", index+1, utxo.TxHash, utxo.Value)
		}
	}
}
