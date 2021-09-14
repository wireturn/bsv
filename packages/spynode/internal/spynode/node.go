package spynode

import (
	"context"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tokenized/pkg/bitcoin"
	"github.com/tokenized/pkg/logger"
	"github.com/tokenized/pkg/storage"
	"github.com/tokenized/pkg/wire"
	"github.com/tokenized/spynode/internal/handlers"
	"github.com/tokenized/spynode/internal/platform/config"
	"github.com/tokenized/spynode/internal/state"
	internalStorage "github.com/tokenized/spynode/internal/storage"
	"github.com/tokenized/spynode/pkg/client"

	"github.com/pkg/errors"
)

const (
	SubSystem = "SpyNode" // For logger
)

type TxCount struct {
	tx    *wire.MsgTx
	count int
}

// OutputFetcher provides a method to fetch transaction outputs from an external source.
type OutputFetcher interface {
	GetOutputs(context.Context, []wire.OutPoint) ([]bitcoin.UTXO, error)
}

// TxFetcher provides a method to fetch transactions from an external source.
type TxFetcher interface {
	GetTx(context.Context, bitcoin.Hash32) (*wire.MsgTx, error)
}

// Node is the main object for spynode.
type Node struct {
	config          config.Config                      // Configuration
	state           *state.State                       // Non-persistent data
	store           storage.Storage                    // Persistent data
	peers           *internalStorage.PeerRepository    // Peer data
	blocks          *internalStorage.BlockRepository   // Block data
	blockRefeeder   handlers.BlockRefeeder             // Reprocess older blocks
	txs             *internalStorage.TxRepository      // Tx data
	reorgs          *internalStorage.ReorgRepository   // Reorg data
	txTracker       *state.TxTracker                   // Tracks tx requests to ensure all txs are received
	memPool         *state.MemPool                     // Tracks which txs have been received and checked
	messageHandlers map[string]handlers.MessageHandler // Handlers for messages from trusted node
	connection      net.Conn                           // Connection to trusted node
	outgoing        MessageChannel                     // Channel for messages to send to trusted node
	handlers        []client.Handler                   // Receive data and notifications about transactions
	untrustedNodes  []*UntrustedNode                   // Randomized peer connections to monitor for double spends
	addresses       map[string]time.Time               // Recently used peer addresses
	unconfTxChannel handlers.TxChannel                 // Channel for directly handled txs so they don't lock the calling thread
	broadcastLock   sync.Mutex
	broadcastTxs    []TxCount // Txs to transmit to nodes upon connection
	needsRestart    bool
	hardStop        bool
	stopping        bool
	stopped         bool
	scanning        bool
	attempts        int // Count of re-connect attempts without completing handshake.
	lock            sync.Mutex
	untrustedLock   sync.Mutex
	blockLock       sync.Mutex

	txFetcher     TxFetcher
	outputFetcher OutputFetcher

	pushDataHashes []bitcoin.Hash20
	pushDataLock   sync.Mutex

	sendContracts bool
	sendHeaders   bool

	// These counts are used to monitor the number of threads active in specific categories.
	// They are used to stop the incoming threads before stopping the processing threads to
	//   prevent the incoming threads from filling channels and getting locked.
	incomingCount   uint32
	processingCount uint32
	untrustedCount  uint32
}

// NewNode creates a new node.
// See handlers/handlers.go for the listener interface definitions.
func NewNode(config config.Config, store storage.Storage, txFetcher TxFetcher,
	outputFetcher OutputFetcher) *Node {
	result := Node{
		config:         config,
		state:          state.NewState(),
		store:          store,
		peers:          internalStorage.NewPeerRepository(store),
		blocks:         internalStorage.NewBlockRepository(config, store),
		txs:            internalStorage.NewTxRepository(store),
		reorgs:         internalStorage.NewReorgRepository(store),
		txTracker:      state.NewTxTracker(),
		memPool:        state.NewMemPool(),
		handlers:       make([]client.Handler, 0),
		untrustedNodes: make([]*UntrustedNode, 0),
		addresses:      make(map[string]time.Time),
		needsRestart:   false,
		hardStop:       false,
		stopping:       false,
		stopped:        false,
		txFetcher:      txFetcher,
		outputFetcher:  outputFetcher,
	}

	atomic.StoreUint32(&result.incomingCount, 0)
	atomic.StoreUint32(&result.processingCount, 0)
	atomic.StoreUint32(&result.untrustedCount, 0)

	return &result
}

// Client interface
func (node *Node) RegisterHandler(handler client.Handler) {
	node.lock.Lock()
	defer node.lock.Unlock()

	node.handlers = append(node.handlers, handler)
}

func (node *Node) GetTx(ctx context.Context, txid bitcoin.Hash32) (*wire.MsgTx, error) {
	// Check storage
	clientTx, err := internalStorage.FetchTxState(ctx, node.store, txid)
	if err == nil {
		return clientTx.Tx, nil
	}

	if errors.Cause(err) != storage.ErrNotFound {
		return nil, errors.Wrap(err, "fetch tx state")
	}

	// Request from fetcher
	return node.txFetcher.GetTx(ctx, txid)
}

