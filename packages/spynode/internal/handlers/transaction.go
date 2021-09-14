package handlers

import (
	"context"
	"sync"

	"github.com/tokenized/pkg/wire"
	"github.com/tokenized/spynode/pkg/client"

	"github.com/pkg/errors"
)

// TXHandler exists to handle the tx command.
type TXHandler struct {
	ready     StateReady
	txChannel *TxChannel
}

type TxData struct {
	Msg             *wire.MsgTx
	Trusted         bool
	Safe            bool
	ConfirmedHeight int
}

// NewTXHandler returns a new TXHandler with the given Config.
func NewTXHandler(ready StateReady, txChannel *TxChannel) *TXHandler {

	result := TXHandler{
		ready:     ready,
		txChannel: txChannel,
	}
	return &result
}

// Handle implements the handler interface for transaction handler.
func (handler *TXHandler) Handle(ctx context.Context, m wire.Message) ([]wire.Message, error) {
	msg, ok := m.(*wire.MsgTx)
	if !ok {
		return nil, errors.New("Could not assert as *wire.MsgTx")
	}

	// Only notify of transactions when in sync or they might be duplicated, since there isn't a
	// mempool yet.
	if !handler.ready.IsReady() {
		return nil, nil
	}

	handler.txChannel.Add(TxData{Msg: msg, Trusted: true, ConfirmedHeight: -1})
	return nil, nil
}

type TxChannel struct {
	Channel chan TxData
	lock    sync.Mutex
	open    bool
}

func (c *TxChannel) Add(tx TxData) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.open {
		return errors.New("Channel closed")
	}

	c.Channel <- tx
	return nil
}

func (c *TxChannel) Open(count int) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.Channel = make(chan TxData, count)
	c.open = true
	return nil
}

func (c *TxChannel) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.open {
		return errors.New("Channel closed")
	}

	close(c.Channel)
	c.open = false
	return nil
}

type TxUpdateChannel struct {
	Channel chan *client.TxUpdate
	lock    sync.Mutex
	open    bool
}

func (c *TxUpdateChannel) Add(update *client.TxUpdate) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.open {
		return errors.New("Channel closed")
	}

	c.Channel <- update
	return nil
}

func (c *TxUpdateChannel) Open(count int) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.Channel = make(chan *client.TxUpdate, count)
	c.open = true
	return nil
}

func (c *TxUpdateChannel) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.open {
		return errors.New("Channel closed")
	}

	close(c.Channel)
	c.open = false
	return nil
}
