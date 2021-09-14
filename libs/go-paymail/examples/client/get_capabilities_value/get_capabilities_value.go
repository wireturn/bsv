package main

import (
	"log"

	"github.com/tonicpow/go-paymail"
)

func main() {

	// Load the client
	client, err := paymail.NewClient()
	if err != nil {
		log.Fatalf("error loading client: %s", err.Error())
	}

	// Get the capabilities
	var capabilities *paymail.Capabilities
	capabilities, err = client.GetCapabilities("moneybutton.com", paymail.DefaultPort)
	if err != nil {
		log.Fatal("error getting capabilities: " + err.Error())
	}
	log.Println("found capabilities: ", len(capabilities.Capabilities))

	// Get the URL for a capability
	endpoint := capabilities.GetString(paymail.BRFCPki, paymail.BRFCPkiAlternate)
	log.Println("capability endpoint found: ", endpoint)
}
