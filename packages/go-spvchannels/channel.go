package spvchannels

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) getChanelBaseEndpoint() string {
	return fmt.Sprintf("https://%s/api/%s/account", c.cfg.baseURL, c.cfg.version)
}

func (c *Client) getTokenBaseEndpoint(accountid, channelid string) string {
	return fmt.Sprintf("%s/%s/channel/%s/api-token", c.getChanelBaseEndpoint(), accountid, channelid)
}

// ChannelsRequest hold data for get channels request
type ChannelsRequest struct {
	AccountID string `json:"accountid"`
}

// ChannelsReply hold data for get channels reply
type ChannelsReply struct {
	Channels []struct {
		ID          string `json:"id"`
		Href        string `json:"href"`
		PublicRead  bool   `json:"public_read"`
		PublicWrite bool   `json:"public_write"`
		Sequenced   bool   `json:"sequenced"`
		Locked      bool   `json:"locked"`
		Head        int    `json:"head"`
		Retention   struct {
			MinAgeDays int  `json:"min_age_days"`
			MaxAgeDays int  `json:"max_age_days"`
			AutoPrune  bool `json:"auto_prune"`
		} `json:"retention"`
		AccessTokens []struct {
			ID          string `json:"id"`
			Token       string `json:"token"`
			Description string `json:"description"`
			CanRead     bool   `json:"can_read"`
			CanWrite    bool   `json:"can_write"`
		} `json:"access_tokens"`
	} `json:"channels"`
}

// ChannelRequest hold data for get channel request
type ChannelRequest struct {
	AccountID string `json:"accountid"`
	ChannelID string `json:"channelid"`
}

// ChannelReply hold data for get channel reply
type ChannelReply struct {
	ID          string `json:"id"`
	Href        string `json:"href"`
	PublicRead  bool   `json:"public_read"`
	PublicWrite bool   `json:"public_write"`
	Sequenced   bool   `json:"sequenced"`
	Locked      bool   `json:"locked"`
	Head        int    `json:"head"`
	Retention   struct {
		MinAgeDays int  `json:"min_age_days"`
		MaxAgeDays int  `json:"max_age_days"`
		AutoPrune  bool `json:"auto_prune"`
	} `json:"retention"`
	AccessTokens []struct {
		ID          string `json:"id"`
		Token       string `json:"token"`
		Description string `json:"description"`
		CanRead     bool   `json:"can_read"`
		CanWrite    bool   `json:"can_write"`
	} `json:"access_tokens"`
}

// ChannelUpdateRequest hold data for update channel request
type ChannelUpdateRequest struct {
	AccountID   string `json:"accountid"`
	ChannelID   string `json:"channelid"`
	PublicRead  bool   `json:"public_read"`
	PublicWrite bool   `json:"public_write"`
	Locked      bool   `json:"locked"`
}

// ChannelUpdateReply hold data for update channel reply
type ChannelUpdateReply struct {
	PublicRead  bool `json:"public_read"`
	PublicWrite bool `json:"public_write"`
	Locked      bool `json:"locked"`
}

// ChannelDeleteRequest hold data for delete channel request
type ChannelDeleteRequest struct {
	AccountID string `json:"accountid"`
	ChannelID string `json:"channelid"`
}

// ChannelCreateRequest hold data for create channel request
type ChannelCreateRequest struct {
	AccountID   string `json:"accountid"`
	PublicRead  bool   `json:"public_read"`
	PublicWrite bool   `json:"public_write"`
	Sequenced   bool   `json:"sequenced"`
	Retention   struct {
		MinAgeDays int  `json:"min_age_days"`
		MaxAgeDays int  `json:"max_age_days"`
		AutoPrune  bool `json:"auto_prune"`
	} `json:"retention"`
}

// ChannelCreateReply hold data for create channel reply
type ChannelCreateReply struct {
	ID          string `json:"id"`
	Href        string `json:"href"`
	PublicRead  bool   `json:"public_read"`
	PublicWrite bool   `json:"public_write"`
	Sequenced   bool   `json:"sequenced"`
	Locked      bool   `json:"locked"`
	Head        int    `json:"head"`
	Retention   struct {
		MinAgeDays int  `json:"min_age_days"`
		MaxAgeDays int  `json:"max_age_days"`
		AutoPrune  bool `json:"auto_prune"`
	} `json:"retention"`
	AccessTokens []struct {
		ID          string `json:"id"`
		Token       string `json:"token"`
		Description string `json:"description"`
		CanRead     bool   `json:"can_read"`
		CanWrite    bool   `json:"can_write"`
	} `json:"access_tokens"`
}

