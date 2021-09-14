package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/tonicpow/go-paymail"
)

// verifyPubKey will return a response if the pubkey matches the paymail given
//
// Specs: https://bsvalias.org/05-verify-public-key-owner.html
func verifyPubKey(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Get the params submitted via URL request
	params := apirouter.GetParams(req)
	incomingPaymail := params.GetString("paymailAddress")
	incomingPubKey := params.GetString("pubKey")

	// Parse, sanitize and basic validation
	alias, domain, address := paymail.SanitizePaymail(incomingPaymail)
	if len(address) == 0 {
		ErrorResponse(w, req, ErrorInvalidParameter, "invalid paymail: "+incomingPaymail, http.StatusBadRequest)
		return
	} else if domain != paymailDomain {
		ErrorResponse(w, req, ErrorUnknownDomain, "domain unknown: "+domain, http.StatusBadRequest)
		return
	}

	// Basic validation on pubkey
	if len(incomingPubKey) != paymail.PubKeyLength {
		ErrorResponse(w, req, ErrorInvalidPubKey, "invalid pubkey: "+incomingPubKey, http.StatusBadRequest)
		return
	}

	// todo: lookup the paymail address in a data-store, database, etc - get the PubKey (return 404 if not found)

	// todo: add caching for fast responses since the pubkey will not change

	// Find in mock database
	foundPaymail := getPaymailByAlias(alias)
	if foundPaymail == nil {
		ErrorResponse(w, req, ErrorPaymailNotFound, "paymail not found", http.StatusNotFound)
		return
	}

	// Return the response
	apirouter.ReturnResponse(w, req, http.StatusOK, &paymail.Verification{
		BsvAlias: paymail.DefaultBsvAliasVersion,
		Handle:   address,
		PubKey:   incomingPubKey, // todo: should this be the incoming or found pubkey?
		Match:    foundPaymail.PubKey == incomingPubKey,
	})
}
