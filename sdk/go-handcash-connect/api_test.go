package handcash

import (
	"encoding/hex"
	"net/http"
	"testing"

	"github.com/bitcoinschema/go-bitcoin"
	"github.com/stretchr/testify/assert"
)

const (
	testTimestamp = "2020-12-10T16:31:23.304Z"
)

func TestGetSignedRequest(t *testing.T) {
	t.Parallel()

	t.Run("valid auth - empty body", func(t *testing.T) {

		client := newTestClient(&mockHTTPDefaultClient{}, EnvironmentBeta)
		assert.NotNil(t, client)

		// These values are used (GOSEC complains about the token)
		token := "68d8fadc95324afa853f00923e0b" + "86f06a76ceb7a6afbb1784e0dde8f43989a0"
		method := http.MethodGet
		timestamp := testTimestamp
		endpoint := endpointProfileCurrent

		// privateKey := "5JcTmjpJkfkcRnf3W2qTvauC4mczNsnUY3SLm6EcKQDS3Gj2wGh" // 68d8fadc95324afa853f00923e0b86f06a76ceb7a6afbb1784e0dde8f439

		signedRequest, err := client.getSignedRequest(method, endpoint, token, nil, timestamp)
		assert.NoError(t, err)
		assert.NotNil(t, signedRequest)

		assert.Equal(t,
			"0275e7081e5b6e73c94998098e075c0ed888d1eb33c721ee38ee741648b108c90d",
			signedRequest.Headers.OauthPublicKey,
		)
		assert.Equal(t,
			"30450221009b613aa82657e28471406d3a390688bc2dceece75cf73d89088d447cc9d1f5c502200a1ff4f02f5dfdb7f48b51b6dd2a0f586f591b07a89d5ad846392fca8ae0c856",
			signedRequest.Headers.OauthSignature,
		)
		assert.Equal(t, testTimestamp, signedRequest.Headers.OauthTimestamp)
		assert.Equal(t, client.Environment.APIURL+endpointProfileCurrent, signedRequest.URI)
		assert.Equal(t, http.MethodGet, signedRequest.Method)
	})

	t.Run("valid auth - with body", func(t *testing.T) {

		client := newTestClient(&mockHTTPDefaultClient{}, EnvironmentBeta)
		assert.NotNil(t, client)

		// These values are used (GOSEC complains about the token)
		token := "68d8fadc95324afa853f00923e0b" + "86f06a76ceb7a6afbb1784e0dde8f43989a0"
		method := http.MethodGet
		timestamp := testTimestamp
		endpoint := endpointProfileCurrent

		// privateKey := "5JcTmjpJkfkcRnf3W2qTvauC4mczNsnUY3SLm6EcKQDS3Gj2wGh" // 68d8fadc95324afa853f00923e0b86f06a76ceb7a6afbb1784e0dde8f439

		var customBodyContents = struct {
			Name    string  `json:"name"`
			Number  int     `json:"number"`
			Float   float64 `json:"float"`
			Boolean bool    `json:"boolean"`
		}{Name: "TestName", Number: 123, Float: 123.123, Boolean: true}

		signedRequest, err := client.getSignedRequest(method, endpoint, token, customBodyContents, timestamp)
		assert.NoError(t, err)
		assert.NotNil(t, signedRequest)

		assert.Equal(t,
			"0275e7081e5b6e73c94998098e075c0ed888d1eb33c721ee38ee741648b108c90d",
			signedRequest.Headers.OauthPublicKey,
		)
		assert.Equal(t,
			"304402200901cc9d77c01d8ab89583eca396f4b753154cb007173ad212e1175e643cf5e1022026f8ab3e03eb1891bb05a24b742d40192fa8b40d2578c5ffc36fb70cf0fafb6a",
			signedRequest.Headers.OauthSignature,
		)
		assert.Equal(t, testTimestamp, signedRequest.Headers.OauthTimestamp)
		assert.Equal(t, client.Environment.APIURL+endpointProfileCurrent, signedRequest.URI)
		assert.Equal(t, http.MethodGet, signedRequest.Method)
	})

	t.Run("invalid auth token", func(t *testing.T) {
		client := newTestClient(&mockHTTPDefaultClient{}, EnvironmentBeta)
		assert.NotNil(t, client)
		method := http.MethodGet
		endpoint := endpointProfileCurrent
		token := "0"
		timestamp := testTimestamp
		signedRequest, err := client.getSignedRequest(method, endpoint, token, nil, timestamp)
		assert.Error(t, err)
		assert.Nil(t, signedRequest)
	})

	t.Run("valid token, invalid json", func(t *testing.T) {
		client := newTestClient(&mockHTTPDefaultClient{}, EnvironmentBeta)
		assert.NotNil(t, client)
		method := http.MethodGet
		endpoint := endpointProfileCurrent
		token := "68d8fadc95324afa853f00923e0b" + "86f06a76ceb7a6afbb1784e0dde8f43989a0"
		timestamp := testTimestamp
		signedRequest, err := client.getSignedRequest(method, endpoint, token, make(chan int), timestamp)
		assert.Error(t, err)
		assert.Nil(t, signedRequest)
	})
}

