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
	// This is required first to get the corresponding PublicProfile endpoint url
	var capabilities *paymail.Capabilities
	if capabilities, err = client.GetCapabilities("moneybutton.com", paymail.DefaultPort); err != nil {
		log.Fatal("error getting capabilities: " + err.Error())
	}
	log.Println("found capabilities: ", len(capabilities.Capabilities))

	// Extract the PublicProfile URL from the capabilities response
	publicProfileURL := capabilities.GetString(paymail.BRFCPublicProfile, "")

	// Get the public profile
	var profile *paymail.PublicProfile
	if profile, err = client.GetPublicProfile(publicProfileURL, "mrz", "moneybutton.com"); err != nil {
		log.Fatal("error getting profile: " + err.Error())
	}
	log.Printf("found profile: %s : %s", profile.Name, profile.Avatar)
}
