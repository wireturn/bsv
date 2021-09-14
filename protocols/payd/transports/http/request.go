package http

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"

	validator "github.com/theflyingcodr/govalidator"
)

// Bind will attempt to parse an io.reader to json, if this fails an error is returned.
//
// The key reason for this is to validate requests implementing the Validator interface
// this checks the requests for implementing Validator and will evaluate and return
// a validation error to be handled in the middleware.
func Bind(req io.Reader, out interface{}) error {
	if err := json.NewDecoder(req).Decode(&out); err != nil {
		return errors.Wrap(err, "failed to decode response")
	}
	if v, ok := out.(validator.Validator); ok {
		return v.Validate().Err()
	}
	return nil
}
