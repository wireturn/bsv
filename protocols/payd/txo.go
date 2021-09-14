package gopayd

import (
	"context"
	"time"
)

// TxoCreate is used when creating outputs as part of a paymentRequest
// the script keys created are stored and later picked up
// and validated when the user sends a payment.
//
// These are partial txos and will be further hydrated when a transaction
// is sent spending them.
type TxoCreate struct {
	KeyName        string
	DerivationPath string
	LockingScript  string
	Satoshis       uint64
}

// UnspentTxoArgs are used to located an unfulfilled txo.
type UnspentTxoArgs struct {
	Keyname       string `db:"keyname"`
	LockingScript string `db:"lockingscript"`
	Satoshis      uint64 `db:"satoshis"`
}

// UnspentTxo is an unfulfilled txo not yet linked to a transaction.
type UnspentTxo struct {
	KeyName        string
	DerivationPath string
	LockingScript  string
	Satoshis       uint64
	CreatedAt      time.Time
	ModifiedAt     time.Time
}

// TxoWriter is used to add transaction information to a data store.
type TxoWriter interface {
	// TxoCreate will add a partial txo to a data store.
	TxoCreate(ctx context.Context, req TxoCreate) error
	// TxosCreate will add an array of partial txos to a data store.
	TxosCreate(ctx context.Context, req []*TxoCreate) error
}

// TxoReader is used to read tx information from a data store.
type TxoReader interface {
	// PartialTxo will return a txo that has not tet been assigned to a transaction.
	PartialTxo(ctx context.Context, args UnspentTxoArgs) (*UnspentTxo, error)
}