func (node *Node) GetOutputs(ctx context.Context,
	outpoints []wire.OutPoint) ([]bitcoin.UTXO, error) {
	return node.outputFetcher.GetOutputs(ctx, outpoints)
}

func (node *Node) SendTx(ctx context.Context, tx *wire.MsgTx) error {
	if err := node.BroadcastTx(ctx, tx); err != nil {
		return errors.Wrap(err, "broadcast")
	}

	node.unconfTxChannel.Add(handlers.TxData{
		Msg:             tx,
		Trusted:         true,
		Safe:            true,
		ConfirmedHeight: -1,
	})
	return nil
}

func (node *Node) SendTxAndMarkOutputs(ctx context.Context, tx *wire.MsgTx,
	indexes []uint32) error {
	return errors.New("Not implemented")
}

func (node *Node) SubscribePushDatas(ctx context.Context, pushDatas [][]byte) error {
	node.pushDataLock.Lock()
	defer node.pushDataLock.Unlock()

	for _, pd := range pushDatas {
		var hash bitcoin.Hash20
		if len(pd) == 20 {
			copy(hash[:], pd)
		} else {
			copy(hash[:], bitcoin.Hash160(pd))
		}

		node.pushDataHashes = append(node.pushDataHashes, hash)
	}

	return nil
}

func (node *Node) UnsubscribePushDatas(ctx context.Context, pushDatas [][]byte) error {
	node.pushDataLock.Lock()
	defer node.pushDataLock.Unlock()

	for _, pd := range pushDatas {
		hash := pushDataToHash(pd)
		for i, epd := range node.pushDataHashes {
			if hash.Equal(&epd) {
				node.pushDataHashes = append(node.pushDataHashes[:i], node.pushDataHashes[i+1:]...)
				break
			}
		}
	}

	return nil
}

func (node *Node) SubscribeTx(context.Context, bitcoin.Hash32, []uint32) error {
	return errors.New("Not implemented")
}

func (node *Node) UnsubscribeTx(context.Context, bitcoin.Hash32, []uint32) error {
	return errors.New("Not implemented")
}

func (node *Node) SubscribeOutputs(context.Context, []*wire.OutPoint) error {
	return errors.New("Not implemented")
}

func (node *Node) UnsubscribeOutputs(context.Context, []*wire.OutPoint) error {
	return errors.New("Not implemented")
}

func (node *Node) SubscribeContracts(ctx context.Context) error {
	logger.Info(ctx, "Subscribing to contracts")
	node.lock.Lock()
	defer node.lock.Unlock()

	node.sendContracts = true
	return nil
}

func (node *Node) UnsubscribeContracts(ctx context.Context) error {
	logger.Info(ctx, "Unsubscribing from contracts")
	node.lock.Lock()
	defer node.lock.Unlock()

	node.sendContracts = false
	return nil
}

func (node *Node) SubscribeHeaders(ctx context.Context) error {
	logger.Info(ctx, "Subscribing to headers")
	node.lock.Lock()
	defer node.lock.Unlock()

	node.sendHeaders = true
	return nil
}

func (node *Node) UnsubscribeHeaders(ctx context.Context) error {
	logger.Info(ctx, "Unsubscribing from headers")
	node.lock.Lock()
	defer node.lock.Unlock()

	node.sendHeaders = false
	return nil
}

func (node *Node) Ready(ctx context.Context, nextMessageID uint64) error {
	node.lock.Lock()
	defer node.lock.Unlock()

	return nil
}

func (node *Node) NextMessageID() uint64 {
	return 0
}

// SetupRetry configures the maximum connection retries and delay in milliseconds between each
//   attempt.
func (node *Node) SetupRetry(max, delay int) {
	node.config.MaxRetries = max
	node.config.RetryDelay = delay
}

// load loads the data for the node.
// Must be called after adding filter(s), but before Run()
func (node *Node) load(ctx context.Context) error {
	ctx = logger.ContextWithLogSubSystem(ctx, SubSystem)
	if err := node.peers.Load(ctx); err != nil {
		return err
	}

	if err := node.blocks.Load(ctx); err != nil {
		return err
	}
	node.state.SetLastHash(*node.blocks.LastHash())
	logger.Info(ctx, "Loaded blocks to height %d", node.blocks.LastHeight())
	startHeight, exists := node.blocks.Height(&node.config.StartHash)
	if exists {
		node.state.SetStartHeight(startHeight)
		logger.Info(ctx, "Start block height %d", startHeight)
	} else {
		logger.Info(ctx, "Start block not found yet")
	}

	if err := node.txs.Load(ctx); err != nil {
		return err
	}

	node.messageHandlers = handlers.NewTrustedMessageHandlers(ctx, node.config, node.state,
		node.peers, node.blocks, &node.blockRefeeder, node.txs, node.reorgs, node.txTracker,
		node.memPool, &node.unconfTxChannel, node.handlers)
	return nil
}

