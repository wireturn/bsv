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
	// This is required first to get the corresponding VerifyPubKey endpoint url
	var capabilities *paymail.Capabilities
	if capabilities, err = client.GetCapabilities("moneybutton.com", paymail.DefaultPort); err != nil {
		log.Fatal("error getting capabilities: " + err.Error())
	}
	log.Println("found capabilities: ", len(capabilities.Capabilities))

	// Extract the Verify URL from the capabilities response
	verifyURL := capabilities.GetString(paymail.BRFCVerifyPublicKeyOwner, "")

	// Verify the pubkey
	var verification *paymail.Verification
	verification, err = client.VerifyPubKey(verifyURL, "mrz", "moneybutton.com", "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10")
	if err != nil {
		log.Fatal("error getting verification: " + err.Error())
	}
	if verification.Match {
		log.Printf("pubkey: %s matched handle: %s", verification.PubKey[:12]+"...", verification.Handle)
	} else {
		log.Printf("pubkey: %s DID NOT MATCH handle: %s", verification.PubKey[:12]+"...", verification.Handle)
	}
}
