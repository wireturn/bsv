package main

import (
	"log"

	"github.com/tonicpow/go-paymail"
)

func main() {

	// Start with a paymail address
	paymailAddress := "MrZ@MoneyButton.com"

	// Validate the paymail address format
	if err := paymail.ValidatePaymail(paymailAddress); err != nil {
		log.Printf("error validating paymail: %s", err.Error())
	} else {
		log.Println("paymail format is valid!")
	}
}
