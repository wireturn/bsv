package spynode

import (
	"bytes"
	"context"
	"fmt"

	"github.com/tokenized/pkg/bitcoin"
	"github.com/tokenized/pkg/logger"
	"github.com/tokenized/pkg/storage"
	"github.com/tokenized/pkg/wire"
	"github.com/tokenized/specification/dist/golang/actions"
	"github.com/tokenized/specification/dist/golang/protocol"
	"github.com/tokenized/spynode/internal/handlers"
	handlerstorage "github.com/tokenized/spynode/internal/storage"
	"github.com/tokenized/spynode/pkg/client"

	"github.com/pkg/errors"
)

// HandleTx processes a tx through spynode as if it came from the network.
// Used to feed "response" txs directly back through spynode.
func (node *Node) HandleTx(ctx context.Context, tx *wire.MsgTx) error {
	return node.unconfTxChannel.Add(handlers.TxData{Msg: tx, Trusted: true, Safe: true,
		ConfirmedHeight: -1})
}

func (node *Node) processUnconfirmedTx(ctx context.Context, tx handlers.TxData) error {
	hash := tx.Msg.TxHash()

	if tx.ConfirmedHeight != -1 {
		return errors.New("Process unconfirmed tx with height")
	}

	node.txTracker.Remove(ctx, *hash)

	// The mempool is needed to track which transactions have been sent to listeners and to check
	//   for attempted double spends.
	conflicts, trusted, added := node.memPool.AddTransaction(ctx, tx.Msg, tx.Trusted)
	if !added {
		return nil // Already saw this tx
	}

	// logger.Debug(ctx, "Tx mempool (added %t) (flagged trusted %t) (received trusted %t) : %s",
	// 	added, trusted, tx.Trusted, hash.String())

	if trusted {
		// Was marked trusted in the mempool by a tx inventory from the trusted node.
		// tx.Trusted means the tx itself was received from the trusted node.
		tx.Trusted = trusted
	}

	if len(conflicts) > 0 {
		logger.Warn(ctx, "Found %d conflicts with %s", len(conflicts), hash)
		// Notify of attempted double spend
		for _, conflict := range conflicts {
			isRelevant, err := node.txs.MarkUnsafe(ctx, conflict)
			if err != nil {
				return errors.Wrap(err, "Failed to check tx repo")
			}
			if !isRelevant {
				continue // Only send for txs that previously matched filters.
			}

			txState, err := handlerstorage.FetchTxState(ctx, node.store, conflict)
			if err != nil {
				continue
			}

			txState.State.UnSafe = true
			txState.State.Safe = false

			if err := handlerstorage.SaveTxState(ctx, node.store, txState); err != nil {
				return errors.Wrap(err, "save tx state")
			}

			update := &client.TxUpdate{
				TxID:  *txState.Tx.TxHash(),
				State: txState.State,
			}

			// Notify of tx conflict
			for _, handler := range node.handlers {
				handler.HandleTxUpdate(ctx, update)
			}
		}
	}

	if !node.IsRelevant(ctx, tx.Msg) {
		if _, err := node.txs.Remove(ctx, *hash, -1); err != nil {
			return errors.Wrap(err, "Failed to remove from tx repo")
		}
		return nil // Filter out
	}

	logger.Info(ctx, "Tx is relevant : %s", hash)

	// We have to succesfully add to tx repo because it is protected by a lock and will prevent
	// processing the same tx twice at the same time.
	added, newlySafe, err := node.txs.Add(ctx, *hash, tx.Trusted, tx.Safe, -1)
	if err != nil {
		return errors.Wrap(err, "add to tx repo")
	}
	if !added {
		logger.Info(ctx, "Tx already added : %s", hash)
		return nil // tx already processed
	}

	// logger.Debug(ctx, "Tx repo (added %t) (newly safe %t) : %s", added, newlySafe, hash.String())

	txState, err := handlerstorage.FetchTxState(ctx, node.store, *hash)
	if err != nil {
		if errors.Cause(err) != storage.ErrNotFound {
			return errors.Wrap(err, "fetch tx state")
		}

		logger.Info(ctx, "Creating new tx state : %s", hash)

		// Create new tx state
		txState = &client.Tx{
			Tx: tx.Msg,
		}

		if err := fetchSpentOutputs(ctx, node.store, node.outputFetcher, txState); err != nil {
			return errors.Wrap(err, "fetch outputs")
		}
	} else {
		logger.Info(ctx, "Updating tx state : %s", hash)
	}

	txState.State.Safe = tx.Safe || newlySafe
	if txState.State.MerkleProof == nil {
		txState.State.UnconfirmedDepth = 1
	}

	if len(conflicts) > 0 {
		// Unsafe
		txState.State.UnSafe = true
		txState.State.Safe = false
	}

	if err := handlerstorage.SaveTxState(ctx, node.store, txState); err != nil {
		return errors.Wrap(err, "save tx state")
	}

	logger.Info(ctx, "Saved tx state : %s", hash)

	// Notify of new tx
	for _, handler := range node.handlers {
		handler.HandleTx(ctx, txState)
	}

	return nil
}

