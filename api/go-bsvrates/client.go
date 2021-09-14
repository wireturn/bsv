package bsvrates

import (
	"context"
	"net/http"
	"time"

	"github.com/mrz1836/go-preev"
	"github.com/mrz1836/go-whatsonchain"
)

// httpInterface is used for the http client (mocking heimdall)
type httpInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// preevInterface is an interface for the Preev Client
type preevInterface interface {
	GetPair(ctx context.Context, pairID string) (pair *preev.Pair, err error)
	GetPairs(ctx context.Context) (pairList *preev.PairList, err error)
	GetTicker(ctx context.Context, pairID string) (ticker *preev.Ticker, err error)
	GetTickers(ctx context.Context) (tickerList *preev.TickerList, err error)
}

// whatsOnChainInterface is an interface for the WOC Client
type whatsOnChainInterface interface {
	GetExchangeRate() (rate *whatsonchain.ExchangeRate, err error)
}

// coinPaprikaInterface is an interface for the Coin Paprika Client
type coinPaprikaInterface interface {
	GetBaseAmountAndCurrencyID(currency string, amount float64) (string, float64)
	GetMarketPrice(ctx context.Context, coinID string) (response *TickerResponse, err error)
	GetPriceConversion(ctx context.Context, baseCurrencyID, quoteCurrencyID string, amount float64) (response *PriceConversionResponse, err error)
	IsAcceptedCurrency(currency string) bool
	GetHistoricalTickers(ctx context.Context, coinID string, start, end time.Time, limit int, quote tickerQuote, interval tickerInterval) (response *HistoricalResponse, err error)
}

// Client is the parent struct that contains the provider clients and list of providers to use
type Client struct {
	CoinPaprika  coinPaprikaInterface  // Coin Paprika client
	Preev        preevInterface        // Preev Client
	Providers    []Provider            // List of providers to use (in order for fail-over)
	WhatsOnChain whatsOnChainInterface // WhatsOnChain client
}

// ClientOptions holds all the configuration for connection, dialer and transport
type ClientOptions struct {
	BackOffExponentFactor          float64       `json:"back_off_exponent_factor"`
	BackOffInitialTimeout          time.Duration `json:"back_off_initial_timeout"`
	BackOffMaximumJitterInterval   time.Duration `json:"back_off_maximum_jitter_interval"`
	BackOffMaxTimeout              time.Duration `json:"back_off_max_timeout"`
	DialerKeepAlive                time.Duration `json:"dialer_keep_alive"`
	DialerTimeout                  time.Duration `json:"dialer_timeout"`
	RequestRetryCount              int           `json:"request_retry_count"`
	RequestTimeout                 time.Duration `json:"request_timeout"`
	TransportExpectContinueTimeout time.Duration `json:"transport_expect_continue_timeout"`
	TransportIdleTimeout           time.Duration `json:"transport_idle_timeout"`
	TransportMaxIdleConnections    int           `json:"transport_max_idle_connections"`
	TransportTLSHandshakeTimeout   time.Duration `json:"transport_tls_handshake_timeout"`
	UserAgent                      string        `json:"user_agent"`
}

// ToPreevOptions will convert the current options to Preev Options
func (c *ClientOptions) ToPreevOptions() (options *preev.Options) {
	options = preev.ClientDefaultOptions()
	options.UserAgent = c.UserAgent + " using " + options.UserAgent
	options.BackOffExponentFactor = c.BackOffExponentFactor
	options.BackOffInitialTimeout = c.BackOffInitialTimeout
	options.BackOffMaximumJitterInterval = c.BackOffMaximumJitterInterval
	options.BackOffMaxTimeout = c.BackOffMaxTimeout
	options.DialerKeepAlive = c.DialerKeepAlive
	options.DialerTimeout = c.DialerTimeout
	options.RequestRetryCount = c.RequestRetryCount
	options.RequestTimeout = c.RequestTimeout
	options.TransportExpectContinueTimeout = c.TransportExpectContinueTimeout
	options.TransportIdleTimeout = c.TransportIdleTimeout
	options.TransportMaxIdleConnections = c.TransportMaxIdleConnections
	options.TransportTLSHandshakeTimeout = c.TransportTLSHandshakeTimeout
	return
}

// ToWhatsOnChainOptions will convert the current options to WOC Options
func (c *ClientOptions) ToWhatsOnChainOptions() (options *whatsonchain.Options) {
	options = whatsonchain.ClientDefaultOptions()
	options.UserAgent = c.UserAgent + " using " + options.UserAgent
	options.BackOffExponentFactor = c.BackOffExponentFactor
	options.BackOffInitialTimeout = c.BackOffInitialTimeout
	options.BackOffMaximumJitterInterval = c.BackOffMaximumJitterInterval
	options.BackOffMaxTimeout = c.BackOffMaxTimeout
	options.DialerKeepAlive = c.DialerKeepAlive
	options.DialerTimeout = c.DialerTimeout
	options.RequestRetryCount = c.RequestRetryCount
	options.RequestTimeout = c.RequestTimeout
	options.TransportExpectContinueTimeout = c.TransportExpectContinueTimeout
	options.TransportIdleTimeout = c.TransportIdleTimeout
	options.TransportMaxIdleConnections = c.TransportMaxIdleConnections
	options.TransportTLSHandshakeTimeout = c.TransportTLSHandshakeTimeout
	return
}

// DefaultClientOptions will return a clientOptions struct with the default settings.
// Useful for starting with the default and then modifying as needed
func DefaultClientOptions() (clientOptions *ClientOptions) {
	return &ClientOptions{
		BackOffExponentFactor:          2.0,
		BackOffInitialTimeout:          2 * time.Millisecond,
		BackOffMaximumJitterInterval:   2 * time.Millisecond,
		BackOffMaxTimeout:              10 * time.Millisecond,
		DialerKeepAlive:                20 * time.Second,
		DialerTimeout:                  5 * time.Second,
		RequestRetryCount:              2,
		RequestTimeout:                 10 * time.Second,
		TransportExpectContinueTimeout: 3 * time.Second,
		TransportIdleTimeout:           20 * time.Second,
		TransportMaxIdleConnections:    10,
		TransportTLSHandshakeTimeout:   5 * time.Second,
		UserAgent:                      defaultUserAgent,
	}
}

// NewClient creates a new client for requests
func NewClient(clientOptions *ClientOptions, customHTTPClient *http.Client, providers ...Provider) (client *Client) {
	client = new(Client)

	// No providers? (Use the default set for now)
	if len(providers) == 0 {
		client.Providers = defaultProviders
	} else {
		client.Providers = providers
	}

	// Set default options if none are provided
	if clientOptions == nil {
		clientOptions = DefaultClientOptions()
	}

	// Create a client for Coin Paprika
	client.CoinPaprika = createPaprikaClient(clientOptions, customHTTPClient)

	// Create a client for Preev
	client.Preev = preev.NewClient(clientOptions.ToPreevOptions(), customHTTPClient)

	// Create a client for WhatsOnChain
	client.WhatsOnChain = whatsonchain.NewClient(whatsonchain.NetworkMain, clientOptions.ToWhatsOnChainOptions(), customHTTPClient)

	return
}
