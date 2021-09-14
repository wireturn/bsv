package main

import (
	"log"

	"github.com/tonicpow/go-paymail"
)

func main() {

	// Start with a domain name
	domainName := "MoneyButton.com"

	// Validate the domain name
	if err := paymail.ValidateDomain(domainName); err != nil {
		log.Printf("error validating domain: %s", err.Error())
	} else {
		log.Println("domain format is valid!")
	}
}