// Run runs the node.
// Doesn't stop until there is a failure or Stop() is called.
func (node *Node) Run(ctx context.Context) error {
	ctx = logger.ContextWithLogSubSystem(ctx, SubSystem)

	var err error = nil
	if err = node.load(ctx); err != nil {
		return err
	}

	initial := true
	for !node.isStopping() {
		if node.attempts != 0 {
			time.Sleep(time.Duration(node.config.RetryDelay) * time.Millisecond)
		}
		if node.attempts > node.config.MaxRetries {
			logger.Error(ctx, "SpyNodeAborted trusted connection to %s", node.config.NodeAddress)
		}
		node.attempts++

		if initial {
			logger.Verbose(ctx, "Connecting to %s", node.config.NodeAddress)
		} else {
			logger.Verbose(ctx, "Re-connecting to %s", node.config.NodeAddress)
		}
		initial = false
		if err = node.connect(ctx); err != nil {
			logger.Error(ctx, "SpyNodeFailed trusted connection to %s : %s",
				node.config.NodeAddress, err)
			continue
		}

		node.outgoing.Open(100)
		node.unconfTxChannel.Open(100)

		// Queue version message to start handshake
		version := buildVersionMsg(node.config.UserAgent, int32(node.blocks.LastHeight()))
		node.outgoing.Add(version)

		go func() {
			atomic.AddUint32(&node.incomingCount, 1)                // increment
			defer atomic.AddUint32(&node.incomingCount, ^uint32(0)) // decrement
			node.monitorIncoming(ctx)
			logger.Verbose(ctx, "Monitor incoming finished")
		}()

		go func() {
			atomic.AddUint32(&node.processingCount, 1)                // increment
			defer atomic.AddUint32(&node.processingCount, ^uint32(0)) // decrement
			node.monitorRequestTimeouts(ctx)
			logger.Verbose(ctx, "Monitor request timeouts finished")
		}()

		go func() {
			atomic.AddUint32(&node.processingCount, 1)                // increment
			defer atomic.AddUint32(&node.processingCount, ^uint32(0)) // decrement
			node.sendOutgoing(ctx)
			logger.Verbose(ctx, "Send outgoing finished")
		}()

		go func() {
			atomic.AddUint32(&node.processingCount, 1)                // increment
			defer atomic.AddUint32(&node.processingCount, ^uint32(0)) // decrement
			node.processBlocks(ctx)
			logger.Verbose(ctx, "Process blocks finished")
		}()

		go func() {
			atomic.AddUint32(&node.processingCount, 1)                // increment
			defer atomic.AddUint32(&node.processingCount, ^uint32(0)) // decrement
			node.processUnconfirmedTxs(ctx)
			logger.Verbose(ctx, "Process unconfirmed txs finished")
		}()

		go func() {
			atomic.AddUint32(&node.incomingCount, 1)                // increment
			defer atomic.AddUint32(&node.incomingCount, ^uint32(0)) // decrement
			node.checkTxDelays(ctx)
			logger.Verbose(ctx, "Check tx delays finished")
		}()

		if node.config.UntrustedCount == 0 {
			logger.Verbose(ctx, "Monitor untrusted not started")
		} else {
			go func() {
				atomic.AddUint32(&node.incomingCount, 1)                // increment
				defer atomic.AddUint32(&node.incomingCount, ^uint32(0)) // decrement
				node.monitorUntrustedNodes(ctx)
				logger.Verbose(ctx, "Monitor untrusted finished")
			}()
		}

		// Send empty accept register message to handlers so they know spynode is started.
		for _, handler := range node.handlers {
			handler.HandleMessage(ctx, &client.AcceptRegister{})
		}

		// Block until goroutines finish as a result of Stop()

		// Phased shutdown
		for !node.isStopping() {
			time.Sleep(100 * time.Millisecond)
		}

		logger.Info(ctx, "Stopping")

		node.txTracker.Stop() // This will reduce network messages

		node.lock.Lock()
		if node.connection != nil {
			node.connection.Close() // This should stop the monitorIncoming thread
			node.connection = nil
		}
		node.lock.Unlock()

		// Wait for incoming threads to stop.
		// We have to be sure that we stop writing to channels before we stop reading from channels or
		//   we can get stuck in a lock trying to write to a full channel.
		waitCount := 0
		for {
			time.Sleep(100 * time.Millisecond)
			incomingCount := atomic.LoadUint32(&node.incomingCount)
			if incomingCount == 0 {
				break
			}
			if waitCount > 30 { // 3 seconds
				logger.Info(ctx, "Waiting for incoming to stop : %d", incomingCount)
				waitCount = 0
			}
			waitCount++
		}
		logger.Info(ctx, "Incoming threads stopped")

		// Close the channels to stop the processing threads.
		node.outgoing.Close()
		node.unconfTxChannel.Close()

		// Wait for processing threads to stop.
		waitCount = 0
		for {
			time.Sleep(100 * time.Millisecond)
			processingCount := atomic.LoadUint32(&node.processingCount)
			if processingCount == 0 {
				break
			}
			if waitCount > 30 { // 3 seconds
				logger.Info(ctx, "Waiting for processing to stop : %d", processingCount)
				waitCount = 0
			}
			waitCount++
		}
		logger.Info(ctx, "Processing threads stopped")

		// Save block repository
		logger.Verbose(ctx, "Saving")
		node.blocks.Save(ctx)
		node.txs.Save(ctx)
		node.peers.Save(ctx)

		node.lock.Lock()
		if !node.needsRestart || node.hardStop {
			node.lock.Unlock()
			break
		}

		logger.Verbose(ctx, "Restarting")
		node.needsRestart = false
		node.stopping = false
		node.lock.Unlock()
		node.state.Reset()
	}

	node.lock.Lock()
	node.stopped = true
	node.lock.Unlock()
	logger.Verbose(ctx, "Stopped")

	return err
}

