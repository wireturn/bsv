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

// mockHTTPGetPayment for mocking requests
type mockHTTPGetPayment struct{}

// Do is a mock http request
func (m *mockHTTPGetPayment) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Beta
	if req.URL.String() == environments[EnvironmentBeta].APIURL+endpointGetPaymentRequest {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"transactionId":"4eb7ab228ab9a23831b5b788e3f0eb5bed6dcdbb6d9d808eaba559c49afb9b0a","note":"Test description","type":"send","time":1608222315,"satoshiFees":113,"satoshiAmount":5301,"fiatExchangeRate":188.6311183023935,"fiatCurrencyCode":"USD","participants":[{"type":"user","alias":"mrz@moneybutton.com","displayName":"MrZ","profilePictureUrl":"https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon","responseNote":""}],"attachments":[],"appAction":"like"}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPInvalidPaymentData for mocking requests
type mockHTTPInvalidPaymentData struct{}

// Do is a mock http request
func (m *mockHTTPInvalidPaymentData) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	resp.StatusCode = http.StatusOK
	resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"invalid":"payment"}`)))

	// Default is valid
	return resp, nil
}

func TestClient_GetPayment(t *testing.T) {
	t.Parallel()

	t.Run("missing auth token", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetPayment{}, EnvironmentBeta)
		assert.NotNil(t, client)
		payment, err := client.GetPayment(context.Background(), "", "")
		assert.Error(t, err)
		assert.Nil(t, payment)
	})

	t.Run("missing transaction id", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetPayment{}, EnvironmentBeta)
		assert.NotNil(t, client)
		payment, err := client.GetPayment(context.Background(), "000000", "")
		assert.Error(t, err)
		assert.Nil(t, payment)
	})

	t.Run("invalid auth token", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetPayment{}, EnvironmentBeta)
		assert.NotNil(t, client)
		payment, err := client.GetPayment(context.Background(), "0", "000000")
		assert.Error(t, err)
		assert.Nil(t, payment)
	})

	t.Run("bad request", func(t *testing.T) {
		client := newTestClient(&mockHTTPBadRequest{}, EnvironmentBeta)
		assert.NotNil(t, client)
		payment, err := client.GetPayment(context.Background(), "000000", "000000")
		assert.Error(t, err)
		assert.Nil(t, payment)
	})

	t.Run("invalid payment data", func(t *testing.T) {
		client := newTestClient(&mockHTTPInvalidPaymentData{}, EnvironmentBeta)
		assert.NotNil(t, client)
		payment, err := client.GetPayment(context.Background(), "000000", "000000")
		assert.Error(t, err)
		assert.Nil(t, payment)
	})

	t.Run("valid payment response", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetPayment{}, EnvironmentBeta)
		assert.NotNil(t, client)
		payment, err := client.GetPayment(context.Background(), "000000", "000000")
		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.Equal(t, "4eb7ab228ab9a23831b5b788e3f0eb5bed6dcdbb6d9d808eaba559c49afb9b0a", payment.TransactionID)
		assert.Equal(t, uint64(1608222315), payment.Time)
		assert.Equal(t, PaymentSend, payment.Type)
		assert.Equal(t, AppActionLike, payment.AppAction)
		assert.Equal(t, CurrencyUSD, payment.FiatCurrencyCode)
		assert.Equal(t, 188.6311183023935, payment.FiatExchangeRate)
		assert.Equal(t, "Test description", payment.Note)
		assert.Equal(t, uint64(5301), payment.SatoshiAmount)
		assert.Equal(t, uint64(113), payment.SatoshiFees)
		assert.Equal(t, 1, len(payment.Participants))
		assert.Equal(t, ParticipantUser, payment.Participants[0].Type)
		assert.Equal(t, "", payment.Participants[0].ResponseNote)
		assert.Equal(t, "mrz@moneybutton.com", payment.Participants[0].Alias)
		assert.Equal(t, "MrZ", payment.Participants[0].DisplayName)
		assert.Equal(t, "https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon", payment.Participants[0].ProfilePictureURL)
	})
}
