package client

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/tokenized/metrics"
	"github.com/tokenized/pkg/bitcoin"
	"github.com/tokenized/pkg/logger"
	"github.com/tokenized/pkg/wire"

	"github.com/pkg/errors"
)

var (
	Endian = binary.LittleEndian

	// RemoteClientVersion is the current version of the communication protocol
	RemoteClientVersion = uint8(0)
)

// RemoteClient is a client for interacting with the spynode service.
type RemoteClient struct {
	conn net.Conn

	config        *Config
	nextMessageID uint64

	// Session
	hash             bitcoin.Hash32    // for generating session key
	serverSessionKey bitcoin.PublicKey // for this session
	sessionKey       bitcoin.Key

	handlers    []Handler
	handlerLock sync.Mutex

	// Process handlers in separate thread to prevent conflict between listener thread calling
	// handlers directly listener thread waiting for responses from within handler.
	handlerChannel       chan MessagePayload
	handlerChannelIsOpen bool
	handlerChannelLock   sync.Mutex

	// Requests
	sendTxRequests []*sendTxRequest
	getTxRequests  []*getTxRequest
	headerRequests []*headerRequest
	requestLock    sync.Mutex

	accepted, ready bool
	lock            sync.Mutex
	wait            sync.WaitGroup

	closeRequested     bool
	closeRequestedLock sync.Mutex

	listenErrChannel *chan error
}

type sendTxRequest struct {
	txid     bitcoin.Hash32
	response *Message
	lock     sync.Mutex
}

type getTxRequest struct {
	txid     bitcoin.Hash32
	response *Message
	lock     sync.Mutex
}

type headerRequest struct {
	height   int
	response *Message
	lock     sync.Mutex
}

// NewRemoteClient creates a remote client.
// Note: If the connection type is not "full" then it will auto-connect when a function is called to
// communicate with the spynode service. Make sure `Close` is called before application end so that
// the connection can be closed and the listen thread completed.
func NewRemoteClient(config *Config) (*RemoteClient, error) {
	result := &RemoteClient{
		config:        config,
		nextMessageID: 1,
	}

	return result, nil
}

// SetupRetry sets the maximum connection retry attempts and delay in milliseconds before failing.
// This can also be set from the config.
func (c *RemoteClient) SetupRetry(max, delay int) {
	c.config.MaxRetries = max
	c.config.RetryDelay = delay
}

func (c *RemoteClient) RegisterHandler(h Handler) {
	c.handlerLock.Lock()
	c.handlers = append(c.handlers, h)
	c.handlerLock.Unlock()
}

func (c *RemoteClient) IsAccepted(ctx context.Context) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.accepted
}

// SubscribePushDatas subscribes to transactions containing the specified push datas.
func (c *RemoteClient) SubscribePushDatas(ctx context.Context, pushDatas [][]byte) error {
	m := &SubscribePushData{
		PushDatas: pushDatas,
	}

	logger.Info(ctx, "Sending subscribe push data message")
	return c.sendMessage(ctx, m)
}

// UnsubscribePushDatas unsubscribes to transactions containing the specified push datas.
func (c *RemoteClient) UnsubscribePushDatas(ctx context.Context, pushDatas [][]byte) error {
	m := &UnsubscribePushData{
		PushDatas: pushDatas,
	}

	logger.Info(ctx, "Sending unsubscribe push data message")
	return c.sendMessage(ctx, m)
}

// SubscribeTx subscribes to information for a specific transaction. Indexes are the indexes of the
// outputs that need to be monitored for spending.
func (c *RemoteClient) SubscribeTx(ctx context.Context, txid bitcoin.Hash32,
	indexes []uint32) error {

	m := &SubscribeTx{
		TxID:    txid,
		Indexes: indexes,
	}

	logger.Info(ctx, "Sending subscribe tx message")
	return c.sendMessage(ctx, m)
}