func (node *Node) isStopped() bool {
	node.lock.Lock()
	defer node.lock.Unlock()

	return node.stopped
}

func (node *Node) isStopping() bool {
	node.lock.Lock()
	defer node.lock.Unlock()

	return node.stopping
}

// Stop closes the connection and causes Run() to return.
func (node *Node) Stop(ctx context.Context) error {
	ctx = logger.ContextWithLogSubSystem(ctx, SubSystem)

	node.lock.Lock()
	stopped := node.stopped
	node.lock.Unlock()
	if stopped {
		logger.Verbose(ctx, "Already stopped")
		return nil
	}

	node.hardStop = true
	err := node.requestStop(ctx)
	count := 0
	for !node.isStopped() {
		time.Sleep(100 * time.Millisecond)
		if count > 30 { // 3 seconds
			logger.Info(ctx, "Waiting for spynode to stop")
			count = 0
		}
		count++
	}
	return err
}

func (node *Node) requestStop(ctx context.Context) error {
	logger.Verbose(ctx, "Requesting stop")

	node.lock.Lock()
	defer node.lock.Unlock()

	if node.stopped {
		logger.Verbose(ctx, "Already stopped")
		return nil
	} else if node.stopping {
		logger.Verbose(ctx, "Already stopping")
		return nil
	}

	node.stopping = true
	return nil
}

func (node *Node) OutgoingCount() int {
	node.untrustedLock.Lock()
	defer node.untrustedLock.Unlock()

	result := 0
	for _, untrusted := range node.untrustedNodes {
		if untrusted.IsReady() {
			result++
		}
	}
	return result
}

// BroadcastTx broadcasts a tx to the network.
func (node *Node) BroadcastTx(ctx context.Context, tx *wire.MsgTx) error {
	ctx = logger.ContextWithLogSubSystem(ctx, SubSystem)
	logger.Info(ctx, "Broadcasting tx : %s", tx.TxHash())

	if node.isStopping() { // TODO Resolve issue when node is restarting
		return errors.New("Node inactive")
	}

	// Wait for ready
	for i := 0; i < 100; i++ {
		if node.isStopping() { // TODO Resolve issue when node is restarting
			return errors.New("Node inactive")
		}
		if node.state.IsReady() {
			break
		}

		time.Sleep(250 * time.Millisecond)
	}

	count := 1

	// Send to trusted node
	if !node.queueOutgoing(tx) {
		return errors.New("Node inactive")
	}

	// Send to untrusted nodes
	node.untrustedLock.Lock()
	for _, untrusted := range node.untrustedNodes {
		if untrusted.IsReady() {
			if err := untrusted.BroadcastTxs(ctx, []*wire.MsgTx{tx}); err != nil {
				logger.Warn(ctx, "Failed to broadcast tx to untrusted : %s", err)
			} else {
				count++
			}
		}
	}
	node.untrustedLock.Unlock()

	if count < node.config.ShotgunCount {
		node.broadcastLock.Lock()
		node.broadcastTxs = append(node.broadcastTxs, TxCount{tx: tx, count: count})
		node.broadcastLock.Unlock()
	}
	return nil
}

func (node *Node) BroadcastIsComplete(ctx context.Context) bool {
	ctx = logger.ContextWithLogSubSystem(ctx, SubSystem)
	if node.isStopping() {
		return true
	}

	node.broadcastLock.Lock()
	count := len(node.broadcastTxs)
	node.broadcastLock.Unlock()

	logger.Info(ctx, "%d broadcast txs remaining", count)
	return count == 0
}

func (node *Node) IsReady(ctx context.Context) bool {
	return node.state.IsReady()
}

// BroadcastTxUntrustedOnly broadcasts a tx to the network.
func (node *Node) BroadcastTxUntrustedOnly(ctx context.Context, tx *wire.MsgTx) error {
	ctx = logger.ContextWithLogSubSystem(ctx, SubSystem)
	logger.Info(ctx, "Broadcasting tx to untrusted only : %s", tx.TxHash())

	if node.isStopping() { // TODO Resolve issue when node is restarting
		return errors.New("Node inactive")
	}

	// Send to untrusted nodes
	node.untrustedLock.Lock()
	defer node.untrustedLock.Unlock()

	for _, untrusted := range node.untrustedNodes {
		if untrusted.IsReady() {
			untrusted.BroadcastTxs(ctx, []*wire.MsgTx{tx})
		}
	}
	return nil
}

