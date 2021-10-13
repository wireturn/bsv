package main

import (
	"log"
	"os"

	"github.com/dotwallet/dotwallet-go-sdk"
)

// For more information: https://developers.dotwallet.com/documents/en/#get-user-receive-address
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

	// Get the receiving address for a known user_id
	var wallets *dotwallet.Wallets
	if wallets, err = c.UserReceiveAddress(
		os.Getenv("DOT_WALLET_USER_ID"),
		dotwallet.CoinTypeBSV,
	); err != nil {
		log.Fatalln(err)
	}

	// Show the wallet info
	log.Println(
		"primary wallet:", wallets.PrimaryWallet.Address,
		"paymail:", wallets.PrimaryWallet.Paymail,
	)
}
