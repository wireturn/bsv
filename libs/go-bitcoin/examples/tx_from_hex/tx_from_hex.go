package main

import (
	"log"

	"github.com/bitcoinschema/go-bitcoin"
)

func main() {

	// Example raw tx
	exampleTx := "0100000001760595866e99c1ce920197844740f5598b34763878696371d41b3a7c0a65b0b7000000006b483045022100eea3d606bd1627be6459a9de4860919225db74843d2fc7f4e7caa5e01f42c2d0022017978d9c6a0e934955a70e7dda71d68cb614f7dd89eb7b9d560aea761834ddd4412102ea87d1fd77d169bd56a71e700628113d0f8dfe57faa0ba0e55a36f9ce8e10be3ffffffff03f4010000000000001976a9147a1980655efbfec416b2b0c663a7b3ac0b6a25d288ac00000000000000001a006a07707265666978310c6578616d706c65206461746102133700000000000000001c006a0770726566697832116d6f7265206578616d706c65206461746100000000"

	rawTx, err := bitcoin.TxFromHex(exampleTx)
	if err != nil {
		log.Printf("error occurred: %s", err.Error())
		return
	}

	log.Printf("tx id: %s", rawTx.GetTxID())
}