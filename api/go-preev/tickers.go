package preev

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// GetTickers this endpoint returns last broadcasted data for all available pairs.
// todo: omitted optional parameters: source, include flat (as they do not work for /tickers
//
// For more information: https://preev.pro/api/
func (c *Client) GetTickers(ctx context.Context) (tickerList *TickerList, err error) {

	var resp string
	// https://api.preev.pro/v1/tickers
	if resp, err = c.request(
		ctx, fmt.Sprintf("%s/tickers", apiEndpoint),
	); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &tickerList)
	return
}

// GetTicker this endpoint returns last broadcasted data for the specified pair
// Optional parameters are not being used at this time - as the response changes dramatically
// todo: omitted optional parameters: source, include flat
//
// For more information: https://preev.pro/api/
func (c *Client) GetTicker(ctx context.Context, pairID string) (ticker *Ticker, err error) {

	var resp string
	// https://api.preev.pro/v1/tickers/<pair_id>
	if resp, err = c.request(
		ctx, fmt.Sprintf("%s/tickers/%s", apiEndpoint, pairID),
	); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &ticker)
	return
}

// GetTickerHistory this endpoint returns the historical data of the specified pair.
// Default time period is last 24 hours if no start or end is specified.
//
// Params    Type            Default  Description
// start     Unix timestamp  0	      Time of the first ticker to be returned
// end	     Unix timestamp  0        Time of the last ticker to be returned
// interval	 Integer         1	      Interval in minutes over which to aggregate the data.
// todo: omitted optional parameters: source, include flat
//
// For more information: https://preev.pro/api/
func (c *Client) GetTickerHistory(ctx context.Context, pairID string,
	start, end, interval int64) (tickers []*Ticker, err error) {

	// Start creating the endpoint
	var endpoint *url.URL
	if endpoint, err = url.Parse(
		fmt.Sprintf("%s/tickers/%s/historical", apiEndpoint, pairID),
	); err != nil {
		return
	}

	// Set the query values
	var query url.Values
	if query, err = url.ParseQuery(endpoint.RawQuery); err != nil {
		return
	}

	// Got a start time?
	if start > 0 {
		query.Add("start", strconv.FormatInt(start, 10))
	}

	// Got an end time?
	if end > 0 {
		query.Add("end", strconv.FormatInt(end, 10))
	}

	// Got an interval?
	if interval > 0 {
		query.Add("interval", strconv.FormatInt(interval, 10))
	}

	// Set the query
	endpoint.RawQuery = query.Encode()

	var resp string
	// https://api.preev.pro/v1/tickers/<pair_id>/historical
	if resp, err = c.request(ctx, endpoint.String()); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &tickers)
	return
}
