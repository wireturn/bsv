package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	lathos "github.com/theflyingcodr/lathos/errs"
)

type sqliteStore struct {
	db *sqlx.DB
}

// NewSQLiteStore will setup and return a sql list data store.
func NewSQLiteStore(db *sqlx.DB) *sqliteStore {
	return &sqliteStore{db: db}
}

func (s *sqliteStore) newTx(ctx context.Context) (*sqlx.Tx, error) {
	ctxx := TxFromContext(ctx)
	if ctxx != nil {
		if ctxx.Tx == nil {
			t, err := s.db.BeginTxx(ctx, nil)
			if err != nil {
				return nil, err
			}
			ctxx.Tx = t
		}
		return ctxx.Tx, nil
	}
	return s.db.BeginTxx(ctx, nil)
}

// commit a transaction, if there is a context based tx
// this will not commit - we wait on the context to close it.
func commit(ctx context.Context, tx *sqlx.Tx) error {
	ctxx := TxFromContext(ctx)
	if ctxx != nil {
		if ctxx.Tx != nil {
			return nil
		}
	}
	return tx.Commit()
}

// rollback a transaction, if there is a context based tx
// this will not rollback - we wait on the context to close it.
func rollback(ctx context.Context, tx *sqlx.Tx) error {
	ctxx := TxFromContext(ctx)
	if ctxx != nil {
		if ctxx.Tx != nil {
			return nil
		}
	}
	return tx.Rollback()
}

func handleNamedExec(tx db, sql string, args interface{}) error {
	res, err := tx.NamedExec(sql, args)
	if err != nil {
		return errors.Wrap(err, "failed to run exec")
	}
	return handleExecRows(res)
}

func handleExecRows(res sql.Result) error {
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to read rows affected")
	}
	if ra <= 0 {
		return errors.Wrap(err, "exec did not affect rows")
	}
	return nil
}

// nolint:deadcode,unused // wip
func dbErr(err error, errCode, message string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return lathos.NewErrNotFound(errCode, message)
	}
	return errors.WithMessage(err, message)
}

// nolint:deadcode,unused // wip
func dbErrf(err error, errCode, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return lathos.NewErrNotFound(errCode, fmt.Sprintf(format, args...))
	}
	return errors.WithMessage(err, fmt.Sprintf(format, args...))
}

type db interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

type execKey int

// nolint:gochecknoglobals // this variable is fine as it's used for context & is private
var exec execKey

// Tx wraps the transaction used in context.
type Tx struct {
	*sqlx.Tx
}

// WithTxContext will add a new empty transaction to the provided context.
func WithTxContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, exec, &Tx{})
}

// TxFromContext will return a context based transaction if found.
func TxFromContext(ctx context.Context) *Tx {
	if tx, ok := ctx.Value(exec).(*Tx); ok {
		return tx
	}
	return nil
}

// Transacter is used to implement the DBTransacter interface for
// managing db transactions in other layers.
type Transacter struct {
}

// WithTx will add a TX to the provided context.
func (t *Transacter) WithTx(ctx context.Context) context.Context {
	return WithTxContext(ctx)
}

// Commit will commit a distributed transaction.
func (t *Transacter) Commit(ctx context.Context) error {
	tx := TxFromContext(ctx)
	if tx.Tx != nil {
		return tx.Tx.Commit()
	}
	return nil
}

// Rollback will revert the tx and rollback any changes.
func (t *Transacter) Rollback(ctx context.Context) error {
	tx := TxFromContext(ctx)
	if tx.Tx != nil {
		return tx.Tx.Rollback()
	}
	return nil
}
