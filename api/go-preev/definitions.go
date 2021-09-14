package preev

// PairList is a list of active pairs
type PairList struct {
	BsvUsd *Pair `json:"BSV:USD"`
	// todo: add more active pairs (only 1 pair is active)
}

// Pair is the info about the current pair
type Pair struct {
	Base    string   `json:"base"`
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Quote   string   `json:"quote"`
	Sources *Sources `json:"sources"`
	Status  *Status  `json:"status"`
}

// Sources are the exchange sources of the rates
type Sources struct {
	Bitfinex string `json:"Bitfinex,omitempty"`
	Bittrex  string `json:"Bittrex,omitempty"`
	OkCoin   string `json:"OkCoin,omitempty"`
	Poloniex string `json:"Poloniex,omitempty"`
	// todo: how to automatically detect other sources?
}

// Status is the current status of the pair
type Status struct {
	Active          bool    `json:"active"`
	Balance         float64 `json:"balance"`
	DaysRemaining   float64 `json:"days_remaining"`
	LastFunded      int64   `json:"last_funded"`
	TotalBroadcasts int64   `json:"total_broadcasts"`
	TxSize          int     `json:"tx_size"`
}

// TickerList is a list of active tickers
type TickerList struct {
	BsvUsd *Ticker `json:"BSV:USD"`
	// todo: add more active tickers (only 1 pair/ticker is active)
}

// Ticker is information for a given ticker
type Ticker struct {
	ID        string       `json:"id,omitempty"` // Pair ID
	Timestamp int64        `json:"t"`            // UNIX timestamp
	Tx        *Transaction `json:"tx"`           // Last broadcasted transaction
	Prices    *PriceSource `json:"p"`            // Price(s)
}

// Transaction is the last broadcasted transaction
type Transaction struct {
	Hash      string `json:"h"` // Transaction hash
	Timestamp int64  `json:"t"` // UNIX timestamp of transaction
}

// PriceSource is the sources of prices
type PriceSource struct {
	Binance  *Price `json:"binance,omitempty"`
	Bitfinex *Price `json:"bitfinex,omitempty"`
	Bitstamp *Price `json:"bitstamp,omitempty"`
	Bittrex  *Price `json:"bittrex,omitempty"`
	Coinbase *Price `json:"coinbase,omitempty"`
	Gemini   *Price `json:"gemini,omitempty"`
	Kraken   *Price `json:"kraken,omitempty"`
	Okcoin   *Price `json:"okcoin,omitempty"`
	Poloniex *Price `json:"poloniex,omitempty"`
	Ppi      *Price `json:"ppi,omitempty"`
}

// Price is the price and volume for a given PriceSource
type Price struct {
	LastPrice float64 `json:"l"` // Last price
	Volume    int64   `json:"v"` // 24hr volume
}

// APIError is an error message returned from the API
type APIError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
