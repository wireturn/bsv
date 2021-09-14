package http

import "net/http"

// HTTPClient interfaces the Do(*http.Request) function to allow for easy mocking.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}
