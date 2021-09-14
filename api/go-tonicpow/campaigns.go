package tonicpow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// permitFields will remove fields that cannot be used
func (c *Campaign) permitFields() {
	c.AdvertiserProfileID = 0
}

// CreateCampaign will make a new campaign for the associated advertiser profile
//
// For more information: https://docs.tonicpow.com/#b67e92bf-a481-44f6-a31d-26e6e0c521b1
func (c *Client) CreateCampaign(campaign *Campaign) (*StandardResponse, error) {

	// Basic requirements
	if campaign.AdvertiserProfileID == 0 {
		return nil, fmt.Errorf("missing required attribute: %s", fieldAdvertiserProfileID)
	} else if len(campaign.Title) == 0 {
		return nil, fmt.Errorf("missing required attribute: %s", fieldTitle)
	} else if len(campaign.Description) == 0 {
		return nil, fmt.Errorf("missing required attribute: %s", fieldDescription)
	} else if len(campaign.TargetURL) == 0 {
		return nil, fmt.Errorf("missing required attribute: %s", fieldTargetURL)
	}

	// Fire the Request
	var response *StandardResponse
	var err error
	if response, err = c.Request(
		http.MethodPost,
		"/"+modelCampaign,
		campaign, http.StatusCreated,
	); err != nil {
		return response, err
	}

	// Convert model response
	return response, json.Unmarshal(response.Body, &campaign)
}

// GetCampaign will get an existing campaign by ID
// This will return an Error if the campaign is not found (404)
//
// For more information: https://docs.tonicpow.com/#b827446b-be34-4678-b347-33c4f63dbf9e
func (c *Client) GetCampaign(campaignID uint64) (campaign *Campaign,
	response *StandardResponse, err error) {

	// Must have an ID
	if campaignID == 0 {
		err = fmt.Errorf("missing required attribute: %s", fieldID)
		return
	}

	// Fire the Request
	if response, err = c.Request(
		http.MethodGet,
		fmt.Sprintf("/%s/details/?%s=%d", modelCampaign, fieldID, campaignID),
		nil, http.StatusOK,
	); err != nil {
		return
	}

	// Convert model response
	err = json.Unmarshal(response.Body, &campaign)
	return
}

// GetCampaignBySlug will get an existing campaign by slug
// This will return an Error if the campaign is not found (404)
//
// For more information: https://docs.tonicpow.com/#b827446b-be34-4678-b347-33c4f63dbf9e
func (c *Client) GetCampaignBySlug(slug string) (campaign *Campaign,
	response *StandardResponse, err error) {

	// Must have a slug
	if len(slug) == 0 {
		err = fmt.Errorf("missing required attribute: %s", fieldSlug)
		return
	}

	// Fire the Request
	if response, err = c.Request(
		http.MethodGet,
		fmt.Sprintf("/%s/details/?%s=%s", modelCampaign, fieldSlug, slug),
		nil, http.StatusOK,
	); err != nil {
		return
	}

	// Convert model response
	err = json.Unmarshal(response.Body, &campaign)
	return
}

// UpdateCampaign will update an existing campaign
//
// For more information: https://docs.tonicpow.com/#665eefd6-da42-4ca9-853c-fd8ca1bf66b2
func (c *Client) UpdateCampaign(campaign *Campaign) (response *StandardResponse, err error) {

	// Basic requirements
	if campaign.ID == 0 {
		err = fmt.Errorf("missing required attribute: %s", fieldID)
		return
	}

	// Permit fields
	campaign.permitFields()

	// Fire the Request
	if response, err = c.Request(
		http.MethodPut,
		"/"+modelCampaign,
		campaign, http.StatusOK,
	); err != nil {
		return
	}

	err = json.Unmarshal(response.Body, &campaign)
	return
}

// CampaignsFeed will return a feed of active campaigns
// This will return an Error if no campaigns are found (404)
//
// For more information: https://docs.tonicpow.com/#b3fe69d3-24ba-4c2a-a485-affbb0a738de
func (c *Client) CampaignsFeed(feedType FeedType) (feed string, response *StandardResponse, err error) {

	// Fire the Request
	if response, err = c.Request(
		http.MethodGet,
		fmt.Sprintf("/%s/feed/?%s=%s", modelCampaign, fieldFeedType, feedType),
		nil, http.StatusOK,
	); err != nil {
		return
	}

	feed = string(response.Body)
	return
}

// ListCampaigns will return a list of campaigns
// This will return an Error if the campaign is not found (404)
//
// For more information: https://docs.tonicpow.com/#c1b17be6-cb10-48b3-a519-4686961ff41c
func (c *Client) ListCampaigns(page, resultsPerPage int, sortBy, sortOrder, searchQuery string,
	minimumBalance uint64, includeExpired bool) (results *CampaignResults, response *StandardResponse, err error) {

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
		fmt.Sprintf(
			"/%s/list?%s=%d&%s=%d&%s=%s&%s=%s&%s=%s&%s=%d&%s=%t",
			modelCampaign,
			fieldCurrentPage, page,
			fieldResultsPerPage, resultsPerPage,
			fieldSortBy, sortBy,
			fieldSortOrder, sortOrder,
			fieldSearchQuery, searchQuery,
			fieldMinimumBalance, minimumBalance,
			fieldExpired, includeExpired,
		),
		nil, http.StatusOK,
	); err != nil {
		return
	}

	err = json.Unmarshal(response.Body, &results)
	return
}

// ListCampaignsByURL will return a list of campaigns using the target url
// This will return an Error if the url is not found (404)
//
// For more information: https://docs.tonicpow.com/#30a15b69-7912-4e25-ba41-212529fba5ff
func (c *Client) ListCampaignsByURL(targetURL string, page, resultsPerPage int,
	sortBy, sortOrder string) (results *CampaignResults, response *StandardResponse, err error) {

	// Must have a value
	if len(targetURL) == 0 {
		err = fmt.Errorf("missing required attribute: %s", fieldTargetURL)
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
		fmt.Sprintf("/%s/list?%s=%s&%s=%d&%s=%d&%s=%s&%s=%s",
			modelCampaign,
			fieldTargetURL, targetURL,
			fieldCurrentPage, page,
			fieldResultsPerPage, resultsPerPage,
			fieldSortBy, sortBy,
			fieldSortOrder, sortOrder,
		),
		nil, http.StatusOK,
	); err != nil {
		return
	}

	err = json.Unmarshal(response.Body, &results)
	return
}
