package mattercloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// AddressBalance this endpoint retrieves balance for a specific address.
//
// For more information: https://developers.mattercloud.net/#get-balance
func (c *Client) AddressBalance(address string) (balance *Balance, err error) {

	var resp string

	// GET https://api.mattercloud.net/api/v3/main/address/<address>/balance
	if resp, err = c.Request("address/"+address+"/balance", http.MethodGet, nil); err != nil {
		return
	}

	// Check for error
	if c.LastRequest.StatusCode != http.StatusOK {
		var apiError APIInternalError
		if err = json.Unmarshal([]byte(resp), &apiError); err != nil {
			return
		}
		err = fmt.Errorf("error: %s", apiError.ErrorMessage)
		return
	}

	// Process the response
	err = json.Unmarshal([]byte(resp), &balance)
	return
}

// AddressBalanceBatch this endpoint retrieves balances for multiple addresses at same time
//
// For more information: https://developers.mattercloud.net/#get-balance-batch
func (c *Client) AddressBalanceBatch(addresses []string) (balances []*Balance, err error) {

	// Check addresses
	if len(addresses) == 0 {
		err = fmt.Errorf("missing addresses")
		return
	}

	// Marshall into JSON
	var data []byte
	if data, err = json.Marshal(&AddressList{Addrs: strings.Join(addresses, ",")}); err != nil {
		return
	}

	var resp string

	// POST https://api.mattercloud.net/api/v3/main/address/balance
	if resp, err = c.Request("address/balance", http.MethodPost, data); err != nil {
		return
	}

	// Check for error
	if c.LastRequest.StatusCode != http.StatusOK {
		var apiError APIInternalError
		if err = json.Unmarshal([]byte(resp), &apiError); err != nil {
			return
		}
		err = fmt.Errorf("error: %s", apiError.ErrorMessage)
		return
	}

	// Process the response
	err = json.Unmarshal([]byte(resp), &balances)
	return
}

// AddressUtxos this endpoint retrieves utxos for a specific address
//
// For more information: https://developers.mattercloud.net/#get-utxos
func (c *Client) AddressUtxos(address string) (utxos []*UnspentTransaction, err error) {

	var resp string

	// GET https://api.mattercloud.net/api/v3/main/address/<address>/utxo
	if resp, err = c.Request("address/"+address+"/utxo", http.MethodGet, nil); err != nil {
		return
	}

	// Check for error
	if c.LastRequest.StatusCode != http.StatusOK {
		var apiError APIInternalError
		if err = json.Unmarshal([]byte(resp), &apiError); err != nil {
			return
		}
		err = fmt.Errorf("error: %s", apiError.ErrorMessage)
		return
	}

	// Process the response
	err = json.Unmarshal([]byte(resp), &utxos)
	return
}

// AddressUtxosBatch this endpoint retrieves utxos for multiple addresses
//
// For more information: https://developers.mattercloud.net/#get-utxos-batch
func (c *Client) AddressUtxosBatch(addresses []string) (utxos []*UnspentTransaction, err error) {

	// Check addresses
	if len(addresses) == 0 {
		err = fmt.Errorf("missing addresses")
		return
	}

	// Marshall into JSON
	var data []byte
	if data, err = json.Marshal(&AddressList{Addrs: strings.Join(addresses, ",")}); err != nil {
		return
	}

	var resp string

	// POST https://api.mattercloud.net/api/v3/main/address/utxo
	if resp, err = c.Request("address/utxo", http.MethodPost, data); err != nil {
		return
	}

	// Check for error
	if c.LastRequest.StatusCode != http.StatusOK {
		var apiError APIInternalError
		if err = json.Unmarshal([]byte(resp), &apiError); err != nil {
			return
		}
		err = fmt.Errorf("error: %s", apiError.ErrorMessage)
		return
	}

	// Process the response
	err = json.Unmarshal([]byte(resp), &utxos)
	return
}

// AddressHistory this endpoint retrieves history for a specific address
//
// For more information: https://developers.mattercloud.net/#get-history
func (c *Client) AddressHistory(address string) (history *History, err error) {

	var resp string

	// GET https://api.mattercloud.net/api/v3/main/address/<address>/history
	if resp, err = c.Request("address/"+address+"/history", http.MethodGet, nil); err != nil {
		return
	}

	// Check for error
	if c.LastRequest.StatusCode != http.StatusOK {
		var apiError APIInternalError
		if err = json.Unmarshal([]byte(resp), &apiError); err != nil {
			return
		}
		err = fmt.Errorf("error: %s", apiError.ErrorMessage)
		return
	}

	// Process the response
	err = json.Unmarshal([]byte(resp), &history)
	return
}

// AddressHistoryBatch this endpoint retrieves history for multiple addresses
//
// For more information: https://developers.mattercloud.net/#get-history-batch
func (c *Client) AddressHistoryBatch(addresses []string) (history *History, err error) {

	// Check addresses
	if len(addresses) == 0 {
		err = fmt.Errorf("missing addresses")
		return
	}

	// Marshall into JSON
	var data []byte
	if data, err = json.Marshal(&AddressList{Addrs: strings.Join(addresses, ",")}); err != nil {
		return
	}

	var resp string

	// POST https://api.mattercloud.net/api/v3/main/address/history
	if resp, err = c.Request("address/history", http.MethodPost, data); err != nil {
		return
	}

	// Check for error
	if c.LastRequest.StatusCode != http.StatusOK {
		var apiError APIInternalError
		if err = json.Unmarshal([]byte(resp), &apiError); err != nil {
			return
		}
		err = fmt.Errorf("error: %s", apiError.ErrorMessage)
		return
	}

	// Process the response
	err = json.Unmarshal([]byte(resp), &history)
	return
}
