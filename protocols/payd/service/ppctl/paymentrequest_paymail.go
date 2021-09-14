package ppctl

import (
	"context"

	"github.com/pkg/errors"
	validator "github.com/theflyingcodr/govalidator"
	gopaymail "github.com/tonicpow/go-paymail"

	gopayd "github.com/libsv/payd"
	"github.com/libsv/payd/config"
)

type paymailOutputs struct {
	cfg    *config.Paymail
	rdrwtr gopayd.PaymailReaderWriter
	txoWtr gopayd.TxoWriter
}

// NewPaymailOutputs will setup and return a new paymailOutputs service that implements a paymentRequestOutputer.
func NewPaymailOutputs(cfg *config.Paymail, rdrwtr gopayd.PaymailReaderWriter, txoWtr gopayd.TxoWriter) *paymailOutputs {
	return &paymailOutputs{
		cfg:    cfg,
		rdrwtr: rdrwtr,
		txoWtr: txoWtr,
	}
}

// CreateOutputs will generate paymail outputs for the current server paymail address.
func (p *paymailOutputs) CreateOutputs(ctx context.Context, args gopayd.OutputsCreate) ([]*gopayd.Output, error) {
	addr, err := gopaymail.ValidateAndSanitisePaymail(p.cfg.Address, p.cfg.IsBeta)
	if err != nil {
		// convert to known type for the global error handler.
		return nil, validator.ErrValidation{
			"paymailAddress": []string{err.Error()},
		}
	}
	oo, err := p.rdrwtr.OutputsCreate(ctx, gopayd.P2POutputCreateArgs{
		Domain: addr.Domain,
		Alias:  addr.Alias,
	}, gopayd.P2PPayment{Satoshis: args.Satoshis})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create outputs for Alias %s", addr.Alias)
	}
	txos := make([]*gopayd.TxoCreate, 0, len(oo))
	for _, o := range oo {
		txos = append(txos, &gopayd.TxoCreate{
			KeyName:        p.cfg.Address,
			DerivationPath: "paymail",
			LockingScript:  o.Script,
			Satoshis:       args.Satoshis,
		})
	}
	return oo, errors.Wrap(p.txoWtr.TxosCreate(ctx, txos), "failed to store paymail outputs")
}
