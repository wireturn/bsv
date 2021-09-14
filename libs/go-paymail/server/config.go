package server

import "github.com/tonicpow/go-paymail"

// Basic configuration for the server
const (
	paymailAPIVersion       = "v1"                                                                                          // Version of API
	paymailDomain           = "test.com"                                                                                    // This is the primary domain for the paymail service
	senderValidationEnabled = false                                                                                         // Turn on if all address resolution requests need a valid signature
	serviceURL              = "https://" + paymailDomain + "/" + paymailAPIVersion + "/" + paymail.DefaultServiceName + "/" // This is appended to all URLs
)
