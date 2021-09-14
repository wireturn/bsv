package tonicpow

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetCurrentRate will get a current rate for the given currency (using default currency amount)
//
// For more information: https://docs.tonicpow.com/#71b8b7fc-317a-4e68-bd2a-5b0da012361c
func (c *Client) GetCurrentRate(currency string,
	customAmount float64) (rate *Rate, response *StandardResponse, err error) {

	// Currency is required
	if len(currency) == 0 {
		err = fmt.Errorf("missing required attribute: %s", fieldCurrency)
		return
	}

	// Fire the Request
	if response, err = c.Request(
		http.MethodGet,
		fmt.Sprintf("/%s/%s?%s=%f", modelRates, currency, fieldAmount, customAmount),
		nil, http.StatusOK,
	); err != nil {
		return
	}

	err = json.Unmarshal(response.Body, &rate)
	return
}
