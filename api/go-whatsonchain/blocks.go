package whatsonchain

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetBlockByHash this endpoint retrieves block details with given hash.
//
// For more information: https://developers.whatsonchain.com/#get-by-hash
func (c *Client) GetBlockByHash(hash string) (blockInfo *BlockInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/hash/<hash>
	if resp, err = c.request(fmt.Sprintf("%s%s/block/hash/%s", apiEndpoint, c.Network, hash), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &blockInfo)
	return
}

// GetBlockByHeight this endpoint retrieves block details with given block height.
//
// For more information: https://developers.whatsonchain.com/#get-by-height
func (c *Client) GetBlockByHeight(height int64) (blockInfo *BlockInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/height/<height>
	if resp, err = c.request(fmt.Sprintf("%s%s/block/height/%d", apiEndpoint, c.Network, height), http.MethodGet, nil); err != nil {
		return
	}

	// Added this condition (for bad block heights, 200 success but empty result - no JSON response)
	if len(resp) > 0 {
		err = json.Unmarshal([]byte(resp), &blockInfo)
	}
	return
}

// GetBlockPages If the block has more that 1000 transactions the page URIs will
// be provided in the pages element when getting a block by hash or height.
//
// For more information: https://developers.whatsonchain.com/#get-block-pages
func (c *Client) GetBlockPages(hash string, page int) (txList BlockPagesInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/hash/<hash>/page/1
	if resp, err = c.request(fmt.Sprintf("%s%s/block/hash/%s/page/%d", apiEndpoint, c.Network, hash, page), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &txList)
	return
}

// GetHeaderByHash this endpoint retrieves block header details with given hash.
//
// For more information: https://developers.whatsonchain.com/#get-header-by-hash
func (c *Client) GetHeaderByHash(hash string) (headerInfo *BlockInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/<hash>/header
	if resp, err = c.request(fmt.Sprintf("%s%s/block/%s/header", apiEndpoint, c.Network, hash), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &headerInfo)
	return
}

// GetHeaders this endpoint retrieves last 10 block headers.
//
// For more information: https://developers.whatsonchain.com/#get-headers
func (c *Client) GetHeaders() (blockHeaders []*BlockInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/headers
	if resp, err = c.request(fmt.Sprintf("%s%s/block/headers", apiEndpoint, c.Network), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &blockHeaders)
	return
}
