package service

import (
	"context"
	"time"

	"github.com/pkg/errors"

	gopayd "github.com/libsv/payd"
	"github.com/libsv/payd/config"

	"github.com/speps/go-hashids"
)

// invoice represents a purchase order system or other such system that a merchant would use
// to receive orders from customers.
// This could be a Pos system or online retailer etc.
// The invoice system would create an invoice / PO and then the protocol
// server would be sent this invoice for lookup.
// This invoicing system is separate to the protocol server itself but added here
// as a very basic example.
type invoice struct {
	store gopayd.InvoiceReaderWriter
	cfg   *config.Server
}

// NewInvoice will setup and return a new invoice service.
func NewInvoice(cfg *config.Server, store gopayd.InvoiceReaderWriter) *invoice {
	return &invoice{
		cfg:   cfg,
		store: store}
}

// Invoice will return an invoice by paymentID.
func (i *invoice) Invoice(ctx context.Context, args gopayd.InvoiceArgs) (*gopayd.Invoice, error) {
	if err := args.Validate().Err(); err != nil {
		return nil, err
	}
	inv, err := i.store.Invoice(ctx, args)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get invoice with id %s", args.PaymentID)
	}
	return inv, err
}

// Invoices will return all currently stored invoices.
func (i *invoice) Invoices(ctx context.Context) ([]gopayd.Invoice, error) {
	ii, err := i.store.Invoices(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get invoices")
	}
	return ii, nil
}

// Create will add a new invoice to the system.
func (i *invoice) Create(ctx context.Context, req gopayd.InvoiceCreate) (*gopayd.Invoice, error) {
	if err := req.Validate().Err(); err != nil {
		return nil, err
	}
	hd := hashids.NewData()
	hd.Alphabet = hashids.DefaultAlphabet
	hd.Salt = i.cfg.Hostname
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	id, err := h.Encode([]int{time.Now().Nanosecond()})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.PaymentID = id
	inv, err := i.store.Create(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return inv, nil
}

// Delete will permanently remove an invoice from the system.
func (i *invoice) Delete(ctx context.Context, args gopayd.InvoiceArgs) error {
	if err := args.Validate().Err(); err != nil {
		return err
	}
	return errors.WithMessagef(i.store.Delete(ctx, args),
		"failed to delete invoice with ID %s", args.PaymentID)
}
