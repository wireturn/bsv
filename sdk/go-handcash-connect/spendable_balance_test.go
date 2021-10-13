package handcash

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockHTTPGetSpendableBalance for mocking requests
type mockHTTPGetSpendableBalance struct{}

// Do is a mock http request
func (m *mockHTTPGetSpendableBalance) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Beta
	if req.URL.String() == environments[EnvironmentBeta].APIURL+endpointGetSpendableBalanceRequest {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"spendableSatoshiBalance":1424992,"spendableFiatBalance":2.7792,"currencyCode":"USD"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPInvalidSpendableBalanceData for mocking requests
type mockHTTPInvalidSpendableBalanceData struct{}

// Do is a mock http request
func (m *mockHTTPInvalidSpendableBalanceData) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	resp.StatusCode = http.StatusOK
	resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"invalid":"currencyCode"}`)))

	// Default is valid
	return resp, nil
}

func TestClient_GetSpendableBalance(t *testing.T) {
	t.Parallel()

	t.Run("missing auth token", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetSpendableBalance{}, EnvironmentBeta)
		assert.NotNil(t, client)
		balance, err := client.GetSpendableBalance(context.Background(), "", "")
		assert.Error(t, err)
		assert.Nil(t, balance)
	})

	t.Run("missing currency code", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetSpendableBalance{}, EnvironmentBeta)
		assert.NotNil(t, client)
		balance, err := client.GetSpendableBalance(context.Background(), "000000", "")
		assert.Error(t, err)
		assert.Nil(t, balance)
	})

	t.Run("invalid auth token", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetSpendableBalance{}, EnvironmentBeta)
		assert.NotNil(t, client)
		balance, err := client.GetSpendableBalance(context.Background(), "0", "USD")
		assert.Error(t, err)
		assert.Nil(t, balance)
	})

	t.Run("invalid currency code", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetSpendableBalance{}, EnvironmentBeta)
		assert.NotNil(t, client)
		balance, err := client.GetSpendableBalance(context.Background(), "0", "FOO")
		assert.Error(t, err)
		assert.Nil(t, balance)
	})

	t.Run("bad request", func(t *testing.T) {
		client := newTestClient(&mockHTTPBadRequest{}, EnvironmentBeta)
		assert.NotNil(t, client)
		balance, err := client.GetSpendableBalance(context.Background(), "000000", "USD")
		assert.Error(t, err)
		assert.Nil(t, balance)
	})

	t.Run("invalid spendable balance data", func(t *testing.T) {
		client := newTestClient(&mockHTTPInvalidSpendableBalanceData{}, EnvironmentBeta)
		assert.NotNil(t, client)
		balance, err := client.GetSpendableBalance(context.Background(), "000000", "USD")
		assert.Error(t, err)
		assert.Nil(t, balance)
	})

	t.Run("valid spendable balance response", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetSpendableBalance{}, EnvironmentBeta)
		assert.NotNil(t, client)
		balance, err := client.GetSpendableBalance(context.Background(), "000000", "USD")
		assert.NoError(t, err)
		assert.NotNil(t, balance)
		assert.Equal(t, CurrencyUSD, balance.CurrencyCode)
		assert.Equal(t, uint64(1424992), balance.SpendableSatoshiBalance)
		assert.Equal(t, float64(2.7792), balance.SpendableFiatBalance)
	})
}
