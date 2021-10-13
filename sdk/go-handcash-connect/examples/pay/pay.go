package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/tonicpow/go-handcash-connect"
)

func main() {

	// Create a new client (Beta ENV)
	client := handcash.NewClient(nil, nil, handcash.EnvironmentBeta)

	// Example attachment (Optional)
	attachmentPayload := `{"some":"data"}`
	dynamic := make(map[string]interface{})
	_ = json.Unmarshal([]byte(attachmentPayload), &dynamic)

	// Payment parameters
	params := &handcash.PayParameters{
		AppAction:   handcash.AppActionLike,
		Description: "Thanks dude!",
		Attachment: &handcash.Attachment{
			Format: handcash.AttachmentFormatJSON,
			Value:  dynamic,
		},
		Receivers: []*handcash.Payment{{
			Amount:       0.01,
			CurrencyCode: handcash.CurrencyUSD,
			To:           "mrz@moneybutton.com",
		}},
	}

	// Make a payment request
	payment, err := client.Pay(context.Background(), os.Getenv("AUTH_TOKEN"), params)
	if err != nil {
		log.Fatalln("error: ", err)
	}
	log.Println("payment: ", payment)
}
