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

func TestUnitMessageHead(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock MessageHead": {
			request: `{
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA"
			}`,
			reply: "",
			err:   nil,
			code:  http.StatusOK,
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

			var req MessageHeadRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			err := c.MessageHead(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}
		})
	}
}

func TestUnitMessageWrite(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock MessageWrite": {
			request: `{
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"message": "Hello this is a message"
			}`,
			reply: `{
				"sequence": 1,
				"received": "2021-08-31T18:43:07.855547Z",
				"content_type": "application/json",
				"payload": "SGVsbG8gdGhpcyBpcyBhIG1lc3NhZ2U="
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

			var req MessageWriteRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.MessageWrite(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp MessageWriteReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestUnitMessages(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock Messages": {
			request: `{
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"unread": true
			}`,
			reply: `[
				{
					"sequence": 1,
					"received": "2021-08-31T17:40:50.618865Z",
					"content_type": "application/json;charset=utf-8",
					"payload": "ZnJvbSBvd25lcg=="
				},
				{
					"sequence": 2,
					"received": "2021-08-31T17:41:54.480861Z",
					"content_type": "application/json;charset=utf-8",
					"payload": "ZnJvbSBvd25lcg=="
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

			var req MessagesRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			resp, err := c.Messages(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}

			var expectedResp MessagesReply
			if err := json.Unmarshal([]byte(test.reply), &expectedResp); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			assert.Equal(t, *resp, expectedResp)
		})
	}
}

func TestUnitMessageMark(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock MessageMark": {
			request: `{
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"sequence": 1,
				"older": true,
				"read": true
			}`,
			reply: "",
			err:   nil,
			code:  http.StatusOK,
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

			var req MessageMarkRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			err := c.MessageMark(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}
		})
	}
}

func TestUnitMessageDelete(t *testing.T) {
	tests := map[string]struct {
		request string
		reply   string
		err     error
		code    int
	}{
		"Mock MessageDelete": {
			request: `{
				"channelid": "H3mNdK-IL_-5OdLG4jymMwlJCW7NlhsNhxd_XrnKlv7J4hyR6EH2NIOaPmWlU7Rs0Zkgv_1yD0qcW7h29BGxbA",
				"sequence": 1
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

			var req MessageDeleteRequest
			if err := json.Unmarshal([]byte(test.request), &req); err != nil {
				assert.Fail(t, "error unmarshalling test json", err)
			}
			err := c.MessageDelete(context.Background(), req)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
				return
			}
		})
	}
}
