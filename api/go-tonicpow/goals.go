package tonicpow

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// permitFields will remove fields that cannot be used
func (g *Goal) permitFields() {
	g.CampaignID = 0
}

// CreateGoal will make a new goal
//
// For more information: https://docs.tonicpow.com/#29a93e9b-9726-474c-b25e-92586200a803
func (c *Client) CreateGoal(goal *Goal) (*StandardResponse, error) {

	// Basic requirements
	if goal.CampaignID == 0 {
		return nil, fmt.Errorf(fmt.Sprintf("missing required attribute: %s", fieldCampaignID))
	} else if len(goal.Name) == 0 {
		return nil, fmt.Errorf(fmt.Sprintf("missing required attribute: %s", fieldName))
	}

	// Fire the Request
	response, err := c.Request(
		http.MethodPost,
		"/"+modelGoal,
		goal, http.StatusCreated,
	)
	if err != nil {
		return response, err
	}

	return response, json.Unmarshal(response.Body, &goal)
}

// GetGoal will get an existing goal
// This will return an Error if the goal is not found (404)
//
// For more information: https://docs.tonicpow.com/#48d7bbc8-5d7b-4078-87b7-25f545c3deaf
func (c *Client) GetGoal(goalID uint64) (goal *Goal, response *StandardResponse, err error) {

	// Must have an ID
	if goalID == 0 {
		err = fmt.Errorf("missing required attribute: %s", fieldID)
		return
	}

	// Fire the Request
	if response, err = c.Request(
		http.MethodGet,
		fmt.Sprintf("/%s/details/%d", modelGoal, goalID),
		nil, http.StatusOK,
	); err != nil {
		return
	}

	err = json.Unmarshal(response.Body, &goal)
	return
}

// UpdateGoal will update an existing goal
//
// For more information: https://docs.tonicpow.com/#395f5b7d-6a5d-49c8-b1ae-abf7f90b42a2
func (c *Client) UpdateGoal(goal *Goal) (*StandardResponse, error) {

	// Basic requirements
	if goal.ID == 0 {
		return nil, fmt.Errorf("missing required attribute: %s", fieldID)
	}

	// Permit fields
	goal.permitFields()

	// Fire the Request
	response, err := c.Request(
		http.MethodPut,
		"/"+modelGoal,
		goal, http.StatusOK,
	)
	if err != nil {
		return response, err
	}

	return response, json.Unmarshal(response.Body, &goal)
}

// DeleteGoal will delete an existing goal
//
// For more information: https://docs.tonicpow.com/#38605b65-72c9-4fc8-87a7-bc644bc89a96
func (c *Client) DeleteGoal(goalID uint64) (bool, *StandardResponse, error) {

	// Basic requirements
	if goalID == 0 {
		return false, nil, fmt.Errorf("missing required attribute: %s", fieldID)
	}

	// Fire the Request
	response, err := c.Request(
		http.MethodDelete,
		fmt.Sprintf("/%s?%s=%d", modelGoal, fieldID, goalID),
		nil, http.StatusOK,
	)
	if err != nil {
		return false, response, err
	}

	// Flag for deleted if no error and good response
	return true, response, err
}
