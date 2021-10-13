package main

import (
	"net/http"
	"os"

	"github.com/dotwallet/dotwallet-go-sdk"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// getAuthorizeURL will return a new url for starting the user oauth2 authorization process
//
// For more information: https://developers.dotwallet.com/documents/en/#user-authorization
func getAuthorizeURL(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Create the DotWallet client
	c, err := dotwallet.NewClient(
		dotwallet.WithCredentials(
			os.Getenv("DOT_WALLET_CLIENT_ID"),
			os.Getenv("DOT_WALLET_CLIENT_SECRET"),
		),
		dotwallet.WithRedirectURI(os.Getenv("DOT_WALLET_REDIRECT_URI")),
	)
	if err != nil {
		ErrorResponse(w, req, "auth-400", err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a state (this needs to be stored in memory or cache for verification later)
	state, _ := c.NewState()
	globalStates[state] = true // This should be replaced with Redis or Caching

	// Return the URL to the user
	apirouter.ReturnResponse(w, req, http.StatusOK, c.GetAuthorizeURL(state, []string{dotwallet.ScopeUserInfo}))
}

// loginPage will return an example HTML page for using the User Authentication Example
func loginPage(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	err := redirectTemplate.Execute(w, req)
	if err != nil {
		ErrorResponse(w, req, "template-error", err.Error(), http.StatusBadRequest)
	}
}

// userLogin will attempt to log in a user after then came from the authorization_url with a code
func userLogin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Get the parameters
	params := apirouter.GetParams(req)
	code := params.GetString("code")
	state := params.GetString("state")

	// Check for required parameters
	if len(code) == 0 {
		ErrorResponse(w, req, "auth-code-1", "missing auth code", http.StatusBadRequest)
		return
	}
	if len(state) == 0 {
		ErrorResponse(w, req, "auth-state-1", "missing auth state", http.StatusBadRequest)
		return
	}

	// Validate the state value given
	_, ok := globalStates[state]
	if !ok {
		ErrorResponse(w, req, "auth-state-2", "state is not valid", http.StatusUnprocessableEntity)
		return
	}

	// Create the DotWallet client
	c, err := dotwallet.NewClient(
		dotwallet.WithCredentials(
			os.Getenv("DOT_WALLET_CLIENT_ID"),
			os.Getenv("DOT_WALLET_CLIENT_SECRET"),
		),
		dotwallet.WithRedirectURI(os.Getenv("DOT_WALLET_REDIRECT_URI")),
	)
	if err != nil {
		ErrorResponse(w, req, "auth-400", err.Error(), http.StatusBadRequest)
		return
	}

	// Get the user token
	var userToken *dotwallet.DotAccessToken
	userToken, err = c.GetUserToken(code)
	if err != nil {
		ErrorResponse(w, req, "auth-login-1", err.Error(), http.StatusUnauthorized)
		return
	}

	// Get the user info
	var user *dotwallet.User
	user, err = c.UserInfo(userToken)
	if err != nil {
		ErrorResponse(w, req, "get-user-info", err.Error(), http.StatusExpectationFailed)
		return
	}

	// Show user info!?
	apirouter.ReturnResponse(w, req, http.StatusOK, user.ID)
}
