package gopayd

import (
	"context"
	"time"

	validator "github.com/theflyingcodr/govalidator"
	"gopkg.in/guregu/null.v3"
)

// Invoice stores information related to a payment.
type Invoice struct {
	PaymentID         string      `json:"paymentID" db:"paymentID"`
	Satoshis          uint64      `json:"satoshis" db:"satoshis"`
	PaymentReceivedAt null.Time   `json:"paymentReceivedAt" db:"paymentReceivedAt"`
	RefundTo          null.String `json:"refundTo" db:"refundTo"`
}

// InvoiceCreate is used to create a new invoice.
type InvoiceCreate struct {
	// PaymentID is the unique identifier for a payment.
	PaymentID string `json:"-" db:"paymentID"`
	Satoshis  uint64 `json:"satoshis" db:"satoshis"`
}

// Validate will check that InvoiceCreate params match expectations.
func (i InvoiceCreate) Validate() validator.ErrValidation {
	return validator.New().
		Validate("satoshis", validator.MinUInt64(i.Satoshis, 546))
}

// InvoiceUpdate can be used to update an invoice after it has been created.
type InvoiceUpdate struct {
	PaymentReceivedAt time.Time   `db:"paymentReceviedAt"`
	RefundTo          null.String `db:"refundTo"`
}

// InvoiceUpdateArgs are used to identify the invoice to update.
type InvoiceUpdateArgs struct {
	PaymentID string
}

// InvoiceArgs contains argument/s to return a single invoice.
type InvoiceArgs struct {
	PaymentID string `param:"paymentID" db:"paymentID"`
}

// Validate will check that invoice arguments match expectations.
func (i *InvoiceArgs) Validate() validator.ErrValidation {
	return validator.New().Validate("paymentID", validator.Length(i.PaymentID, 1, 30))
}

// InvoiceService defines a service for managing invoices.
type InvoiceService interface {
	Invoice(ctx context.Context, args InvoiceArgs) (*Invoice, error)
	Invoices(ctx context.Context) ([]Invoice, error)
	Create(ctx context.Context, req InvoiceCreate) (*Invoice, error)
	Delete(ctx context.Context, args InvoiceArgs) error
}

// InvoiceReaderWriter can be implemented to support storing and retrieval of invoices.
type InvoiceReaderWriter interface {
	InvoiceWriter
	InvoiceReader
}

// InvoiceWriter defines a data store used to write invoice data.
type InvoiceWriter interface {
	// Create will persist a new Invoice in the data store.
	Create(ctx context.Context, req InvoiceCreate) (*Invoice, error)
	// Update will update an invoice matching the provided args with the requested changes.
	Update(ctx context.Context, args InvoiceUpdateArgs, req InvoiceUpdate) (*Invoice, error)
	Delete(ctx context.Context, args InvoiceArgs) error
}

// InvoiceReader defines a data store used to read invoice data.
type InvoiceReader interface {
	// Invoice will return an invoice that matches the provided args.
	Invoice(ctx context.Context, args InvoiceArgs) (*Invoice, error)
	// Invoices returns all currently stored invoices TODO: update to support search args
	Invoices(ctx context.Context) ([]Invoice, error)
}