// UnsubscribeTx unsubscribes to information for a specific transaction.
func (c *RemoteClient) UnsubscribeTx(ctx context.Context, txid bitcoin.Hash32,
	indexes []uint32) error {

	m := &UnsubscribeTx{
		TxID:    txid,
		Indexes: indexes,
	}

	logger.Info(ctx, "Sending unsubscribe tx message")
	return c.sendMessage(ctx, m)
}

func (c *RemoteClient) SubscribeOutputs(ctx context.Context, outputs []*wire.OutPoint) error {
	m := &SubscribeOutputs{
		Outputs: outputs,
	}

	logger.Info(ctx, "Sending subscribe outputs message")
	return c.sendMessage(ctx, m)
}

func (c *RemoteClient) UnsubscribeOutputs(ctx context.Context, outputs []*wire.OutPoint) error {
	m := &UnsubscribeOutputs{
		Outputs: outputs,
	}

	logger.Info(ctx, "Sending unsubscribe outputs message")
	return c.sendMessage(ctx, m)
}

// SubscribeHeaders subscribes to information on new block headers.
func (c *RemoteClient) SubscribeHeaders(ctx context.Context) error {
	m := &SubscribeHeaders{}

	logger.Info(ctx, "Sending subscribe headers message")
	return c.sendMessage(ctx, m)
}

// UnsubscribeHeaders unsubscribes to information on new block headers.
func (c *RemoteClient) UnsubscribeHeaders(ctx context.Context) error {
	m := &UnsubscribeHeaders{}

	logger.Info(ctx, "Sending unsubscribe headers message")
	return c.sendMessage(ctx, m)
}

// SubscribeContracts subscribes to information on contracts.
func (c *RemoteClient) SubscribeContracts(ctx context.Context) error {
	m := &SubscribeContracts{}

	logger.Info(ctx, "Sending subscribe contracts message")
	return c.sendMessage(ctx, m)
}

// UnsubscribeContracts unsubscribes to information on contracts.
func (c *RemoteClient) UnsubscribeContracts(ctx context.Context) error {
	m := &UnsubscribeContracts{}

	logger.Info(ctx, "Sending unsubscribe contracts message")
	return c.sendMessage(ctx, m)
}

// Ready tells the spynode the client is ready to start receiving updates. Call this after
// connecting and subscribing to all relevant push data.
func (c *RemoteClient) Ready(ctx context.Context, nextMessageID uint64) error {
	if nextMessageID == 0 {
		nextMessageID = 1 // first message id is 1
	}

	m := &Ready{
		NextMessageID: nextMessageID,
	}

	c.lock.Lock()
	c.nextMessageID = nextMessageID
	c.ready = true
	c.lock.Unlock()

	logger.Info(ctx, "Sending ready message (next message %d)", nextMessageID)
	if err := c.sendMessage(ctx, m); err != nil {
		return err
	}
	return nil
}

func (c *RemoteClient) NextMessageID() uint64 {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.nextMessageID
}

func (c *RemoteClient) SendTx(ctx context.Context, tx *wire.MsgTx) error {
	return c.SendTxAndMarkOutputs(ctx, tx, nil)
}

