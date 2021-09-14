package bsvrates

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gojektech/heimdall/v6"
	"github.com/gojektech/heimdall/v6/httpclient"
	"github.com/shopspring/decimal"
)

// todo: Use the official CoinPaprika package when:
//		1) They support "price-converter"
//		2) They swap HTTPClient for an interface

// coinPaprikaBaseURL is the main url for the service
const coinPaprikaBaseURL = "https://api.coinpaprika.com/v1/"

// usd is a const for the dollar
const usd = "usd"

// List of accepted known currencies (works for CoinPaprika only)
const (
	AUDCurrencyID = "aud-australian-dollar"
	BRLCurrencyID = "brl-brazil-real"
	CADCurrencyID = "cad-canadian-dollar"
	CHFCurrencyID = "chf-swiss-franc"
	CNYCurrencyID = "cny-yuan-renminbi"
	EURCurrencyID = "eur-euro"
	GBPCurrencyID = "gbp-pound-sterling"
	JPYCurrencyID = "jpy-japanese-yen"
	KRWCurrencyID = "krw-south-korea-won"
	MXNCurrencyID = "mxn-mexican-peso"
	NEWCurrencyID = "new-zealand-dollar"
	NOKCurrencyID = "nok-norwegian-krone"
	PLNCurrencyID = "pln-polish-zloty"
	RUBCurrencyID = "rub-russian-ruble"
	SEKCurrencyID = "sek-swedish-krona"
	TRYCurrencyID = "try-turkish-lira"
	TWDCurrencyID = "twd-taiwan-new-dollar"
	USDCurrencyID = usd + "-us-dollars"
	ZARCurrencyID = "zar-south-african-rand"
)

// PriceConversionResponse is the result returned from Coin Paprika conversion request
type PriceConversionResponse struct {
	Amount                float64      `json:"amount"`
	BaseCurrencyID        string       `json:"base_currency_id"`
	BaseCurrencyName      string       `json:"base_currency_name"`
	BasePriceLastUpdated  string       `json:"base_price_last_updated"`
	LastRequest           *lastRequest `json:"-"` // is the raw information from the last request
	Price                 float64      `json:"price"`
	QuoteCurrencyID       string       `json:"quote_currency_id"`
	QuoteCurrencyName     string       `json:"quote_currency_name"`
	QuotePriceLastUpdated string       `json:"quote_price_last_updated"`
}

// TickerResponse is the result returned from Coin Paprika ticker request
type TickerResponse struct {
	BetaValue         float64      `json:"beta_value"`
	CirculatingSupply int64        `json:"circulating_supply"`
	ID                string       `json:"id"`
	LastRequest       *lastRequest `json:"-"` // is the raw information from the last request
	LastUpdated       string       `json:"last_updated"`
	MaxSupply         int64        `json:"max_supply"`
	Name              string       `json:"name"`
	Quotes            *currency    `json:"quotes"`
	Rank              int          `json:"rank"`
	Symbol            string       `json:"symbol"`
	TotalSupply       int64        `json:"total_supply"`
}

// currency is the parent struct of a quote for a Ticker request
type currency struct {
	USD *quote `json:"USD"`
}

// quote is the JSON for each currency that is returned from the Ticker request
type quote struct {
	Price              float64 `json:"price"`
	Volume24h          float64 `json:"volume_24h"`
	Volume24hChange24h float64 `json:"volume_24h_change_24h"`
	MarketCap          int64   `json:"market_cap"`
	AthPrice           float64 `json:"ath_price"`
	AthDate            string  `json:"ath_date"`
}

// Accepted currencies (countries / fiat)
var acceptedCurrenciesCoinPaprika = []string{
	"aud",
	"brl",
	"cad",
	"chf",
	"cny",
	"eur",
	"gbp",
	"jpy",
	"krw",
	"mxn",
	"new",
	"nok",
	"pln",
	"rub",
	"sek",
	"try",
	"twd",
	usd,
	"zar",
}

