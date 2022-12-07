package paymail

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// BRFCSpec is a full BRFC specification document
//
// See more: http://bsvalias.org/01-brfc-specifications.html
type BRFCSpec struct {
	Alias      string `json:"alias,omitempty"`      // Alias is used in the list of capabilities
	Author     string `json:"author"`               // Free-form, could include a name, alias, paymail address, GitHub/social media handle, etc.
	ID         string `json:"id"`                   // Public BRFC ID
	Supersedes string `json:"supersedes,omitempty"` // A BRFC ID (or list of IDs) that this document supersedes
	Title      string `json:"title"`                // Title of the brfc
	URL        string `json:"url,omitempty"`        // Public URL to view the specification
	Version    string `json:"version"`              // No set format; could be a sequence number, publication date, or any other scheme
	Valid      bool   `json:"valid"`                // Validated the ID -> (title,author,version)
}

// LoadBRFCs will load the known "default" specifications into structs from JSON
//
// additionSpecifications is appended to the default specs
// BRFCKnownSpecifications is a local constant of JSON to pre-load known BRFC ids
func (c *ClientOptions) LoadBRFCs(additionalSpecifications string) (err error) {

	// Load the default specs
	if err = json.Unmarshal([]byte(BRFCKnownSpecifications), &c.BRFCSpecs); err != nil {
		// This error case should never occur since the JSON is hardcoded, but good practice anyway
		return
	}

	// No additional specs to process
	if len(additionalSpecifications) == 0 {
		return
	}

	// Process the additional specifications
	var tempSpecs []*BRFCSpec
	if err = json.Unmarshal([]byte(additionalSpecifications), &tempSpecs); err != nil {
		return
	}

	// Add the specs to the existing specifications
	for _, spec := range tempSpecs {

		// Validate the spec before adding
		var valid bool
		var id string
		if valid, id, err = spec.Validate(); err != nil {
			return
		} else if !valid {
			err = fmt.Errorf("brfc: [%s] is invalid - id returned: %s vs %s", spec.Title, id, spec.ID)
			return
		}

		// Add to existing list
		c.BRFCSpecs = append(c.BRFCSpecs, spec)
	}

	return
}

// Generate will generate a new BRFC ID from the given specification
//
// See more: http://bsvalias.org/01-02-brfc-id-assignment.html
func (b *BRFCSpec) Generate() error {

	// Validate the title (only required field)
	if len(b.Title) == 0 {
		b.ID = ""
		return fmt.Errorf("invalid brfc title, length: 0")
	}

	// Start a new SHA256 hash
	h := sha256.New()

	// Append all values (trim leading & trailing whitespace)
	_, _ = h.Write([]byte(strings.TrimSpace(b.Title) + strings.TrimSpace(b.Author) + strings.TrimSpace(b.Version)))

	// Start the double SHA256
	h2 := sha256.New()

	// Write the first SHA256 result
	_, _ = h2.Write(h.Sum(nil))

	// Create the final double SHA256
	doubleHash := h2.Sum(nil)

	// Reverse the order
	for i, j := 0, len(doubleHash)-1; i < j; i, j = i+1, j-1 {
		doubleHash[i], doubleHash[j] = doubleHash[j], doubleHash[i]
	}

	// Hex encode the value
	hexDoubleHash := make([]byte, hex.EncodedLen(len(doubleHash)))
	hex.Encode(hexDoubleHash, doubleHash)

	// Check that the ID length is valid
	// this error case was never hit as long as title is len() > 0
	/*
		if len(hexDoubleHash) < 12 {
			b.ID = ""
			return fmt.Errorf("failed to generate a valid id, length was %d", len(hexDoubleHash))
		}
	*/

	// Extract the ID and set (first 12 characters)
	b.ID = string(hexDoubleHash[:12])

	return nil
}

// Validate will check if the BRFC is valid or not (and set b.Valid)
//
// Returns the ID that was generated to compare against the existing id
// Returns valid bool for convenience, but also sets b.Valid = true
func (b *BRFCSpec) Validate() (valid bool, id string, err error) {

	// Copy and generate (copying ensures that running Generate() will not override the existing ID)
	tempBRFC := new(BRFCSpec)
	*tempBRFC = *b

	// Start by invalidating the BRFC
	b.Valid = false

	// Run the generate method to return an ID
	if err = tempBRFC.Generate(); err != nil {
		return
	}

	// Set the ID generated (for external comparison etc)
	id = tempBRFC.ID

	// Test if the ID generated matches what was set previously
	if tempBRFC.ID == b.ID {
		valid = true
		b.Valid = valid
	}

	return
}