// SendTxAndMarkOutputs sends a tx message to the bitcoin network. It is synchronous meaning it
// will wait for a response before returning.
func (c *RemoteClient) SendTxAndMarkOutputs(ctx context.Context, tx *wire.MsgTx,
	indexes []uint32) error {
	start := time.Now()
	defer metrics.Elapsed(ctx, start, "SpyNodeClient.SendTxAndMarkOutputs")

	// Register with listener for response
	request := &sendTxRequest{
		txid: *tx.TxHash(),
	}

	c.requestLock.Lock()
	c.sendTxRequests = append(c.sendTxRequests, request)
	c.requestLock.Unlock()

	logger.Info(ctx, "Sending send tx message : %s", tx.TxHash())
	m := &SendTx{Tx: tx}
	if err := c.sendMessage(ctx, m); err != nil {
		return err
	}

	// Wait for response
	timeout := start.Add(time.Duration(c.config.RequestTimeout) * time.Millisecond)
	for time.Now().Before(timeout) {
		request.lock.Lock()
		if request.response != nil {
			request.lock.Unlock()

			// Remove
			c.requestLock.Lock()
			for i, r := range c.sendTxRequests {
				if r == request {
					c.sendTxRequests = append(c.sendTxRequests[:i], c.sendTxRequests[i+1:]...)
					break
				}
			}
			c.requestLock.Unlock()

			switch msg := request.response.Payload.(type) {
			case *Reject:
				return errors.Wrap(ErrReject, msg.Message)
			case *Accept:
				return nil
			default:
				return fmt.Errorf("Unknown response : %d", request.response.Payload.Type())
			}
		}
		request.lock.Unlock()

		time.Sleep(1 * time.Millisecond)
	}

	return ErrTimeout
}

// GetTx requests a tx from the bitcoin network. It is synchronous meaning it will wait for a
// response before returning.
func (c *RemoteClient) GetTx(ctx context.Context, txid bitcoin.Hash32) (*wire.MsgTx, error) {
	start := time.Now()
	defer metrics.Elapsed(ctx, start, "SpyNodeClient.GetTx")

	// Register with listener for response tx
	request := &getTxRequest{
		txid: txid,
	}

	c.requestLock.Lock()
	c.getTxRequests = append(c.getTxRequests, request)
	c.requestLock.Unlock()

	logger.Info(ctx, "Sending get tx message : %s", txid)
	m := &GetTx{TxID: txid}
	if err := c.sendMessage(ctx, m); err != nil {
		return nil, err
	}

	// Wait for response
	timeout := start.Add(time.Duration(c.config.RequestTimeout) * time.Millisecond)
	for time.Now().Before(timeout) {
		request.lock.Lock()
		if request.response != nil {
			request.lock.Unlock()
			// Remove
			c.requestLock.Lock()
			for i, r := range c.getTxRequests {
				if r == request {
					c.getTxRequests = append(c.getTxRequests[:i], c.getTxRequests[i+1:]...)
					break
				}
			}
			c.requestLock.Unlock()

			switch msg := request.response.Payload.(type) {
			case *Reject:
				return nil, errors.Wrap(ErrReject, msg.Message)
			case *BaseTx:
				return msg.Tx, nil
			default:
				return nil, fmt.Errorf("Unknown response : %d", request.response.Payload.Type())
			}
		}
		request.lock.Unlock()

		time.Sleep(1 * time.Millisecond)
	}

	return nil, ErrTimeout
}

func (c *RemoteClient) GetOutputs(ctx context.Context,
	outpoints []wire.OutPoint) ([]bitcoin.UTXO, error) {
	start := time.Now()
	defer metrics.Elapsed(ctx, start, "SpyNodeClient.GetOutputs")

	outputs := make([]*wire.TxOut, len(outpoints))
	for i, outpoint := range outpoints {
		if outputs[i] != nil {
			continue // already fetched this output
		}

		tx, err := c.GetTx(ctx, outpoint.Hash)
		if err != nil {
			return nil, errors.Wrap(err, "get tx")
		}

		if int(outpoint.Index) >= len(tx.TxOut) {
			return nil, errors.Wrap(err, "invalid index")
		}
		outputs[i] = tx.TxOut[outpoint.Index]

		// Check if other outpoints have the same txid.
		for j := range outpoints[i+1:] {
			if outpoints[j].Hash.Equal(&outpoint.Hash) {
				if int(outpoint.Index) >= len(tx.TxOut) {
					return nil, errors.Wrap(err, "invalid index")
				}
				outputs[j] = tx.TxOut[outpoint.Index]
			}
		}
	}

	result := make([]bitcoin.UTXO, len(outputs))
	for i, output := range outputs {
		result[i] = bitcoin.UTXO{
			Hash:          outpoints[i].Hash,
			Index:         outpoints[i].Index,
			Value:         output.Value,
			LockingScript: output.PkScript,
		}
	}

	return result, nil
}

