package spv

import (
	"context"

	"github.com/libsv/go-bt/v2"
)

// VerifyPayment verifies whether or not the txs supplied via the supplied spv.Envelope are valid
func (v *verifier) VerifyPayment(ctx context.Context, initialPayment *Envelope) (bool, error) {
	if initialPayment == nil {
		return false, ErrNilInitialPayment
	}

	// The tip tx is the transaction we're trying to verify, and it should not have a supplied
	// Merkle Proof.
	if initialPayment.IsAnchored() {
		return false, ErrTipTxConfirmed
	}

	valid, err := v.verifyTxs(ctx, initialPayment)
	if err != nil {
		return false, err
	}

	return valid, nil
}

func (v *verifier) verifyTxs(ctx context.Context, payment *Envelope) (bool, error) {
	// If at the beginning or middle of the tx chain and tx is unconfirmed, fail and error.
	if !payment.IsAnchored() && (payment.Parents == nil || len(payment.Parents) == 0) {
		return false, ErrNoConfirmedTransaction
	}

	// Recurse back to the anchor transactions of the transaction chain and verify forward towards
	// the tip transaction. This way, we check that the first transactions in the chain are anchored
	// to the blockchain through a valid Merkle Proof.
	for parentTxID, parent := range payment.Parents {
		if parent.TxID == "" {
			parent.TxID = parentTxID
		}

		valid, err := v.verifyTxs(ctx, parent)
		if err != nil {
			return false, err
		}
		if !valid {
			return false, nil
		}
	}

	// If a Merkle Proof is provided, assume we are at the anchor/beginning of the tx chain.
	// Verify and return the result.
	if payment.IsAnchored() {
		return v.verifyTxAnchor(ctx, payment)
	}

	tx, err := bt.NewTxFromString(payment.RawTx)
	if err != nil {
		return false, err
	}

	// We must verify the tx or else we can not know if any of it's child txs are valid.
	return v.verifyUnconfirmedTx(tx, payment)
}

func (v *verifier) verifyTxAnchor(ctx context.Context, payment *Envelope) (bool, error) {
	proofTxID := payment.Proof.TxOrID
	if len(proofTxID) != 64 {
		proofTx, err := bt.NewTxFromString(payment.Proof.TxOrID)
		if err != nil {
			return false, err
		}

		proofTxID = proofTx.TxID()
	}

	// If the txid of the Merkle Proof doesn't match the txid provided in the spv.Envelope,
	// fail and error
	if proofTxID != payment.TxID {
		return false, ErrTxIDMismatch
	}

	valid, _, err := v.VerifyMerkleProofJSON(ctx, payment.Proof)
	if err != nil {
		return false, err
	}

	return valid, nil
}

func (v *verifier) verifyUnconfirmedTx(tx *bt.Tx, payment *Envelope) (bool, error) {
	// If no tx inputs have been provided, fail and error
	if len(tx.Inputs) == 0 {
		return false, ErrNoTxInputsToVerify
	}

	for _, input := range tx.Inputs {
		parent, ok := payment.Parents[input.PreviousTxIDStr()]
		if !ok {
			return false, ErrNotAllInputsSupplied
		}

		parentTx, err := bt.NewTxFromString(parent.RawTx)
		if err != nil {
			return false, err
		}

		// If the input is indexing an output that is out of bounds, fail and error
		if int(input.PreviousTxOutIndex) > len(parentTx.Outputs)-1 {
			return false, ErrInputRefsOutOfBoundsOutput
		}

		output := parentTx.Outputs[int(input.PreviousTxOutIndex)]

		// TODO: verify script using input and previous output
		_ = output
	}

	return true, nil
}
