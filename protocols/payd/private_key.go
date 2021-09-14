package gopayd

import (
	"context"
	"time"

	"github.com/libsv/go-bk/bip32"
)

// PrivateKey describes a named private key.
type PrivateKey struct {
	// Name of the private key.
	Name string `db:"name"`
	// Xprv is the private key.
	Xprv string `db:"xprv"`
	// CreatedAt is the date/time when the key was stored.
	CreatedAt time.Time `db:"createdAt"`
}

// KeyArgs defines all arguments required to get a key.
type KeyArgs struct {
	// Name is the name of the key to return.
	Name string `db:"name"`
}

// PrivateKeyService can be implemented to get and create PrivateKeys.
type PrivateKeyService interface {
	// Create will create a new private key if it doesn't exist already.
	Create(ctx context.Context, keyName string) error
	// PrivateKey will return a private key.
	PrivateKey(ctx context.Context, keyName string) (*bip32.ExtendedKey, error)
}

// PrivateKeyReader reads private info from a data store.
type PrivateKeyReader interface {
	// PrivateKey can be used to return an existing private key.
	PrivateKey(ctx context.Context, args KeyArgs) (*PrivateKey, error)
}

// PrivateKeyWriter will add private key to the datastore.
type PrivateKeyWriter interface {
	// PrivateKeyCreate will add a new private key to the data store.
	PrivateKeyCreate(ctx context.Context, req PrivateKey) (*PrivateKey, error)
}

// PrivateKeyReaderWriter describes a data store that can be implemented to get and store private keys.
type PrivateKeyReaderWriter interface {
	PrivateKeyReader
	PrivateKeyWriter
}
