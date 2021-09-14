package main

import (
	"log"

	"github.com/bitcoinschema/go-aip"
)

func main() {
	out, _, a, err := aip.SignOpReturnData(
		"54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd",
		aip.BitcoinECDSA,
		[][]byte{[]byte("some op_return data")},
	)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}
	log.Printf("address: %s", a.AlgorithmSigningComponent)
	log.Printf("signature: %s", a.Signature)
	log.Printf("output: %s", out.GetLockingScriptHexString())
}