// Scan opens a lot of connections at once to try to find peers.
func (node *Node) Scan(ctx context.Context, connections int) error {
	ctx = logger.ContextWithLogSubSystem(ctx, SubSystem)

	if err := node.load(ctx); err != nil {
		return err
	}

	if err := node.scan(ctx, connections, 1); err != nil {
		return err
	}

	return node.peers.Save(ctx)
}

// AddPeer adds a peer to the database with a specific score.
func (node *Node) AddPeer(ctx context.Context, address string, score int) error {
	ctx = logger.ContextWithLogSubSystem(ctx, SubSystem)
	if err := node.load(ctx); err != nil {
		return err
	}
	if _, err := node.peers.Add(ctx, address); err != nil {
		return err
	}

	if !node.peers.UpdateScore(ctx, address, int32(score)) {
		return errors.New("Failed to update score")
	}

	return node.peers.Save(ctx)
}

// sendOutgoing waits for and sends outgoing messages
//
// This is a blocking function that will run forever, so it should be run
// in a goroutine.
func (node *Node) sendOutgoing(ctx context.Context) {
	// Wait for outgoing messages on channel
	for msg := range node.outgoing.Channel {
		node.lock.Lock()
		connection := node.connection
		node.lock.Unlock()

		if connection == nil {
			continue // Keep clearing the channel
		}

		tx, ok := msg.(*wire.MsgTx)
		if ok {
			logger.Verbose(ctx, "Sending Tx : %s", tx.TxHash().String())
		}

		if err := sendAsync(ctx, connection, msg, wire.BitcoinNet(node.config.Net)); err != nil {
			logger.Error(ctx, "SpyNodeFailed to send %s message : %s", msg.Command(), err)
			// don't break out of the loop because we need to continue emptying the channel.
			// otherwise any processing adding to the channel can get locked on a write.
			node.restart(ctx)
		}
	}
}

// handleMessage Processes an incoming message
func (node *Node) handleMessage(ctx context.Context, msg wire.Message) error {
	if node.isStopping() {
		return nil
	}

	handler, ok := node.messageHandlers[msg.Command()]
	if !ok {
		// no handler for this command
		return nil
	}

	responses, err := handler.Handle(ctx, msg)
	if err != nil {
		logger.Warn(ctx, "Failed to handle [%s] message : %s", msg.Command(), err)
		return nil
	}

	// Queue messages to be sent in response
	for _, response := range responses {
		if !node.queueOutgoing(response) {
			break
		}
	}

	return nil
}

// CleanupBlock is called when a block is being processed.
// Implements handlers.BlockProcessor interface
// It is responsible for any cleanup as a result of a block.
func (node *Node) CleanupBlock(ctx context.Context, txids []*bitcoin.Hash32) error {
	// logger.Debug(ctx, "Cleaning up after block")

	node.txTracker.RemoveList(ctx, txids)

	node.untrustedLock.Lock()
	defer node.untrustedLock.Unlock()

	for _, untrusted := range node.untrustedNodes {
		untrusted.CleanupBlock(ctx, txids)
	}

	return nil
}

func (node *Node) connect(ctx context.Context) error {
	conn, err := net.Dial("tcp", node.config.NodeAddress)
	if err != nil {
		return err
	}

	node.lock.Lock()
	node.connection = conn
	node.lock.Unlock()
	node.state.MarkConnected()
	node.peers.UpdateTime(ctx, node.config.NodeAddress)
	return nil
}

// monitorIncoming processes incoming messages.
//
// This is a blocking function that will run forever, so it should be run
// in a goroutine.
func (node *Node) monitorIncoming(ctx context.Context) {
	for !node.isStopping() {
		if err := node.check(ctx); err != nil {
			logger.Error(ctx, "SpyNodeAborted check : %s", err.Error())
			node.requestStop(ctx)
			break
		}

		node.lock.Lock()
		connection := node.connection
		node.lock.Unlock()

		if connection == nil {
			break
		}

		// read new messages, blocking
		_, msg, _, err := wire.ReadMessageN(connection, wire.ProtocolVersion,
			wire.BitcoinNet(node.config.Net))
		if err != nil {
			wireError, ok := errors.Cause(err).(*wire.MessageError)
			if ok {
				switch wireError.Type {
				case wire.MessageErrorUnknownCommand:
					logger.Verbose(ctx, wireError.Error())
					continue
				case wire.MessageErrorConnectionClosed:
					if !node.isStopping() {
						logger.Warn(ctx, "SpyNodeFailed : %s", wireError)
						node.restart(ctx)
					}
					return
				default:
					logger.Warn(ctx, "SpyNodeFailed read message (wireError) : %s", wireError)
					node.restart(ctx)
					return
				}
			} else {
				logger.Warn(ctx, "SpyNodeFailed to read message : %s", err)
				node.restart(ctx)
				return
			}
		}

		if err := node.handleMessage(ctx, msg); err != nil {
			logger.Error(ctx, "SpyNodeAborted handling [%s] message : %s", msg.Command(),
				err.Error())
			node.restart(ctx)
			return
		}
		if msg.Command() == "reject" {
			reject, ok := msg.(*wire.MsgReject)
			if ok {
				logger.Warn(ctx, "(%s) Reject message : %s - %s", node.config.NodeAddress,
					reject.Reason, reject.Hash.String())
			}
		}
	}
}

