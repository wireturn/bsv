package validate

import "fmt"

// ValidationError is key and message of the corresponding error
type ValidationError struct {

	// Key is the Field name, key name
	Key string

	// Message is the error message
	Message string
}

// ValidationError returns a string of key + message
func (v *ValidationError) Error() string {
	return v.Key + " " + v.Message
}

// ValidationErrors is a slice of validation errors
type ValidationErrors []ValidationError

// ValidationError returns the list of errors from the slice of errors
func (v ValidationErrors) Error() (errors string) {

	// No errors?
	if len(v) == 0 {
		return
	}

	// Get the first error
	errors = v[0].Error()

	// Add x other errors
	if len(v) > 1 {
		errors += fmt.Sprintf(" and %d other errors", len(v)-1)
	}

	return
}
