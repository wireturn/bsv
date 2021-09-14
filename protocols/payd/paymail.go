package gopayd

import (
	"context"
)

// P2PTransactionArgs is used to get a transaction.
type P2PTransactionArgs struct {
	Alias     string
	Domain    string
	PaymentID string
	TxHex     string
}

// P2PTransaction defines a peer to peer transaction.
type P2PTransaction struct {
	TxHex    string
	Metadata P2PTransactionMetadata
}

// P2PTransactionMetadata contains potentially optional
// metadata that can be sent along with a p2p payment.
type P2PTransactionMetadata struct {
	Note      string // A human readable bit of information about the payment
	PubKey    string // Public key to validate the signature (if signature is given)
	Sender    string // The paymail of the person that originated the transaction
	Signature string
}

// P2PCapabilityArgs is used to retrieve information from a
// p2p server.
// https://bsvalias.org/02-02-capability-discovery.html
type P2PCapabilityArgs struct {
	Domain string
	BrfcID string
}

// P2POutputCreateArgs is used to locate a signer when sending a p2p payment.
type P2POutputCreateArgs struct {
	Domain string
	Alias  string
}

// P2PPayment contains the amount of satoshis to send peer to peer.
type P2PPayment struct {
	Satoshis uint64
}

// PaymailReader reads paymail information from a datastore.
type PaymailReader interface {
	Capability(ctx context.Context, args P2PCapabilityArgs) (string, error)
}

// PaymailWriter writes to a paymail datastore.
type PaymailWriter interface {
	OutputsCreate(ctx context.Context, args P2POutputCreateArgs, req P2PPayment) ([]*Output, error)
}

// PaymailReaderWriter combines the reader and writer interfaces.
type PaymailReaderWriter interface {
	PaymailReader
	PaymailWriter
}