// TokenRequest hold data for get token request
type TokenRequest struct {
	AccountID string `json:"accountid"`
	ChannelID string `json:"channelid"`
	TokenID   string `json:"tokenid"`
}

// TokenReply hold data for get token reply
type TokenReply struct {
	ID          string `json:"id"`
	Token       string `json:"token"`
	Description string `json:"description"`
	CanRead     bool   `json:"can_read"`
	CanWrite    bool   `json:"can_write"`
}

// TokenDeleteRequest hold data for delete token request
type TokenDeleteRequest struct {
	AccountID string `json:"accountid"`
	ChannelID string `json:"channelid"`
	TokenID   string `json:"tokenid"`
}

// TokensRequest hold data for get tokens request
type TokensRequest struct {
	AccountID string `json:"accountid"`
	ChannelID string `json:"channelid"`
}

// TokensReply hold data for get tokens reply
type TokensReply []TokenReply

// TokenCreateRequest hold data for create token request
type TokenCreateRequest struct {
	AccountID   string `json:"accountid"`
	ChannelID   string `json:"channelid"`
	Description string `json:"description"`
	CanRead     bool   `json:"can_read"`
	CanWrite    bool   `json:"can_write"`
}

// TokenCreateReply hold data for create token reply
type TokenCreateReply struct {
	ID          string `json:"id"`
	Token       string `json:"token"`
	Description string `json:"description"`
	CanRead     bool   `json:"can_read"`
	CanWrite    bool   `json:"can_write"`
}

// Channels get the list of channels
func (c *Client) Channels(ctx context.Context, r ChannelsRequest) (*ChannelsReply, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s/channel/list", c.getChanelBaseEndpoint(), r.AccountID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	res := ChannelsReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Channel get the channel
func (c *Client) Channel(ctx context.Context, r ChannelRequest) (*ChannelReply, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s/channel/%s", c.getChanelBaseEndpoint(), r.AccountID, r.ChannelID),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res := ChannelReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ChannelUpdate update the channel
func (c *Client) ChannelUpdate(ctx context.Context, r ChannelUpdateRequest) (*ChannelUpdateReply, error) {
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/%s/channel/%s", c.getChanelBaseEndpoint(), r.AccountID, r.ChannelID),
		bytes.NewBuffer(payload),
	)

	if err != nil {
		return nil, err
	}

	res := ChannelUpdateReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ChannelDelete delete the channel
func (c *Client) ChannelDelete(ctx context.Context, r ChannelDeleteRequest) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/%s/channel/%s", c.getChanelBaseEndpoint(), r.AccountID, r.ChannelID),
		nil,
	)

	if err != nil {
		return err
	}

	return c.sendRequest(req, nil)
}

// ChannelCreate create a channel
func (c *Client) ChannelCreate(ctx context.Context, r ChannelCreateRequest) (*ChannelCreateReply, error) {
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/%s/channel", c.getChanelBaseEndpoint(), r.AccountID),
		bytes.NewBuffer(payload),
	)

	if err != nil {
		return nil, err
	}

	res := ChannelCreateReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Token get the token
func (c *Client) Token(ctx context.Context, r TokenRequest) (*TokenReply, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s", c.getTokenBaseEndpoint(r.AccountID, r.ChannelID), r.TokenID),
		nil,
	)

	if err != nil {
		return nil, err
	}

	res := TokenReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// TokenDelete delate the token
func (c *Client) TokenDelete(ctx context.Context, r TokenDeleteRequest) error {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/%s", c.getTokenBaseEndpoint(r.AccountID, r.ChannelID), r.TokenID),
		nil,
	)

	if err != nil {
		return err
	}

	return c.sendRequest(req, nil)
}

// Tokens get the list of tokens
func (c *Client) Tokens(ctx context.Context, r TokensRequest) (*TokensReply, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.getTokenBaseEndpoint(r.AccountID, r.ChannelID), nil)
	if err != nil {
		return nil, err
	}

	res := TokensReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// TokenCreate create a token
func (c *Client) TokenCreate(ctx context.Context, r TokenCreateRequest) (*TokenCreateReply, error) {
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.getTokenBaseEndpoint(r.AccountID, r.ChannelID), bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, err
	}

	res := TokenCreateReply{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
