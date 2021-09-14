package main

import (
	"log"

	"github.com/bitcoinschema/go-bob"
	magic "github.com/bitcoinschema/go-map"
)

func main() {

	tape := bob.Tape{
		Cell: []bob.Cell{
			{S: magic.Prefix},
			{S: magic.Set},
			{S: "app"},
			{S: "myapp"},
		},
	}

	tx, err := magic.NewFromTape(&tape)
	if err != nil {
		log.Fatalf("failed to create new MAP from tape")
	}

	log.Printf("cmd: [%s] key: [%s] value: [%s]", tx[magic.Cmd], "app", tx["app"])
}
