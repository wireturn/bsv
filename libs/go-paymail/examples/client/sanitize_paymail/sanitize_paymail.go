package main

import (
	"log"

	"github.com/tonicpow/go-paymail"
)

func main() {

	// Start with a paymail address
	paymailAddress := "MrZ@MoneyButton.com"

	// Sanitize the address, extract the parts
	alias, domain, address := paymail.SanitizePaymail(paymailAddress)
	log.Printf("alias: %s domain: %s address: %s", alias, domain, address)
}
