package main

import (
	"log"

	"github.com/bitcoinschema/go-bap"
)

func main() {
	examplePrivateKey := "xprv9s21ZrQH143K2beTKhLXFRWWFwH8jkwUssjk3SVTiApgmge7kNC3jhVc4NgHW8PhW2y7BCDErqnKpKuyQMjqSePPJooPJowAz5BVLThsv6c"
	exampleIdKey := "8bafa4ca97d770276253585cb2a49da1775ec7aeed3178e346c8c1b55eaf5ca2"

	tx, err := bap.CreateIdentity(examplePrivateKey, exampleIdKey, 0)
	if err != nil {
		log.Fatalf("failed to create identity: %s", err.Error())
	}

	log.Printf("tx created: %s", tx.ToString())
}
