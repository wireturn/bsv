package ppctl

import (
	"context"

	"github.com/pkg/errors"

	gopayd "github.com/libsv/payd"
)

type txStatusService struct {
	rdr gopayd.TxStatusReader
}

// NewTxStatusService will setup and return a new TxStatusService for retrieving Tx status information
// and enforcing business rules around this.
func NewTxStatusService(rdr gopayd.TxStatusReader) *txStatusService {
	return &txStatusService{rdr: rdr}
}

// Status will return the current broadcast status for a transaction.
func (t *txStatusService) Status(ctx context.Context, args gopayd.TxStatusArgs) (*gopayd.TxStatus, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	resp, err := t.rdr.Status(ctx, args)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get status for transaction %s", args.TxID)
	}
	return resp, nil
}
