package paymail

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/pkg/errors"
	lathos "github.com/theflyingcodr/lathos/errs"
	gopaymail "github.com/tonicpow/go-paymail"

	gopayd "github.com/libsv/payd"
	"github.com/libsv/payd/config"
	"github.com/libsv/payd/data/paymail/models"
)

// TODO: test this code.
// powmailClient is a wrapper interface for tonicpow.go-paymail
// to allow easier unit testing of this code.
type powmailClient interface {
	GetSRVRecord(service, protocol, domainName string) (srv *net.SRV, err error)
	GetCapabilities(target string, port int) (response *gopaymail.Capabilities, err error)
	GetP2PPaymentDestination(p2pURL, alias, domain string, paymentRequest *gopaymail.PaymentRequest) (response *gopaymail.PaymentDestination, err error)
	VerifyPubKey(verifyURL, alias, domain, pubKey string) (response *gopaymail.Verification, err error)
}

type paymail struct {
	mu     sync.Mutex
	cli    powmailClient
	addr   *gopaymail.SanitisedPaymail
	cstore map[string]*gopaymail.Capabilities
}

// NewPaymail will setup and return a new paymail data store used
// to create and send paymail transactions.
// Currently we support one paymail address per server.
func NewPaymail(cfg *config.Paymail, cli powmailClient) *paymail {
	addr, err := gopaymail.ValidateAndSanitisePaymail(cfg.Address, cfg.IsBeta)
	if err != nil {
		log.Fatalf("failed to validate wallet payment address %s", cfg.Address)
	}
	return &paymail{
		addr:   addr,
		cli:    cli,
		cstore: map[string]*gopaymail.Capabilities{},
	}
}

// Capability will return a capability or a notfound error if it could not be found.
func (p *paymail) Capability(ctx context.Context, args gopayd.P2PCapabilityArgs) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	c, ok := p.cstore[args.Domain]
	if ok && c.Has(args.BrfcID, "") {
		return c.GetString(args.BrfcID, ""), nil
	}
	srv, err := p.cli.GetSRVRecord(gopaymail.DefaultServiceName, gopaymail.DefaultProtocol, args.Domain)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get srv record for %s", args.Domain)
	}
	cp, err := p.cli.GetCapabilities(srv.Target, int(srv.Port))
	if err != nil {
		return "", errors.Wrapf(err, "failed to get capabilities for %s", args.Domain)
	}
	p.cstore[args.Domain] = cp
	if cp.Has(args.BrfcID, "") {
		return cp.GetString(args.BrfcID, ""), nil
	}
	return "", lathos.NewErrNotFound("N001",
		fmt.Sprintf("brfcID [%s] not found for domain [%s]", args.BrfcID, args.Domain))
}

// OutputsCreate will create outputs for the provided payment information. Args are used to gather capability information
// a lathos.NotFound error may be returned if the paymail or brfc doesn't exist.
func (p *paymail) OutputsCreate(ctx context.Context, args gopayd.P2POutputCreateArgs, req gopayd.P2PPayment) ([]*gopayd.Output, error) {
	url, err := p.Capability(ctx, gopayd.P2PCapabilityArgs{
		Domain: args.Domain,
		BrfcID: gopaymail.BRFCP2PPaymentDestination,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get BRFCP2P Payment Destination for domain %s", args.Domain)
	}
	resp, err := p.cli.GetP2PPaymentDestination(url, args.Alias, args.Domain, &gopaymail.PaymentRequest{Satoshis: req.Satoshis})
	if err != nil {
		if err.Error() == "paymail address not found" {
			return nil, lathos.NewErrNotFound("N003", err.Error())
		}
		return nil, errors.Wrapf(err, "failed to generate paymail outputs for alias %s", args.Alias)
	}
	return models.OutputsToPayd(resp.Outputs), nil
}
