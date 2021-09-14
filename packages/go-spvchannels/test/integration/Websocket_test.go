// +build integration

package integration

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	spv "github.com/libsv/go-spvchannels"
	"github.com/stretchr/testify/assert"
)

// getChannelTokens create a new channel with 2 tokens
// Return chanelid, token1, token2 and error if any
func getChannelTokens() (string, string, string, error) {
	restClient := getRestClient()

	// Create a new channel
	replyCreateChannel, err := createChannel(restClient)
	if err != nil {
		return "", "", "", err
	}
	channelid := replyCreateChannel.ID
	token1 := replyCreateChannel.AccessTokens[0].Token

	// Create and additional token on the channel
	r := spv.TokenCreateRequest{
		AccountID:   accountid,
		ChannelID:   replyCreateChannel.ID,
		Description: "Create Second Token for websocket test",
		CanRead:     true,
		CanWrite:    true,
	}

	reply, err := restClient.TokenCreate(context.Background(), r)
	if err != nil {
		return "", "", "", err
	}
	token2 := reply.Token
	return channelid, token1, token2, nil
}

// TestWebsocket run 2 goroutines :
//     One open a websocket client, listen and process messages. The process exit when it fully received N messages
//     Other keep sending messages to spv channel server, exit only when the first one received enough messages
func TestWebsocket(t *testing.T) {
	channelid, token1, token2, err := getChannelTokens()
	assert.NoError(t, err)
	maxReceive := uint64(10)
	totalReceive := uint64(0)

	var wg sync.WaitGroup

	// Websocket client routine ---------------------------------------------------------
	wg.Add(1)
	go func() {
		ws := spv.NewWSClient(
			spv.WithBaseURL(baseURL),
			spv.WithVersion(version),
			spv.WithToken(token1),
			spv.WithChannelID(channelid),
			spv.WithWebsocketCallBack(func(t int, msg []byte, err error) error {
				atomic.AddUint64(&totalReceive, 1)
				return nil
			}),
			spv.WithMaxNotified(maxReceive),
			spv.WithInsecure(),
		)

		err := ws.Run()
		assert.NoError(t, err)
		assert.Equal(t, ws.NbNotified(), maxReceive)
		wg.Done()
	}()

	// Message writer client routine ------------------------------------------------------
	wg.Add(1)
	go func() {
		restClient := spv.NewClient(
			spv.WithBaseURL(baseURL),
			spv.WithVersion(version),
			spv.WithUser(duser),
			spv.WithPassword(dpassword),
			spv.WithToken(token2),
			spv.WithInsecure(),
		)

		for atomic.LoadUint64(&totalReceive) < maxReceive {
			r := spv.MessageWriteRequest{
				ChannelID: channelid,
				Message:   "Some random message",
			}

			_, err := restClient.MessageWrite(context.Background(), r)
			assert.NoError(t, err)
		}
		wg.Done()
	}()

	wg.Wait()
}
