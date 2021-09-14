package gopayd

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

// Transaction defines a single transaction.
type Transaction struct {
	PaymentID string    `db:"paymentid"`
	TxID      string    `db:"txid"`
	TxHex     string    `db:"txhex"`
	CreatedAt time.Time `db:"createdat"`
	Outputs   []Txo     `db:"-"`
}

// Txo defines a single txo and can be returned from the data store.
type Txo struct {
	Outpoint       string      `db:"outpoint"`
	TxID           string      `db:"txid"`
	Vout           int         `db:"vout"`
	KeyName        null.String `db:"keyname"`
	DerivationPath null.String `db:"derivationpath"`
	LockingScript  string      `db:"lockingscript"`
	Satoshis       uint64      `db:"satoshis"`
	SpentAt        null.Time   `db:"spentat"`
	SpendingTxID   null.String `db:"spendingtxid"`
	CreatedAt      time.Time   `db:"createdat"`
	ModifiedAt     time.Time   `db:"modifiedat"`
}

// CreateTransaction is used to insert a tx into the data store.
type CreateTransaction struct {
	PaymentID string       `db:"paymentID"`
	TxID      string       `db:"txid"`
	TxHex     string       `db:"txhex"`
	Outputs   []*UpdateTxo `db:"-"`
}

// UpdateTxo is used to update a single txo in the data store.
type UpdateTxo struct {
	Outpoint       string      `db:"outpoint"`
	TxID           string      `db:"txid"`
	Vout           int         `db:"vout"`
	KeyName        null.String `db:"keyname"`
	DerivationPath null.String `db:"derivationpath"`
	LockingScript  string      `db:"lockingscript"`
	Satoshis       uint64      `db:"satoshis"`
}

// SpendTxo can be used to update a transaction out with information
// on when it was spent and by what transaction.
type SpendTxo struct {
	SpentAt      *time.Time
	SpendingTxID string
}

// SpendTxoArgs are used to identify the transaction output to mark as spent.
type SpendTxoArgs struct {
	Outpoint string
}

// TxoArgs is used to get a single txo.
type TxoArgs struct {
	Outpoint string
}
