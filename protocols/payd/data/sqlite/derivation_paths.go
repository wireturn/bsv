package sqlite

import (
	"context"

	"github.com/pkg/errors"

	gopayd "github.com/libsv/payd"
)

const (
	sqlDerivationPathExists = `
	SELECT EXISTS(
	    SELECT derivationpath FROM txos WHERE derivationpath = $1 AND keyname = $2 
	    )
	`
)

// DerivationPathExists will return true / false if the supplied derivation path exists or not.
func (s *sqliteStore) DerivationPathExists(ctx context.Context, args gopayd.DerivationExistsArgs) (bool, error) {
	var exists int
	if err := s.db.GetContext(ctx, &exists, sqlDerivationPathExists, args.Path, args.KeyName); err != nil {
		return false, errors.WithStack(err)
	}
	return exists == 1, nil
}
