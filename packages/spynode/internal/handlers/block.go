package handlers

import (
	"context"

	"github.com/tokenized/pkg/logger"
	"github.com/tokenized/pkg/wire"
	"github.com/tokenized/spynode/internal/state"

	"github.com/pkg/errors"
)

// BlockHandler exists to handle the block command.
type BlockHandler struct {
	state         *state.State
	blockRefeeder *BlockRefeeder
}

// NewBlockHandler returns a new BlockHandler with the given Config.
func NewBlockHandler(state *state.State, blockRefeeder *BlockRefeeder) *BlockHandler {
	result := BlockHandler{
		state:         state,
		blockRefeeder: blockRefeeder,
	}
	return &result
}

// Handle implements the Handler interface for a block handler.
func (handler *BlockHandler) Handle(ctx context.Context, m wire.Message) ([]wire.Message, error) {
	block, ok := m.(*wire.MsgParseBlock)
	if ok {
		hash := block.Header.BlockHash()

		logger.Verbose(ctx, "Received block : %s", hash)

		if handler.blockRefeeder != nil && handler.blockRefeeder.SetBlock(*hash, block) {
			return nil, nil
		}

		if !handler.state.AddBlock(hash, block) {
			logger.Warn(ctx, "Block not requested : %s", hash)
		}
	} else {
		block, ok := m.(*wire.MsgBlock)
		if !ok {
			return nil, errors.New("Could not assert as *wire.MsgParseBlock or *wire.MsgParseBlock")
		}

		hash := block.Header.BlockHash()

		logger.Verbose(ctx, "Received block : %s", hash)

		if handler.blockRefeeder != nil && handler.blockRefeeder.SetBlock(*hash, block) {
			return nil, nil
		}

		if !handler.state.AddBlock(hash, block) {
			logger.Warn(ctx, "Block not requested : %s", hash)
		}
	}

	return nil, nil
}