// GetHeaders requests a header from the bitcoin network. It is synchronous meaning it will wait for
// a response before returning.
func (c *RemoteClient) GetHeaders(ctx context.Context, height, count int) (*Headers, error) {
	start := time.Now()
	defer metrics.Elapsed(ctx, start, "SpyNodeClient.GetHeaders")

	// Register with listener for response tx
	request := &headerRequest{
		height: height,
	}

	c.requestLock.Lock()
	c.headerRequests = append(c.headerRequests, request)
	c.requestLock.Unlock()

	logger.Info(ctx, "Sending get header message : height %d, count %d", height, count)
	m := &GetHeaders{
		RequestHeight: int32(height),
		MaxCount:      uint32(count),
	}
	if err := c.sendMessage(ctx, m); err != nil {
		return nil, err
	}

	// Wait for response
	timeout := start.Add(time.Duration(c.config.RequestTimeout) * time.Millisecond)
	for time.Now().Before(timeout) {
		request.lock.Lock()
		if request.response != nil {
			request.lock.Unlock()
			// Remove
			c.requestLock.Lock()
			for i, r := range c.headerRequests {
				if r == request {
					c.headerRequests = append(c.headerRequests[:i], c.headerRequests[i+1:]...)
					break
				}
			}
			c.requestLock.Unlock()

			switch msg := request.response.Payload.(type) {
			case *Reject:
				return nil, errors.Wrap(ErrReject, msg.Message)
			case *Headers:
				return msg, nil
			default:
				return nil, fmt.Errorf("Unknown response : %d", request.response.Payload.Type())
			}
		}
		request.lock.Unlock()

		time.Sleep(1 * time.Millisecond)
	}

	return nil, ErrTimeout
}

func (c *RemoteClient) BlockHash(ctx context.Context, height int) (*bitcoin.Hash32, error) {
	headers, err := c.GetHeaders(ctx, height, 1)
	if err != nil {
		return nil, errors.Wrap(err, "get headers")
	}

	if len(headers.Headers) == 0 {
		return nil, errors.New("No headers returned")
	}

	return headers.Headers[0].BlockHash(), nil
}

// sendMessage wraps and sends a message to the server.
func (c *RemoteClient) sendMessage(ctx context.Context, payload MessagePayload) error {
	if c.config.ConnectionType != ConnectionTypeFull {
		// Connect if not already connected
		if err := c.Connect(ctx); err != nil {
			return errors.Wrap(err, "connect")
		}
	}

	c.lock.Lock()

	if c.conn == nil {
		c.lock.Unlock()
		return ErrNotConnected
	}

	conn := c.conn
	c.lock.Unlock()

	message := &Message{
		Payload: payload,
	}

	// TODO Possibly add streaming encryption here. --ce

	if err := message.Serialize(conn); err != nil {
		return errors.Wrap(err, "send message")
	}

	return nil
}

func (c *RemoteClient) IsConnected() bool {
	c.lock.Lock()
	result := c.conn != nil
	c.lock.Unlock()

	return result
}

// SetListenerErrorChannel sets a channel that will receive an error when the listener returns.
func (c *RemoteClient) SetListenerErrorChannel(channel *chan error) {
	c.listenErrChannel = channel
}

