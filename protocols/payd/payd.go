package gopayd

import (
	"context"
)

// Transacter can be implemented to provide context based transactions.
type Transacter interface {
	WithTx(ctx context.Context) context.Context
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
