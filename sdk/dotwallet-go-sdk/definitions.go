package dotwallet

import (
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// Package configuration defaults
	apiVersion                   string = "v1"                           // Version of the API
	defaultHost                         = "https://api.ddpurse.com"      // Default host for API endpoints
	defaultHTTPTimeout                  = 10 * time.Second               // Default timeout for all GET requests in seconds
	defaultRefreshTokenExpiresIn        = 7 * 24 * time.Hour             // Default is 7 days (from documentation)
	defaultRetryCount            int    = 2                              // Default retry count for HTTP requests
	defaultUserAgent                    = "dotwallet-go-sdk: " + version // Default user agent
	version                      string = "v0.0.3"                       // dotwallet-go-sdk version

	// Grants
	grantTypeAuthorizationCode = "authorization_code"
	grantTypeClientCredentials = "client_credentials"
	grantTypeRefreshToken      = "refresh_token"

	// Endpoints
	authorizeURI          = "/" + apiVersion + "/oauth2/authorize"
	getAccessTokenURI     = "/" + apiVersion + "/oauth2/get_access_token"
	getUserInfo           = "/" + apiVersion + "/user/get_user_info"
	getUserReceiveAddress = "/" + apiVersion + "/user/get_user_receive_address"

	// Headers
	headerAuthorization = "Authorization"

	// ScopeUserInfo is for getting user info
	ScopeUserInfo = "user.info"

	// ScopeAutoPayBSV is for auto-pay with a BSV balance
	ScopeAutoPayBSV = "autopay.bsv"

	// ScopeAutoPayBTC is for auto-pay with a BTC balance
	ScopeAutoPayBTC = "autopay.btc"

	// ScopeAutoPayETH is for auto-pay with a ETH balance
	ScopeAutoPayETH = "autopay.eth"
)

// These are used for the accepted coin types in regard to wallet functions
const (
	CoinTypeBSV coinType = "BSV" // BitcoinSV
	CoinTypeBTC coinType = "BTC" // BitcoinCore
	CoinTypeETH coinType = "ETH" // Ethereum
)

// coinType is used for determining the coin_type for wallet functions
type coinType string

// String is the string version of coin_type
func (c coinType) String() string {
	return string(c)
}

// StandardResponse is the standard fields returned on all responses from Request()
type StandardResponse struct {
	Body       []byte          `json:"-"` // Body of the response request
	Error      *Error          `json:"-"` // API error response
	StatusCode int             `json:"-"` // Status code returned on the request
	Tracing    resty.TraceInfo `json:"-"` // Trace information if enabled on the request
}

// genericResponse is the generic part of the response body
type genericResponse struct {
	Code    int    `json:"code"` // If there is an error, this will be a value
	Message string `json:"msg"`  // If there is an error, this will be the error message
}

/*
Example
{
    "code": 74012,
    "msg": "Client authentication failed",
    "data": null,
    "req_id": "zacd8b1b05b12a36d45fvfc20a4b97c5"
}
*/

// Error is the universal Error response from the API
//
// For more information: https://developers.dotwallet.com/documents/en/#errors
type Error struct {
	genericResponse
	Data      interface{} `json:"data"`
	Method    string      `json:"method"`
	RequestID string      `json:"req_id"`
	URL       string      `json:"url"`
}

// getAccessTokenRequest is used for the access token request
//
// For more information: https://developers.dotwallet.com/documents/en/#application-authorization
type getAccessTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

// getDotUserTokenRequest is used for the GetDotUserToken request
//
// For more information: https://developers.dotwallet.com/documents/en/#user-authorization
type getDotUserTokenRequest struct {
	getAccessTokenRequest
	Code        string `json:"code"`         // User code given from the oauth2 handshake
	RedirectURI string `json:"redirect_uri"` // The redirect URI set for the application
}

// refreshDotUserTokenRequest is used for the RefreshDotUserToken request
//
// For more information: https://developers.dotwallet.com/documents/en/#user-authorization
type refreshDotUserTokenRequest struct {
	getAccessTokenRequest
	RefreshToken string `json:"refresh_token"` // Refresh token which was given upon first auth_token generation
}

// DotAccessToken is the access token information
//
// For more information: https://developers.dotwallet.com/documents/en/#user-authorization
type DotAccessToken struct {
	AccessToken           string   `json:"access_token"`                       // Access token from the API
	ExpiresAt             int64    `json:"expires_at,omitempty"`               // Friendly unix time from UTC when the access_token expires
	ExpiresIn             int64    `json:"expires_in"`                         // Seconds from now that the token expires
	RefreshToken          string   `json:"refresh_token,omitempty"`            // Refresh token for the user
	RefreshTokenExpiresIn int64    `json:"refresh_token_expires_in,omitempty"` // Seconds from now that the token expires
	RefreshTokenExpiresAt int64    `json:"refresh_token_expires_at,omitempty"` // Friendly unix time from UTC when the refresh_token expires
	Scopes                []string `json:"scopes,omitempty"`                   // Scopes for the user token
	TokenType             string   `json:"token_type"`                         // Token type
}

// accessTokenResponse is the response from creating the new access token
//
// For more information: https://developers.dotwallet.com/documents/en/#user-authorization
type accessTokenResponse struct {
	genericResponse
	Data struct {
		AccessToken  string `json:"access_token"`            // Access token from the API
		ExpiresIn    int64  `json:"expires_in"`              // Seconds from now that the token expires
		RefreshToken string `json:"refresh_token,omitempty"` // Refresh token for the user
		Scope        string `json:"scope,omitempty"`         // Scopes for the user token
		TokenType    string `json:"token_type"`              // Token type
	}
}

// userResponse is the response from the user info request
//
// For more information: https://developers.dotwallet.com/documents/en/#user-info
type userResponse struct {
	genericResponse
	Data struct {
		User
	}
}

// User is the DotWallet user profile information
//
// For more information: https://developers.dotwallet.com/documents/en/#user-info
type User struct {
	Avatar           string            `json:"avatar"`
	ID               string            `json:"id"`
	Nickname         string            `json:"nickname"`
	WebWalletAddress *webWalletAddress `json:"web_wallet_address"`
}

// webWalletAddress is the user's wallet addresses
type webWalletAddress struct {
	BSVRegular string `json:"bsv_regular"`
	BTCRegular string `json:"btc_regular"`
	ETHRegular string `json:"eth_regular"`
}

// userReceiveRequest is used for the user receive address request
//
// For more information: https://developers.dotwallet.com/documents/en/#get-user-receive-address
type userReceiveRequest struct {
	UserID   string   `json:"user_id"`
	CoinType coinType `json:"coin_type"`
}

// userReceiveAddressResponse is the response from the user receive address request
//
// For more information: https://developers.dotwallet.com/documents/en/#get-user-receive-address
type userReceiveAddressResponse struct {
	genericResponse
	Data struct {
		Wallets
	}
}

// Wallets is the user's wallet information
type Wallets struct {
	AutopayWallet *walletInfo `json:"autopay_wallet"`
	PrimaryWallet *walletInfo `json:"primary_wallet"`
}

// walletInfo is the user's wallet information
type walletInfo struct {
	Address     string `json:"address"`
	CoinType    string `json:"coin_type"`
	Paymail     string `json:"paymail,omitempty"`
	UserID      string `json:"user_id"`
	WalletIndex int64  `json:"wallet_index"`
	WalletType  string `json:"wallet_type"`
}
