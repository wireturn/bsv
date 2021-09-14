package gopayd

import (
	"context"
)

// OutputsCreate will be used to create n outputs that match
// the required satoshis and are split into the set denomination.
// Each output should have a different locking script.
type OutputsCreate struct {
	Satoshis     uint64
	Denomination uint64
}

// Output message used in BIP270.
// See https://github.com/moneybutton/bips/blob/master/bip-0270.mediawiki#output
type Output struct {
	// Amount is the number of satoshis to be paid.
	Amount uint64 `json:"amount"`
	// Script is a locking script where payment should be sent, formatted as a hexadecimal string.
	Script string `json:"script"`
	// Description, an optional description such as "tip" or "sales tax". Maximum length is 100 chars.
	Description string `json:"description"`
}

// PaymentRequestOutputer will create outputs that equal the amount of request satoshis.
type PaymentRequestOutputer interface {
	CreateOutputs(ctx context.Context, req OutputsCreate) ([]*Output, error)
}
