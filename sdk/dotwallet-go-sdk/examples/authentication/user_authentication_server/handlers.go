package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// Handlers are used to isolate loading the routes (used for testing)
func Handlers() *httprouter.Router {

	// Create a new router
	r := apirouter.New()

	// Turned off all CORs - should be accessed outside a browser
	r.CrossOriginEnabled = false
	r.CrossOriginAllowCredentials = false
	r.CrossOriginAllowOriginAll = false

	// Register basic server routes
	registerBasicRoutes(r)

	// Register auth routes
	registerAuthRoutes(r)

	// Return the router
	return r.HTTPRouter.Router
}

// registerBasicRoutes will register basic server related routes
func registerBasicRoutes(router *apirouter.Router) {

	// Set the health request (used for load balancers)
	router.HTTPRouter.GET("/health", router.RequestNoLogging(health))
	router.HTTPRouter.OPTIONS("/health", router.SetCrossOriginHeaders)
	router.HTTPRouter.HEAD("/health", router.SetCrossOriginHeaders)

	// Set the 404 handler (any request not detected)
	router.HTTPRouter.NotFound = http.HandlerFunc(notFound)

	// Set the method not allowed
	router.HTTPRouter.MethodNotAllowed = http.HandlerFunc(methodNotAllowed)
}

// registerAuthRoutes will register all example related routes
func registerAuthRoutes(router *apirouter.Router) {

	// Example login page (root page for the example)
	router.HTTPRouter.GET(
		"/",
		router.Request(loginPage),
	)

	// Get the authorization url
	router.HTTPRouter.GET(
		"/dot_wallet_auth",
		router.Request(getAuthorizeURL),
	)

	// Login as a user to DotWallet
	router.HTTPRouter.POST(
		"/dot_wallet_login",
		router.Request(userLogin),
	)
}
