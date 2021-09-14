package main

import (
	"log"

	"github.com/tonicpow/go-paymail"
)

func main() {

	// Start with a BRFC specification
	existingBRFC := &paymail.BRFCSpec{
		Author:  "MrZ",
		ID:      "e898079d7d1a",
		Title:   "New BRFC",
		Version: "1",
	}

	// Validate the BRFC ID
	if valid, id, err := existingBRFC.Validate(); err != nil {
		log.Fatalf("error validating BRFC id: %s", err.Error())
	} else if !valid {
		log.Fatalf("brfc is invalid: %s", id)
	} else if valid {
		log.Printf("brfc: %s is valid", existingBRFC.ID)
	}
}
