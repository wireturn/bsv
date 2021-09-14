/*
Package main is an example package using go-polynym
*/
package main

import (
	"log"

	"github.com/mrz1836/go-polynym"
)

func main() {

	// Start a new client
	client := polynym.NewClient(nil)

	// Resolve a handle or paymail
	resp, err := polynym.GetAddress(client, "mrz@relayx.io")
	if err != nil {
		log.Fatal(err.Error())
	}

	// Success
	log.Println("address: ", resp.Address)
}