// Connect connects to the spynode service if it isn't already connected and also starts the
// listiner thread.
func (c *RemoteClient) Connect(ctx context.Context) error {
	if isNewConnection, err := c.connect(ctx); err != nil {
		return err
	} else if !isNewConnection {
		return nil
	}

	// Start listener thread
	c.wait.Add(1)
	go func() {
		logger.Info(ctx, "Spynode client listening")
		if c.listenErrChannel != nil {
			*c.listenErrChannel <- c.listen(ctx)
		} else {
			if err := c.listen(ctx); err != nil {
				logger.Warn(ctx, "Listener finished with error : %s", err)
			}
		}
		logger.Info(ctx, "Spynode client finished listening")
		c.wait.Done()
	}()

	// Start handler thread
	c.handlerChannelLock.Lock()
	if c.handlerChannelIsOpen {
		logger.Error(ctx, "Handler channel already open")
	}
	c.handlerChannelIsOpen = true
	c.handlerChannel = make(chan MessagePayload, 1000)
	c.handlerChannelLock.Unlock()

	c.wait.Add(1)
	go func() {
		logger.Info(ctx, "Spynode client handler running")
		if err := c.handle(ctx); err != nil {
			logger.Warn(ctx, "Spynode client handler finished with error : %s", err)
		}
		logger.Info(ctx, "Spynode client handler finished")
		c.wait.Done()
	}()

	// Start ping thread
	c.wait.Add(1)
	go func() {
		logger.Info(ctx, "Spynode client pinging")
		if err := c.ping(ctx); err != nil {
			logger.Warn(ctx, "Pinger finished with error : %s", err)
			c.closeRequestedLock.Lock()
			c.closeRequested = true
			c.closeRequestedLock.Unlock()
		}
		logger.Info(ctx, "Spynode client finished pinging")
		c.wait.Done()
	}()

	return nil
}

func (c *RemoteClient) Close(ctx context.Context) {
	c.closeRequestedLock.Lock()
	c.closeRequested = true
	c.closeRequestedLock.Unlock()

	c.close(ctx)
	c.wait.Wait() // Wait for listen thread to finish

	// Clear close requested flag
	c.closeRequestedLock.Lock()
	c.closeRequested = false
	c.closeRequestedLock.Unlock()
}

func (c *RemoteClient) connect(ctx context.Context) (bool, error) {
	c.closeRequestedLock.Lock()
	c.closeRequested = false
	c.closeRequestedLock.Unlock()

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.conn != nil {
		return false, nil // already connected
	}

	start := time.Now()
	defer metrics.Elapsed(ctx, start, "SpyNodeClient.connect")

	var connectErr error
	for i := 0; i <= c.config.MaxRetries; i++ {
		if i > 0 {
			// Delay, then retry
			logger.Info(ctx, "Delaying %d milliseconds before dial retry %d", c.config.RetryDelay,
				i)
			time.Sleep(time.Millisecond * time.Duration(c.config.RetryDelay))
		}

		// Check if we are trying to close
		c.closeRequestedLock.Lock()
		stop := c.closeRequested
		c.closeRequestedLock.Unlock()
		if stop {
			return false, connectErr
		}

		logger.Info(ctx, "Connecting to spynode service")

		if err := c.generateSession(); err != nil {
			return false, errors.Wrap(err, "session")
		}

		var dialer net.Dialer
		conn, err := dialer.DialContext(ctx, "tcp", c.config.ServerAddress)
		if err != nil {
			logger.Warn(ctx, "Spynode service dial failed : %s", err)
			connectErr = err
			continue
		}

		// Create and sign register message
		register := &Register{
			Version:          RemoteClientVersion,
			Key:              c.config.ClientKey.PublicKey(),
			Hash:             c.hash,
			StartBlockHeight: c.config.StartBlockHeight,
			ConnectionType:   c.config.ConnectionType,
		}

		sigHash, err := register.SigHash()
		if err != nil {
			conn.Close()
			return false, errors.Wrap(err, "sig hash")
		}

		register.Signature, err = c.config.ClientKey.Sign(sigHash.Bytes())
		if err != nil {
			conn.Close()
			return false, errors.Wrap(err, "sign")
		}

		message := Message{Payload: register}
		if err := message.Serialize(conn); err != nil {
			conn.Close()
			return false, errors.Wrap(err, "send register")
		}

		c.conn = conn

		return true, nil
	}

	return false, connectErr
}

