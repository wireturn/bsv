package spynode

import (
	"context"
	"fmt"
	"time"

	"github.com/tokenized/pkg/bitcoin"
	"github.com/tokenized/pkg/logger"
	"github.com/tokenized/pkg/storage"
	"github.com/tokenized/pkg/wire"
	handlersstorage "github.com/tokenized/spynode/internal/storage"
	"github.com/tokenized/spynode/pkg/client"

	"github.com/pkg/errors"
)

var (
	ErrBlockNotNextBlock = errors.New("Not next block")
	ErrBlockNotAdded     = errors.New("Block not added")
)

func (node *Node) processBlocks(ctx context.Context) error {

	for !node.isStopping() {

		block, height, refeederActive := node.blockRefeeder.GetBlock()
		if refeederActive {
			if block != nil {
				node.provideBlock(ctx, block, height)

				if height != 0 { // block header still active
					if node.blocks.LastHeight() == height {
						logger.Info(ctx, "Refeed complete at block %d", height)
						node.blockRefeeder.Clear(height)
					} else {
						logger.Info(ctx, "Refeed setting next block %d", height+1)
						nextHash, err := node.blocks.Hash(ctx, height+1)
						if err != nil {
							return errors.Wrap(err, "get next hash")
						}
						node.blockRefeeder.Increment(height+1, *nextHash)
					}
				}
			}

			hash := node.blockRefeeder.GetBlockToRequest()
			if hash != nil {
				getBlocks := wire.NewMsgGetData()
				getBlocks.AddInvVect(wire.NewInvVect(wire.InvTypeBlock, hash))
				if !node.queueOutgoing(getBlocks) {
					return nil
				}
			}

			time.Sleep(200 * time.Millisecond)
			continue
		}

		// Blocks are fed into the state when received by the block handler, then pulled out and
		//   processed here.
		block = node.state.NextBlock()
		if block == nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		if err := node.ProcessBlock(ctx, block); err != nil {
			c := errors.Cause(err)
			if c != ErrBlockNotNextBlock && c != ErrBlockNotAdded {
				header := block.GetHeader()
				logger.Warn(ctx, "Failed to process block : %s : %s",
					header.BlockHash().String(), err)
				return err
			}
		}

		// Request more blocks if necessary
		// TODO Send some requests to other nodes --ce
		getBlocks := wire.NewMsgGetData() // Block request message

		for {
			requestHash, _ := node.state.GetNextBlockToRequest()
			if requestHash == nil {
				break
			}

			// logger.Debug(ctx, "Requesting block : %s", requestHash.String())
			getBlocks.AddInvVect(wire.NewInvVect(wire.InvTypeBlock, requestHash))
			if len(getBlocks.InvList) == wire.MaxInvPerMsg {
				// Start new get data (block request) message
				if !node.queueOutgoing(getBlocks) {
					return nil
				}
				getBlocks = wire.NewMsgGetData()
			}
		}

		// Add any non-full requests.
		if len(getBlocks.InvList) > 0 {
			if !node.queueOutgoing(getBlocks) {
				return nil
			}
		}
	}

	return nil
}

// provideBlock feeds the block to the handlers.
func (node *Node) provideBlock(ctx context.Context, block wire.Block, height int) error {

	h := block.GetHeader()
	headers := &client.Headers{
		StartHeight: uint32(height),
		Headers:     []*wire.BlockHeader{&h},
	}

	for _, handler := range node.handlers {
		handler.HandleHeaders(ctx, headers)
	}

	// logger.Debug(ctx, "Providing block %d (%d tx) : %s", height, block.GetTxCount(), hash)
	merkleTree := wire.NewMerkleTree(true)
	var txs []*wire.MsgTx

	for {
		tx, err := block.GetNextTx()
		if err != nil {
			return errors.Wrap(err, "get next tx")
		}
		if tx == nil {
			break // parsed all txs
		}

		txid := tx.TxHash()

		if !node.IsRelevant(ctx, tx) {
			merkleTree.AddHash(*txid)
			continue
		}

		// Add to txs for block
		if _, _, err := node.txs.Add(ctx, *txid, true, true, height); err != nil {
			return errors.Wrap(err, "add to tx repo")
		}

		merkleTree.AddMerkleProof(*txid)
		txs = append(txs, tx)
	}

	// Check merkle root hash
	merkleRootHash, merkleProofs := merkleTree.FinalizeMerkleProofs()
	if !merkleRootHash.Equal(&h.MerkleRoot) {
		// TODO Remove confirmations for all txs in block? --ce
		logger.Warn(ctx, "Invalid merkle root hash : calculated %s, header %s", merkleRootHash,
			h.MerkleRoot)
		return errors.New("Invalid merkle root hash")
	}

	// Send updates for relevant txs.
	for i, tx := range txs {
		txState, err := handlersstorage.FetchTxState(ctx, node.store, *tx.TxHash())
		if err != nil {
			if errors.Cause(err) != storage.ErrNotFound {
				return errors.Wrap(err, "fetch tx state")
			}

			txState := &client.Tx{
				Tx: tx,
				State: client.TxState{
					Safe:             true,
					UnSafe:           false,
					Cancelled:        false,
					UnconfirmedDepth: 0,
					MerkleProof:      convertMerkleProof(merkleProofs[i], h),
				},
			}

			if err := fetchSpentOutputs(ctx, node.store, node.outputFetcher, txState); err != nil {
				return errors.Wrap(err, "fetch outputs")
			}

			if err := handlersstorage.SaveTxState(ctx, node.store, txState); err != nil {
				return errors.Wrap(err, "save tx state")
			}
		}

		// Send tx to handlers
		for _, handler := range node.handlers {
			handler.HandleTx(ctx, txState)
		}
	}

	return nil
}

