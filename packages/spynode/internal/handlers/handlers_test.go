package handlers

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/tokenized/pkg/bitcoin"
	"github.com/tokenized/pkg/logger"
	"github.com/tokenized/pkg/storage"
	"github.com/tokenized/pkg/wire"
	"github.com/tokenized/spynode/internal/platform/config"
	"github.com/tokenized/spynode/internal/state"
	handlerStorage "github.com/tokenized/spynode/internal/storage"
	"github.com/tokenized/spynode/pkg/client"

	"github.com/pkg/errors"
)

func TestHandlers(test *testing.T) {
	testBlockCount := 9264
	reorgDepth := 4

	// Setup context
	logConfig := logger.NewConfig(true, false, "")
	// logConfig.IsText = true
	ctx := logger.ContextWithLogConfig(context.Background(), logConfig)

	// For logging to test from within functions
	ctx = context.WithValue(ctx, 999, test)
	// Use this to get the test value from within non-test code.
	// testValue := ctx.Value(999)
	// test, ok := testValue.(*testing.T)
	// if ok {
	// test.Logf("Test Debug Message")
	// }

	// Setup storage
	store := storage.NewMockStorage()

	// Setup config
	startHash := "0000000000000000000000000000000000000000000000000000000000000000"
	config, err := config.NewConfig(bitcoin.MainNet, true, "test", "Tokenized Test", startHash, 8,
		2000, 10, 10, 1000, true)
	if err != nil {
		test.Errorf("Failed to create config : %v", err)
	}

	// Setup state
	st := state.NewState()
	st.SetStartHeight(1)

	// Create peer repo
	peerRepo := handlerStorage.NewPeerRepository(store)
	if err := peerRepo.Load(ctx); err != nil {
		test.Errorf("Failed to initialize peer repo : %v", err)
	}

	// Create block repo
	t := uint32(time.Now().Unix())
	blockRepo := handlerStorage.NewBlockRepository(config, store)
	if err := blockRepo.Initialize(ctx, t); err != nil {
		test.Errorf("Failed to initialize block repo : %v", err)
	}

	// Create tx repo
	txRepo := handlerStorage.NewTxRepository(store)
	// Clear any pre-existing data
	for i := 0; i <= testBlockCount; i++ {
		txRepo.ClearBlock(ctx, i)
	}

	// Create reorg repo
	reorgRepo := handlerStorage.NewReorgRepository(store)

	// TxTracker
	txTracker := state.NewTxTracker()

	// Create mempool
	memPool := state.NewMemPool()

	// Setup handlers
	testHandler := &TestHandler{test: test, state: st, blocks: blockRepo, txs: txRepo, height: 0,
		txTracker: txTracker}
	handlers := []client.Handler{testHandler}

	unconfTxChannel := TxChannel{}
	unconfTxChannel.Open(100)

	// Create handlers
	testMessageHandlers := NewTrustedMessageHandlers(ctx, config, st, peerRepo, blockRepo, nil,
		txRepo, reorgRepo, txTracker, memPool, &unconfTxChannel, handlers)

	test.Logf("Testing Blocks")

	// Build a bunch of headers
	blocks := make([]*wire.MsgBlock, 0, testBlockCount)
	txs := make([]*wire.MsgTx, 0, testBlockCount)
	headersMsg := wire.NewMsgHeaders()
	zeroHash, _ :=
		bitcoin.NewHash32FromStr("0000000000000000000000000000000000000000000000000000000000000000")
	previousHash, err := blockRepo.Hash(ctx, 0)
	if err != nil {
		test.Errorf("Failed to get genesis hash : %s", err)
		return
	}
	st.SetLastHash(*previousHash)
	headerMsgCount := 0
	for i := 0; i < testBlockCount; i++ {
		height := i

		// Create coinbase tx to make a valid block
		tx := wire.NewMsgTx(1)
		outpoint := wire.NewOutPoint(zeroHash, 0xffffffff)
		script := make([]byte, 5)
		script[0] = 4 // push 4 bytes
		// Push 4 byte height
		script[1] = byte((height >> 24) & 0xff)
		script[2] = byte((height >> 16) & 0xff)
		script[3] = byte((height >> 8) & 0xff)
		script[4] = byte((height >> 0) & 0xff)
		input := wire.NewTxIn(outpoint, script)
		tx.AddTxIn(input)
		txs = append(txs, tx)

		merkleRoot := tx.TxHash()
		header := wire.NewBlockHeader(1, previousHash, merkleRoot, 0, 0)
		header.Timestamp = time.Unix(int64(t), 0)
		t += 600
		block := wire.NewMsgBlock(header)
		if err := block.AddTransaction(tx); err != nil {
			test.Errorf("Failed to add tx to block (%d) : %s", height, err)
		}

		blocks = append(blocks, block)
		if headerMsgCount > 1000 {
			// Send headers to handlers
			if err := handleMessage(ctx, testMessageHandlers, headersMsg); err != nil {
				test.Errorf("Failed to process headers message : %v", err)
			}

			headersMsg = wire.NewMsgHeaders()
		}
		if err := headersMsg.AddBlockHeader(header); err != nil {
			test.Errorf("Failed to add header to headers message : %v", err)
		}
		headerMsgCount++
		hash := header.BlockHash()
		previousHash = hash
	}

	if headerMsgCount > 0 {
		// Send headers to handlers
		if err := handleMessage(ctx, testMessageHandlers, headersMsg); err != nil {
			test.Errorf("Failed to process headers message : %v", err)
		}
	}

	// Send corresponding blocks
	if err := sendBlocks(ctx, testMessageHandlers, blocks[:len(blocks)-2], 0,
		testHandler); err != nil {
		test.Errorf("Failed to send block messages : %v", err)
	}

	verify(ctx, test, blocks[:len(blocks)-2], blockRepo, len(blocks)-2)

	test.Logf("Testing Reorg")

	// Cause a reorg
	reorgHeadersMsg := wire.NewMsgHeaders()
	reorgBlocks := make([]*wire.MsgBlock, 0, testBlockCount)
	hash := blocks[testBlockCount-reorgDepth].Header.BlockHash()
	previousHash = hash
	test.Logf("Reorging to (%d) : %s", (testBlockCount-reorgDepth)+1, previousHash)
	for i := testBlockCount - reorgDepth; i < testBlockCount; i++ {
		height := (testBlockCount - reorgDepth) + 1 + i

		// Create coinbase tx to make a valid block
		tx := wire.NewMsgTx(1)
		outpoint := wire.NewOutPoint(zeroHash, 0xffffffff)
		script := make([]byte, 5)
		script[0] = 4 // push 4 bytes
		// Push 4 byte height
		script[1] = byte((height >> 24) & 0xff)
		script[2] = byte((height >> 16) & 0xff)
		script[3] = byte((height >> 8) & 0xff)
		script[4] = byte((height >> 0) & 0xff)
		input := wire.NewTxIn(outpoint, script)
		tx.AddTxIn(input)
		txs = append(txs, tx)

		merkleRoot := tx.TxHash()
		header := wire.NewBlockHeader(int32(wire.ProtocolVersion), previousHash, merkleRoot, 0, 1)
		block := wire.NewMsgBlock(header)
		if err := block.AddTransaction(tx); err != nil {
			test.Errorf(fmt.Sprintf("Failed to add tx to block (%d)", height), err)
		}

		reorgBlocks = append(reorgBlocks, block)
		if err := reorgHeadersMsg.AddBlockHeader(header); err != nil {
			test.Errorf("Failed to add header to reorg headers message : %v", err)
		}
		hash := header.BlockHash()
		previousHash = hash
	}

	// Send reorg headers to handlers
	if err := handleMessage(ctx, testMessageHandlers, reorgHeadersMsg); err != nil {
		test.Errorf("Failed to process reorg headers message : %v", err)
	}

	// Send corresponding reorg blocks
	if err := sendBlocks(ctx, testMessageHandlers, reorgBlocks, (testBlockCount-reorgDepth)+1,
		testHandler); err != nil {
		test.Errorf("Failed to send reorg block messages : %v", err)
	}

	// Check reorg
	activeReorg, err := reorgRepo.GetActive(ctx)
	if err != nil {
		test.Errorf("Failed to get active reorg : %v", err)
	}

	if activeReorg == nil {
		test.Errorf("Failed: No active reorg found")
	}

	err = reorgRepo.ClearActive(ctx)
	if err != nil {
		test.Errorf("Failed to clear active reorg : %v", err)
	}

	activeReorg, err = reorgRepo.GetActive(ctx)
	if err != nil {
		test.Errorf("Failed to get active reorg after clear : %v", err)
	}
	if activeReorg != nil {
		test.Errorf("Failed: Active reorg was not cleared")
	}

	// Update headers array for reorg
	blocks = blocks[:(testBlockCount-reorgDepth)+1]
	for _, hash := range reorgBlocks {
		blocks = append(blocks, hash)
	}

	test.Logf("Block count %d = %d", len(blocks), testBlockCount+1)
	verify(ctx, test, blocks, blockRepo, testBlockCount+1)
}

