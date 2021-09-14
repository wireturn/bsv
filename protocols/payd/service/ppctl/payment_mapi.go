package ppctl

import (
	"context"

	gopayd "github.com/libsv/payd"
	"github.com/pkg/errors"
)

type paymentMapiService struct {
	sender gopayd.PaymentSender
}

// NewPaymentMapiSender will setup and return a new mapi payment service.
func NewPaymentMapiSender(sender gopayd.PaymentSender) *paymentMapiService {
	return &paymentMapiService{sender: sender}
}

// CreatePayment will inform the merchant of a new payment being made,
// this payment will then be transmitted to the network and and acknowledgement sent to the user.
func (p *paymentMapiService) Send(ctx context.Context, args gopayd.SendTransactionArgs, req gopayd.CreatePayment) error {
	// Broadcast the transaction.
	return errors.WithStack(p.sender.Send(ctx, args, req))
}
