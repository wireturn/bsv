package ppctl

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/libsv/go-bc/spv"
	"github.com/libsv/go-bt/v2"
	"github.com/pkg/errors"
	validator "github.com/theflyingcodr/govalidator"
	"github.com/theflyingcodr/lathos"
	"github.com/theflyingcodr/lathos/errs"
	"gopkg.in/guregu/null.v3"

	gopayd "github.com/libsv/payd"
	"github.com/libsv/payd/config"
	"github.com/libsv/payd/errcodes"
)

// payment is a layer on top of the payment services of which we currently support:
// * wallet payments, that are handled by the wallet and transmitted to the network
// * paymail payments, that use the paymail protocol for making the payments.
type payment struct {
	cfg       *config.Wallet
	store     gopayd.PaymentWriter
	txoRdr    gopayd.TxoReader
	invStore  gopayd.InvoiceReaderWriter
	sender    gopayd.PaymentSender
	txrunner  gopayd.Transacter
	envVerify spv.PaymentVerifier
}

// NewPayment will create and return a new payment service.
func NewPayment(cfg *config.Wallet, store gopayd.PaymentWriter, txoRdr gopayd.TxoReader, invStore gopayd.InvoiceReaderWriter, sender gopayd.PaymentSender, txrunner gopayd.Transacter, envVerify spv.PaymentVerifier) *payment {
	return &payment{
		cfg:       cfg,
		store:     store,
		txoRdr:    txoRdr,
		invStore:  invStore,
		sender:    sender,
		txrunner:  txrunner,
		envVerify: envVerify,
	}
}

// CreatePayment will setup a new payment and return the result.
func (p *payment) CreatePayment(ctx context.Context, args gopayd.CreatePaymentArgs, req gopayd.CreatePayment) (*gopayd.PaymentACK, error) {
	if err := validator.New().Validate("paymentID", validator.NotEmpty(args.PaymentID)); err.Err() != nil {
		return nil, err
	}

	// Validate the SPV Envelope before processing transaction
	ok, err := p.envVerify.VerifyPayment(ctx, req.SPVEnvelope)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to validate spv envelope for paymentID %s", args.PaymentID)
	}
	if !ok {
		return nil, errors.WithStack(errors.New("failed to verify merkle proof in spv envelope"))
	}

	// get the invoice for the paymentID to check total satoshis required.
	inv, err := p.invStore.Invoice(ctx, gopayd.InvoiceArgs{PaymentID: args.PaymentID})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get invoice to validate output total for paymentID %s.", args.PaymentID)
	}
	if !inv.PaymentReceivedAt.IsZero() {
		return nil, errs.NewErrDuplicate(errcodes.ErrDuplicatePayment, fmt.Sprintf("payment already received for paymentID '%s'", args.PaymentID))
	}
	pa := &gopayd.PaymentACK{
		Payment: &req,
	}
	// get and attempt to store transaction before processing payment.
	tx, err := bt.NewTxFromString(req.Transaction)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse transaction for paymentID %s", args.PaymentID)
	}
	// TODO: validate the transaction inputs
	outputTotal := uint64(0)
	txos := make([]*gopayd.UpdateTxo, 0)
	// iterate outputs and gather the total satoshis for our known outputs
	for i, o := range tx.Outputs {
		sk, err := p.txoRdr.PartialTxo(ctx, gopayd.UnspentTxoArgs{
			LockingScript: o.LockingScript.String(),
			Satoshis:      o.Satoshis,
			Keyname:       keyname,
		})
		if err != nil {
			// script isn't known to us, could be a change utxo, skip and carry on
			if lathos.IsNotFound(err) {
				continue
			}
			return nil, errors.Wrapf(err, "failed to get store output for paymentID %s", args.PaymentID)
		}
		// has the payment expired - createdAt of the txo is the date the paymentRequest was received
		// so used this as our base.
		if sk.CreatedAt.Add(time.Hour * time.Duration(p.cfg.PaymentExpiryHours)).Before(time.Now().UTC()) {
			return nil, errs.NewErrUnprocessable(errcodes.ErrExpiredPayment, fmt.Sprintf("paymentRequest '%s' has expired, request a new payment", args.PaymentID))
		}
		// push new txo onto list for persistence later
		txos = append(txos, &gopayd.UpdateTxo{
			Outpoint:       fmt.Sprintf("%s%d", tx.TxID(), i),
			TxID:           tx.TxID(),
			Vout:           i,
			KeyName:        null.StringFrom(keyname),
			DerivationPath: null.StringFrom(sk.DerivationPath),
			LockingScript:  sk.LockingScript,
			Satoshis:       o.Satoshis,
		})
		outputTotal += o.Satoshis
	}
	// if it doesn't fully pay the invoice, reject it
	if outputTotal < inv.Satoshis {
		pa.Error = 1
		pa.Memo = "Outputs do not fully pay invoice for paymentID " + args.PaymentID
		return pa, nil
	}
	ctx = p.txrunner.WithTx(ctx)
	// Store utxos and set invoice to paid.
	if _, err = p.store.StoreUtxos(ctx, gopayd.CreateTransaction{
		PaymentID: inv.PaymentID,
		TxID:      tx.TxID(),
		TxHex:     req.Transaction,
		Outputs:   txos,
	}); err != nil {
		log.Error(err)
		pa.Error = 1
		pa.Memo = err.Error()
		return nil, errors.Wrapf(err, "failed to complete payment for paymentID %s", args.PaymentID)
	}
	if _, err := p.invStore.Update(ctx, gopayd.InvoiceUpdateArgs{PaymentID: args.PaymentID}, gopayd.InvoiceUpdate{
		RefundTo: req.RefundTo,
	}); err != nil {
		log.Error(err)
		pa.Error = 1
		pa.Memo = err.Error()
		return nil, errors.Wrapf(err, "failed to update invoice payment for paymentID %s", args.PaymentID)
	}
	// Broadcast the transaction.
	if err := p.sender.Send(ctx, gopayd.SendTransactionArgs{TxID: tx.TxID()}, req); err != nil {
		log.Error(err)
		pa.Error = 1
		pa.Memo = err.Error()
		return pa, errors.Wrapf(err, "failed to send payment for paymentID %s", args.PaymentID)
	}
	return pa, errors.WithStack(p.txrunner.Commit(ctx))
}