// restart triggers a disconnection and reconnection from the node.
func (node *Node) restart(ctx context.Context) {
	if node.isStopping() {
		return
	}
	node.lock.Lock()
	node.needsRestart = true
	node.lock.Unlock()
	node.requestStop(ctx)
}

// queueOutgoing adds the message to the queue to be sent. It returns true if the message was added.
func (node *Node) queueOutgoing(msg wire.Message) bool {
	if node.isStopping() {
		return false
	}
	err := node.outgoing.Add(msg)
	return err == nil
}

// TransmitMessage interface
func (node *Node) TransmitMessage(msg wire.Message) bool {
	return node.queueOutgoing(msg)
}

// check checks the state of spynode and performs state related actions.
func (node *Node) check(ctx context.Context) error {
	if !node.state.VersionReceived() {
		return nil // Still performing handshake
	}

	if !node.state.HandshakeComplete() {
		// Send header request to kick off sync
		headerRequest, err := buildHeaderRequest(ctx, node.state.ProtocolVersion(), node.blocks,
			node.state, 0, 50)
		if err != nil {
			return err
		}

		if node.queueOutgoing(headerRequest) {
			logger.Verbose(ctx, "Requesting headers")
			node.state.MarkHeadersRequested()
			node.state.SetHandshakeComplete()
		}
	}

	// Check sync
	if node.state.IsReady() {
		node.attempts = 0

		if !node.state.SentSendHeaders() {
			// Send sendheaders message to get headers instead of block inventories.
			sendheaders := wire.NewMsgSendHeaders()
			if node.queueOutgoing(sendheaders) {
				node.state.SetSentSendHeaders()
			}
		}

		if !node.state.AddressesRequested() {
			addresses := wire.NewMsgGetAddr()
			if node.queueOutgoing(addresses) {
				node.state.SetAddressesRequested()
			}
		}

		if node.config.RequestMempool && !node.state.MemPoolRequested() {
			// Send mempool request
			// This tells the peer to send inventory of all tx in their mempool.
			mempool := wire.NewMsgMemPool()
			if node.queueOutgoing(mempool) {
				node.state.SetMemPoolRequested()
			}
		} else {
			if !node.state.WasInSync() {
				node.reorgs.ClearActive(ctx)
				node.state.SetWasInSync()
			}

			if !node.state.NotifiedSync() {
				// TODO Add method to wait for mempool to sync
				for _, handler := range node.handlers {
					handler.HandleInSync(ctx)
				}
				node.state.SetNotifiedSync()
			}
		}

		if err := node.txTracker.Check(ctx, node.memPool, node); err != nil {
			return err
		}
		if node.isStopping() {
			return nil
		}
	} else if node.state.HeadersRequested() == nil && node.state.TotalBlockRequestCount() < 5 {
		// Request more headers
		headerRequest, err := buildHeaderRequest(ctx, node.state.ProtocolVersion(), node.blocks,
			node.state, 1, 50)
		if err != nil {
			return err
		}

		if node.queueOutgoing(headerRequest) {
			logger.Verbose(ctx, "Requesting headers after : %s",
				headerRequest.BlockLocatorHashes[0])
			node.state.MarkHeadersRequested()
		}
	}

	return nil
}

// monitorRequestTimeouts monitors for request timeouts.
//
// This is a blocking function that will run forever, so it should be run
// in a goroutine.
func (node *Node) monitorRequestTimeouts(ctx context.Context) {
	for !node.isStopping() {
		node.sleepUntilStop(50) // Only check every 5 seconds

		if err := node.state.CheckTimeouts(); err != nil {
			logger.Error(ctx, "SpyNodeFailed timeouts : %s", err)
			node.restart(ctx)
			break
		}
	}
}

// checkTxDelays monitors txs for when they have passed the safe tx delay without seeing a
//   conflicting tx.
//
// This is a blocking function that will run forever, so it should be run
// in a goroutine.
func (node *Node) checkTxDelays(ctx context.Context) {
	logger.Info(ctx, "Safe tx delay : %d ms", node.config.SafeTxDelay)
	for !node.isStopping() {
		time.Sleep(100 * time.Millisecond)

		if !node.state.IsReady() {
			continue
		}

		// Get newly safe txs
		cutoffTime := time.Now().Add(time.Millisecond * -time.Duration(node.config.SafeTxDelay))
		txids, err := node.txs.GetNewSafe(ctx, node.memPool, cutoffTime)
		if err != nil {
			logger.Error(ctx, "SpyNodeFailed GetNewSafe : %s", err)
			node.restart(ctx)
			break
		}

		for _, txid := range txids {
			txState, err := internalStorage.FetchTxState(ctx, node.store, txid)
			if err != nil {
				logger.Error(ctx, "SpyNodeFailed fetch tx state : %s", err)
				continue
			}

			if txState.State.UnSafe || txState.State.Cancelled {
				continue
			}

			txState.State.Safe = true

			if err := internalStorage.SaveTxState(ctx, node.store, txState); err != nil {
				logger.Error(ctx, "SpyNodeFailed save tx state : %s", err)
				continue
			}

			// Send update
			update := &client.TxUpdate{
				TxID:  txid,
				State: txState.State,
			}
			for _, handler := range node.handlers {
				handler.HandleTxUpdate(ctx, update)
			}
		}
	}
}

