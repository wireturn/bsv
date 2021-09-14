package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/tokenized/pkg/bitcoin"
	"github.com/tokenized/pkg/storage"
	"github.com/tokenized/spynode/pkg/client"
)

const (
	txStatePath = "spynode/txs/state"
)

// SaveTxState saves a tx state to storage.
func SaveTxState(ctx context.Context, store storage.Storage, tx *client.Tx) error {
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return errors.Wrap(err, "serialize")
	}

	path := fmt.Sprintf("%s/%s", txStatePath, tx.Tx.TxHash())
	if err := store.Write(ctx, path, buf.Bytes(), nil); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

// FetchTxState fetches a tx state from storage.
func FetchTxState(ctx context.Context, store storage.Storage,
	txid bitcoin.Hash32) (*client.Tx, error) {

	path := fmt.Sprintf("%s/%s", txStatePath, txid)
	b, err := store.Read(ctx, path)
	if err != nil {
		return nil, errors.Wrap(err, "read")
	}

	var result client.Tx
	if err := result.Deserialize(bytes.NewReader(b)); err != nil {
		return nil, errors.Wrap(err, "deserialize")
	}

	return &result, nil
}