func (c *RemoteClient) close(ctx context.Context) {
	c.lock.Lock()
	if c.conn != nil {
		logger.Info(ctx, "Closing spynode connection")
		c.conn.Close()
		c.conn = nil
	}
	c.lock.Unlock()

	c.handlerChannelLock.Lock()
	if c.handlerChannelIsOpen {
		c.handlerChannelIsOpen = false
		close(c.handlerChannel)
	}
	c.handlerChannelLock.Unlock()
}

func (c *RemoteClient) addHandlerMessage(ctx context.Context, msg MessagePayload) {
	c.handlerChannelLock.Lock()
	if c.handlerChannelIsOpen {
		c.handlerChannel <- msg
	}
	c.handlerChannelLock.Unlock()
}

// ping sends pings to keep the connection alive.
func (c *RemoteClient) ping(ctx context.Context) error {
	sinceLastPing := 0
	for {
		c.closeRequestedLock.Lock()
		stop := c.closeRequested
		c.closeRequestedLock.Unlock()

		if stop {
			return nil
		}

		c.lock.Lock()
		conn := c.conn
		c.lock.Unlock()

		if conn == nil {
			return nil // connection closed
		}

		sinceLastPing++
		if sinceLastPing >= 500 {
			timeStamp := uint64(time.Now().UnixNano())
			message := &Message{
				Payload: &Ping{TimeStamp: timeStamp},
			}
			if err := message.Serialize(conn); err != nil {
				return errors.Wrap(err, "send message")
			}
			sinceLastPing = 0
			logger.InfoWithFields(ctx, []logger.Field{
				logger.Float64("timestamp", float64(timeStamp)/1000000000.0),
			}, "Sent ping")
		}

		time.Sleep(200 * time.Millisecond)
	}
}

