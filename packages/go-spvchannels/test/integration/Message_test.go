// +build integration

package integration

import (
	"context"
	"testing"

	spv "github.com/libsv/go-spvchannels"
	"github.com/stretchr/testify/assert"
)

func getClientWithTokenAndChainID() (spv.Client, string) {
	c := getRestClient()
	replyCreateChannel, _ := createChannel(c)

	newClient := spv.NewClient(
		spv.WithBaseURL(baseURL),
		spv.WithVersion(version),
		spv.WithUser(duser),
		spv.WithPassword(dpassword),
		spv.WithToken(replyCreateChannel.AccessTokens[0].Token),
		spv.WithInsecure(),
	)

	//cfg.Token = replyCreateChannel.AccessTokens[0].Token
	return *newClient, replyCreateChannel.ID
}

func writeMessage(client spv.Client, channelid string) (*spv.MessageWriteReply, error) {

	r := spv.MessageWriteRequest{
		ChannelID: channelid,
		Message:   "Hello, this is a message",
	}

	reply, err := client.MessageWrite(context.Background(), r)
	return reply, err
}

func TestMessageIntegration(t *testing.T) {

	t.Run("TestMessageHead", func(t *testing.T) {
		client, channelid := getClientWithTokenAndChainID()

		r := spv.MessageHeadRequest{
			ChannelID: channelid,
		}

		err := client.MessageHead(context.Background(), r)
		assert.NoError(t, err)
	})

	t.Run("TestMessageWrite", func(t *testing.T) {
		client, channelid := getClientWithTokenAndChainID()

		reply, err := writeMessage(client, channelid)
		assert.NoError(t, err)
		assert.True(t, reply.Sequence > 0)
		assert.NotEmpty(t, reply.Payload)
	})

	t.Run("TestMessageWrite", func(t *testing.T) {
		client, channelid := getClientWithTokenAndChainID()

		reply, err := writeMessage(client, channelid)
		assert.NoError(t, err)
		assert.True(t, reply.Sequence > 0)
		assert.NotEmpty(t, reply.Payload)
	})

	t.Run("TestMessages", func(t *testing.T) {
		client, channelid := getClientWithTokenAndChainID()

		_, err := writeMessage(client, channelid)
		assert.NoError(t, err)

		r := spv.MessagesRequest{
			ChannelID: channelid,
			UnRead:    false,
		}

		reply, err := client.Messages(context.Background(), r)
		assert.NoError(t, err)
		assert.True(t, len(*reply) > 0)
		assert.True(t, (*reply)[0].Sequence > 0)
		assert.NotEmpty(t, (*reply)[0].Payload)
	})

	t.Run("TestMessageMark", func(t *testing.T) {
		client, channelid := getClientWithTokenAndChainID()
		replyWriteMessage, _ := writeMessage(client, channelid)

		r := spv.MessageMarkRequest{
			ChannelID: channelid,
			Sequence:  replyWriteMessage.Sequence,
			Older:     false,
			Read:      false,
		}

		err := client.MessageMark(context.Background(), r)
		assert.NoError(t, err)
	})

	t.Run("TestMessageDelete", func(t *testing.T) {
		client, channelid := getClientWithTokenAndChainID()
		replyWriteMessage, _ := writeMessage(client, channelid)

		r := spv.MessageDeleteRequest{
			ChannelID: channelid,
			Sequence:  replyWriteMessage.Sequence,
		}

		err := client.MessageDelete(context.Background(), r)
		assert.NoError(t, err)
	})
}
