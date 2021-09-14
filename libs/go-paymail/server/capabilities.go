package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/tonicpow/go-paymail"
)

// Capabilities is the standard response for returning the Paymail capabilities
// FYI - could not use paymail.Capabilities because it uses map[string]interface{}
type Capabilities struct {
	BsvAlias     string              `json:"bsvalias"`     // Version of the bsvalias
	Capabilities *activeCapabilities `json:"capabilities"` // List of the capabilities
}

// activeCapabilities is used to display only the active capabilities of the Paymail server
type activeCapabilities struct {
	ForceSenderValidation bool   `json:"6745385c3fc0"`       // Will force sender to have a signature if enabled
	PaymentDestination    string `json:"paymentDestination"` // Resolve an address aka Payment Destination - Alternate: 759684b1a19a
	PKI                   string `json:"pki"`                // Get public key information - Alternate: 0c4339ef99c2
	PublicProfile         string `json:"f12f968c92d6"`       // Returns a public profile
	VerifyPublicKey       string `json:"a9f510c16bde"`       // Verify a given pubkey
}

// createCapabilities will create a default set of capabilities
func createCapabilities() *Capabilities {
	return &Capabilities{
		BsvAlias: paymail.DefaultBsvAliasVersion,
		Capabilities: &activeCapabilities{
			ForceSenderValidation: senderValidationEnabled,
			PaymentDestination:    serviceURL + "address/{alias}@{domain.tld}",
			PKI:                   serviceURL + "id/{alias}@{domain.tld}",
			PublicProfile:         serviceURL + "public-profile/{alias}@{domain.tld}",
			VerifyPublicKey:       serviceURL + "verify-pubkey/{alias}@{domain.tld}/{pubkey}",
		},
	}
}

// showCapabilities will return the service discovery results for the server
// and list all active capabilities of the Paymail server
//
// Specs: http://bsvalias.org/02-02-capability-discovery.html
func showCapabilities(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	apirouter.ReturnResponse(w, req, http.StatusOK, createCapabilities())
}
