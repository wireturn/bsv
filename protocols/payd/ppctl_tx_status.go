package gopayd

import (
	"context"

	validator "github.com/theflyingcodr/govalidator"
)

// TxStatusArgs is used to locate a specific transaction in order to get its status.
type TxStatusArgs struct {
	TxID string `param:"txId"`
}

// Validate will verify the TxStatusArgs are correct.
func (t TxStatusArgs) Validate() error {
	return validator.New().
		Validate("txId", validator.Length(t.TxID, 64, 64), validator.IsHex(t.TxID)).
		Err()
}

// TxStatus will return the current broadcast status of a transaction.
type TxStatus struct {
	TxID          string `json:"txId"`
	Status        string `json:"status"`
	BlockHash     string `json:"blockHash"`
	BlockHeight   uint64 `json:"blockHeight"`
	Confirmations uint64 `json:"confirmations"`
	Error         int    `json:"error,omitempty"`
}

// TxStatusService is used to enforce business rules / emit events etc in response to
// a txStatus request.
type TxStatusService interface {
	// Status will return the broadcast status of a transaction.
	Status(ctx context.Context, args TxStatusArgs) (*TxStatus, error)
}

// TxStatusReader can be implemented to lookup broadcast status for a transaction,
// this could be interfacing with a node, mAPI or some other system to get the tx broadcast status.
type TxStatusReader interface {
	// Status will return the broadcast status of a transaction.
	Status(ctx context.Context, args TxStatusArgs) (*TxStatus, error)
}
