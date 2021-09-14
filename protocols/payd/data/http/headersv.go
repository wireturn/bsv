package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/libsv/go-bc"
	"github.com/pkg/errors"
)

type headersv struct {
	client HTTPClient
	host   string
}

// NewHeadersv returns a bc.BlockHeaderChain using Headersv.
func NewHeadersv(client HTTPClient, host string) bc.BlockHeaderChain {
	return &headersv{
		client: client,
		host:   host,
	}
}

// BlockHeader returns the header for the provided blockhash.
func (h *headersv) BlockHeader(ctx context.Context, blockHash string) (*bc.BlockHeader, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/v1/chain/header/%s", h.host, blockHash),
		nil,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating request for chain/header/%s", blockHash)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse error message body")
		}

		return nil, fmt.Errorf("block header request: unexpected status code %d\nresponse body:\n%s", resp.StatusCode, body)
	}

	var bh *bc.BlockHeader
	if err = json.NewDecoder(resp.Body).Decode(&bh); err != nil {
		return nil, err
	}

	return bh, nil
}
