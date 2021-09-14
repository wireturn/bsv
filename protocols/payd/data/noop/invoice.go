package noop

import (
	"context"

	gopayd "github.com/libsv/payd"
	"gopkg.in/guregu/null.v3"
)

// invoice is a no-op invoice that returns some stubbed data.
type invoice struct {
}

// NewInvoice will return a new instance of a noop invoice.
func NewInvoice() *invoice {
	return &invoice{}
}

// Invoice will return an invoice that matches the provided args.
func (i *invoice) Invoice(ctx context.Context, args gopayd.InvoiceArgs) (*gopayd.Invoice, error) {
	return &gopayd.Invoice{
		PaymentID: args.PaymentID,
		Satoshis:  10000,
	}, nil
}

// Create will persist a new Invoice in the data store.
func (i *invoice) Create(ctx context.Context, req gopayd.InvoiceCreate) (*gopayd.Invoice, error) {
	return &gopayd.Invoice{
		PaymentID: req.PaymentID,
		Satoshis:  req.Satoshis,
	}, nil
}

// Update will update an invoice and return the result.
func (i *invoice) Update(ctx context.Context, args gopayd.InvoiceUpdateArgs, req gopayd.InvoiceUpdate) (*gopayd.Invoice, error) {
	return &gopayd.Invoice{
		PaymentID:         args.PaymentID,
		Satoshis:          10000,
		PaymentReceivedAt: null.TimeFrom(req.PaymentReceivedAt),
	}, nil
}