func TestGetRequestSignature(t *testing.T) {
	t.Parallel()

	t.Run("valid signature", func(t *testing.T) {
		method := http.MethodGet
		endpoint := endpointProfileCurrent
		timestamp := testTimestamp

		privateKey, err := bitcoin.CreatePrivateKey()
		assert.NoError(t, err)
		assert.NotNil(t, privateKey)

		var signature []byte
		signature, err = getRequestSignature(method, endpoint, nil, timestamp, privateKey)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, len(hex.EncodeToString(signature)))
	})

	t.Run("missing method", func(t *testing.T) {
		method := ""
		endpoint := endpointProfileCurrent
		timestamp := testTimestamp

		privateKey, err := bitcoin.CreatePrivateKey()
		assert.NoError(t, err)
		assert.NotNil(t, privateKey)

		var signature []byte
		signature, err = getRequestSignature(method, endpoint, nil, timestamp, privateKey)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, len(hex.EncodeToString(signature)))
	})

	t.Run("missing endpoint", func(t *testing.T) {
		method := http.MethodGet
		endpoint := ""
		timestamp := testTimestamp

		privateKey, err := bitcoin.CreatePrivateKey()
		assert.NoError(t, err)
		assert.NotNil(t, privateKey)

		var signature []byte
		signature, err = getRequestSignature(method, endpoint, nil, timestamp, privateKey)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, len(hex.EncodeToString(signature)))
	})

	t.Run("missing timestamp", func(t *testing.T) {
		method := http.MethodGet
		endpoint := endpointProfileCurrent
		timestamp := ""

		privateKey, err := bitcoin.CreatePrivateKey()
		assert.NoError(t, err)
		assert.NotNil(t, privateKey)

		var signature []byte
		signature, err = getRequestSignature(method, endpoint, nil, timestamp, privateKey)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, len(hex.EncodeToString(signature)))
	})

	t.Run("custom struct", func(t *testing.T) {
		method := http.MethodGet
		endpoint := endpointProfileCurrent
		timestamp := testTimestamp

		privateKey, err := bitcoin.CreatePrivateKey()
		assert.NoError(t, err)
		assert.NotNil(t, privateKey)

		var customBodyContents = struct {
			Name    string  `json:"name"`
			Number  int     `json:"number"`
			Float   float64 `json:"float"`
			Boolean bool    `json:"boolean"`
		}{Name: "TestName", Number: 123, Float: 123.123, Boolean: true}

		var signature []byte
		signature, err = getRequestSignature(method, endpoint, customBodyContents, timestamp, privateKey)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, len(hex.EncodeToString(signature)))
	})

	t.Run("invalid body / json", func(t *testing.T) {
		method := http.MethodGet
		endpoint := endpointProfileCurrent
		timestamp := testTimestamp

		privateKey, err := bitcoin.CreatePrivateKey()
		assert.NoError(t, err)
		assert.NotNil(t, privateKey)

		var signature []byte
		signature, err = getRequestSignature(method, endpoint, make(chan int), timestamp, privateKey)
		assert.Error(t, err)
		assert.Equal(t, 0, len(signature))
	})
}

