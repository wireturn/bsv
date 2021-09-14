package whatsonchain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// GetChainInfo this endpoint retrieves various state info of the chain for the selected network.
//
// For more information: https://developers.whatsonchain.com/#chain-info
func (c *Client) GetChainInfo() (chainInfo *ChainInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/chain/info
	if resp, err = c.request(fmt.Sprintf("%s%s/chain/info", apiEndpoint, c.Network), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &chainInfo)
	return
}

// GetCirculatingSupply this endpoint retrieves the current circulating supply
//
// For more information: (undocumented) //todo: add link once in documentation
func (c *Client) GetCirculatingSupply() (supply float64, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/circulatingsupply
	if resp, err = c.request(fmt.Sprintf("%s%s/circulatingsupply", apiEndpoint, c.Network), http.MethodGet, nil); err != nil {
		return
	}

	supply, err = strconv.ParseFloat(strings.TrimSpace(resp), 64)
	return
}
