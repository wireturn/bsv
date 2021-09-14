package spynode

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/tokenized/pkg/wire"
	"github.com/tokenized/spynode/internal/state"
	"github.com/tokenized/spynode/internal/storage"

	"github.com/pkg/errors"
)

var (
	// ErrNotFound should be returned if the file was not found.
	ErrChannelClosed = errors.New("Channel Closed")
)

// Send a message requesting headers after the latest seen
func buildHeaderRequest(ctx context.Context, protocol uint32, blocks *storage.BlockRepository,
	requests *state.State, delta, max int) (*wire.MsgGetHeaders, error) {

	getheaders := wire.NewMsgGetHeaders()
	getheaders.ProtocolVersion = protocol

	if requests != nil {
		hash := requests.BlockRequestHash(delta)
		if hash != nil {
			getheaders.AddBlockLocatorHash(hash)
		}
	}

	// Add block hashes in reverse order
	for ; delta <= blocks.LastHeight(); delta *= 2 {
		hash, err := blocks.Hash(ctx, blocks.LastHeight()-delta)
		if err != nil {
			return getheaders, err
		}
		getheaders.AddBlockLocatorHash(hash)
		if len(getheaders.BlockLocatorHashes) > max {
			break
		}
		if blocks.LastHeight() <= delta {
			break
		}
	}

	return getheaders, nil
}

// sendAsync writes a message to a peer.
func sendAsync(ctx context.Context, conn net.Conn, m wire.Message, net wire.BitcoinNet) error {
	// build the message to send
	_, err := wire.WriteMessageN(conn, m, wire.ProtocolVersion, net)
	if err != nil {
		return errors.Wrap(err, "write message")
	}

	return nil
}

func buildVersionMsg(userAgent string, blockHeight int32) *wire.MsgVersion {
	// my local. This doesn't matter, we don't accept inboound connections.
	local := wire.NewNetAddressIPPort(net.IPv4(127, 0, 0, 1), 9333, 0)

	// build the address of the remote
	remote := wire.NewNetAddressIPPort(net.IPv4(127, 0, 0, 1), 8333, 0)

	version := wire.NewMsgVersion(remote, local, nonce(), blockHeight)
	version.UserAgent = buildUserAgent(userAgent)
	version.Services = 0x01

	return version
}

func buildUserAgent(userAgent string) string {
	return fmt.Sprintf("%v", userAgent)
}

func nonce() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}
