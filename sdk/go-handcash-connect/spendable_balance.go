package handcash

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SpendableBalanceResponse is the balance response
type SpendableBalanceResponse struct {
	SpendableSatoshiBalance uint64       `json:"spendableSatoshiBalance"`
	SpendableFiatBalance    float64      `json:"spendableFiatBalance"`
	CurrencyCode            CurrencyCode `json:"currencyCode"`
}

// GetSpendableBalance gets the user's spendable balance from the handcash connect API
func (c *Client) GetSpendableBalance(ctx context.Context, authToken string,
	currencyCode CurrencyCode) (*SpendableBalanceResponse, error) {

	// Make sure we have an auth token
	if len(authToken) == 0 {
		return nil, fmt.Errorf("missing auth token")
	}

	if len(currencyCode) == 0 {
		return nil, fmt.Errorf("missing currency code")
	}

	// Get the signed request
	signed, err := c.getSignedRequest(
		http.MethodGet,
		endpointGetSpendableBalanceRequest,
		authToken,
		&BalanceRequest{CurrencyCode: currencyCode},
		currentISOTimestamp(),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating signed request: %w", err)
	}

	// Convert into bytes
	var params []byte
	if params, err = json.Marshal(
		&BalanceRequest{CurrencyCode: currencyCode},
	); err != nil {
		return nil, err
	}

	// Make the HTTP request
	response := httpRequest(
		ctx,
		c,
		&httpPayload{
			Data:           params,
			ExpectedStatus: http.StatusOK,
			Method:         signed.Method,
			URL:            signed.URI,
		},
		signed,
	)

	if response.Error != nil {
		return nil, response.Error
	}

	spendableBalanceResponse := new(SpendableBalanceResponse)

	if err = json.Unmarshal(response.BodyContents, &spendableBalanceResponse); err != nil {
		return nil, fmt.Errorf("failed unmarshal %w", err)
	} else if spendableBalanceResponse == nil || spendableBalanceResponse.CurrencyCode == "" {
		return nil, fmt.Errorf("failed to get balance")
	}

	return spendableBalanceResponse, nil
}
