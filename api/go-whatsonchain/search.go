package whatsonchain

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetExplorerLinks this endpoint identifies whether the posted query text is a block hash, txid or address and
// responds with WoC links. Ideal for extending customized search in apps.
//
// For more information: https://developers.whatsonchain.com/#get-history
func (c *Client) GetExplorerLinks(query string) (results SearchResults, err error) {

	// Start the post data
	stringVal := fmt.Sprintf(`{"query":"%s"}`, query)
	postData := []byte(stringVal)

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/search/links
	if resp, err = c.request(fmt.Sprintf("%s%s/search/links", apiEndpoint, c.Network), http.MethodPost, postData); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &results)
	return
}
