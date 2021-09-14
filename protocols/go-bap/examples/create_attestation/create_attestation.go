package main

import (
	"log"

	"github.com/bitcoinschema/go-bap"
)

func main() {
	// examplePrivateKey := "xprv9s21ZrQH143K2beTKhLXFRWWFwH8jkwUssjk3SVTiApgmge7kNC3jhVc4NgHW8PhW2y7BCDErqnKpKuyQMjqSePPJooPJowAz5BVLThsv6c"
	exampleIdKey := "8bafa4ca97d770276253585cb2a49da1775ec7aeed3178e346c8c1b55eaf5ca2"

	exampleAttributeName := "legal-name"
	exampleAttributeValue := "John Adams"
	exampleIdentityAttributeSecret := "e2c6fb4063cc04af58935737eaffc938011dff546d47b7fbb18ed346f8c4d4fa"

	tx, err := bap.CreateAttestation(
		exampleIdKey,
		"127d0ab318252b4622d8eac61407359a4cab7c1a5d67754b5bf9db910eaf052c",
		exampleAttributeName,
		exampleAttributeValue,
		exampleIdentityAttributeSecret,
	)
	if err != nil {
		log.Fatalf("failed to create attestation: %s", err.Error())
	}

	log.Printf("attestation tx created: %s", tx.ToString())
}
