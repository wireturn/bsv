package client

import (
	"context"

	"github.com/tokenized/pkg/bitcoin"
	"github.com/tokenized/pkg/wire"

	"github.com/pkg/errors"
)

var (
	ErrUnknownMessageType = errors.New("Unknown Message Type")
	ErrInvalid            = errors.New("Invalid")
	ErrIncomplete         = errors.New("Incomplete") // merkle proof is incomplete
	ErrWrongHash          = errors.New("Wrong Hash") // Non-matching merkle root hash
	ErrNotConnected       = errors.New("Not Connected")
	ErrWrongKey           = errors.New("Wrong Key") // The wrong key was provided during auth
	ErrBadSignature       = errors.New("Bad Signature")
	ErrTimeout            = errors.New("Timeout")
	ErrReject             = errors.New("Reject")
)

// Handler provides an interface for handling data from the spynode client.
type Handler interface {
	HandleTx(context.Context, *Tx)
	HandleTxUpdate(context.Context, *TxUpdate)
	HandleHeaders(context.Context, *Headers)
	HandleInSync(context.Context)

	// HandleMessage handles all other client data messages
	HandleMessage(context.Context, MessagePayload)
}

// Client is the interface for interacting with a spynode.
type Client interface {
	RegisterHandler(Handler)

	// Subscribe to notifications about any transactions whose output scripts contain these push
	// datas.
	SubscribePushDatas(context.Context, [][]byte) error
	UnsubscribePushDatas(context.Context, [][]byte) error

	// Subscribe to notifications about the transaction with the specified hash and any transactions
	// spending these outputs.
	SubscribeTx(context.Context, bitcoin.Hash32, []uint32) error
	UnsubscribeTx(context.Context, bitcoin.Hash32, []uint32) error

	// Subscribe to notifications about any transactions spending these outputs.
	// Note: This should be used for Tokenized requests with the contract output(s) so the response
	// transaction from the contract is seen.
	SubscribeOutputs(context.Context, []*wire.OutPoint) error
	UnsubscribeOutputs(context.Context, []*wire.OutPoint) error

	// Subscribe to notifications about any transactions whose output scripts contain Tokenized
	// "contract-wide" actions.
	SubscribeContracts(context.Context) error
	UnsubscribeContracts(context.Context) error

	// Subscribe to all new block headers.
	SubscribeHeaders(context.Context) error
	UnsubscribeHeaders(context.Context) error

	GetTx(context.Context, bitcoin.Hash32) (*wire.MsgTx, error)
	GetOutputs(context.Context, []wire.OutPoint) ([]bitcoin.UTXO, error)

	// Send a tx to the network.
	SendTx(context.Context, *wire.MsgTx) error

	// Send a tx to the network and subscribe to the outputs specified by indexes.
	SendTxAndMarkOutputs(context.Context, *wire.MsgTx, []uint32) error

	GetHeaders(context.Context, int, int) (*Headers, error)
	BlockHash(context.Context, int) (*bitcoin.Hash32, error)

	// Notify the service to activate the notification message feed.
	Ready(context.Context, uint64) error

	// Returns the next message that is expected.
	NextMessageID() uint64
}

// PushDataSubscriber subscribes to monitoring for transactions containing the push datas.
type PushDataSubscriber interface {
	SubscribePushDatas(context.Context, [][]byte) error
}

func SubscribeAddresses(ctx context.Context, ras []bitcoin.RawAddress,
	subscriber PushDataSubscriber) error {

	pds := make([][]byte, 0, len(ras))
	for _, ra := range ras {
		if ra.IsEmpty() {
			continue
		}

		hashes, err := ra.Hashes()
		if err != nil {
			return errors.Wrap(err, "address hashes")
		}

		for _, hash := range hashes {
			pds = append(pds, hash[:])
		}
	}

	return subscriber.SubscribePushDatas(ctx, pds)
}

func SubscribeAddress(ctx context.Context, ra bitcoin.RawAddress,
	subscriber PushDataSubscriber) error {

	hashes, err := ra.Hashes()
	if err != nil {
		return errors.Wrap(err, "address hashes")
	}

	var pds [][]byte
	for _, hash := range hashes {
		pds = append(pds, hash[:])
	}

	return subscriber.SubscribePushDatas(ctx, pds)
}