func TestGetRequestSignatureHash(t *testing.T) {
	t.Parallel()

	t.Run("valid signature hash", func(t *testing.T) {
		method := http.MethodGet
		endpoint := endpointProfileCurrent
		timestamp := testTimestamp

		hash, err := getRequestSignatureHash(method, endpoint, nil, timestamp)
		assert.NoError(t, err)
		assert.Equal(t, 32, len(hash))
		assert.Equal(t, "aaa21242e579564ec36a3c5108cbd215661fb61cfb8cf17dbf00c074f4561378", hex.EncodeToString(hash))
	})

	t.Run("missing method", func(t *testing.T) {
		method := ""
		endpoint := endpointProfileCurrent
		timestamp := testTimestamp

		hash, err := getRequestSignatureHash(method, endpoint, nil, timestamp)
		assert.NoError(t, err)
		assert.Equal(t, 32, len(hash))
		assert.Equal(t, "1f1917cb12ecd2ef0d245ac193b348f91db912e93d692b504a94ca850ac412f8", hex.EncodeToString(hash))
	})

	t.Run("missing endpoint", func(t *testing.T) {
		method := http.MethodGet
		endpoint := ""
		timestamp := testTimestamp

		hash, err := getRequestSignatureHash(method, endpoint, nil, timestamp)
		assert.NoError(t, err)
		assert.Equal(t, 32, len(hash))
		assert.Equal(t, "b7b2f37bcd5d28ebd048f56160787a7a86b23e6233167e638e288332dedbdb3b", hex.EncodeToString(hash))
	})

	t.Run("missing timestamp", func(t *testing.T) {
		method := http.MethodGet
		endpoint := endpointProfileCurrent
		timestamp := ""

		hash, err := getRequestSignatureHash(method, endpoint, nil, timestamp)
		assert.NoError(t, err)
		assert.Equal(t, 32, len(hash))
		assert.Equal(t, "e4740d26cbe8defe7842304461e08d311aaba77fa7f22e80283a3d7c4ced26cf", hex.EncodeToString(hash))
	})

	t.Run("custom body - basic struct - valid json", func(t *testing.T) {
		method := http.MethodGet
		endpoint := endpointProfileCurrent
		timestamp := testTimestamp

		var customBodyContents = struct {
			Name    string  `json:"name"`
			Number  int     `json:"number"`
			Float   float64 `json:"float"`
			Boolean bool    `json:"boolean"`
		}{Name: "TestName", Number: 123, Float: 123.123, Boolean: true}

		hash, err := getRequestSignatureHash(method, endpoint, customBodyContents, timestamp)
		assert.NoError(t, err)
		assert.Equal(t, 32, len(hash))
		assert.Equal(t, "73a73a79325098f309881d0103e189b73f1a8a247440c93d36852ca036ab7dc5", hex.EncodeToString(hash))
	})

	t.Run("custom body - advanced struct - valid json", func(t *testing.T) {
		method := http.MethodGet
		endpoint := endpointProfileCurrent
		timestamp := testTimestamp

		// todo: nested structs, arrays of structs

		var customBodyContents = struct {
			Name    string   `json:"name"`
			Number  int      `json:"number"`
			Float   float64  `json:"float"`
			Boolean bool     `json:"boolean"`
			Parts   []string `json:"parts"`
		}{Name: "TestName", Number: 123, Float: 123.123, Boolean: true, Parts: []string{"one", "two"}}

		hash, err := getRequestSignatureHash(method, endpoint, customBodyContents, timestamp)
		assert.NoError(t, err)
		assert.Equal(t, 32, len(hash))
		assert.Equal(t, "287e66e5b4ffe062015c078895d70120638080bde96f9a48022add1f4c69f78e", hex.EncodeToString(hash))
	})

	t.Run("invalid body - produces error", func(t *testing.T) {
		method := http.MethodGet
		endpoint := endpointProfileCurrent
		timestamp := testTimestamp

		hash, err := getRequestSignatureHash(method, endpoint, make(chan int), timestamp)
		assert.Error(t, err)
		assert.Equal(t, 0, len(hash))
	})
}
