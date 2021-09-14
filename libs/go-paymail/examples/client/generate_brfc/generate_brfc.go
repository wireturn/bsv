package main

import (
	"log"

	"github.com/tonicpow/go-paymail"
)

func main() {

	// Start with a new BRFC specification
	newBRFC := &paymail.BRFCSpec{
		Author:  "MrZ",
		Title:   "New BRFC",
		Version: "1",
	}

	// Generate the BRFC ID
	if err := newBRFC.Generate(); err != nil {
		log.Fatalf("error generating BRFC id: %s", err.Error())
	}
	log.Printf("id generated: %s", newBRFC.ID)
}
