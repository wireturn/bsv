package main

import (
	"log"
	"os"
	"time"

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

	// Previously, we got a "code" from the authorize_url callback to the redirect_uri
	var token *dotwallet.DotAccessToken
	if token, err = c.GetUserToken(os.Getenv("DOT_WALLET_CODE")); err != nil {
		log.Fatalln(err)
	}

	// Show that we got an application token
	log.Println(
		"token:", token.AccessToken,
		"type:", token.TokenType,
		"expires:", token.ExpiresAt,
		"refresh:", token.RefreshToken,
	)

	// Sleep for a few seconds
	log.Println("waiting 5 seconds...")
	time.Sleep(5 * time.Second)

	// Trigger a refresh of the token
	if token, err = c.RefreshUserToken(token); err != nil {
		log.Fatalln(err)
	}

	// Show that we got an application token
	log.Println(
		"token:", token.AccessToken,
		"type:", token.TokenType,
		"expires:", token.ExpiresAt,
		"refresh:", token.RefreshToken,
	)
}
