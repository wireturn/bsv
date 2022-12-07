package paymail

import (
	"fmt"
	"strings"
	"time"

	"github.com/mrz1836/go-sanitize"
	"github.com/mrz1836/go-validate"
)

// SanitisedPaymail contains elements of a sanitized paymail address.
// All elements are lowercased.
type SanitisedPaymail struct{
	Alias, Domain, Address string
}

// ValidateAndSanitisePaymail will take a paymail address or handle,
// convert to a paymail address if it's a handle,
// validate that address, then sanitize it if it is valid.
// If the address or handle isn't valid, an error will be returned.
func ValidateAndSanitisePaymail(paymail string, isBeta bool) (*SanitisedPaymail, error){
	h := ConvertHandle(paymail, isBeta)
	if err := ValidatePaymail(h); err != nil{
		return nil, err
	}
	a,d,ad := SanitizePaymail(h)
	return &SanitisedPaymail{
		Alias:  a,
		Domain:  d,
		Address: ad,
	}, nil
}

// SanitizePaymail will take an input and return the sanitized version (alias@domain.tld)
//
// Alias is the first part of the address (alias @)
// Domain is the lowercase sanitized version (domain.tld)
// Address is the full sanitized paymail address (alias@domain.tld)
func SanitizePaymail(paymailAddress string) (alias, domain, address string) {

	// Sanitize the paymail address
	address = sanitize.Email(paymailAddress, false)

	// Split the email parts (alias @ domain)
	parts := strings.Split(address, "@")

	// Sanitize the domain name (force to lowercase, remove www.)
	if len(parts) > 1 {
		domain, _ = sanitize.Domain(parts[1], false, true)
	}

	// Set the alias (lowercase, no spaces)
	alias = strings.TrimSpace(strings.ToLower(parts[0]))

	// Paymail address does not meet the basic requirement of an email address
	// Since we don't return an error, we will return an empty result
	if len(alias) == 0 || len(domain) == 0 {
		address = ""
	}

	return
}

// ValidatePaymail will do a basic validation on the paymail format (email address format)
//
// This will not check to see if the paymail address is active via the provider
func ValidatePaymail(paymailAddress string) error {

	// Validate the format for the paymail address (paymail addresses follow conventional email requirements)
	if _, err := validate.IsValidEmail(paymailAddress, false); err != nil {
		return fmt.Errorf("paymail address failed format validation: %w", err)
	}

	return nil
}

// ValidateDomain will do a basic validation on the domain format
//
// This will not check to see if the domain is an active paymail provider
// This will not check DNS records to make sure the domain is active
func ValidateDomain(domain string) error {

	// Check for a real domain (require at least one period)
	if !strings.Contains(domain, ".") {
		return fmt.Errorf("domain name is invalid: %s", domain)
	} else if !validate.IsValidHost(domain) { // Basic DNS check (not a REAL domain name check)
		return fmt.Errorf("domain name failed host check: %s", domain)
	}

	return nil
}

// ConvertHandle will convert a $handle or 1handle to a paymail address
//
// For HandCash: $handle = handle@handcash.io or handle@beta.handcash.io
// For RelayX:   1handle = handle@relayx.io
func ConvertHandle(handle string, isBeta bool) string {
	if strings.HasPrefix(handle, "$") {
		if isBeta {
			return strings.ToLower(strings.Replace(handle, "$", "", -1)) + "@beta.handcash.io"
		}
		return strings.ToLower(strings.Replace(handle, "$", "", -1)) + "@handcash.io"
	} else if strings.HasPrefix(handle, "1") && len(handle) < 25 && !strings.Contains(handle, "@") {
		return strings.ToLower(strings.Replace(handle, "1", "", -1)) + "@relayx.io"
	}
	return handle
}

// ValidateTimestamp will test if the timestamp is valid
//
// This is used to validate the "dt" parameter in resolve_address.go
// Allowing 3 minutes before/after for
func ValidateTimestamp(timestamp string) error {

	// Parse the time using the RFC3339 layout
	dt, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return err
	}

	// Timestamp cannot be more than 2 minutes in the past
	// Specs: http://bsvalias.org/04-02-sender-validation.html
	if dt.Before(time.Now().UTC().Add(-2 * time.Minute)) {
		return fmt.Errorf("timestamp: %s is in the past", timestamp)
	} else if dt.After(time.Now().UTC().Add(2 * time.Minute)) {
		return fmt.Errorf("timestamp: %s is in the future", timestamp)
	}

	return nil
}
