package service

import (
	"context"

	"github.com/pkg/errors"

	gopayd "github.com/libsv/payd"
)

type balance struct {
	store gopayd.BalanceReader
}

// NewBalance will setup and return the current balance of the wallet.
func NewBalance(store gopayd.BalanceReader) *balance {
	return &balance{store: store}
}

// Balance will return the current wallet balance.
func (b *balance) Balance(ctx context.Context) (*gopayd.Balance, error) {
	resp, err := b.store.Balance(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get balance")
	}
	return resp, nil
}