// fetchSpentOutputs fetches the outputs spent by this tx.
func fetchSpentOutputs(ctx context.Context, store storage.Storage, outputFetcher OutputFetcher,
	tx *client.Tx) error {

	tx.Outputs = make([]*wire.TxOut, len(tx.Tx.TxIn))

	var toFetch []wire.OutPoint
	for i, input := range tx.Tx.TxIn {
		if input.PreviousOutPoint.Index == wire.MaxPrevOutIndex {
			// coinbase tx. there is no spent output. these are new coins
			tx.Outputs[i] = wire.NewTxOut(0, nil)
			continue
		}

		inputTx, err := handlerstorage.FetchTxState(ctx, store, input.PreviousOutPoint.Hash)
		if errors.Cause(err) == storage.ErrNotFound {
			toFetch = append(toFetch, input.PreviousOutPoint)
			continue
		} else if err != nil {
			return errors.Wrap(err, "fetch tx")
		}

		if int(input.PreviousOutPoint.Index) >= len(inputTx.Tx.TxOut) {
			// Invalid previous outpoint
			tx.Outputs[i] = wire.NewTxOut(0, nil)
			continue
		}

		tx.Outputs[i] = inputTx.Tx.TxOut[input.PreviousOutPoint.Index]
	}

	if len(toFetch) == 0 {
		return nil
	}

	utxos, err := outputFetcher.GetOutputs(ctx, toFetch)
	if err != nil {
		return errors.Wrap(err, "fetch outputs")
	}

	utxoIndex := 0
	for i := range tx.Outputs {
		if tx.Outputs[i] != nil {
			continue
		}

		if utxoIndex >= len(utxos) {
			return fmt.Errorf("Not enough outputs returned : %d/%d", len(utxos), len(toFetch))
		}

		tx.Outputs[i] = wire.NewTxOut(utxos[utxoIndex].Value, utxos[utxoIndex].LockingScript)
		utxoIndex++
	}

	return nil
}

// processUnconfirmedTxs pulls txs from the unconfirmed tx channel and processes them.
func (node *Node) processUnconfirmedTxs(ctx context.Context) {
	for tx := range node.unconfTxChannel.Channel {
		if err := node.processUnconfirmedTx(ctx, tx); err != nil {
			logger.Error(ctx, "SpyNodeAborted to process unconfirmed tx : %s : %s", err,
				tx.Msg.TxHash().String())
			node.requestStop(ctx)
			break
		}
	}
}

func (node *Node) IsSubscribedToContracts(ctx context.Context) bool {
	node.lock.Lock()
	defer node.lock.Unlock()

	return node.sendContracts
}

func (node *Node) IsRelevant(ctx context.Context, tx *wire.MsgTx) bool {
	if node.IsSubscribedToContracts(ctx) && checkContracts(ctx, tx, node.config.IsTest) {
		logger.Info(ctx, "Contract found : %s", tx.TxHash())
		return true
	}

	node.pushDataLock.Lock()
	defer node.pushDataLock.Unlock()

	// Check the hashes from each output against this client.
	for _, output := range tx.TxOut {
		r := bytes.NewReader(output.PkScript)
		for {
			_, pushdata, err := bitcoin.ParsePushDataScript(r)
			if err != nil {
				if err == bitcoin.ErrNotPushOp { // ignore non push op codes
					continue
				}
				break
			}

			hash := pushDataToHash(pushdata)
			for _, pd := range node.pushDataHashes {
				if pd.Equal(&hash) {
					return true
				}
			}
		}
	}

	// Note: To support P2PK addresses we would need to track UTXOs since the signature scripts
	// would only be the signature.
	for _, input := range tx.TxIn {
		r := bytes.NewReader(input.SignatureScript)
		for {
			_, pushdata, err := bitcoin.ParsePushDataScript(r)
			if err != nil {
				if err == bitcoin.ErrNotPushOp { // ignore non push op codes
					continue
				}
				break
			}

			hash := pushDataToHash(pushdata)
			for _, pd := range node.pushDataHashes {
				if pd.Equal(&hash) {
					return true
				}
			}
		}
	}

	return false
}

func pushDataToHash(b []byte) bitcoin.Hash20 {
	var hash bitcoin.Hash20
	if len(b) == bitcoin.Hash20Size {
		copy(hash[:], b)
	} else {
		copy(hash[:], bitcoin.Hash160(b))
	}

	return hash
}

// checkContracts returns true if the tx contains a Tokenized "contract wide" op return.
// This includes contract formations and asset creations so can be used to index contract and asset
// information.
func checkContracts(ctx context.Context, tx *wire.MsgTx, isTest bool) bool {
	for _, output := range tx.TxOut {
		action, err := protocol.Deserialize(output.PkScript, isTest)
		if err != nil {
			continue
		}

		switch action.(type) {
		case *actions.ContractFormation, *actions.AssetCreation:
			return true
		default:
			continue
		}
	}

	return false
}
