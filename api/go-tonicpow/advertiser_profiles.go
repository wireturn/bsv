package tonicpow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// permitFields will remove fields that cannot be used
func (a *AdvertiserProfile) permitFields() {
	a.PublicGUID = ""
	a.UserID = 0
}

// GetAdvertiserProfile will get an existing advertiser profile
// This will return an Error if the profile is not found (404)
//
// For more information: https://docs.tonicpow.com/#b3a62d35-7778-4314-9321-01f5266c3b51
func (c *Client) GetAdvertiserProfile(profileID uint64) (profile *AdvertiserProfile,
	response *StandardResponse, err error) {

	// Must have an ID
	if profileID == 0 {
		err = fmt.Errorf("missing field: %s", fieldID)
		return
	}

	// Fire the Request
	if response, err = c.Request(
		http.MethodGet,
		fmt.Sprintf("/%s/details/%d", modelAdvertiser, profileID),
		nil, http.StatusOK,
	); err != nil {
		return
	}

	// Convert model response
	err = json.Unmarshal(response.Body, &profile)
	return
}

// UpdateAdvertiserProfile will update an existing profile
//
// For more information: https://docs.tonicpow.com/#0cebd1ff-b1ce-4111-aff6-9d586f632a84
func (c *Client) UpdateAdvertiserProfile(profile *AdvertiserProfile) (*StandardResponse, error) {

	// Basic requirements
	if profile.ID == 0 {
		return nil, fmt.Errorf("missing required attribute: %s", fieldID)
	}

	// Permit fields
	profile.permitFields()

	// Fire the Request
	response, err := c.Request(
		http.MethodPut,
		"/"+modelAdvertiser,
		profile, http.StatusOK,
	)
	if err != nil {
		return response, err
	}

	// Convert model response
	return response, json.Unmarshal(response.Body, &profile)
}

// ListCampaignsByAdvertiserProfile will return a list of campaigns
//
// For more information: https://docs.tonicpow.com/#98017e9a-37dd-4810-9483-b6c400572e0c
func (c *Client) ListCampaignsByAdvertiserProfile(profileID uint64, page, resultsPerPage int,
	sortBy, sortOrder string) (campaigns *CampaignResults, response *StandardResponse, err error) {

	// Basic requirements
	if profileID == 0 {
		err = fmt.Errorf("missing required attribute: %s", fieldAdvertiserProfileID)
		return
	}

	// Do we know this field?
	if len(sortBy) > 0 {
		if !isInList(strings.ToLower(sortBy), campaignSortFields) {
			err = fmt.Errorf("sort by %s is not valid", sortBy)
			return
		}
	} else {
		sortBy = SortByFieldCreatedAt
		sortOrder = SortOrderDesc
	}

	// Fire the Request
	if response, err = c.Request(
		http.MethodGet,
		fmt.Sprintf("/%s/%s/%d?%s=%d&%s=%d&%s=%s&%s=%s", modelAdvertiser, modelCampaign, profileID,
			fieldCurrentPage, page,
			fieldResultsPerPage, resultsPerPage,
			fieldSortBy, sortBy,
			fieldSortOrder, sortOrder,
		),
		nil, http.StatusOK,
	); err != nil {
		return
	}

	// Convert model response
	err = json.Unmarshal(response.Body, &campaigns)
	return
}

// ListAppsByAdvertiserProfile will return a list of apps
//
// For more information: https://docs.tonicpow.com/#9c9fa8dc-3017-402e-8059-136b0eb85c2e
func (c *Client) ListAppsByAdvertiserProfile(profileID uint64, page, resultsPerPage int,
	sortBy, sortOrder string) (apps *AppResults, response *StandardResponse, err error) {

	// Basic requirements
	if profileID == 0 {
		err = fmt.Errorf("missing required attribute: %s", fieldAdvertiserProfileID)
		return
	}

	// Do we know this field?
	if len(sortBy) > 0 {
		if !isInList(strings.ToLower(sortBy), appSortFields) {
			err = fmt.Errorf("sort by %s is not valid", sortBy)
			return
		}
	} else {
		sortBy = SortByFieldCreatedAt
		sortOrder = SortOrderDesc
	}

	// Fire the Request
	if response, err = c.Request(
		http.MethodGet,
		fmt.Sprintf(
			"/%s/%s/?%s=%d&%s=%d&%s=%d&%s=%s&%s=%s",
			modelAdvertiser, modelApp,
			fieldID, profileID,
			fieldCurrentPage, page,
			fieldResultsPerPage, resultsPerPage,
			fieldSortBy, sortBy,
			fieldSortOrder, sortOrder,
		),
		nil, http.StatusOK,
	); err != nil {
		return
	}

	// Convert model response
	err = json.Unmarshal(response.Body, &apps)
	return
}