// Scan opens a lot of connections at once to try to find peers.
func (node *Node) scan(ctx context.Context, connections, uncheckedCount int) error {
	if node.scanning {
		return nil
	}
	node.scanning = true

	ctx = logger.ContextWithLogTrace(ctx, "scan")

	peers, err := node.peers.GetUnchecked(ctx)
	if err != nil {
		return err
	}
	logger.Verbose(ctx, "Found %d peers with no score", len(peers))
	if len(peers) < uncheckedCount {
		return nil // Not enough unchecked peers to run a scan
	}

	logger.Verbose(ctx, "Scanning %d peers", connections)

	count := 0
	unodes := make([]*UntrustedNode, 0, connections)
	wg := sync.WaitGroup{}
	var address string
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	for !node.isStopping() && count < connections && len(peers) > 0 {
		// Pick peer randomly
		random := seed.Intn(len(peers))
		address = peers[random].Address

		// Remove this address and try again
		peers = append(peers[:random], peers[random+1:]...)

		// Attempt connection
		newNode := NewUntrustedNode(address, node.config, node.state, node.store, node.peers,
			node.blocks, node.txs, node.memPool, &node.unconfTxChannel, node.handlers, node, true)
		unodes = append(unodes, newNode)
		wg.Add(1)
		go func() {
			defer wg.Done()
			newNode.Run(ctx)
		}()
		count++
	}

	node.sleepUntilStop(100) // Wait for handshake

	for _, unode := range unodes {
		unode.Stop(ctx)
	}

	logger.Verbose(ctx, "Waiting for %d scanning nodes to stop", len(unodes))
	wg.Wait()
	node.scanning = false
	logger.Verbose(ctx, "Finished scanning")
	return nil
}

// monitorUntrustedNodes monitors untrusted nodes.
// Attempt to keep the specified number running.
// Watch for when they become inactive and replace them.
//
// This is a blocking function that will run forever, so it should be run
// in a goroutine.
func (node *Node) monitorUntrustedNodes(ctx context.Context) {
	for !node.isStopping() {
		if !node.state.IsReady() {
			node.sleepUntilStop(5)
			continue
		}

		node.broadcastLock.Lock()
		broadcastCount := len(node.broadcastTxs)
		node.broadcastLock.Unlock()

		if broadcastCount == 0 {
			node.scan(ctx, 1000, 1)
		}
		if node.isStopping() {
			break
		}

		node.untrustedLock.Lock()
		if !node.state.IsReady() {
			node.untrustedLock.Unlock()
			node.sleepUntilStop(5)
			continue
		}

		// Check for inactive
		for {
			removed := false
			for i, untrusted := range node.untrustedNodes {
				if !untrusted.IsActive() {
					// Remove
					node.untrustedNodes = append(node.untrustedNodes[:i], node.untrustedNodes[i+1:]...)
					removed = true
					break
				}
			}

			if !removed {
				break
			}
		}

		count := len(node.untrustedNodes)
		verifiedCount := 0
		for _, untrusted := range node.untrustedNodes {
			if untrusted.untrustedState.IsReady() {
				verifiedCount++
			}
		}

		desiredCount := node.config.UntrustedCount

		node.untrustedLock.Unlock()

		var txs []*wire.MsgTx
		sentCount := 0
		node.broadcastLock.Lock()
		if len(node.broadcastTxs) > 0 {
			desiredCount = node.config.ShotgunCount
			txs = make([]*wire.MsgTx, 0, len(node.broadcastTxs))
			for _, btx := range node.broadcastTxs {
				txs = append(txs, btx.tx)
			}
		}
		node.broadcastLock.Unlock()

		if verifiedCount < desiredCount {
			logger.Debug(ctx, "Untrusted connections : %d", verifiedCount)
		}

		if count < desiredCount/2 {
			// Try for peers with a good score
			for !node.isStopping() && count < desiredCount/2 {
				if node.addUntrustedNode(ctx, 5, txs) {
					count++
					sentCount++
				} else {
					break
				}
			}
		}

		// Try for peers with a score above zero
		for !node.isStopping() && count < desiredCount {
			if node.addUntrustedNode(ctx, 1, txs) {
				count++
				sentCount++
			} else {
				break
			}
		}

		if node.isStopping() {
			break
		}

		if sentCount > 0 {
			for _, tx := range txs {
				node.broadcastLock.Lock()
				for i, btx := range node.broadcastTxs {
					if tx == btx.tx {
						node.broadcastTxs[i].count += sentCount
						if node.broadcastTxs[i].count > node.config.ShotgunCount {
							// tx has been sent to enough nodes. remove it
							node.broadcastTxs = append(node.broadcastTxs[:i],
								node.broadcastTxs[i+1:]...)
						}
						break
					}
				}
				node.broadcastLock.Unlock()
			}
		}

		node.sleepUntilStop(20) // Only check every 2 seconds
	}

	// Stop all
	node.untrustedLock.Lock()
	for _, untrusted := range node.untrustedNodes {
		logger.Verbose(ctx, "Stopping untrusted node : %s", untrusted.address)
		untrusted.Stop(ctx)
	}
	node.untrustedLock.Unlock()

	waitCount := 0
	for {
		time.Sleep(100 * time.Millisecond)
		untrustedCount := atomic.LoadUint32(&node.untrustedCount)
		if untrustedCount == 0 {
			break
		}
		if waitCount > 30 { // 3 seconds
			logger.Info(ctx, "Waiting for %d untrusted nodes to stop", untrustedCount)

			node.untrustedLock.Lock()
			for _, untrusted := range node.untrustedNodes {
				if untrusted.IsActive() {
					logger.Info(ctx, "Waiting for untrusted node : %s", untrusted.address)
				}
			}
			node.untrustedLock.Unlock()

			waitCount = 0
		}
		waitCount++
	}
}

