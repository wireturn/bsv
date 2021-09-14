package main

import (
	"fmt"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {

	// Create a client
	client := whatsonchain.NewClient(whatsonchain.NetworkMain, nil, nil)

	// Get the balance for multiple addresses
	balances, _ := client.BulkTransactionDetailsProcessor(
		&whatsonchain.TxHashes{TxIDs: []string{
			"cc84bf6aa5f0c3ab7e1e7f71bc40325576d0561bd07908ff8354308fcba7b4f0",
			"a33d408055fd8b2ac571a7d2016cf9f572d6a8cf5d905c0858d57818025c363a",
			"7677542b511bdf4d445dad6e835dd921f7fbe25833613479022ff1803007562e",
			"766d9f2b7da5f13aa679736d7a172b1b26de984838a7e9a2302a99c1a4c908fd",
			"60b22f4cf81b5e1aec080529096fd3cc99dd7eae09626c088f13768934aa7a4d",
			"c1b38c534773fdc5d2a600858e1c02572f309b40ef6182fa9149756ac7be15b1",
			"10b44f6f8a739a6223f911f7b52cd40c7b6b2abecfc287c2f9f11af3f8f7ed61",
		}},
	)

	for _, record := range balances {
		fmt.Printf(
			"tx: %s outputs: %d \n",
			record.TxID,
			len(record.Vout),
		)
	}
}