// GetSatoshi will convert the price into Satoshi's (integer value)
func (p PriceConversionResponse) GetSatoshi() (satoshi int64, err error) {

	switch {
	case math.IsNaN(p.Price):
		fallthrough
	case math.IsInf(p.Price, 1):
		fallthrough
	case math.IsInf(p.Price, -1):
		return 0, fmt.Errorf("invalid price conversion")
	}

	satoshiDecimal := decimal.NewFromFloat(p.Price).Mul(decimal.NewFromInt(1e8))

	// there are no sub-Satoshis so get the ceiling so that you never underpay
	satoshi = satoshiDecimal.Ceil().IntPart()
	return
}

// PaprikaClient is the client for Coin Paprika
type PaprikaClient struct {
	HTTPClient httpInterface // carries out the http operations (heimdall client)
	UserAgent  string
}

// lastRequest is used to track what was submitted via the Coin Paprika Request
type lastRequest struct {
	Method     string `json:"method"`      // method is the HTTP method used
	PostData   string `json:"post_data"`   // postData is the post data submitted if POST/PUT request
	StatusCode int    `json:"status_code"` // statusCode is the last code from the request
	URL        string `json:"url"`         // url is the url used for the request
}

// createPaprikaClient will make a new http client based on the options provided
func createPaprikaClient(options *ClientOptions, customHTTPClient *http.Client) (c *PaprikaClient) {

	// Create a client
	c = new(PaprikaClient)

	// Is there a custom HTTP client to use?
	if customHTTPClient != nil {
		c.HTTPClient = customHTTPClient
		return
	}

	// Set options (either default or user modified)
	if options == nil {
		options = DefaultClientOptions()
	}

	// Set the user agent
	c.UserAgent = options.UserAgent

	// dial is the net dialer for clientDefaultTransport
	dial := &net.Dialer{KeepAlive: options.DialerKeepAlive, Timeout: options.DialerTimeout}

	// clientDefaultTransport is the default transport struct for the HTTP client
	clientDefaultTransport := &http.Transport{
		DialContext:           dial.DialContext,
		ExpectContinueTimeout: options.TransportExpectContinueTimeout,
		IdleConnTimeout:       options.TransportIdleTimeout,
		MaxIdleConns:          options.TransportMaxIdleConnections,
		Proxy:                 http.ProxyFromEnvironment,
		TLSHandshakeTimeout:   options.TransportTLSHandshakeTimeout,
	}

	// Determine the strategy for the http client (no retry enabled)
	if options.RequestRetryCount <= 0 {
		c.HTTPClient = httpclient.NewClient(
			httpclient.WithHTTPTimeout(options.RequestTimeout),
			httpclient.WithHTTPClient(&http.Client{
				Transport: clientDefaultTransport,
				Timeout:   options.RequestTimeout,
			}),
		)
	} else { // Retry enabled
		// Create exponential back-off
		backOff := heimdall.NewExponentialBackoff(
			options.BackOffInitialTimeout,
			options.BackOffMaxTimeout,
			options.BackOffExponentFactor,
			options.BackOffMaximumJitterInterval,
		)

		c.HTTPClient = httpclient.NewClient(
			httpclient.WithHTTPTimeout(options.RequestTimeout),
			httpclient.WithRetrier(heimdall.NewRetrier(backOff)),
			httpclient.WithRetryCount(options.RequestRetryCount),
			httpclient.WithHTTPClient(&http.Client{
				Transport: clientDefaultTransport,
				Timeout:   options.RequestTimeout,
			}),
		)
	}

	return
}

// GetBaseAmountAndCurrencyID will return an ID and default amount
func (p *PaprikaClient) GetBaseAmountAndCurrencyID(currency string, amount float64) (string, float64) {

	// Most are two decimal places
	if amount <= 0 {
		amount = 0.01 // Default
	}

	// Switch on known currencies
	switch currency {
	case "aud":
		return AUDCurrencyID, amount
	case "brl":
		return BRLCurrencyID, amount
	case "cad":
		return CADCurrencyID, amount
	case usd:
		return USDCurrencyID, amount
	case "eur":
		return EURCurrencyID, amount
	case "jpy":
		if amount < 1 {
			amount = 1
		}
		return JPYCurrencyID, amount
	case "mxn":
		return MXNCurrencyID, amount
	case "new":
		return NEWCurrencyID, amount
	case "nok":
		return NOKCurrencyID, amount
	case "pln":
		return PLNCurrencyID, amount
	case "gbp":
		return GBPCurrencyID, amount
	case "rub":
		return RUBCurrencyID, amount
	case "zar":
		return ZARCurrencyID, amount
	case "krw":
		if amount < 1 {
			amount = 1
		}
		return KRWCurrencyID, amount
	case "sek":
		return SEKCurrencyID, amount
	case "chf":
		return CHFCurrencyID, amount
	case "twd":
		return TWDCurrencyID, amount
	case "try":
		return TRYCurrencyID, amount
	case "cny":
		return CNYCurrencyID, amount
	default:
		return "", 0.00
	}
}

