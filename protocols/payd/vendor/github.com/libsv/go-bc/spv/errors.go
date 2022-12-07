package spv

import "github.com/pkg/errors"

var (
	// ErrNoTxInputs returns if an envelope is attempted to be created from a transaction that has no inputs
	ErrNoTxInputs = errors.New("provided tx has no inputs to build envelope from")

	// ErrPaymentNotVerified returns if a transaction in the tree provided was missed during verification
	ErrPaymentNotVerified = errors.New("a tx was missed during validation")

	// ErrTipTxConfirmed returns if the tip transaction is already confirmed
	ErrTipTxConfirmed = errors.New("tip transaction must be unconfirmed")

	// ErrNoConfirmedTransaction returns if a path from tip to beginning/anchor contains no confirmed transcation
	ErrNoConfirmedTransaction = errors.New("not confirmed/anchored tx(s) provided")

	// ErrTxIDMismatch returns if they key value pair of a transactions input has a mismatch in txID
	ErrTxIDMismatch = errors.New("input and proof ID mismatch")

	// ErrNotAllInputsSupplied returns if an unconfirmed transaction in envelope contains inputs which are not
	// present in the parent envelope
	ErrNotAllInputsSupplied = errors.New("a tx input missing in parent envelope")

	// ErrNoTxInputsToVerify returns if a transaction has no inputs
	ErrNoTxInputsToVerify = errors.New("a tx has no inputs to verify")

	// ErrNilInitialPayment returns if a transaction has no inputs
	ErrNilInitialPayment = errors.New("initial payment cannot be nil")

	// ErrInputRefsOutOfBoundsOutput returns if a transaction has no inputs
	ErrInputRefsOutOfBoundsOutput = errors.New("tx input index into output is out of bounds")
)
