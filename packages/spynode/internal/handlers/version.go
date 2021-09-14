package handlers

import (
	"context"
	"errors"

	"github.com/tokenized/pkg/logger"
	"github.com/tokenized/pkg/wire"
	"github.com/tokenized/spynode/internal/state"
)

// VersionHandler exists to handle the Version command.
type VersionHandler struct {
	state   *state.State
	address string
}

// NewVersionHandler returns a new VersionHandler with the given Config.
func NewVersionHandler(state *state.State, address string) *VersionHandler {
	result := VersionHandler{state: state, address: address}
	return &result
}

// Verifies version message and sends acknowledge
func (handler *VersionHandler) Handle(ctx context.Context, m wire.Message) ([]wire.Message, error) {
	msg, ok := m.(*wire.MsgVersion)
	if !ok {
		return nil, errors.New("Could not assert as *wire.MsgVersion")
	}

	logger.Verbose(ctx, "(%s) Version : %s protocol %d, blocks %d", handler.address, msg.UserAgent,
		msg.ProtocolVersion, msg.LastBlock)
	handler.state.SetVersionReceived()

	// Return a version acknowledge
	// TODO Verify the version is compatible
	return []wire.Message{wire.NewMsgVerAck()}, nil
}
