package main

import (
	"log"
	"os"

	"github.com/dotwallet/dotwallet-go-sdk"
)

// For more information: https://developers.dotwallet.com/documents/en/#application-authorization
func main() {

	// Create the DotWallet client
	c, err := dotwallet.NewClient(
		dotwallet.WithCredentials(
			os.Getenv("DOT_WALLET_CLIENT_ID"),
			os.Getenv("DOT_WALLET_CLIENT_SECRET"),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Update the application token (get or update)
	if err = c.UpdateApplicationAccessToken(); err != nil {
		log.Fatalln(err)
	}

	// Show that we got an application token
	t := c.Token()
	log.Println(
		"token:", t.AccessToken,
		"type:", t.TokenType,
		"expires:", t.ExpiresAt,
	)
}
