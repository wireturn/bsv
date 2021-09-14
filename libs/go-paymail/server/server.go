package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mrz1836/go-logger"
)

// Start will run the Paymail server
//
// This is just a basic example - all options should be set via config or ENV
func Start(port, timeoutInSeconds int) {

	// Load the server
	logger.Data(2, logger.DEBUG, "starting go paymail server...", logger.MakeParameter("port", port))
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),                      // Address to run the server on
		Handler:      Handlers(),                                    // Load all the routes
		ReadTimeout:  time.Duration(timeoutInSeconds) * time.Second, // Basic default timeout for read requests
		WriteTimeout: time.Duration(timeoutInSeconds) * time.Second, // Basic default timeout for write requests
	}
	logger.Fatalln(srv.ListenAndServe())
}
