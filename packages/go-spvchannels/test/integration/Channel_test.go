// +build integration

package integration

import (
	"context"

	"testing"

	spv "github.com/libsv/go-spvchannels"
	"github.com/stretchr/testify/assert"
)

func TestChannelIntegration(t *testing.T) {

	t.Run("TestChannels", func(t *testing.T) {
		client := getRestClient()
		_, err := createChannel(client)
		assert.NoError(t, err)

		reply, err := getChannels(client, accountid)
		assert.NoError(t, err)
		assert.NotEmpty(t, reply.Channels)
	})

	t.Run("TestChannel", func(t *testing.T) {
		client := getRestClient()
		replyCreateChannel, _ := createChannel(client)

		reply, err := getChannel(client, accountid, replyCreateChannel.ID)
		assert.NoError(t, err)
		assert.Equal(t, reply.ID, replyCreateChannel.ID)
	})

	t.Run("TestChannelUpdate", func(t *testing.T) {
		client := getRestClient()
		replyCreateChannel, _ := createChannel(client)

		r := spv.ChannelUpdateRequest{
			AccountID:   accountid,
			ChannelID:   replyCreateChannel.ID,
			PublicRead:  false,
			PublicWrite: false,
			Locked:      false,
		}

		reply, err := client.ChannelUpdate(context.Background(), r)
		assert.NoError(t, err)
		assert.Equal(t, reply.PublicRead, r.PublicRead)
		assert.Equal(t, reply.PublicWrite, r.PublicWrite)
		assert.Equal(t, reply.Locked, r.Locked)
	})

	t.Run("TestChannelDelete", func(t *testing.T) {
		client := getRestClient()
		replyCreateChannel, _ := createChannel(client)

		replyGetChannelsBefore, _ := getChannels(client, accountid)

		r := spv.ChannelDeleteRequest{
			AccountID: accountid,
			ChannelID: replyCreateChannel.ID,
		}
		err := client.ChannelDelete(context.Background(), r)
		assert.NoError(t, err)
		replyGetChannelsAfter, _ := getChannels(client, accountid)

		assert.Equal(t, len(replyGetChannelsBefore.Channels), len(replyGetChannelsAfter.Channels)+1)
	})

	t.Run("TestChannelCreate", func(t *testing.T) {
		client := getRestClient()
		reply, err := createChannel(client)
		assert.NotNil(t, reply)
		assert.NoError(t, err)
		assert.Equal(t, len(reply.AccessTokens), 1)
		assert.NotEmpty(t, reply.ID)
		assert.NotEmpty(t, reply.AccessTokens[0].Token)
	})
}

func TestChannelTokenIntegration(t *testing.T) {

	t.Run("TestToken", func(t *testing.T) {
		client := getRestClient()
		replyCreateChannel, _ := createChannel(client)

		r := spv.TokenRequest{
			AccountID: accountid,
			ChannelID: replyCreateChannel.ID,
			TokenID:   replyCreateChannel.AccessTokens[0].ID,
		}

		reply, err := client.Token(context.Background(), r)
		assert.NoError(t, err)
		assert.Equal(t, reply.ID, replyCreateChannel.AccessTokens[0].ID)
		assert.Equal(t, reply.Token, replyCreateChannel.AccessTokens[0].Token)
		assert.Equal(t, reply.Description, replyCreateChannel.AccessTokens[0].Description)
		assert.Equal(t, reply.CanRead, replyCreateChannel.AccessTokens[0].CanRead)
		assert.Equal(t, reply.CanWrite, replyCreateChannel.AccessTokens[0].CanWrite)
	})

	t.Run("TestTokenDelete", func(t *testing.T) {
		client := getRestClient()
		replyCreateChannel, _ := createChannel(client)

		r := spv.TokenDeleteRequest{
			AccountID: accountid,
			ChannelID: replyCreateChannel.ID,
			TokenID:   replyCreateChannel.AccessTokens[0].ID,
		}

		err := client.TokenDelete(context.Background(), r)
		assert.NoError(t, err)

		r2 := spv.TokensRequest{
			AccountID: accountid,
			ChannelID: replyCreateChannel.ID,
		}

		// Token list after deleting the only one is empty
		reply, _ := client.Tokens(context.Background(), r2)
		assert.Equal(t, len(*reply), 0)
	})

	t.Run("TestTokens", func(t *testing.T) {
		client := getRestClient()
		replyCreateChannel, _ := createChannel(client)

		r := spv.TokensRequest{
			AccountID: accountid,
			ChannelID: replyCreateChannel.ID,
		}

		reply, err := client.Tokens(context.Background(), r)
		assert.NoError(t, err)
		assert.Equal(t, len(*reply), 1)
		assert.Equal(t, (*reply)[0].ID, replyCreateChannel.AccessTokens[0].ID)
		assert.Equal(t, (*reply)[0].Token, replyCreateChannel.AccessTokens[0].Token)
		assert.Equal(t, (*reply)[0].Description, replyCreateChannel.AccessTokens[0].Description)
		assert.Equal(t, (*reply)[0].CanRead, replyCreateChannel.AccessTokens[0].CanRead)
		assert.Equal(t, (*reply)[0].CanWrite, replyCreateChannel.AccessTokens[0].CanWrite)
	})

	t.Run("TestTokenCreate", func(t *testing.T) {
		client := getRestClient()
		replyCreateChannel, _ := createChannel(client)

		r := spv.TokenCreateRequest{
			AccountID:   accountid,
			ChannelID:   replyCreateChannel.ID,
			Description: "TestTokenCreate",
			CanRead:     true,
			CanWrite:    true,
		}

		reply, err := client.TokenCreate(context.Background(), r)
		assert.NoError(t, err)
		assert.NotEmpty(t, reply.ID)
		assert.NotEmpty(t, reply.Token)
		assert.Equal(t, reply.Description, r.Description)
		assert.Equal(t, reply.CanRead, r.CanRead)
		assert.Equal(t, reply.CanWrite, r.CanWrite)
	})
}
