package tonicpow

// AdvertiserProfile is the advertiser_profile model (child of User)
//
// For more information: https://docs.tonicpow.com/#2f9ec542-0f88-4671-b47c-d0ee390af5ea
type AdvertiserProfile struct {
	HomepageURL         string `json:"homepage_url"`
	IconURL             string `json:"icon_url"`
	PublicGUID          string `json:"public_guid"`
	Name                string `json:"name"`
	ID                  uint64 `json:"id,omitempty"`
	LinkServiceDomainID uint64 `json:"link_service_domain_id"`
	UserID              uint64 `json:"user_id"`
	DomainVerified      bool   `json:"domain_verified"`
	Unlisted            bool   `json:"unlisted"`
}

// AdvertiserResults is the page response for advertiser profile results from listing
type AdvertiserResults struct {
	Advertisers    []*AdvertiserProfile `json:"advertisers"`
	CurrentPage    int                  `json:"current_page"`
	Results        int                  `json:"results"`
	ResultsPerPage int                  `json:"results_per_page"`
}

// App is the app model (child of advertiser_profile)
//
// For more information: (todo)
type App struct {
	AdvertiserProfileID uint64 `json:"advertiser_profile_id"`
	ID                  uint64 `json:"id"`
	Name                string `json:"name"`
	UserID              uint64 `json:"user_id"`
	WebhookURL          string `json:"webhook_url"`
}

// AppResults is the page response for app results from listing
type AppResults struct {
	Apps           []*App `json:"apps"`
	CurrentPage    int    `json:"current_page"`
	Results        int    `json:"results"`
	ResultsPerPage int    `json:"results_per_page"`
}

// Campaign is the campaign model (child of AdvertiserProfile)
//
// For more information: https://docs.tonicpow.com/#5aca2fc7-b3c8-445b-aa88-f62a681f8e0c
type Campaign struct {
	Goals                 []*Goal               `json:"goals"`
	Images                []*CampaignImage      `json:"images"`
	CreatedAt             string                `json:"created_at"`
	LastEventAt           string                `json:"last_event_at"`
	Currency              string                `json:"currency"`
	Description           string                `json:"description"`
	ExpiresAt             string                `json:"expires_at"`
	FundingAddress        string                `json:"funding_address"`
	ImageURL              string                `json:"image_url"`
	PublicGUID            string                `json:"public_guid"`
	Slug                  string                `json:"slug"`
	TargetURL             string                `json:"target_url"`
	Title                 string                `json:"title"`
	TxID                  string                `json:"-"`
	AdvertiserProfile     *AdvertiserProfile    `json:"advertiser_profile"`
	Balance               float64               `json:"balance"`
	BalanceAlertThreshold float64               `json:"balance_alert_threshold"`
	PayPerClickRate       float64               `json:"pay_per_click_rate"`
	AdvertiserProfileID   uint64                `json:"advertiser_profile_id"`
	BalanceSatoshis       uint64                `json:"balance_satoshis"`
	ID                    uint64                `json:"id,omitempty"`
	LinksCreated          uint64                `json:"links_created"`
	LinkServiceDomainID   uint64                `json:"link_service_domain_id"`
	PaidClicks            uint64                `json:"paid_clicks"`
	PaidConversions       uint64                `json:"paid_conversions"`
	PayoutMode            int                   `json:"payout_mode"`
	Requirements          *CampaignRequirements `json:"requirements"`
	BotProtection         bool                  `json:"bot_protection"`
	ContributeEnabled     bool                  `json:"contribute_enabled"`
	DomainVerified        bool                  `json:"domain_verified"`
	Unlisted              bool                  `json:"unlisted"`
	MatchDomain           bool                  `json:"match_domain"`
}

// CampaignImage is the structure of the image metadata
type CampaignImage struct {
	Height   int    `json:"height"`
	MimeType string `json:"mime_type"`
	URL      string `json:"url"`
	Width    int    `json:"width"`
}

// CampaignRequirements is the structure for "requirements"
//
// DO NOT CHANGE ORDER - malign
//
type CampaignRequirements struct {
	VisitorCountries    []string `json:"visitor_countries"`
	DotWallet           bool     `json:"dotwallet"`
	Facebook            bool     `json:"facebook"`
	Google              bool     `json:"google"`
	HandCash            bool     `json:"handcash"`
	KYC                 bool     `json:"kyc"`
	MoneyButton         bool     `json:"moneybutton"`
	PPCBid              bool     `json:"ppc_bid"`
	Relay               bool     `json:"relay"`
	Twitter             bool     `json:"twitter"`
	VisitorRestrictions bool     `json:"visitor_restrictions"`
}

// CampaignResults is the page response for campaign results from listing
type CampaignResults struct {
	Campaigns      []*Campaign `json:"campaigns"`
	CurrentPage    int         `json:"current_page"`
	Results        int         `json:"results"`
	ResultsPerPage int         `json:"results_per_page"`
}

// Conversion is the response of getting a conversion
//
// For more information: https://docs.tonicpow.com/#75c837d5-3336-4d87-a686-d80c6f8938b9
type Conversion struct {
	Amount           float64 `json:"amount,omitempty"`
	CampaignID       uint64  `json:"campaign_id"`
	CustomDimensions string  `json:"custom_dimensions"`
	GoalID           uint64  `json:"goal_id"`
	GoalName         string  `json:"goal_name,omitempty"`
	ID               uint64  `json:"id,omitempty"`
	PayoutAfter      string  `json:"payout_after,omitempty"`
	Status           string  `json:"status"`
	StatusData       string  `json:"status_data"`
	TxID             string  `json:"tx_id"`
	UserID           uint64  `json:"user_id"`
}

// Goal is the goal model (child of Campaign)
//
// For more information: https://docs.tonicpow.com/#316b77ab-4900-4f3d-96a7-e67c00af10ca
type Goal struct {
	CampaignID      uint64  `json:"campaign_id"`
	Description     string  `json:"description"`
	ID              uint64  `json:"id,omitempty"`
	LastConvertedAt string  `json:"last_converted_at"`
	MaxPerPromoter  int16   `json:"max_per_promoter"`
	MaxPerVisitor   int16   `json:"max_per_visitor"`
	Name            string  `json:"name"`
	PayoutRate      float64 `json:"payout_rate"`
	Payouts         int     `json:"payouts"`
	PayoutType      string  `json:"payout_type"`
	Title           string  `json:"title"`
}

// Rate is the rate results
//
// For more information: https://docs.tonicpow.com/#fb00736e-61b9-4ec9-acaf-e3f9bb046c89
type Rate struct {
	Currency        string  `json:"currency"`
	CurrencyAmount  float64 `json:"currency_amount"`
	PriceInSatoshis int64   `json:"price_in_satoshis"`
}
