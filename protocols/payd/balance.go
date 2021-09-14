package gopayd

import (
	"context"
)

// Balance contains the current balance of unspent outputs.
type Balance struct {
	Satoshis uint64 `json:"satoshis" db:"satoshis"`
}

// BalanceService is used to enforce balance buiness rules.
type BalanceService interface {
	Balance(ctx context.Context) (*Balance, error)
}

// BalanceReader is used to read balance info from a datastore.
type BalanceReader interface {
	Balance(ctx context.Context) (*Balance, error)
}