func handleMessage(ctx context.Context, messageHandlers map[string]MessageHandler,
	msg wire.Message) error {

	h, ok := messageHandlers[msg.Command()]
	if !ok {
		// no handler for this command
		return nil
	}

	_, err := h.Handle(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}

func sendBlocks(ctx context.Context, messageHandlers map[string]MessageHandler,
	blocks []*wire.MsgBlock, startHeight int, handler *TestHandler) error {

	logger.Info(ctx, "sending blocks")

	for i, block := range blocks {
		// Convert from MsgBlock to MsgParseBlock for the handler.
		var buf bytes.Buffer
		if err := block.BtcEncode(&buf, wire.ProtocolVersion); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to encode block (%d) message",
				startHeight+i))
		}

		parseBlock := &wire.MsgParseBlock{}
		if err := parseBlock.BtcDecode(bytes.NewReader(buf.Bytes()),
			wire.ProtocolVersion); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to decode block (%d) message",
				startHeight+i))
		}

		// Send block to message handlers
		if err := handleMessage(ctx, messageHandlers, parseBlock); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to process block (%d) message",
				startHeight+i))
		}

		hash, _ := handler.state.GetNextBlockToRequest()
		if hash != nil {
			logger.Info(ctx, "requesting block : %s", hash)
		}

		if i%8 == 0 {
			if err := handler.ProcessBlocks(ctx); err != nil {
				return errors.Wrap(err, "process blocks")
			}
		}
	}

	if err := handler.ProcessBlocks(ctx); err != nil {
		return errors.Wrap(err, "process blocks")
	}

	return nil
}

