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

	// Check the SSL certificate
	var valid bool
	if valid, err = client.CheckSSL("moneybutton.com"); err != nil {
		log.Fatal("error getting SSL certificate: " + err.Error())
	} else if !valid {
		log.Fatal("SSL certificate validation failed")
	}
	log.Println("valid SSL certificate found for:", "moneybutton.com")
}
