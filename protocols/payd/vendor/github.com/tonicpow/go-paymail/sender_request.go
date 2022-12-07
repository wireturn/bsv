package paymail

import (
	"fmt"

	"github.com/bitcoinschema/go-bitcoin"
)

/*
Example:
{
    "senderName": "FirstName LastName",
    "senderHandle": "<alias>@<domain.tld>",
    "dt": "2013-10-21T13:28:06.419Z",
    "amount": 550,
    "purpose": "message to receiver",
    "signature": "<compact Bitcoin message signature>"
}
*/

// SenderRequest is the request body for the basic address resolution
//
// This is required to make a basic resolution request, and Dt and SenderHandle are required
type SenderRequest struct {
	Amount       uint64 `json:"amount,omitempty"`     // The amount, in Satoshis, that the sender intends to transfer to the receiver
	Dt           string `json:"dt"`                   // (required) ISO-8601 formatted timestamp; see notes
	Purpose      string `json:"purpose,omitempty"`    // Human-readable description of the purpose of the payment
	SenderHandle string `json:"senderHandle"`         // (required) Sender's paymail handle
	SenderName   string `json:"senderName,omitempty"` // Human-readable sender display name
	Signature    string `json:"signature,omitempty"`  // Compact Bitcoin message signature; http://bsvalias.org/04-01-basic-address-resolution.html#signature-field
}

// Verify will verify the given components in the ResolveAddress() request
//
// Source: https://github.com/moneybutton/paymail-client/blob/master/src/VerifiableMessage.js
// Specs: http://bsvalias.org/04-01-basic-address-resolution.html#signature-field
func (s *SenderRequest) Verify(keyAddress, signature string) error {

	// Basic checks before trying the signature verification
	if len(keyAddress) == 0 {
		return fmt.Errorf("missing key address")
	} else if len(signature) == 0 {
		return fmt.Errorf("missing a signature to verify")
	}

	// Concatenate & verify the message
	return bitcoin.VerifyMessage(keyAddress, signature, fmt.Sprintf("%s%d%s%s", s.SenderHandle, s.Amount, s.Dt, s.Purpose))
}

// Sign will sign the given components in the ResolveAddress() request
//
// Source: https://github.com/moneybutton/paymail-client/blob/master/src/VerifiableMessage.js
// Specs: http://bsvalias.org/04-01-basic-address-resolution.html#signature-field
// Additional Specs: http://bsvalias.org/04-02-sender-validation.html
func (s *SenderRequest) Sign(privateKey string) (string, error) {

	// Basic checks before trying to sign the request
	if len(privateKey) == 0 {
		return "", fmt.Errorf("missing private key")
	} else if len(s.Dt) == 0 {
		return "", fmt.Errorf("missing dt")
	} else if len(s.SenderHandle) == 0 {
		return "", fmt.Errorf("missing senderHandle")
	}

	// Concatenate & sign message
	return bitcoin.SignMessage(
		privateKey,
		fmt.Sprintf("%s%d%s%s", s.SenderHandle, s.Amount, s.Dt, s.Purpose),
		false,
	)
}
