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

// mockHTTPPay for mocking requests
type mockHTTPPay struct{}

// Do is a mock http request
func (m *mockHTTPPay) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Beta
	if req.URL.String() == environments[EnvironmentBeta].APIURL+endpointGetPayRequest {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"transactionId":"05d7df52a1c58cabada16709469e6940342cb13e8cfa3c7e1438d7ea84765787","note":"Thanks dude!","type":"send","time":1608226019,"satoshiFees":127,"satoshiAmount":5372,"fiatExchangeRate":186.15198556884275,"fiatCurrencyCode":"USD","participants":[{"type":"user","alias":"mrz@moneybutton.com","displayName":"MrZ","profilePictureUrl":"https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon","responseNote":""}],"attachments":[{"value":{"some":"data"},"format":"json"}],"appAction":"like","rawTransactionHex":"01000000018598fbea559e4a59772361994f800adb63bab592e276de7ebd5805ecc639b3b8010000006a47304402200fc98489e2bbba5cb7f8cea970c0037585d42618ef60d172179307b4446854a802206be468ffd31f97c6e01a6549be50241d42633e32ba4e06ff4b2565ec897232a2412103c1fbc71737d3820890535112ac99b2471d6bacbd8a7e7825c65863a67b1d0c7effffffff03000000000000000012006a0f7b22736f6d65223a2264617461227dfc140000000000001976a914b7ce7a4c1350f1cb9dcaecca10d48f064be9197f88ac57020000000000001976a9145233794b8bdf2fd7f809b11da081189d2e79000c88ac00000000"}`)))
	}

	// Default is valid
	return resp, nil
}

func TestClient_Pay(t *testing.T) {
	t.Parallel()

	t.Run("missing auth token", func(t *testing.T) {
		client := newTestClient(&mockHTTPPay{}, EnvironmentBeta)
		assert.NotNil(t, client)
		payment, err := client.Pay(context.Background(), "", nil)
		assert.Error(t, err)
		assert.Nil(t, payment)
	})

	t.Run("missing payment parameters", func(t *testing.T) {
		client := newTestClient(&mockHTTPPay{}, EnvironmentBeta)
		assert.NotNil(t, client)
		payment, err := client.Pay(context.Background(), "000000", nil)
		assert.Error(t, err)
		assert.Nil(t, payment)
	})

	t.Run("missing payment receivers", func(t *testing.T) {
		client := newTestClient(&mockHTTPPay{}, EnvironmentBeta)
		assert.NotNil(t, client)

		payParams := &PayParameters{
			AppAction:   AppActionLike,
			Description: "Test description",
		}

		payment, err := client.Pay(context.Background(), "000000", payParams)
		assert.Error(t, err)
		assert.Nil(t, payment)
	})

	t.Run("invalid auth token", func(t *testing.T) {
		client := newTestClient(&mockHTTPPay{}, EnvironmentBeta)
		assert.NotNil(t, client)

		payParams := &PayParameters{
			AppAction:   AppActionLike,
			Description: "Test description",
			Receivers: []*Payment{{
				Amount:       0.01,
				CurrencyCode: CurrencyUSD,
				To:           "mrz@moneybutton.com",
			}},
		}

		payment, err := client.Pay(context.Background(), "0", payParams)
		assert.Error(t, err)
		assert.Nil(t, payment)
	})

	t.Run("bad request", func(t *testing.T) {
		client := newTestClient(&mockHTTPBadRequest{}, EnvironmentBeta)
		assert.NotNil(t, client)
		payParams := &PayParameters{
			AppAction:   AppActionLike,
			Description: "Test description",
			Receivers: []*Payment{{
				Amount:       0.01,
				CurrencyCode: CurrencyUSD,
				To:           "mrz@moneybutton.com",
			}},
		}

		payment, err := client.Pay(context.Background(), "000000", payParams)
		assert.Error(t, err)
		assert.Nil(t, payment)
	})

	t.Run("invalid payment data", func(t *testing.T) {
		client := newTestClient(&mockHTTPInvalidPaymentData{}, EnvironmentBeta)
		assert.NotNil(t, client)
		payParams := &PayParameters{
			AppAction:   AppActionLike,
			Description: "Test description",
			Receivers: []*Payment{{
				Amount:       0.01,
				CurrencyCode: CurrencyUSD,
				To:           "mrz@moneybutton.com",
			}},
		}

		payment, err := client.Pay(context.Background(), "000000", payParams)
		assert.Error(t, err)
		assert.Nil(t, payment)
	})

	t.Run("valid payment", func(t *testing.T) {
		client := newTestClient(&mockHTTPPay{}, EnvironmentBeta)
		assert.NotNil(t, client)
		payParams := &PayParameters{
			AppAction:   AppActionLike,
			Description: "Test description",
			Receivers: []*Payment{{
				Amount:       0.01,
				CurrencyCode: CurrencyUSD,
				To:           "mrz@moneybutton.com",
			}},
		}

		payment, err := client.Pay(context.Background(), "000000", payParams)
		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.Equal(t, "05d7df52a1c58cabada16709469e6940342cb13e8cfa3c7e1438d7ea84765787", payment.TransactionID)
		assert.Equal(t, "01000000018598fbea559e4a59772361994f800adb63bab592e276de7ebd5805ecc639b3b8010000006a47304402200fc98489e2bbba5cb7f8cea970c0037585d42618ef60d172179307b4446854a802206be468ffd31f97c6e01a6549be50241d42633e32ba4e06ff4b2565ec897232a2412103c1fbc71737d3820890535112ac99b2471d6bacbd8a7e7825c65863a67b1d0c7effffffff03000000000000000012006a0f7b22736f6d65223a2264617461227dfc140000000000001976a914b7ce7a4c1350f1cb9dcaecca10d48f064be9197f88ac57020000000000001976a9145233794b8bdf2fd7f809b11da081189d2e79000c88ac00000000", payment.RawTransactionHex)
		assert.Equal(t, uint64(1608226019), payment.Time)
		assert.Equal(t, PaymentSend, payment.Type)
		assert.Equal(t, AppActionLike, payment.AppAction)
		assert.Equal(t, CurrencyUSD, payment.FiatCurrencyCode)
		assert.Equal(t, 186.15198556884275, payment.FiatExchangeRate)
		assert.Equal(t, "Thanks dude!", payment.Note)
		assert.Equal(t, uint64(5372), payment.SatoshiAmount)
		assert.Equal(t, uint64(127), payment.SatoshiFees)
		assert.Equal(t, 1, len(payment.Participants))
		assert.Equal(t, ParticipantUser, payment.Participants[0].Type)
		assert.Equal(t, "", payment.Participants[0].ResponseNote)
		assert.Equal(t, "mrz@moneybutton.com", payment.Participants[0].Alias)
		assert.Equal(t, "MrZ", payment.Participants[0].DisplayName)
		assert.Equal(t, "https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon", payment.Participants[0].ProfilePictureURL)
	})
}
