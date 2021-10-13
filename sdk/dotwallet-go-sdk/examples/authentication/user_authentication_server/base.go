package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// ServerError is the standard error response from a paymail server
type ServerError struct {
	Code    string `json:"code"`    // Shows the corresponding code
	Message string `json:"message"` // Shows the error message returned by the server
}

// ErrorResponse is a standard way to return errors to the client
func ErrorResponse(w http.ResponseWriter, req *http.Request, code, message string, statusCode int) {
	apirouter.ReturnResponse(w, req, statusCode, &ServerError{Code: code, Message: message})
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
