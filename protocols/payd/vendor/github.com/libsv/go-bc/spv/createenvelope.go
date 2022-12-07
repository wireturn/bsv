package spv

import (
	"context"
	"fmt"

	"github.com/libsv/go-bt/v2"
	"github.com/pkg/errors"
)

// CreateEnvelope builds and returns an spv.Envelope for the provided tx.
func (c *creator) CreateEnvelope(ctx context.Context, tx *bt.Tx) (*Envelope, error) {
	if len(tx.Inputs) == 0 {
		return nil, ErrNoTxInputs
	}

	envelope := &Envelope{
		TxID:    tx.TxID(),
		RawTx:   tx.String(),
		Parents: make(map[string]*Envelope),
	}

	for _, input := range tx.Inputs {
		pTxID := input.PreviousTxIDStr()

		// If we already have added the tx to the parent envelope, there's no point in
		// redoing the same work
		if _, ok := envelope.Parents[pTxID]; ok {
			continue
		}

		// Check the store for a Merkle Proof for the current input.
		mp, err := c.mpc.MerkleProof(ctx, pTxID)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get merkle proof for tx %s", pTxID)
		}
		if mp != nil {
			// If a Merkle Proof exists, build and return an spv.Envelope
			envelope.Parents[pTxID] = &Envelope{
				TxID:  pTxID,
				Proof: mp,
			}

			// Skip getting the tx data as we have everything we need for verifying the current tx.
			continue
		}

		// If no merkle proof was found for the input, build a *bt.Tx from its TxID and recursively
		// call this function building envelopes for inputs without proofs, until a parent with a
		// Merkle Proof is found.
		pTx, err := c.txc.Tx(ctx, pTxID)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get tx %s", pTxID)
		}
		if pTx == nil {
			return nil, fmt.Errorf("could not find tx %s", pTxID)
		}

		pEnvelope, err := c.CreateEnvelope(ctx, pTx)
		if err != nil {
			return nil, err
		}

		envelope.Parents[pTxID] = pEnvelope
	}

	return envelope, nil
}
