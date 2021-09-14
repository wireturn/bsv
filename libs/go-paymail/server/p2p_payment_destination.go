package server

import (
	"net/http"

	"github.com/bitcoinschema/go-bitcoin"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/tonicpow/go-paymail"
)

/*
Incoming Data Object Example:
{
  "satoshis": 1000100,
}
*/

// p2pDestination will return a output script(s) for a destination (used with SendP2PTransaction)
//
// Specs: https://docs.moneybutton.com/docs/paymail-07-p2p-payment-destination.html
func p2pDestination(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Get the params & paymail address submitted via URL request
	params := apirouter.GetParams(req)
	incomingPaymail := params.GetString("paymailAddress")

	// Start the PaymentRequest
	paymentRequest := &paymail.PaymentRequest{
		Satoshis: params.GetUint64("satoshis"),
	}

	// Parse, sanitize and basic validation
	alias, domain, paymailAddress := paymail.SanitizePaymail(incomingPaymail)
	if len(paymailAddress) == 0 {
		ErrorResponse(w, req, ErrorInvalidParameter, "invalid paymail: "+incomingPaymail, http.StatusBadRequest)
		return
	} else if domain != paymailDomain {
		ErrorResponse(w, req, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest)
		return
	}

	// Did we get some satoshis?
	if paymentRequest.Satoshis == 0 {
		ErrorResponse(w, req, ErrorMissingSatoshis, "missing parameter: satoshis", http.StatusBadRequest)
		return
	}

	// todo: lookup the paymail address in a data-store, database, etc - get the PubKey (return 404 if not found)

	// Find in mock database
	foundPaymail := getPaymailByAlias(alias)
	if foundPaymail == nil {
		ErrorResponse(w, req, ErrorPaymailNotFound, "paymail not found", http.StatusNotFound)
		return
	}

	// Start the script
	output := &paymail.PaymentOutput{
		Satoshis: paymentRequest.Satoshis,
	}

	// Generate the script
	// todo: multiple scripts if you want to break apart (IE: over X satoshis, break apart into multiple outputs)
	var err error
	if output.Script, err = bitcoin.ScriptFromAddress(foundPaymail.LastAddress); err != nil {
		ErrorResponse(w, req, ErrorScript, "error generating script: "+err.Error(), http.StatusNotFound)
		return
	}

	// todo: Generate a unique reference (stored in data-store once created)
	reference := "1234567890"

	// Return the response
	apirouter.ReturnResponse(w, req, http.StatusOK, &paymail.PaymentDestination{
		Outputs:   []*paymail.PaymentOutput{output},
		Reference: reference,
	})
}
