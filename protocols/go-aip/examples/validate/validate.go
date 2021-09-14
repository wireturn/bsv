package main

import (
	"log"

	"github.com/bitcoinschema/go-aip"
)

func main() {
	a, err := aip.Sign(
		"54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd",
		aip.BitcoinECDSA,
		"example message",
	)
	if err != nil {
		log.Fatalf("error occurred: %s", err.Error())
	}
	if _, err = a.Validate(); err == nil {
		log.Printf("signature is valid: %s", a.Signature)
	} else {
		log.Fatalf("signature failed validation: %s", err.Error())
	}
}