func (node *Node) ProcessBlock(ctx context.Context, block wire.Block) error {
	node.blockLock.Lock()
	defer node.blockLock.Unlock()

	header := block.GetHeader()
	hash := header.BlockHash()
	start := time.Now()

	if node.blocks.Contains(hash) {
		height, _ := node.blocks.Height(hash)
		logger.Warn(ctx, "Already have block (%d) : %s", height, hash)
		return ErrBlockNotAdded
	}

	if header.PrevBlock != *node.blocks.LastHash() {
		// Ignore this as it can happen when there is a reorg.
		logger.Warn(ctx, "Not next block : %s", hash)
		logger.Warn(ctx, "Previous hash : %s", header.PrevBlock)
		return ErrBlockNotNextBlock // Unknown or out of order block
	}

	// Validate
	if !block.IsMerkleRootValid() {
		logger.Warn(ctx, "Invalid merkle hash for block %s", hash)
		return ErrBlockNotAdded
	}

	// Add to repo
	if err := node.blocks.Add(ctx, &header); err != nil {
		return errors.Wrap(err, "add block")
	}

	// If we are in sync we can save after every block
	if node.state.IsReady() {
		if err := node.blocks.Save(ctx); err != nil {
			return errors.Wrap(err, "save blocks")
		}
	}

	// Get unconfirmed "relevant" txs
	var unconfirmed []bitcoin.Hash32
	var err error
	// This locks the tx repo so that propagated txs don't interfere while a block is being
	//   processed.
	unconfirmed, err = node.txs.GetUnconfirmed(ctx)
	if err != nil {
		return errors.Wrap(err, "get unconfirmed txs")
	}

	// Send block notification
	height := node.blocks.LastHeight()
	h := block.GetHeader()
	headers := &client.Headers{
		StartHeight: uint32(height),
		Headers:     []*wire.BlockHeader{&h},
	}

	for _, handler := range node.handlers {
		handler.HandleHeaders(ctx, headers)
	}

	// Notify Tx for block and tx listeners
	logger.Verbose(ctx, "Processing block %d (%d tx) : %s", height, block.GetTxCount(), hash)
	defer logger.Elapsed(ctx, start, fmt.Sprintf("Processed block : %s", hash))
	inUnconfirmed := false
	txids := make([]*bitcoin.Hash32, 0, block.GetTxCount())
	merkleTree := wire.NewMerkleTree(true)
	var txs []*wire.MsgTx
	var txsIsNew []bool
	var txsIsSafe []bool
	for {
		tx, err := block.GetNextTx()
		if err != nil {
			node.txs.ReleaseUnconfirmed(ctx)
			return errors.Wrap(err, "get next tx")
		}
		if tx == nil {
			break // parsed all txs
		}

		txid := tx.TxHash()
		txids = append(txids, txid)

		// Remove from unconfirmed. Only matching are in unconfirmed.
		inUnconfirmed, unconfirmed = removeHash(*txid, unconfirmed)

		// Remove from mempool
		inMemPool := false
		if node.state.IsReady() {
			inMemPool = node.memPool.RemoveTransaction(*txid)
		}

		if inUnconfirmed {
			// Already seen and marked relevant
			merkleTree.AddMerkleProof(*txid)
			txs = append(txs, tx)
			txsIsNew = append(txsIsNew, false)
			txsIsSafe = append(txsIsSafe, true)

		} else if !inMemPool {
			// Not seen yet
			isSafe := true

			// Transaction wasn't in the mempool.
			// Check for transactions in the mempool with conflicting inputs (double spends).
			if conflicting := node.memPool.Conflicting(tx); len(conflicting) > 0 {
				isSafe = false
				for _, confHash := range conflicting {
					if containsHash(confHash, unconfirmed) {
						// Only send for txs that previously matched filters.

						// Mark cancelled
						txState, err := handlersstorage.FetchTxState(ctx, node.store, *txid)
						if err != nil {
							node.txs.ReleaseUnconfirmed(ctx)
							return errors.Wrap(err, "fetch tx state")
						}

						txState.State.UnSafe = true
						txState.State.Cancelled = true

						if err := handlersstorage.SaveTxState(ctx, node.store, txState); err != nil {
							node.txs.ReleaseUnconfirmed(ctx)
							return errors.Wrap(err, "save tx state")
						}

						// Send update
						update := &client.TxUpdate{
							TxID:  *txid,
							State: txState.State,
						}
						for _, handler := range node.handlers {
							handler.HandleTxUpdate(ctx, update)
						}
					}
				}
			}

			if node.IsRelevant(ctx, tx) {
				// Add to txs for block
				if _, _, err := node.txs.Add(ctx, *txid, true, true, height); err != nil {
					node.txs.ReleaseUnconfirmed(ctx)
					return errors.Wrap(err, "add to tx repo")
				}

				merkleTree.AddMerkleProof(*txid)
				txs = append(txs, tx)
				txsIsNew = append(txsIsNew, true)
				txsIsSafe = append(txsIsSafe, isSafe)

			} else {
				if _, err := node.txs.Remove(ctx, *txid, height); err != nil {
					node.txs.ReleaseUnconfirmed(ctx)
					return errors.Wrap(err, "remove from tx repo")
				}
			}
		}

		merkleTree.AddHash(*txid)
	}

	// Check merkle root hash
	merkleRootHash, merkleProofs := merkleTree.FinalizeMerkleProofs()
	if !merkleRootHash.Equal(&header.MerkleRoot) {
		// TODO Remove confirmations for all txs in block? --ce
		logger.Warn(ctx, "Invalid merkle root hash : calculated %s, header %s", merkleRootHash,
			header.MerkleRoot)
		node.txs.ReleaseUnconfirmed(ctx)
		return errors.New("Invalid merkle root hash")
	}

	// Send updates for relevant txs.
	for i, tx := range txs {
		if txsIsNew[i] {
			txState := &client.Tx{
				Tx: tx,
				State: client.TxState{
					Safe:             txsIsSafe[i],
					UnSafe:           !txsIsSafe[i],
					Cancelled:        false,
					UnconfirmedDepth: 0,
					MerkleProof:      convertMerkleProof(merkleProofs[i], header),
				},
			}

			if err := fetchSpentOutputs(ctx, node.store, node.outputFetcher, txState); err != nil {
				node.txs.ReleaseUnconfirmed(ctx)
				return errors.Wrap(err, "fetch outputs")
			}

			if err := handlersstorage.SaveTxState(ctx, node.store, txState); err != nil {
				node.txs.ReleaseUnconfirmed(ctx)
				return errors.Wrap(err, "save tx state")
			}

			// Send new tx
			for _, handler := range node.handlers {
				handler.HandleTx(ctx, txState)
			}

		} else {
			txState, err := handlersstorage.FetchTxState(ctx, node.store, *tx.TxHash())
			if err != nil {
				node.txs.ReleaseUnconfirmed(ctx)
				return errors.Wrap(err, "fetch tx state")
			}

			txState.State.MerkleProof = convertMerkleProof(merkleProofs[i], header)
			txState.State.UnconfirmedDepth = 0
			if !txState.State.UnSafe && txsIsSafe[i] {
				txState.State.Safe = true
				txState.State.UnSafe = false
			} else {
				txState.State.Safe = false
				txState.State.UnSafe = true
			}

			if err := handlersstorage.SaveTxState(ctx, node.store, txState); err != nil {
				node.txs.ReleaseUnconfirmed(ctx)
				return errors.Wrap(err, "save tx state")
			}

			// Send update
			update := &client.TxUpdate{
				TxID:  *tx.TxHash(),
				State: txState.State,
			}
			for _, handler := range node.handlers {
				handler.HandleTxUpdate(ctx, update)
			}

		}
	}

	// Perform any block cleanup
	if err := node.CleanupBlock(ctx, txids); err != nil {
		logger.Warn(ctx, "Failed clean up after block : %s", hash)
		node.txs.ReleaseUnconfirmed(ctx) // Release unconfirmed
		return err
	}

	if !node.state.IsReady() {
		if node.state.IsPendingSync() && node.state.BlockRequestsEmpty() {
			node.state.SetInSync()
			logger.Info(ctx, "Blocks in sync at height %d", node.blocks.LastHeight())
		}
	}

	if err := node.txs.FinalizeUnconfirmed(ctx, unconfirmed); err != nil {
		return err
	}

	return nil
}

func convertMerkleProof(mp *wire.MerkleProof, header wire.BlockHeader) *client.MerkleProof {
	result := &client.MerkleProof{
		Index:       uint64(mp.Index),
		Path:        mp.Path,
		BlockHeader: header,
	}

	result.DuplicatedIndexes = make([]uint64,
		len(mp.DuplicatedIndexes))
	for i, value := range mp.DuplicatedIndexes {
		result.DuplicatedIndexes[i] = uint64(value)
	}

	return result
}

func containsHash(hash bitcoin.Hash32, list []bitcoin.Hash32) bool {
	for _, listhash := range list {
		if hash.Equal(&listhash) {
			return true
		}
	}
	return false
}

func removeHash(hash bitcoin.Hash32, list []bitcoin.Hash32) (bool, []bitcoin.Hash32) {
	for i, listhash := range list {
		if hash.Equal(&listhash) {
			return true, append(list[:i], list[i+1:]...)
		}
	}
	return false, list
}
