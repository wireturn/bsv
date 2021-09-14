package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	gopayd "github.com/libsv/payd"
	"github.com/pkg/errors"
	lathos "github.com/theflyingcodr/lathos/errs"
)

const (
	sqlCreateInvoice = `
	INSERT INTO invoices(paymentID, satoshis)
	VALUES(:paymentID, :satoshis)
	`

	sqlInvoiceByPayID = `
	SELECT paymentId,satoshis,paymentReceivedAt
	FROM invoices
	WHERE paymentId = :paymentID
	`

	sqlInvoices = `
	SELECT paymentId,satoshis,paymentReceivedAt
	FROM invoices
	`

	sqlInvoiceUpdate = `
		UPDATE invoices 
		SET paymentReceivedAt = :paymentReceivedAt, refundTo = :refundTo
		WHERE paymentID = :paymentID
	`

	sqlInvoiceDelete = `
	DELETE FROM invoices 
	WHERE paymentID = :paymentID
	`
)

// Invoice will return an invoice that matches the provided args.
func (s *sqliteStore) Invoice(ctx context.Context, args gopayd.InvoiceArgs) (*gopayd.Invoice, error) {
	var resp gopayd.Invoice
	if err := s.db.GetContext(ctx, &resp, sqlInvoiceByPayID, args.PaymentID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, lathos.NewErrNotFound("N0001", fmt.Sprintf("invoice with paymentID %s not found", args.PaymentID))
		}
		return nil, errors.Wrapf(err, "failed to get invoice with paymentID %s", args.PaymentID)
	}
	return &resp, nil
}

// Invoice will return an invoice that matches the provided args.
func (s *sqliteStore) Invoices(ctx context.Context) ([]gopayd.Invoice, error) {
	var resp []gopayd.Invoice
	if err := s.db.SelectContext(ctx, &resp, sqlInvoices); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, lathos.NewErrNotFound("N0002", "no invoices found")
		}
		return nil, errors.Wrapf(err, "failed to get invoices")
	}
	return resp, nil
}

// Create will persist a new Invoice in the data store.
func (s *sqliteStore) Create(ctx context.Context, req gopayd.InvoiceCreate) (*gopayd.Invoice, error) {
	tx, err := s.newTx(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create new invoice with paymentID %s", req.PaymentID)
	}
	defer func() {
		_ = rollback(ctx, tx)
	}()
	if err := handleNamedExec(tx, sqlCreateInvoice, req); err != nil {
		return nil, errors.Wrap(err, "failed to insert invoice for ")
	}
	var resp gopayd.Invoice
	if err := tx.Get(&resp, sqlInvoiceByPayID, req.PaymentID); err != nil {
		return nil, errors.Wrapf(err, "failed to get new invoice with paymentID %s after creation", req.PaymentID)
	}
	if err := commit(ctx, tx); err != nil {
		return nil, errors.Wrapf(err, "failed to commit transaction when creating invoice with paymentID %s", req.PaymentID)
	}
	return &resp, nil
}

// Update will update an invoice to mark it paid and return the result.
func (s *sqliteStore) Update(ctx context.Context, args gopayd.InvoiceUpdateArgs, req gopayd.InvoiceUpdate) (*gopayd.Invoice, error) {
	tx, err := s.newTx(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update invoice with paymentID %s", args.PaymentID)
	}
	defer func() {
		_ = rollback(ctx, tx)
	}()
	req.PaymentReceivedAt = time.Now().UTC()
	resp, err := s.txUpdateInvoicePaid(tx, args, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update invoice")
	}
	if err := commit(ctx, tx); err != nil {
		return nil, errors.Wrapf(err, "failed to commit transaction when updating invoice with paymentID %s", args.PaymentID)
	}
	return resp, nil
}

func (s *sqliteStore) Delete(ctx context.Context, args gopayd.InvoiceArgs) error {
	tx, err := s.newTx(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to delete invoice with paymentID %s", args.PaymentID)
	}
	defer func() {
		_ = rollback(ctx, tx)
	}()
	if _, err := s.Invoice(ctx, args); err != nil {
		return errors.WithMessagef(err, "failed to find key with id %s to delete", args.PaymentID)
	}
	if err := handleNamedExec(tx, sqlInvoiceDelete, args); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lathos.NewErrNotFound("N0003", fmt.Sprintf("invoice with ID %s not found", args.PaymentID))
		}
		return errors.Wrapf(err, "failed to delete invoice for paymentID %s", args.PaymentID)
	}
	if err := commit(ctx, tx); err != nil {
		return errors.Wrapf(err, "failed to commit transaction when deleting invoice with paymentID %s", args.PaymentID)
	}
	return nil
}

// txUpdateInvoicePaid takes a db object / transaction and adds a transaction to the data store
// along with utxos, returning the updated invoice.
// This method can be used with other methods in the store allowing
// multiple methods to be ran in the same db transaction.
func (s *sqliteStore) txUpdateInvoicePaid(tx db, args gopayd.InvoiceUpdateArgs, req gopayd.InvoiceUpdate) (*gopayd.Invoice, error) {
	req.PaymentReceivedAt = time.Now().UTC()
	if err := handleNamedExec(tx, sqlInvoiceUpdate, map[string]interface{}{
		"paymentReceivedAt": req.PaymentReceivedAt,
		"refundTo":          req.RefundTo,
		"paymentID":         args.PaymentID,
	}); err != nil {
		return nil, errors.Wrapf(err, "failed to update invoice for paymentID %s", args.PaymentID)
	}
	var resp gopayd.Invoice
	if err := tx.Get(&resp, sqlInvoiceByPayID, args.PaymentID); err != nil {
		return nil, errors.Wrapf(err, "failed to get invoice with paymentID %s after update", args.PaymentID)
	}
	return &resp, nil
}
