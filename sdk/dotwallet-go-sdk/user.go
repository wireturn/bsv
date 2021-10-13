package dotwallet

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// UserInfo can obtain information authorized by DotWallet users via their user access_token
//
// For more information: https://developers.dotwallet.com/documents/en/#user-info
func (c *Client) UserInfo(userToken *DotAccessToken) (*User, error) {

	// Make the request
	response, err := c.Request(
		http.MethodPost,
		getUserInfo,
		nil,
		http.StatusOK,
		userToken,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response
	resp := new(userResponse)
	if err = json.Unmarshal(
		response.Body, &resp,
	); err != nil {
		return nil, err
	}

	// Error?
	if resp.Code > 0 {
		return nil, fmt.Errorf(resp.Message)
	}

	return &resp.Data.User, nil
}

// UserReceiveAddress will return the wallets for the user based on the given coin_type
// Currently supported: BSV, BTC and ETH
// Paymail is currently only supported on BSV
//
// For more information: https://developers.dotwallet.com/documents/en/#get-user-receive-address
func (c *Client) UserReceiveAddress(userID string, coinType coinType) (*Wallets, error) {

	// Make the request
	response, err := c.Request(
		http.MethodPost,
		getUserReceiveAddress,
		&userReceiveRequest{
			CoinType: coinType,
			UserID:   userID,
		},
		http.StatusOK,
		c.Token(),
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response
	resp := new(userReceiveAddressResponse)
	if err = json.Unmarshal(
		response.Body, &resp,
	); err != nil {
		return nil, err
	}

	// Error?
	if resp.Code > 0 {
		return nil, fmt.Errorf(resp.Message)
	}

	return &resp.Data.Wallets, nil
}