// GetPriceConversion returns a response of the conversion price from Coin Paprika
//
// See: https://api.coinpaprika.com/#tag/Tools/paths/~1price-converter/get
func (p *PaprikaClient) GetPriceConversion(ctx context.Context, baseCurrencyID, quoteCurrencyID string,
	amount float64) (response *PriceConversionResponse, err error) {

	// Set the api url
	// price-converter?base_currency_id=usd-us-dollars&quote_currency_id=bsv-bitcoin-sv&amount=0.01
	reqURL := fmt.Sprintf(
		"%sprice-converter?base_currency_id=%s&quote_currency_id=%s&amount=%f",
		coinPaprikaBaseURL, baseCurrencyID, quoteCurrencyID, amount,
	)

	// Start the request
	var req *http.Request
	if req, err = http.NewRequestWithContext(
		ctx, http.MethodGet, reqURL, nil,
	); err != nil {
		return
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")

	// Change the header (user agent is in case they block default Go user agents)
	req.Header.Set("User-Agent", p.UserAgent)

	// Start the response
	response = new(PriceConversionResponse)
	response.LastRequest = new(lastRequest)
	response.LastRequest.Method = http.MethodGet
	response.LastRequest.URL = reqURL

	// Fire the request
	var resp *http.Response
	if resp, err = p.HTTPClient.Do(req); err != nil {
		if resp != nil {
			response.LastRequest.StatusCode = resp.StatusCode
		}
		return
	}

	// Close the body
	defer func() {
		_ = resp.Body.Close()
	}()

	// Set the status
	response.LastRequest.StatusCode = resp.StatusCode

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response from provider: %d", resp.StatusCode)
		return
	}

	// Try and decode the response
	err = json.NewDecoder(resp.Body).Decode(&response)
	return
}

// GetMarketPrice returns a response of the market price from Coin Paprika
//
// See: https://api.coinpaprika.com/#operation/getTickersById
func (p *PaprikaClient) GetMarketPrice(ctx context.Context, coinID string) (response *TickerResponse, err error) {

	// Set the api url
	// tickers/:coin_id
	reqURL := fmt.Sprintf("%stickers/%s", coinPaprikaBaseURL, coinID)

	// Start the request
	var req *http.Request
	if req, err = http.NewRequestWithContext(
		ctx, http.MethodGet, reqURL, nil,
	); err != nil {
		return
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")

	// Change the header (user agent is in case they block default Go user agents)
	req.Header.Set("User-Agent", p.UserAgent)

	// Start the response
	response = new(TickerResponse)
	response.LastRequest = new(lastRequest)
	response.LastRequest.Method = http.MethodGet
	response.LastRequest.URL = reqURL

	// Fire the request
	var resp *http.Response
	if resp, err = p.HTTPClient.Do(req); err != nil {
		if resp != nil {
			response.LastRequest.StatusCode = resp.StatusCode
		}
		return
	}

	// Close the body
	defer func() {
		_ = resp.Body.Close()
	}()

	// Set the status
	response.LastRequest.StatusCode = resp.StatusCode

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response from provider: %d", resp.StatusCode)
		return
	}

	// Try and decode the response
	err = json.NewDecoder(resp.Body).Decode(&response)
	return
}

// tickerQuote is the quote value to return the price
type tickerQuote string

// Quote types
const (
	TickerQuoteUSD tickerQuote = "usd"
	TickerQuoteBTC tickerQuote = "btc"
)

// String is the string version of tickerQuote
func (t tickerQuote) String() string {
	return string(t)
}

// tickerInterval is the interval of time
type tickerInterval string

