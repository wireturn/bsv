package handlers

import (
	"context"
	"fmt"

	"github.com/tokenized/pkg/logger"
	"github.com/tokenized/pkg/wire"
	"github.com/tokenized/spynode/internal/platform/config"
	"github.com/tokenized/spynode/internal/state"
	"github.com/tokenized/spynode/internal/storage"
	"github.com/tokenized/spynode/pkg/client"

	"github.com/pkg/errors"
)

// HeadersHandler exists to handle the headers command.
type HeadersHandler struct {
	config   config.Config
	state    *state.State
	blocks   *storage.BlockRepository
	txs      *storage.TxRepository
	reorgs   *storage.ReorgRepository
	handlers []client.Handler
}

// NewHeadersHandler returns a new HeadersHandler with the given Config.
func NewHeadersHandler(config config.Config, state *state.State,
	blockRepo *storage.BlockRepository, txRepo *storage.TxRepository,
	reorgs *storage.ReorgRepository, handlers []client.Handler) *HeadersHandler {

	result := HeadersHandler{
		config:   config,
		state:    state,
		blocks:   blockRepo,
		txs:      txRepo,
		reorgs:   reorgs,
		handlers: handlers,
	}
	return &result
}

// Implements the Handler interface.
// Headers are in order from lowest block height, to highest
func (handler *HeadersHandler) Handle(ctx context.Context, m wire.Message) ([]wire.Message, error) {
	message, ok := m.(*wire.MsgHeaders)
	if !ok {
		return nil, errors.New("Could not assert as *wire.Msginv")
	}

	response := []wire.Message{}
	modified := false
	addedCount := 0
	// logger.Debug(ctx, "Received %d headers", len(message.Headers))

	lastHash := handler.state.LastHash()
	newHeight := handler.blocks.LastHeight() + handler.state.TotalBlockRequestCount()

	if !handler.state.IsReady() && (len(message.Headers) == 0 || (len(message.Headers) == 1 &&
		lastHash.Equal(message.Headers[0].BlockHash()))) {

		logger.Info(ctx, "Headers in sync at height %d", handler.blocks.LastHeight())
		handler.state.SetPendingSync() // We are in sync
		if handler.state.StartHeight() == -1 {
			handler.state.SetInSync()
			logger.Error(ctx, "Headers in sync before start block found")
		} else if handler.state.BlockRequestsEmpty() {
			handler.state.SetInSync()
			logger.Info(ctx, "Blocks in sync at height %d", handler.blocks.LastHeight())
		}
		handler.state.ClearHeadersRequested()
		handler.blocks.Save(ctx) // Save when we get in sync
		return response, nil
	}

	// Process headers
	getBlocks := wire.NewMsgGetData()
	for _, header := range message.Headers {
		if len(header.PrevBlock) == 0 {
			continue
		}

		hash := header.BlockHash()

		if lastHash.Equal(&header.PrevBlock) {
			if len(message.Headers) < 10 {
				logger.InfoWithFields(ctx, []logger.Field{
					logger.Stringer("hash", hash),
				}, "Adding header")
			}
			request, err := handler.checkStartHeight(ctx, header)
			if err != nil {
				return response, errors.Wrap(err, "check start height")
			}
			if request { // block should be processed
				// Request it if it isn't already requested.
				sendRequest, err := handler.state.AddBlockRequest(&header.PrevBlock, hash)
				if err != nil {
					if err == state.ErrWrongPreviousHash {
						logger.Warn(ctx, "Wrong previous hash : %s", header.PrevBlock)
					}
				} else if sendRequest {
					// logger.Debug(ctx, "Requesting block : %s", hash)
					getBlocks.AddInvVect(wire.NewInvVect(wire.InvTypeBlock, hash))
					if len(getBlocks.InvList) == wire.MaxInvPerMsg {
						// Start new get data (blocks) message
						response = append(response, getBlocks)
						getBlocks = wire.NewMsgGetData()
					}
				}
			}

			addedCount++
			newHeight++
			lastHash = *hash
			modified = true
			continue
		}

		if hash.Equal(&lastHash) {
			continue // Already latest header
		}

		logger.InfoWithFields(ctx, []logger.Field{
			logger.Stringer("hash", hash),
			logger.Stringer("previous_hash", header.PrevBlock),
			logger.Stringer("last_hash", lastHash),
		}, "Header not next")

		// Check if we already have this block
		if handler.blocks.Contains(hash) || handler.state.BlockIsRequested(hash) ||
			handler.state.BlockIsToBeRequested(hash) {
			continue
		}

		// Check if this block is a reorg in pending blocks
		if handler.state.BlockIsRequested(&header.PrevBlock) ||
			handler.state.BlockIsToBeRequested(&header.PrevBlock) {
			logger.Info(ctx, "Reorg in pending blocks")
			handler.state.ClearBlockRequestsAfter(ctx, header.PrevBlock)

			// Request it if it isn't already requested.
			sendRequest, err := handler.state.AddBlockRequest(&header.PrevBlock, hash)
			if err != nil {
				if errors.Cause(err) == state.ErrWrongPreviousHash {
					logger.Warn(ctx, "Wrong previous hash : %s", header.PrevBlock)
				}
			} else if sendRequest {
				// logger.Debug(ctx, "Requesting block : %s", hash)
				getBlocks.AddInvVect(wire.NewInvVect(wire.InvTypeBlock, hash))
				if len(getBlocks.InvList) == wire.MaxInvPerMsg {
					// Start new get data (blocks) message
					response = append(response, getBlocks)
					getBlocks = wire.NewMsgGetData()
				}
			}

			lastHash = *hash
			modified = true
			addedCount = 1
			newHeight = handler.blocks.LastHeight() + handler.state.TotalBlockRequestCount()
			continue
		}

		// Check for a reorg in processed blocks
		reorgHeight, exists := handler.blocks.Height(&header.PrevBlock)
		if exists {
			if reorgHeight == handler.blocks.LastHeight() {
				logger.Info(ctx, "Reorg on latest block")
				handler.state.ClearInSync()
				handler.state.ClearBlockRequests(ctx)
				continue
			}

			logger.Info(ctx, "Reorging to height %d", reorgHeight)
			handler.state.ClearInSync()
			handler.state.ClearBlockRequests(ctx)

			// Call reorg listener for all blocks above reorg height.
			reorg := storage.Reorg{
				BlockHeight: reorgHeight,
			}

			for height := handler.blocks.LastHeight(); height > reorgHeight; height-- {
				// Add block to reorg
				revertHeader, err := handler.blocks.Header(ctx, height)
				if err != nil {
					return response, errors.Wrap(err, "Failed to get reverted block header")
				}

				reorgBlock := storage.ReorgBlock{
					Header: *revertHeader,
				}

				revertTxs, err := handler.txs.GetBlock(ctx, height)
				if err != nil {
					return response, errors.Wrap(err, "Failed to get reverted txs")
				}
				for _, txid := range revertTxs {
					reorgBlock.TxIds = append(reorgBlock.TxIds, txid)
				}

				reorg.Blocks = append(reorg.Blocks, reorgBlock)

				if len(revertTxs) > 0 {
					if err := handler.txs.RemoveBlock(ctx, height); err != nil {
						return response, errors.Wrap(err, "Failed to remove reverted txs")
					}
				} else {
					if err := handler.txs.ReleaseBlock(ctx, height); err != nil {
						return response, errors.Wrap(err, "Failed to remove reverted txs")
					}
				}
			}

			if len(reorg.Blocks) > 0 {
				logger.Info(ctx, "Removed %d blocks", len(reorg.Blocks))
				if err := handler.reorgs.Save(ctx, &reorg); err != nil {
					return response, errors.Wrap(err, "save reorg")
				}

				// Revert block repository
				if err := handler.blocks.Revert(ctx, reorgHeight); err != nil {
					return response, errors.Wrap(err, "revert blocks")
				}
			} else {
				logger.Info(ctx, "No blocks removed")
			}

			// Assert this header is now next
			newLastHash := handler.blocks.LastHash()
			if newLastHash == nil || !newLastHash.Equal(&header.PrevBlock) {
				return response, fmt.Errorf("Revert failed to produce correct last hash : %s",
					newLastHash)
			}
			handler.state.SetLastHash(*newLastHash)

			// Add this header after the new top block
			request, err := handler.checkStartHeight(ctx, header)
			if err != nil {
				return response, errors.Wrap(err, "check start height")
			}
			if request {
				// Request it if it isn't already requested.
				sendRequest, err := handler.state.AddBlockRequest(&header.PrevBlock, hash)
				if err != nil {
					if errors.Cause(err) == state.ErrWrongPreviousHash {
						logger.Warn(ctx, "Wrong previous hash : %s", header.PrevBlock)
					}
				} else if sendRequest {
					// logger.Debug(ctx, "Requesting block : %s", hash)
					getBlocks.AddInvVect(wire.NewInvVect(wire.InvTypeBlock, hash))
					if len(getBlocks.InvList) == wire.MaxInvPerMsg {
						// Start new get data (blocks) message
						response = append(response, getBlocks)
						getBlocks = wire.NewMsgGetData()
					}
				}
			}

			lastHash = *hash
			modified = true
			addedCount = 1
			newHeight = handler.blocks.LastHeight() + handler.state.TotalBlockRequestCount()
			continue
		}

		// Ignore unknown blocks as they might happen when there is a reorg.
		logger.Verbose(ctx, "Unknown header : %s", hash)
		logger.Verbose(ctx, "Previous hash : %s", header.PrevBlock)
		return nil, nil //errors.New(fmt.Sprintf("Unknown header : %s", hash))
	}

	// Add any non-full requests.
	if len(getBlocks.InvList) > 0 {
		response = append(response, getBlocks)
	}

	if modified {
		handler.state.ClearHeadersRequested()
	}
	if addedCount > 0 {
		logger.Info(ctx, "Added %d headers to height %d", addedCount, newHeight)
	}
	return response, nil
}

func (handler HeadersHandler) checkStartHeight(ctx context.Context,
	header *wire.BlockHeader) (bool, error) {

	startHeight := handler.state.StartHeight()
	if startHeight == -1 {
		hash := header.BlockHash()
		// Check if it is the start block
		if handler.config.StartHash.Equal(hash) {
			startHeight = handler.blocks.LastHeight() + 1
			handler.state.SetStartHeight(startHeight)
			handler.state.SetLastHash(header.PrevBlock)
			logger.Verbose(ctx, "Found start block at height %d", startHeight)
		} else {
			err := handler.blocks.Add(ctx, header) // Just add hashes before the start block
			if err != nil {
				return false, errors.Wrap(err, "add header")
			}
			handler.state.SetLastHash(*hash)
			return false, nil
		}
	}

	return true, nil
}
