package server

import (
	"net/http"

	apirouter "github.com/mrz1836/go-api-router"
	"github.com/tonicpow/go-paymail"
)

// Error codes for server response errors
const (
	ErrorInvalidDt           = "invalid-dt"
	ErrorInvalidParameter    = "invalid-parameter"
	ErrorInvalidPubKey       = "invalid-pubkey"
	ErrorInvalidSenderHandle = "invalid-sender-handle"
	ErrorInvalidSignature    = "invalid-signature"
	ErrorMissingHex          = "missing-hex"
	ErrorMissingReference    = "missing-reference"
	ErrorMissingSatoshis     = "missing-satoshis"
	ErrorPaymailNotFound     = "not-found"
	ErrorScript              = "script-error"
	ErrorUnknownDomain       = "unknown-domain"
)

// ErrorResponse is a standard way to return errors to the client
//
// Specs: http://bsvalias.org/99-01-recommendations.html
func ErrorResponse(w http.ResponseWriter, req *http.Request, code, message string, statusCode int) {
	apirouter.ReturnResponse(w, req, statusCode, &paymail.ServerError{Code: code, Message: message})
}