// Interval types
const (
	TickerInterval5m   tickerInterval = "5m"
	TickerInterval10m  tickerInterval = "10m"
	TickerInterval15m  tickerInterval = "15m"
	TickerInterval30m  tickerInterval = "30m"
	TickerInterval45m  tickerInterval = "45m"
	TickerInterval1h   tickerInterval = "1h"
	TickerInterval2h   tickerInterval = "2h"
	TickerInterval3h   tickerInterval = "3h"
	TickerInterval6h   tickerInterval = "6h"
	TickerInterval12h  tickerInterval = "12h"
	TickerInterval24h  tickerInterval = "24h"
	TickerInterval1d   tickerInterval = "1d"
	TickerInterval7d   tickerInterval = "7d"
	TickerInterval14d  tickerInterval = "14d"
	TickerInterval30d  tickerInterval = "30d"
	TickerInterval90d  tickerInterval = "90d"
	TickerInterval365d tickerInterval = "365d"
)

// String is the string version of tickerInterval
func (t tickerInterval) String() string {
	return string(t)
}

// This is the max amount of results that can be returned
const (
	maxHistoricalLimit = 5000
)

// HistoricalResponse is the response returned from the request
type HistoricalResponse struct {
	LastRequest *lastRequest      `json:"-"`
	Results     HistoricalResults `json:"-"`
}

// HistoricalResults is the results returned by the historical ticker request
type HistoricalResults []*HistoricalTicker

// HistoricalTicker is the ticker struct for historical request
type HistoricalTicker struct {
	MarketCap int64   `json:"market_cap"`
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
	Volume24h int64   `json:"volume_24h"`
}

// GetHistoricalTickers will return the historical tickers given the range of time
//
// See: https://api.coinpaprika.com/#tag/Tickers/paths/~1tickers~1{coin_id}~1historical/get
func (p *PaprikaClient) GetHistoricalTickers(ctx context.Context, coinID string, start, end time.Time, limit int,
	quote tickerQuote, interval tickerInterval) (response *HistoricalResponse, err error) {

	// Validate "start" time
	if start.IsZero() {
		err = fmt.Errorf("start time cannot be zero")
		return
	}

	// Send end if zero
	if end.IsZero() {
		end = time.Now().UTC()
	}

	// Test if start is after end
	if start.After(end) || start == end {
		err = fmt.Errorf("start time must be before end time")
		return
	}

	// Check for "max" limit (set default if not set)
	if limit > maxHistoricalLimit {
		limit = maxHistoricalLimit
	} else if limit <= 0 {
		limit = maxHistoricalLimit
	}

	// Set the api url
	// tickers/:coin_id/historical?start=
	reqURL := fmt.Sprintf(
		"%stickers/%s/historical?start=%d&end=%d&limit=%d&quote=%s&interval=%s",
		coinPaprikaBaseURL,
		coinID,
		start.Unix(),
		end.Unix(),
		limit,
		quote,
		interval,
	)

	// Start the request
	var req *http.Request
	if req, err = http.NewRequestWithContext(
		ctx, http.MethodGet, reqURL, nil,
	); err != nil {
		return
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")

	// Change the header (user agent is in case they block default Go user agents)
	req.Header.Set("User-Agent", p.UserAgent)

	// Start the response
	response = new(HistoricalResponse)
	response.LastRequest = new(lastRequest)
	response.LastRequest.Method = http.MethodGet
	response.LastRequest.URL = reqURL

	// Fire the request
	var resp *http.Response
	if resp, err = p.HTTPClient.Do(req); err != nil {
		if resp != nil {
			response.LastRequest.StatusCode = resp.StatusCode
		}
		return
	}

	// Close the body
	defer func() {
		_ = resp.Body.Close()
	}()

	// Set the status
	response.LastRequest.StatusCode = resp.StatusCode

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response from provider: %d", resp.StatusCode)
		return
	}

	// Try and decode the response
	err = json.NewDecoder(resp.Body).Decode(&response.Results)
	return
}

// IsAcceptedCurrency checks if the currency is accepted or not
func (p *PaprikaClient) IsAcceptedCurrency(currency string) bool {
	for _, val := range acceptedCurrenciesCoinPaprika {
		if strings.EqualFold(currency, val) {
			return true
		}
	}
	return false
}
