package gopayd

import (
	"context"

	"github.com/libsv/go-bt/v2"
)

// FeeReader can be implemented to read fees.
type FeeReader interface {
	// Fees will return fees from a datastore.
	Fees(ctx context.Context) (*bt.FeeQuote, error)
}
