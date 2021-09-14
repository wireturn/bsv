package main

import (
	"context"
	"log"
	"net"

	"github.com/tonicpow/go-paymail"
)

func main() {

	// Load the client
	client, err := paymail.NewClient()
	if err != nil {
		log.Fatalf("error loading client: %s", err.Error())
	}

	// Get the SRV record
	var srv *net.SRV
	srv, err = client.GetSRVRecord(paymail.DefaultServiceName, paymail.DefaultProtocol, "moneybutton.com")
	if err != nil {
		log.Fatal("error getting SRV record: " + err.Error())
	}

	// Found record!
	log.Println("found SRV record:", srv)

	// Validate the record (1 instead of 10, moneybutton deviated from the defaults)
	err = client.ValidateSRVRecord(
		context.Background(), srv, paymail.DefaultPort, 1, paymail.DefaultWeight,
	)
	if err != nil {
		log.Fatal("failed validating SRV record: " + err.Error())
	}
}
