package gopayd

import (
	"context"
)

// DerivationExistsArgs are used to check a derivation path exists for a specific
// master key and path.
type DerivationExistsArgs struct {
	KeyName string `db:"keyName"`
	Path    string `db:"derivationPath"`
}

// DerivationReader can be used to read derivation path data from a data store.
type DerivationReader interface {
	DerivationPathExists(ctx context.Context, args DerivationExistsArgs) (bool, error)
}
