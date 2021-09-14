package tonicpow

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ConversionOps allow functional options to be supplied
// that overwrite default conversion options.
type ConversionOps func(c *conversionOptions)

// conversionOptions holds all the configuration for the conversion
type conversionOptions struct {
	goalID           uint64  // Goal by ID
	goalName         string  // Goal by name
	tncpwSession     string  // tncpw session
	customDimensions string  // (optional) custom dimensions to add to the conversion
	purchaseAmount   float64 // (optional) purchase amount (total for e-commerce)
	delayInMinutes   uint64  // (optional) delay the conversion x minutes (before processing, allowing cancellation)
	tonicPowUserID   uint64  // (optional) trigger a conversion for a specific user
}

// validate will check the options before processing
func (o *conversionOptions) validate() error {
	if o.goalID == 0 && len(o.goalName) == 0 {
		return fmt.Errorf("missing required attribute(s): %s or %s", fieldID, fieldName)
	} else if o.goalID == 0 && o.tonicPowUserID > 0 {
		return fmt.Errorf("missing required attribute: %s", fieldID)
	} else if o.tonicPowUserID == 0 && len(o.tncpwSession) == 0 {
		return fmt.Errorf("missing required attribute(s): %s or %s", fieldVisitorSessionGUID, fieldUserID)
	}
	return nil
}

// payload will generate the payload given the options
func (o *conversionOptions) payload() map[string]string {
	m := map[string]string{}

	// Set goal id
	if o.goalID > 0 {
		m[fieldGoalID] = fmt.Sprintf("%d", o.goalID)
	}

	// Set goal name
	if len(o.goalName) > 0 {
		m[fieldName] = o.goalName
	}

	// Set tonic pow user
	if o.tonicPowUserID > 0 {
		m[fieldUserID] = fmt.Sprintf("%d", o.tonicPowUserID)
	} else if len(o.tncpwSession) > 0 {
		m[fieldVisitorSessionGUID] = o.tncpwSession
	}

	// Set delay in minutes
	if o.delayInMinutes > 0 {
		m[fieldDelayInMinutes] = fmt.Sprintf("%d", o.delayInMinutes)
	}

	// Set purchase amount
	if o.purchaseAmount > 0 {
		m[fieldAmount] = fmt.Sprintf("%f", o.purchaseAmount)
	}

	// Set custom dimensions
	if len(o.customDimensions) > 0 {
		m[fieldCustomDimensions] = o.customDimensions
	}

	return m
}

// WithGoalID will set a goal ID
func WithGoalID(goalID uint64) ConversionOps {
	return func(c *conversionOptions) {
		c.goalID = goalID
	}
}

// WithGoalName will set a goal name
func WithGoalName(name string) ConversionOps {
	return func(c *conversionOptions) {
		c.goalName = name
	}
}

// WithTncpwSession will set a tncpw_session
func WithTncpwSession(session string) ConversionOps {
	return func(c *conversionOptions) {
		c.tncpwSession = session
	}
}

// WithCustomDimensions will set custom dimensions (string / json)
func WithCustomDimensions(dimensions string) ConversionOps {
	return func(c *conversionOptions) {
		c.customDimensions = dimensions
	}
}

// WithPurchaseAmount will set purchase amount from e-commerce
func WithPurchaseAmount(amount float64) ConversionOps {
	return func(c *conversionOptions) {
		c.purchaseAmount = amount
	}
}

// WithDelay will set a delay in minutes
func WithDelay(minutes uint64) ConversionOps {
	return func(c *conversionOptions) {
		c.delayInMinutes = minutes
	}
}

// WithUserID will set a tonicpow user ID
func WithUserID(userID uint64) ConversionOps {
	return func(c *conversionOptions) {
		c.tonicPowUserID = userID
	}
}

// CreateConversion will fire a conversion for a given goal, if successful it will make a new Conversion
//
// For more information: https://docs.tonicpow.com/#caeffdd5-eaad-4fc8-ac01-8288b50e8e27
func (c *Client) CreateConversion(opts ...ConversionOps) (conversion *Conversion,
	response *StandardResponse, err error) {

	// Start the options
	options := new(conversionOptions)

	// Set the conversion options
	for _, opt := range opts {
		opt(options)
	}

	// Validate options
	if err = options.validate(); err != nil {
		return
	}

	// Fire the Request
	if response, err = c.Request(
		http.MethodPost,
		"/"+modelConversion,
		options.payload(), http.StatusCreated,
	); err != nil {
		return
	}

	err = json.Unmarshal(response.Body, &conversion)
	return
}

// GetConversion will get an existing conversion
// This will return an Error if the goal is not found (404)
//
// For more information: https://docs.tonicpow.com/#fce465a1-d8d5-442d-be22-95169170167e
func (c *Client) GetConversion(conversionID uint64) (conversion *Conversion,
	response *StandardResponse, err error) {

	// Must have an ID
	if conversionID == 0 {
		err = fmt.Errorf("missing required attribute: %s", fieldID)
		return
	}

	// Fire the Request
	if response, err = c.Request(
		http.MethodGet,
		fmt.Sprintf("/%s/details/%d", modelConversion, conversionID),
		nil, http.StatusOK,
	); err != nil {
		return
	}

	err = json.Unmarshal(response.Body, &conversion)
	return
}

// CancelConversion will cancel an existing conversion (if delay was set and > 1 minute remaining)
//
// For more information: https://docs.tonicpow.com/#e650b083-bbb4-4ff7-9879-c14b1ab3f753
func (c *Client) CancelConversion(conversionID uint64, cancelReason string) (conversion *Conversion,
	response *StandardResponse, err error) {

	// Must have an ID
	if conversionID == 0 {
		err = fmt.Errorf("missing required attribute: %s", fieldID)
		return
	}

	// Fire the Request
	if response, err = c.Request(
		http.MethodPut,
		fmt.Sprintf("/%s/cancel", modelConversion),
		map[string]string{
			fieldID:     fmt.Sprintf("%d", conversionID),
			fieldReason: cancelReason,
		},
		http.StatusOK,
	); err != nil {
		return
	}

	err = json.Unmarshal(response.Body, &conversion)
	return
}
