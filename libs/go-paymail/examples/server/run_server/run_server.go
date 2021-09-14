package main

import "github.com/tonicpow/go-paymail/server"

func main() {

	// Run the server on port 3000 and timeout requests after 15 seconds
	server.Start(3000, 15)
}
