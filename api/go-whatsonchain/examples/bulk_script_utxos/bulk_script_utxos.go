package main

import (
	"fmt"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {

	// Create a client
	client := whatsonchain.NewClient(whatsonchain.NetworkMain, nil, nil)

	// Get the balance for multiple addresses
	balances, _ := client.BulkScriptUnspentTransactions(
		&whatsonchain.ScriptsList{Scripts: []string{
			"f814a7c3a40164aacc440871e8b7b14eb6a45f0ca7dcbeaea709edc83274c5e7",
			"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3",
		}},
	)

	for _, record := range balances {
		fmt.Printf(
			"script: %s utxos: %d \n",
			record.Script,
			len(record.Utxos),
		)
	}
}