// listen listens for incoming messages.
func (c *RemoteClient) listen(ctx context.Context) error {
	for {
		c.lock.Lock()
		conn := c.conn
		c.lock.Unlock()

		if conn == nil {
			logger.Info(ctx, "Connection closed")
			return nil // connection closed
		}

		m := &Message{}
		if err := m.Deserialize(conn); err != nil {
			var returnErr error
			if errors.Cause(err) == io.EOF {
				logger.Info(ctx, "Server disconnected")
			} else {
				logger.Warn(ctx, "Failed to read incoming message : %s", err)
				returnErr = err
			}

			c.close(ctx)

			// Check if we are trying to close
			c.closeRequestedLock.Lock()
			stop := c.closeRequested
			c.closeRequestedLock.Unlock()
			if stop {
				return returnErr
			}

			if _, err := c.connect(ctx); err != nil {
				return errors.Wrap(err, "connect")
			}

			continue
		}

		// Handle message
		switch msg := m.Payload.(type) {
		case *AcceptRegister:
			if !msg.Key.Equal(c.serverSessionKey) {
				logger.Error(ctx, "Wrong server session key returned : got %s, want %s", msg.Key,
					c.serverSessionKey)
				c.close(ctx)
				return ErrWrongKey
			}

			sigHash, err := msg.SigHash(c.hash)
			if err != nil {
				logger.Error(ctx, "Failed to create accept sig hash : %s", err)
				c.close(ctx)
				return errors.Wrap(err, "accept sig hash")
			}

			if !msg.Signature.Verify(sigHash.Bytes(), msg.Key) {
				logger.Error(ctx, "Invalid server signature")
				c.close(ctx)
				return ErrBadSignature
			}

			logger.Info(ctx, "Server accepted connection : %+v", msg)
			c.lock.Lock()
			c.accepted = true
			c.lock.Unlock()

			c.addHandlerMessage(ctx, m.Payload)

		case *Tx:
			txid := *msg.Tx.TxHash()
			logger.InfoWithFields(ctx, []logger.Field{
				logger.Stringer("txid", txid),
				logger.Uint64("message_id", msg.ID),
			}, "Received tx")

			if msg.ID == 0 { // non-sequential message (from a request)
				c.requestLock.Lock()
				found := false
				for _, request := range c.getTxRequests {
					if request.txid.Equal(&txid) {
						request.response = m
						found = true
						break
					}
				}
				c.requestLock.Unlock()

				if !found {
					logger.WarnWithFields(ctx, []logger.Field{
						logger.Stringer("txid", txid),
					}, "No matching request found for non-sequential tx")
				}
			} else if c.nextMessageID != msg.ID {
				logger.WarnWithFields(ctx, []logger.Field{
					logger.Uint64("expected_message_id", c.nextMessageID),
					logger.Uint64("message_id", msg.ID),
				}, "Wrong message ID in tx message")
			} else {
				c.nextMessageID = msg.ID + 1
				c.addHandlerMessage(ctx, m.Payload)
			}

		case *TxUpdate:
			logger.InfoWithFields(ctx, []logger.Field{
				logger.Stringer("txid", msg.TxID),
				logger.Uint64("message_id", msg.ID),
			}, "Received tx state")

			if c.nextMessageID != msg.ID {
				logger.WarnWithFields(ctx, []logger.Field{
					logger.Uint64("expected_message_id", c.nextMessageID),
					logger.Uint64("message_id", msg.ID),
				}, "Wrong message ID in tx update message")
			} else {
				c.nextMessageID = msg.ID + 1
				c.addHandlerMessage(ctx, m.Payload)
			}

		case *Headers:
			logger.InfoWithFields(ctx, []logger.Field{
				logger.Int("header_count", len(msg.Headers)),
				logger.Uint32("start_height", msg.StartHeight),
			}, "Received headers")

			requestFound := false
			c.requestLock.Lock()
			for _, request := range c.headerRequests {
				if request.height == int(msg.RequestHeight) {
					request.response = m
					requestFound = true
					break
				}
			}
			c.requestLock.Unlock()

			if !requestFound {
				c.addHandlerMessage(ctx, m.Payload)
			}

		case *InSync:
			logger.Info(ctx, "Received in sync")

			c.addHandlerMessage(ctx, m.Payload)

		case *ChainTip:
			logger.InfoWithFields(ctx, []logger.Field{
				logger.Stringer("hash", msg.Hash),
				logger.Uint32("height", msg.Height),
			}, "Received chain tip")

			c.addHandlerMessage(ctx, m.Payload)

		case *BaseTx:
			txid := *msg.Tx.TxHash()
			logger.InfoWithFields(ctx, []logger.Field{
				logger.Stringer("txid", txid),
			}, "Received base tx")

			c.requestLock.Lock()
			found := false
			for _, request := range c.getTxRequests {
				if request.txid.Equal(&txid) {
					request.response = m
					found = true
					break
				}
			}
			c.requestLock.Unlock()

			if !found {
				logger.WarnWithFields(ctx, []logger.Field{
					logger.Stringer("txid", txid),
				}, "No matching request found for base tx")
			}

		case *Accept:
			if msg.Hash != nil && msg.MessageType == MessageTypeSendTx {
				logger.InfoWithFields(ctx, []logger.Field{
					logger.Stringer("txid", msg.Hash),
				}, "Received accept for send tx")

				found := false
				c.requestLock.Lock()
				for _, request := range c.sendTxRequests {
					request.lock.Lock()
					if request.txid.Equal(msg.Hash) {
						request.response = m
						request.lock.Unlock()
						found = true
						break
					}
					request.lock.Unlock()
				}
				c.requestLock.Unlock()

				if !found {
					logger.WarnWithFields(ctx, []logger.Field{
						logger.Stringer("txid", msg.Hash),
					}, "No matching request found for send tx accept")
				}
			}

		case *Reject:
			c.lock.Lock()
			accepted := c.accepted
			c.lock.Unlock()

			if !accepted {
				// Service rejected registration
				c.close(ctx)
				return errors.Wrap(ErrReject, msg.Message)
			}

			if msg.Hash != nil {
				found := false

				if msg.MessageType == MessageTypeSendTx {
					logger.WarnWithFields(ctx, []logger.Field{
						logger.Stringer("txid", msg.Hash),
					}, "Received reject for send tx : %s", msg.Message)

					c.requestLock.Lock()
					for _, request := range c.sendTxRequests {
						request.lock.Lock()
						if request.txid.Equal(msg.Hash) {
							request.response = m
							request.lock.Unlock()
							found = true
							break
						}
						request.lock.Unlock()
					}
					c.requestLock.Unlock()

					if !found {
						logger.WarnWithFields(ctx, []logger.Field{
							logger.Stringer("txid", msg.Hash),
						}, "No matching request found for send tx reject")
					}
				}

				if msg.MessageType == MessageTypeGetTx {
					logger.WarnWithFields(ctx, []logger.Field{
						logger.Stringer("txid", msg.Hash),
					}, "Received reject for get tx : %s", msg.Message)

					c.requestLock.Lock()
					for _, request := range c.getTxRequests {
						request.lock.Lock()
						if request.txid.Equal(msg.Hash) {
							request.response = m
							request.lock.Unlock()
							found = true
							break
						}
						request.lock.Unlock()
					}
					c.requestLock.Unlock()

					if !found {
						logger.WarnWithFields(ctx, []logger.Field{
							logger.Stringer("txid", msg.Hash),
						}, "No matching request found for get tx reject")
					}
				}
			}

		case *Ping:
			logger.InfoWithFields(ctx, []logger.Field{
				logger.Float64("timestamp", float64(msg.TimeStamp)/1000000000.0),
			}, "Received ping")

		default:
			logger.Error(ctx, "Unknown message type : %d", msg.Type())

		}
	}
}

