package spvchannels

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var c = NewClient(
	WithBaseURL("somedomain"),
	WithInsecure(),
)

func TestUnitChannels(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock Channels": {
			request: `{
				"accountid": "1"
			}`,
			reply: `{
				"channels": [
					{
						"id": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
						"href": "https://localhost:5010/api/v1/channel/H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
						"public_read": true,
						"public_write": true,
						"sequenced": true,
						"locked": false,
						"head": 0,
						"retention": {
							"min_age_days": 0,
							"max_age_days": 99999,
							"auto_prune": true
						},
						"access_tokens": [
							{
								"id": "1",
								"token": "20_j2-GfF6GFk8lnofe7EW5u7DhztfLQmRsa8d8R3CBZCGVU7xS1vhQwqfT-K-P2PLyxkS1wznAbj1VF1U3TFA",
								"description": "Owner",
								"can_read": true,
								"can_write": true
							}
						]
					}
				]
			}`,
			err:  nil,
			code: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req ChannelsRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.Channels(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp ChannelsReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestUnitChannel(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock Channel": {
			request: `{
				"accountid": "1",
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA"
			}`,
			reply: `{
				"id": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"href": "https://localhost:5010/api/v1/channel/H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"public_read": true,
				"public_write": true,
				"sequenced": true,
				"locked": false,
				"head": 0,
				"retention": {
					"min_age_days": 0,
					"max_age_days": 99999,
					"auto_prune": true
				},
				"access_tokens": [
					{
						"id": "1",
						"token": "20_j2-GfF6GFk8lnofe7EW5u7DhztfLQmRsa8d8R3CBZCGVU7xS1vhQwqfT-K-P2PLyxkS1wznAbj1VF1U3TFA",
						"description": "Owner",
						"can_read": true,
						"can_write": true
					}
				]
			}`,
			err:  nil,
			code: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req ChannelRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.Channel(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp ChannelReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestUnitChannelUpdate(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock ChannelUpdate": {
			request: `{
				"accountid": "1",
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"PublicRead": true,
				"PublicWrite": true,
				"Locked": true
			  }`,
			reply: `{
				"PublicRead": true,
				"PublicWrite": true,
				"Locked": true
			  }`,
			err:  nil,
			code: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req ChannelUpdateRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.ChannelUpdate(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp ChannelUpdateReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestUnitChannelDelete(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock ChannelDelete": {
			request: `{
				"accountid": "1",
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA"
			  }`,
			reply: "{}",
			err:   nil,
			code:  http.StatusNoContent,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req ChannelDeleteRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			err := c.ChannelDelete(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}
		})
	}
}

func TestUnitChannelCreate(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock ChannelCreate": {
			request: `{
				"accountid": "1",
				"public_read": true,
				"public_write": true,
				"sequenced": true,
				"retention": {
				  "min_age_days": 0,
				  "max_age_days": 99999,
				  "auto_prune": true
				}
			  }`,
			reply: `{
				"ID": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"Href": "https://localhost:5010/api/v1/channel/H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"PublicRead": true,
				"PublicWrite": true,
				"Sequenced": true,
				"Locked": false,
				"Head": 0,
				"Retention": {
					"MinAgeDays": 0,
					"MaxAgeDays": 99999,
					"AutoPrune": true
				},
				"AccessTokens": [
					{
						"ID": "1",
						"Token": "OEdvoTD3ozLxDfXrko2J3RKNHI7LrGW-sxyYF1aoLUNJI2mcFH9CMQXv3oRPbkcgx0EM3nEhYT61F6T72sPXEA",
						"Description": "Owner",
						"CanRead": true,
						"CanWrite": true
					}
				]
			}`,
			err:  nil,
			code: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req ChannelCreateRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.ChannelCreate(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp ChannelCreateReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestUnitToken(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock Token": {
			request: `{
				"accountid": "1",
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"tokenid": "1"
			  }`,
			reply: `{
				"ID": "1",
				"Token": "20_j2-GfF6GFk8lnofe7EW5u7DhztfLQmRsa8d8R3CBZCGVU7xS1vhQwqfT-K-P2PLyxkS1wznAbj1VF1U3TFA",
				"Description": "Owner",
				"CanRead": true,
				"CanWrite": true
			}`,
			err:  nil,
			code: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req TokenRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.Token(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp TokenReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestUnitTokenDelete(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock TokenDelete": {
			request: `{
				"accountid": "1",
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"tokenid": "1"
			  }`,
			reply: "",
			err:   nil,
			code:  http.StatusNoContent,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req TokenDeleteRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			err := c.TokenDelete(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}
		})
	}
}

func TestUnitTokens(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock Tokens": {
			request: `{
				"accountid": "1",
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA"
			  }`,
			reply: `[
				{
					"ID": "1",
					"Token": "20_j2-GfF6GFk8lnofe7EW5u7DhztfLQmRsa8d8R3CBZCGVU7xS1vhQwqfT-K-P2PLyxkS1wznAbj1VF1U3TFA",
					"Description": "Owner",
					"CanRead": true,
					"CanWrite": true
				}
			]`,
			err:  nil,
			code: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req TokensRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.Tokens(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp TokensReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestUnitTokenCreate(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock TokenCreate": {
			request: `{
				"accountid": "1",
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"description": "Owner",
				"can_read": true,
				"can_write": true
			  }`,
			reply: `{
					"ID": "1",
					"Token": "20_j2-GfF6GFk8lnofe7EW5u7DhztfLQmRsa8d8R3CBZCGVU7xS1vhQwqfT-K-P2PLyxkS1wznAbj1VF1U3TFA",
					"Description": "Owner",
					"CanRead": true,
					"CanWrite": true
				}`,
			err:  nil,
			code: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			c.HTTPClient = &MockClient{
				MockDo: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: test.code,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(strings.Join(strings.Fields(test.reply), "")))),
					}, nil
				},
			}

			var req TokenCreateRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.TokenCreate(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp TokenCreateReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}
