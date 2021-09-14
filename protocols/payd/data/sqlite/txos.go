package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/theflyingcodr/lathos/errs"

	gopayd "github.com/libsv/payd"
)

const (
	sqlTxoCreate = `
	INSERT INTO txos (keyname, derivationpath, lockingscript, satoshis)
	VALUES(:keyname, :derivationpath, :lockingscript, :satoshis)
	`

	sqlPartialTxo = `
	SELECT keyname, derivationpath, lockingscript, satoshis, createdat, modifiedat
	FROM txos
	WHERE lockingscript = $1 AND satoshis = $2 AND keyname = $3 AND outpoint IS NULL
	`

	sqlTxoUpdate = `
	UPDATE txos SET outpoint = :outpoint, vout = :vout, txid = :txid, modifiedat = DATETIME('now')
	WHERE outpoint IS NULL AND lockingscript = :lockingscript AND keyname = :keyname AND satoshis = :satoshis
	`
)

// TxoCreate will store a txo created during payment requests.
func (s *sqliteStore) TxoCreate(ctx context.Context, req gopayd.TxoCreate) error {
	return s.TxosCreate(ctx, []*gopayd.TxoCreate{
		&req,
	})
}

// TxosCreate will store txos created during payment requests.
func (s *sqliteStore) TxosCreate(ctx context.Context, req []*gopayd.TxoCreate) error {
	tx, err := s.newTx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create context when creating a txo")
	}
	defer func() {
		_ = rollback(ctx, tx)
	}()
	if err := handleNamedExec(tx, sqlTxoCreate, req); err != nil {
		return errors.Wrap(err, "failed to insert script keys.")
	}
	return errors.Wrap(commit(ctx, tx), "failed to commit transaction when creating txos.")
}

// PartialTxo will return a txo that has been stored but not yet assigned to a transaction.
func (s *sqliteStore) PartialTxo(ctx context.Context, args gopayd.UnspentTxoArgs) (*gopayd.UnspentTxo, error) {
	var txo gopayd.UnspentTxo
	if err := s.db.GetContext(ctx, &txo, sqlPartialTxo, args.LockingScript, args.Satoshis, args.Keyname); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewErrNotFound("N104",
				fmt.Sprintf("unable to find txo with script '%s', value '%d' and keyname '%s'",
					args.LockingScript, args.Satoshis, args.Keyname))
		}
		return nil, errors.Wrap(err, "failed to read partialTxo")
	}
	return &txo, nil
}
