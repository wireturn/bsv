package sqlite

import (
	"context"

	gopayd "github.com/libsv/payd"
	"github.com/pkg/errors"
)

func (s *sqliteStore) StoreUtxos(ctx context.Context, req gopayd.CreateTransaction) (*gopayd.Transaction, error) {
	tx, err := s.newTx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start transaction when inserting transaction to db")
	}
	defer func() {
		_ = rollback(ctx, tx)
	}()
	resp, err := s.txCreateTransaction(tx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transaction and utxos")
	}
	return resp, errors.Wrapf(commit(ctx, tx),
		"failed to commit transaction when adding tx and outputs for paymentID %s", req.PaymentID)
}

// txCreateTransaction takes a db object / transaction and adds a transaction to the data store
// along with utxos, returning the transaction.
// This method can be used with other methods in the store allowing
// multiple methods to be ran in the same db transaction.
func (s *sqliteStore) txCreateTransaction(tx db, req gopayd.CreateTransaction) (*gopayd.Transaction, error) {
	// insert tx and utxos
	if err := handleNamedExec(tx, sqlTransactionCreate, req); err != nil {
		return nil, errors.Wrap(err, "failed to insert new transaction")
	}
	if err := handleNamedExec(tx, sqlTxoUpdate, req.Outputs); err != nil {
		return nil, errors.Wrap(err, "failed to insert transaction outputs")
	}
	var outTx gopayd.Transaction
	if err := tx.Get(&outTx, sqlTransactionByID, req.TxID); err != nil {
		return nil, errors.Wrapf(err, "failed to get stored transaction for paymentID %s", req.PaymentID)
	}
	var outTxos []gopayd.Txo
	if err := tx.Select(&outTxos, sqlTxosByTxID, req.TxID); err != nil {
		return nil, errors.Wrapf(err, "failed to get stored transaction outputs for paymentID %s", req.PaymentID)
	}
	outTx.Outputs = outTxos
	return &outTx, nil
}
