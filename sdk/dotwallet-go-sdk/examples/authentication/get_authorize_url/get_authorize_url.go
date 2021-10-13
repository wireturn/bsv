package main

import (
	"log"
	"os"

	"github.com/dotwallet/dotwallet-go-sdk"
)

// For more information: https://developers.dotwallet.com/documents/en/#user-authorization
func main() {

	// Create the DotWallet client
	c, err := dotwallet.NewClient(
		dotwallet.WithCredentials(
			os.Getenv("DOT_WALLET_CLIENT_ID"),
			os.Getenv("DOT_WALLET_CLIENT_SECRET"),
		),
		dotwallet.WithRedirectURI(os.Getenv("DOT_WALLET_REDIRECT_URI")),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Create a state UUID for checking later after the oauth2 callback
	state, _ := c.NewState()

	// Show the url (used for oauth2 user login)
	log.Println("url:", c.GetAuthorizeURL(state, []string{dotwallet.ScopeUserInfo}))
}