// addUntrustedNode adds a new untrusted node.
// Returns true if a new node connection was attempted
func (node *Node) addUntrustedNode(ctx context.Context, minScore int32,
	txs []*wire.MsgTx) bool {

	// Get new address
	// Check we aren't already connected and haven't used it recently
	peers, err := node.peers.Get(ctx, minScore)
	if err != nil {
		return false
	}
	logger.Debug(ctx, "Found %d peers with score %d", len(peers), minScore)

	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	var address string
	for {
		if node.isStopping() || len(peers) == 0 {
			return false
		}

		if len(peers) == 1 {
			if node.checkAddress(ctx, peers[0].Address) {
				address = peers[0].Address
				break
			} else {
				return false
			}
		}

		// Pick one randomly
		random := seed.Intn(len(peers))
		if node.checkAddress(ctx, peers[random].Address) {
			address = peers[random].Address
			break
		}

		// Remove this address and try again
		peers = append(peers[:random], peers[random+1:]...)
	}

	// Attempt connection
	newNode := NewUntrustedNode(address, node.config, node.state, node.store, node.peers,
		node.blocks, node.txs, node.memPool, &node.unconfTxChannel, node.handlers, node,
		false)
	if txs != nil {
		newNode.BroadcastTxs(ctx, txs)
	}
	node.untrustedLock.Lock()
	node.untrustedNodes = append(node.untrustedNodes, newNode)
	if node.isStopping() {
		newNode.Stop(ctx)
	}
	node.untrustedLock.Unlock()
	go func() {
		atomic.AddUint32(&node.untrustedCount, 1) // increment
		newNode.Run(ctx)
		atomic.AddUint32(&node.untrustedCount, ^uint32(0)) // decrement
	}()
	return true
}

// checkAddress checks if an address was recently used.
func (node *Node) checkAddress(ctx context.Context, address string) bool {
	lastUsed, exists := node.addresses[address]
	if exists {
		if time.Now().Sub(lastUsed).Minutes() > 10 {
			// Address hasn't been used for a while
			node.addresses[address] = time.Now()
			return true
		}

		// Address was used recently
		return false
	}

	// Add address
	node.addresses[address] = time.Now()
	return true
}

func (node *Node) sleepUntilStop(deciseconds int) {
	for i := 0; i < deciseconds; i++ {
		if node.isStopping() {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (node *Node) RefeedBlocksFromHeight(ctx context.Context, height int) error {
	hash, err := node.Hash(ctx, height)
	if err != nil {
		return errors.Wrap(err, "get hash")
	}

	node.blockRefeeder.SetHeight(height, *hash)
	return nil
}

// ------------------------------------------------------------------------------------------------
// BitcoinHeaders interface
func (node *Node) LastHeight(ctx context.Context) int {
	return node.blocks.LastHeight()
}

func (node *Node) Hash(ctx context.Context, height int) (*bitcoin.Hash32, error) {
	return node.blocks.Hash(ctx, height)
}

func (node *Node) GetHeaders(ctx context.Context, height, maxCount int) (*client.Headers, error) {
	var headers []*wire.BlockHeader
	startHeight := height
	if height == -1 {
		startHeight = node.blocks.LastHeight()
		if startHeight > maxCount {
			startHeight -= maxCount
		} else {
			startHeight = 0
		}
	}
	for i := startHeight; i <= startHeight+maxCount; i++ {
		header, err := node.blocks.Header(ctx, i)
		if err != nil {
			if errors.Cause(err) == internalStorage.ErrInvalidHeight {
				return &client.Headers{}, nil
			}
			return nil, errors.Wrap(err, "header")
		}

		headers = append(headers, header)
	}

	result := &client.Headers{
		RequestHeight: int32(height),
		StartHeight:   uint32(startHeight),
		Headers:       headers,
	}

	return result, nil
}

func (node *Node) BlockHash(ctx context.Context, height int) (*bitcoin.Hash32, error) {
	if height == -1 {
		height = node.blocks.LastHeight()
	}
	return node.blocks.Hash(ctx, height)
}

func (node *Node) Time(ctx context.Context, height int) (uint32, error) {
	return node.blocks.Time(ctx, height)
}