func verify(ctx context.Context, test *testing.T, blocks []*wire.MsgBlock,
	blockRepo *handlerStorage.BlockRepository, testBlockCount int) {

	if blockRepo.LastHeight() != len(blocks) {
		test.Fatalf("Failed: Block repo height %d doesn't match added %d", blockRepo.LastHeight(),
			len(blocks))
	}

	if !blocks[len(blocks)-1].Header.BlockHash().Equal(blockRepo.LastHash()) {
		test.Fatalf("Failed: Block repo last hash doesn't match last added")
	}

	for i := 0; i < testBlockCount; i++ {
		hash := blocks[i].Header.BlockHash()
		height, _ := blockRepo.Height(hash)
		if height != i+1 {
			test.Fatalf("Failed: Block repo height %d should be %d : %s", height, i+1, hash)
		}
	}

	for i := 0; i < testBlockCount; i++ {
		hash, err := blockRepo.Hash(ctx, i+1)
		if err != nil || hash == nil {
			test.Fatalf("Failed: Block repo hash failed at height %d", i+1)
		} else if !hash.Equal(blocks[i].Header.BlockHash()) {
			test.Fatalf("Failed: Block repo hash %d should : %s", i+1, blocks[i].Header.BlockHash())
		}
	}

	// Save repo
	if err := blockRepo.Save(ctx); err != nil {
		test.Fatalf("Failed to save block repo : %v", err)
	}

	// Load repo
	if err := blockRepo.Load(ctx); err != nil {
		test.Fatalf("Failed to load block repo : %v", err)
	}

	if blockRepo.LastHeight() != len(blocks) {
		test.Fatalf("Failed: Block repo height %d doesn't match added %d after reload",
			blockRepo.LastHeight(), len(blocks))
	}

	if !blocks[len(blocks)-1].Header.BlockHash().Equal(blockRepo.LastHash()) {
		test.Fatalf("Failed: Block repo last hash doesn't match last added after reload")
	}

	for i := 0; i < testBlockCount; i++ {
		hash := blocks[i].Header.BlockHash()
		height, _ := blockRepo.Height(hash)
		if height != i+1 {
			test.Fatalf("Failed: Block repo height %d should be %d : %s", height, i+1, hash)
		}
	}

	for i := 0; i < testBlockCount; i++ {
		hash, err := blockRepo.Hash(ctx, i+1)
		if err != nil || hash == nil {
			test.Fatalf("Failed: Block repo hash failed at height %d", i+1)
		} else if !hash.Equal(blocks[i].Header.BlockHash()) {
			test.Fatalf("Failed: Block repo hash %d should : %s", i+1, blocks[i].Header.BlockHash())
		}
	}

	test.Logf("Verified %d blocks", len(blocks))
}

