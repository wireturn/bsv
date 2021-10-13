package handcash

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
{
  "transactionId": "4eb7ab228ab9a23831b5b788e3f0eb5bed6dcdbb6d9d808eaba559c49afb9b0a",
  "note": "Test description",
  "type": "send",
  "time": 1608222315,
  "satoshiFees": 113,
  "satoshiAmount": 5301,
  "fiatExchangeRate": 188.6311183023935,
  "fiatCurrencyCode": "USD",
  "participants": [
    {
      "type": "user",
      "alias": "mrz@moneybutton.com",
      "displayName": "MrZ",
      "profilePictureUrl": "https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon",
      "responseNote": ""
    }
  ],
  "attachments": [],
  "appAction": "like"
}
*/

// GetPayment fetches a payment by transaction id
//
// Specs: https://github.com/HandCash/handcash-connect-sdk-js/blob/master/src/api/http_request_factory.js
func (c *Client) GetPayment(ctx context.Context, authToken,
	transactionID string) (*PaymentResponse, error) {

	// Make sure we have an auth token
	if len(authToken) == 0 {
		return nil, fmt.Errorf("missing auth token")
	} else if len(transactionID) == 0 {
		return nil, fmt.Errorf("missing transaction id")
	}

	// Get the signed request
	signed, err := c.getSignedRequest(
		http.MethodGet,
		endpointGetPaymentRequest,
		authToken,
		&PaymentRequest{TransactionID: transactionID},
		currentISOTimestamp(),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating signed request: %w", err)
	}

	// Convert into bytes
	var params []byte
	if params, err = json.Marshal(
		&PaymentRequest{TransactionID: transactionID},
	); err != nil {
		return nil, err
	}

	// Make the HTTP request
	response := httpRequest(
		ctx,
		c,
		&httpPayload{
			Data:           params,
			ExpectedStatus: http.StatusOK,
			Method:         signed.Method,
			URL:            signed.URI,
		},
		signed,
	)

	// Error in request?
	if response.Error != nil {
		return nil, response.Error
	}

	// Unmarshal pay response
	paymentResponse := new(PaymentResponse)
	if err = json.Unmarshal(response.BodyContents, &paymentResponse); err != nil {
		return nil, fmt.Errorf("failed unmarshal: %w", err)
	} else if paymentResponse == nil || paymentResponse.TransactionID == "" {
		return nil, fmt.Errorf("failed to find payment")
	}
	return paymentResponse, nil
}
