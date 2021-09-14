package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// index basic request to /
func index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	apirouter.ReturnResponse(w, req, http.StatusOK, map[string]interface{}{"message": "Welcome to the Paymail Server ✌(◕‿-)✌"})
}

// health is a basic request to return a health response
func health(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

// notFound handles all 404 requests
func notFound(w http.ResponseWriter, req *http.Request) {
	ErrorResponse(w, req, "request-404", "request not found", http.StatusNotFound)
}

// methodNotAllowed handles all 405 requests
func methodNotAllowed(w http.ResponseWriter, req *http.Request) {
	ErrorResponse(w, req, "method-405", "method "+req.Method+" not allowed", http.StatusMethodNotAllowed)
}