type TestHandler struct {
	test      *testing.T
	state     *state.State
	blocks    *handlerStorage.BlockRepository
	txs       *handlerStorage.TxRepository
	height    int
	txTracker *state.TxTracker
}

// This is called when a block is being processed.
// It is responsible for any cleanup as a result of a block.
func (handler *TestHandler) ProcessBlocks(ctx context.Context) error {

	for {
		block := handler.state.NextBlock()

		if block == nil {
			break
		}

		header := block.GetHeader()
		hash := header.BlockHash()

		if handler.blocks.Contains(hash) {
			height, _ := handler.blocks.Height(hash)
			logger.Warn(ctx, "Already have block (%d) : %s", height, hash)
			return errors.New("block not added")
		}

		if header.PrevBlock != *handler.blocks.LastHash() {
			// Ignore this as it can happen when there is a reorg.
			logger.Warn(ctx, "Not next block : %s", hash)
			logger.Warn(ctx, "Previous hash : %s", header.PrevBlock)
			return errors.New("not next block") // Unknown or out of order block
		}

		// Add to repo
		if err := handler.blocks.Add(ctx, &header); err != nil {
			return err
		}
	}

	return nil
}

// Spynode handler interface
func (handler *TestHandler) HandleHeaders(ctx context.Context, headers *client.Headers) {
	handler.test.Logf("New header (%d) : %s", headers.StartHeight, headers.Headers[0].BlockHash())
}

func (handler *TestHandler) HandleTx(ctx context.Context, tx *client.Tx) {
	handler.test.Logf("Tx : %s", tx.Tx.TxHash())
}

func (handler *TestHandler) HandleTxUpdate(ctx context.Context, update *client.TxUpdate) {
	handler.test.Logf("Tx update : %s", update.TxID)
}

func (handler *TestHandler) HandleInSync(ctx context.Context) {
	handler.test.Logf("In Sync")
}

func (handler *TestHandler) HandleMessage(ctx context.Context, payload client.MessagePayload) {}
