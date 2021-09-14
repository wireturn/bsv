package main

import (
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
	if srv, err = client.GetSRVRecord(paymail.DefaultServiceName, paymail.DefaultProtocol, "moneybutton.com"); err != nil {
		log.Fatal("error getting SRV record: " + err.Error())
	}
	log.Println("found SRV record:", srv)
}
