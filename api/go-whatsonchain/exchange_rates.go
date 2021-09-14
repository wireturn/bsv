package whatsonchain

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetExchangeRate this endpoint provides exchange rate for BSV
//
// For more information: https://developers.whatsonchain.com/#get-exchange-rate
func (c *Client) GetExchangeRate() (rate *ExchangeRate, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/exchangerate
	if resp, err = c.request(fmt.Sprintf("%s%s/exchangerate", apiEndpoint, c.Network), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &rate)
	return
}