func (c *RemoteClient) handle(ctx context.Context) error {
	for msg := range c.handlerChannel {
		switch m := msg.(type) {
		case *AcceptRegister, *ChainTip:
			c.handlerLock.Lock()
			for _, handler := range c.handlers {
				handler.HandleMessage(ctx, msg)
			}
			c.handlerLock.Unlock()

		case *InSync:
			c.handlerLock.Lock()
			for _, handler := range c.handlers {
				handler.HandleInSync(ctx)
			}
			c.handlerLock.Unlock()

		case *Headers:
			c.handlerLock.Lock()
			for _, handler := range c.handlers {
				handler.HandleHeaders(ctx, m)
			}
			c.handlerLock.Unlock()

		case *TxUpdate:
			c.handlerLock.Lock()
			for _, handler := range c.handlers {
				handler.HandleTxUpdate(ctx, m)
			}
			c.handlerLock.Unlock()

		case *Tx:
			c.handlerLock.Lock()
			for _, handler := range c.handlers {
				handler.HandleTx(ctx, m)
			}
			c.handlerLock.Unlock()
		}
	}

	return nil
}

// generateSession generates session keys from root keys.
func (c *RemoteClient) generateSession() error {
	for { // loop through any out of range keys
		var err error

		// Generate random hash
		c.hash, err = bitcoin.GenerateSeedValue()
		if err != nil {
			return errors.Wrap(err, "generate hash")
		}

		// Derive session keys
		c.serverSessionKey, err = bitcoin.NextPublicKey(c.config.ServerKey, c.hash)
		if err != nil {
			if errors.Cause(err) == bitcoin.ErrOutOfRangeKey {
				continue // try with a new hash
			}
			return errors.Wrap(err, "next public key")
		}

		c.sessionKey, err = bitcoin.NextKey(c.config.ClientKey, c.hash)
		if err != nil {
			if errors.Cause(err) == bitcoin.ErrOutOfRangeKey {
				continue // try with a new hash
			}
			return errors.Wrap(err, "next key")
		}

		return nil
	}
}
