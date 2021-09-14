package sqlite

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	gopayd "github.com/libsv/payd"
)

const (
	sqlBalance = `
	SELECT SUM(satoshis) as satoshis
	FROM txos
	WHERE spentat IS NULL 
	`
)

// Balance will return the current account balance.
func (s *sqliteStore) Balance(ctx context.Context) (*gopayd.Balance, error) {
	var resp gopayd.Balance
	if err := s.db.GetContext(ctx, &resp, sqlBalance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &gopayd.Balance{Satoshis: 0}, nil
		}
		return nil, errors.Wrap(err, "failed to get balance")
	}
	return &resp, nil
}
