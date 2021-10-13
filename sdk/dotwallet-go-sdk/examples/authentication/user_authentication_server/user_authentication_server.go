package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

var (
	redirectTemplate *template.Template
	globalStates     = make(map[string]bool) // This should be replaced with Redis or Caching
)

func main() {

	// Get the working directory
	wd, _ := os.Getwd()

	// Load the HTML example page
	var err error
	if redirectTemplate, err = template.ParseFiles(
		filepath.Join(wd, "examples", "authentication", "user_authentication_server", "login.html"),
	); err != nil {
		log.Fatalln(err)
	}

	// Run the server on port 3000 and timeout requests after 15 seconds
	Start(3000, 15)
}

// Start will run the Example Auth server
func Start(port, timeoutInSeconds int) {

	// Load the server
	log.Println("starting go example auth server...")
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),                      // Address to run the server on
		Handler:      Handlers(),                                    // Load all the routes
		ReadTimeout:  time.Duration(timeoutInSeconds) * time.Second, // Basic default timeout for read requests
		WriteTimeout: time.Duration(timeoutInSeconds) * time.Second, // Basic default timeout for write requests
	}
	log.Fatalln(srv.ListenAndServe())
}
