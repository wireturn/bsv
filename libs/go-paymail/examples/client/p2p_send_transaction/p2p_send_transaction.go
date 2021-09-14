package main

import (
	"log"

	"github.com/tonicpow/go-paymail"
)

func main() {

	// Load the client
	client, err := paymail.NewClient()
	if err != nil {
		log.Fatalf("error loading client: %s", err.Error())
	}

	// Get the capabilities
	// This is required first to get the corresponding P2P endpoint urls
	var capabilities *paymail.Capabilities
	if capabilities, err = client.GetCapabilities("moneybutton.com", paymail.DefaultPort); err != nil {
		log.Fatal("error getting capabilities: " + err.Error())
	}
	log.Println("found capabilities: ", len(capabilities.Capabilities))

	// Extract the URL from the capabilities response
	p2pDestinationURL := capabilities.GetString(paymail.BRFCP2PPaymentDestination, "")
	p2pSendURL := capabilities.GetString(paymail.BRFCP2PTransactions, "")

	// Create the basic paymentRequest to achieve a payment destination (how many sats are you planning to send?)
	paymentRequest := &paymail.PaymentRequest{Satoshis: 1000}

	// Get the P2P destination
	var destination *paymail.PaymentDestination
	destination, err = client.GetP2PPaymentDestination(p2pDestinationURL, "satchmo", "moneybutton.com", paymentRequest)
	if err != nil {
		log.Fatal("error getting destination: " + err.Error())
	}
	log.Printf("destination returned reference: %s and outputs: %d", destination.Reference, len(destination.Outputs))

	// Create a new P2P transaction
	rawTransaction := &paymail.P2PTransaction{
		Hex: "replace-with-raw-transaction-hex", // todo: replace with a real transaction
		MetaData: &paymail.P2PMetaData{
			Note:      "Thanks for dinner Satchmo!",
			Sender:    "mrz@moneybutton.com",
			PubKey:    "insert-pubkey-for-sender", // todo: replace with a real pubkey for the Sender
			Signature: "insert-signature-if-txid", // todo: replace with a real signature of the txid by the sender
		},
		Reference: destination.Reference,
	}

	// Send the P2P transaction
	var transaction *paymail.P2PTransactionResponse
	transaction, err = client.SendP2PTransaction(p2pSendURL, "satchmo", "moneybutton.com", rawTransaction)
	if err != nil {
		log.Fatal("error sending transaction: " + err.Error())
	}
	log.Printf("transaction sent: %s", transaction.TxID)
}
