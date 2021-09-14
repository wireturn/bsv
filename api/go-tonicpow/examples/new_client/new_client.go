package main

import (
	"log"
	"os"

	"github.com/tonicpow/go-tonicpow"
)

func main() {

	//
	// Load the api client
	// You can also set the environment or client options
	//
	client, err := tonicpow.NewClient(
		tonicpow.WithAPIKey(os.Getenv("TONICPOW_API_KEY")),
		tonicpow.WithEnvironmentString(os.Getenv("TONICPOW_ENVIRONMENT")),

		// Custom options for loading the TonicPow client
		// tonicpow.WithCustomEnvironment("customEnv", "customAlias", "https://localhost:3002"),
		// tonicpow.WithEnvironment(tonicpow.EnvironmentStaging),
		// tonicpow.WithHTTPTimeout(10*time.Second),
		// tonicpow.WithRequestTracing(),
		// tonicpow.WithRetryCount(3),
		// tonicpow.WithUserAgent("my custom user agent v9.0.9"),

		/*
			// Example adding custom headers to each request
			headers := make(map[string][]string)
			headers["custom_header_1"] = append(headers["custom_header_1"], "value_1")
			headers["custom_header_2"] = append(headers["custom_header_2"], "value_2")
			tonicpow.WithCustomHeaders(headers),
		*/
	)
	if err != nil {
		log.Fatalf("error in NewClient: %s", err.Error())
	}

	// Use your own custom Resty
	// client.WithCustomHTTPClient(resty.New())

	log.Println(
		"client: ", client.GetUserAgent(),
		"environment: ", client.GetEnvironment().Name(),
	)
}
