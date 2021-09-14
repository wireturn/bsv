package main

import (
	"log"

	"github.com/tonicpow/go-paymail"
)

func main() {

	// Start with a handle
	handle := "$Mr-Z"

	// Convert the handle to paymail address
	address := paymail.ConvertHandle(handle, false)
	log.Printf("handle %s was converted to: %s", handle, address)

	// Try another handle
	handle = "1MrZ"

	// Convert the handle to paymail address
	address = paymail.ConvertHandle(handle, false)
	log.Printf("handle %s was converted to: %s", handle, address)
}
