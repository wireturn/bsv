package tonicpow

import (
	"time"
)

const (
	// Package configuration defaults
	apiVersion         string = "v1"
	defaultHTTPTimeout        = 10 * time.Second          // Default timeout for all GET requests in seconds
	defaultRetryCount  int    = 2                         // Default retry count for HTTP requests
	defaultUserAgent          = "go-tonicpow: " + version // Default user agent
	version            string = "v0.6.8"                  // go-tonicpow version

	// Field key names for various model requests
	fieldAdvertiserProfileID = "advertiser_profile_id"
	fieldAmount              = "amount"
	fieldAPIKey              = "api_key"
	fieldCampaignID          = "campaign_id"
	fieldCurrency            = "currency"
	fieldCurrentPage         = "current_page"
	fieldCustomDimensions    = "custom_dimensions"
	fieldDelayInMinutes      = "delay_in_minutes"
	fieldDescription         = "description"
	fieldExpired             = "expired"
	fieldFeedType            = "feed_type"
	fieldGoalID              = "goal_id"
	fieldID                  = "id"
	fieldMinimumBalance      = "minimum_balance"
	fieldName                = "name"
	fieldReason              = "reason"
	fieldResultsPerPage      = "results_per_page"
	fieldSearchQuery         = "query"
	fieldSlug                = "slug"
	fieldSortBy              = "sort_by"
	fieldSortOrder           = "sort_order"
	fieldTargetURL           = "target_url"
	fieldTitle               = "title"
	fieldUserID              = "user_id"
	fieldVisitorSessionGUID  = "tncpw_session"

	// Model names (used for Request endpoints)
	modelAdvertiser string = "advertisers"
	modelApp        string = "apps"
	modelCampaign   string = "campaigns"
	modelConversion string = "conversions"
	modelGoal       string = "goals"
	modelRates      string = "rates"

	// Environment names
	environmentDevelopmentAlias string = "local"
	environmentDevelopmentName  string = "development"
	environmentLiveAlias        string = "production"
	environmentLiveName         string = "live"
	environmentStagingAlias     string = "beta"
	environmentStagingName      string = "staging"

	// Environment API URLs
	developmentURL = "http://localhost:3000/" + apiVersion
	liveAPIURL     = "https://api.tonicpow.com/" + apiVersion
	stagingAPIURL  = "https://api.staging.tonicpow.com/" + apiVersion

	// SortByFieldBalance is for sorting results by field: balance
	SortByFieldBalance string = "balance"

	// SortByFieldCreatedAt is for sorting results by field: created_at
	SortByFieldCreatedAt string = "created_at"

	// SortByFieldName is for sorting results by field: name
	SortByFieldName string = "name"

	// SortByFieldLinksCreated is for sorting results by field: links_created
	SortByFieldLinksCreated string = "links_created"

	// SortByFieldPaidClicks is for sorting results by field: paid_clicks
	SortByFieldPaidClicks string = "paid_clicks"

	// SortByFieldPayPerClick is for sorting results by field: pay_per_click_rate
	SortByFieldPayPerClick string = "pay_per_click_rate"

	// SortOrderAsc is for returning the results in ascending order
	SortOrderAsc string = "asc"

	// SortOrderDesc is for returning the results in descending order
	SortOrderDesc string = "desc"

	// FeedTypeAtom is for using the feed type: Atom
	FeedTypeAtom FeedType = "atom"

	// FeedTypeJSON is for using the feed type: JSON
	FeedTypeJSON FeedType = "json"

	// FeedTypeRSS is for using the feed type: RSS
	FeedTypeRSS FeedType = "rss"
)

var (

	// appSortFields is used for allowing specific fields for sorting
	appSortFields = []string{
		SortByFieldCreatedAt,
		SortByFieldName,
	}

	// campaignSortFields is used for allowing specific fields for sorting
	campaignSortFields = []string{
		SortByFieldBalance,
		SortByFieldCreatedAt,
		SortByFieldLinksCreated,
		SortByFieldPaidClicks,
		SortByFieldPayPerClick,
	}
)

// FeedType is used for the campaign feeds (rss, atom, json)
type FeedType string

// Environment is used for changing the Environment for running client requests
type Environment struct {
	alias  string
	apiURL string
	name   string
}

// Alias will return the Environment's alias
func (e Environment) Alias() string {
	return e.alias
}

// Name will return the Environment's name
func (e Environment) Name() string {
	return e.name
}

// URL will return the Environment's url
func (e Environment) URL() string {
	return e.apiURL
}

// Current environments available
var (
	EnvironmentLive = Environment{
		apiURL: liveAPIURL,
		name:   environmentLiveName,
		alias:  environmentLiveAlias,
	}
	EnvironmentStaging = Environment{
		apiURL: stagingAPIURL,
		name:   environmentStagingName,
		alias:  environmentStagingAlias,
	}
	EnvironmentDevelopment = Environment{
		apiURL: developmentURL,
		name:   environmentDevelopmentName,
		alias:  environmentDevelopmentAlias,
	}
)

// Error is the universal Error response from the API
//
// For more information: https://docs.tonicpow.com/#d7fe13a3-2b6d-4399-8d0f-1d6b8ad6ebd9
type Error struct {
	Code        int         `json:"code"`
	Data        interface{} `json:"data"`
	IPAddress   string      `json:"ip_address"`
	Message     string      `json:"message"`
	Method      string      `json:"method"`
	RequestGUID string      `json:"request_guid"`
	StatusCode  int         `json:"status_code"`
	URL         string      `json:"url"`
}
