package server

import (
	"net/http"

	apirouter "github.com/mrz1836/go-api-router"
	"github.com/newrelic/go-agent/v3/integrations/nrhttprouter"
	"github.com/tonicpow/go-paymail"
)

// Handlers are used to isolate loading the routes (used for testing)
func Handlers() *nrhttprouter.Router {

	// Create a new router
	r := apirouter.New()

	// Turned off all CORs - should be accessed outside a browser
	r.CrossOriginEnabled = false
	r.CrossOriginAllowCredentials = false
	r.CrossOriginAllowOriginAll = false

	// Register basic server routes
	registerBasicRoutes(r)

	// Register paymail routes
	registerPaymailRoutes(r)

	// Return the router
	return r.HTTPRouter
}

// registerBasicRoutes will register basic server related routes
func registerBasicRoutes(router *apirouter.Router) {

	// Set the main index page (navigating to slash)
	router.HTTPRouter.GET("/", router.Request(index))
	// router.HTTPRouter.OPTIONS("/", router.SetCrossOriginHeaders) // Disabled for security

	// Set the health request (used for load balancers)
	router.HTTPRouter.GET("/health", router.RequestNoLogging(health))
	router.HTTPRouter.OPTIONS("/health", router.SetCrossOriginHeaders)
	router.HTTPRouter.HEAD("/health", router.SetCrossOriginHeaders)

	// Set the 404 handler (any request not detected)
	router.HTTPRouter.NotFound = http.HandlerFunc(notFound)

	// Set the method not allowed
	router.HTTPRouter.MethodNotAllowed = http.HandlerFunc(methodNotAllowed)
}

// registerPaymailRoutes will register all paymail related routes
func registerPaymailRoutes(router *apirouter.Router) {

	// Capabilities (service discovery)
	router.HTTPRouter.GET(
		"/.well-known/"+paymail.DefaultServiceName,
		router.Request(showCapabilities),
	)

	// PKI request (public key information)
	router.HTTPRouter.GET(
		"/"+paymailAPIVersion+"/"+paymail.DefaultServiceName+"/id/:paymailAddress",
		router.Request(showPKI),
	)

	// Verify PubKey request (public key verification to paymail address)
	router.HTTPRouter.GET(
		"/"+paymailAPIVersion+"/"+paymail.DefaultServiceName+"/verify-pubkey/:paymailAddress/:pubKey",
		router.Request(verifyPubKey),
	)

	// Payment Destination request (address resolution)
	router.HTTPRouter.POST(
		"/"+paymailAPIVersion+"/"+paymail.DefaultServiceName+"/address/:paymailAddress",
		router.Request(resolveAddress),
	)

	// Public Profile request (returns Name & Avatar)
	router.HTTPRouter.GET(
		"/"+paymailAPIVersion+"/"+paymail.DefaultServiceName+"/public-profile/:paymailAddress",
		router.Request(publicProfile),
	)

	// P2P Destination request (returns output & reference)
	router.HTTPRouter.POST(
		"/"+paymailAPIVersion+"/"+paymail.DefaultServiceName+"/p2p-payment-destination/:paymailAddress",
		router.Request(p2pDestination),
	)

	// P2P Receive Tx request (receives the P2P transaction, broadcasts, returns tx_id)
	router.HTTPRouter.POST(
		"/"+paymailAPIVersion+"/"+paymail.DefaultServiceName+"/receive-transaction/:paymailAddress",
		router.Request(p2pReceiveTx),
	)
}
