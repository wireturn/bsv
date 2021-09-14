package tonicpow

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestCampaign will return a dummy example for tests
func newTestCampaign() *Campaign {
	return &Campaign{
		Goals:               []*Goal{newTestGoal()},
		Images:              []*CampaignImage{newTestCampaignImages()},
		CreatedAt:           "2021-01-01 00:00:01",
		Currency:            "usd",
		Description:         "This is a test campaign",
		FundingAddress:      "124oW4xLDfay1BXmubUG9r64bGCCxnuf4g",
		ImageURL:            "https://res.cloudinary.com/tonicpow/image/upload/c_crop,x_0,y_0,w_2048,h_1024/w_600,h_300,c_limit,g_center/v1611266301/glfwmr0yhyjydeyfhyih.jpg",
		PublicGUID:          "b02e13a7001546b1b7efb9df40ab75e5",
		Slug:                "tonicpow",
		TargetURL:           "https://tonicpow.com",
		Title:               "TonicPow",
		AdvertiserProfileID: testAdvertiserID,
		AdvertiserProfile:   newTestAdvertiserProfile(),
		Balance:             13.37,
		PayPerClickRate:     1,
		BalanceSatoshis:     11333377,
		ID:                  testCampaignID,
		LinksCreated:        1,
		PaidClicks:          1,
		PaidConversions:     1,
		Requirements:        newTestCampaignRequirements(),
		BotProtection:       true,
		ContributeEnabled:   true,
		MatchDomain:         true,
	}
}

// newTestCampaignImages will return a dummy example for tests
func newTestCampaignImages() *CampaignImage {
	return &CampaignImage{
		Height:   300,
		MimeType: "image/jpeg",
		URL:      "https://res.cloudinary.com/tonicpow/image/upload/c_crop,x_0,y_0,w_2048,h_1024/w_600,h_300,c_limit,g_center/v1611266301/glfwmr0yhyjydeyfhyih.jpg",
		Width:    600,
	}
}

// newTestCampaignRequirements will return a dummy example for tests
func newTestCampaignRequirements() *CampaignRequirements {
	return &CampaignRequirements{
		HandCash:    true,
		MoneyButton: true,
	}
}

// newTestCampaignResults will return a dummy example for tests
func newTestCampaignResults(currentPage, resultsPerPage int) *CampaignResults {
	return &CampaignResults{
		Campaigns:      []*Campaign{newTestCampaign()},
		CurrentPage:    currentPage,
		Results:        1,
		ResultsPerPage: resultsPerPage,
	}
}

// newTestCampaignFeedRSS will return a dummy example for tests
func newTestCampaignFeedRSS() string {
	return `<?xml version="1.0" encoding="UTF-8"?><rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
  <channel>
    <title>TonicPow</title>
    <link>https://tonicpow.com</link>
    <description>List of active campaigns, ordered by newest first</description>
    <managingEditor>support@tonicpow.com (tonicpow)</managingEditor>
    <pubDate>Fri, 04 Jun 2021 17:20:28 +0000</pubDate>
    <item>
      <title>TonicPow</title>
      <link>https://tonicpow.com/campaign/tonicpow</link>
      <description>Earn BSV for sharing things you like. Offer BSV for sharing your product or service.</description>
      <author>TonicPow</author>
      <pubDate>Thu, 05 Sep 2019 01:50:06 +0000</pubDate>
    </item>
  </channel>
</rss>`
}

// newTestCampaignFeedAtom will return a dummy example for tests
func newTestCampaignFeedAtom() string {
	return `<?xml version="1.0" encoding="UTF-8"?><feed xmlns="http://www.w3.org/2005/Atom">
  <title>TonicPow</title>
  <id>https://tonicpow.com</id>
  <updated>2021-06-04T17:23:38Z</updated>
  <subtitle>List of active campaigns, ordered by newest first</subtitle>
  <link href="https://tonicpow.com"></link>
  <author>
    <name>tonicpow</name>
    <email>support@tonicpow.com</email>
  </author>
  <entry>
    <title>TonicPow</title>
    <updated>2019-09-05T01:50:06Z</updated>
    <id>tag:tonicpow.com,2019-09-05:/campaign/tonicpow</id>
    <link href="https://tonicpow.com/campaign/tonicpow" rel="alternate"></link>
    <summary type="html">Earn BSV for sharing things you like. Offer BSV for sharing your product or service.</summary>
    <author>
      <name>TonicPow</name>
    </author>
  </entry>
</feed>`
}

// newTestCampaignFeedJSON will return a dummy example for tests
func newTestCampaignFeedJSON() string {
	return `{
  "version": "https://jsonfeed.org/version/1",
  "title": "TonicPow",
  "home_page_url": "https://tonicpow.com",
  "description": "List of active campaigns, ordered by newest first",
  "author": {
    "name": "tonicpow"
  },
  "items": [
    {
      "id": "",
      "url": "https://tonicpow.com/campaign/tonicpow",
      "title": "TonicPow",
      "summary": "Earn BSV for sharing things you like. Offer BSV for sharing your product or service.",
      "date_published": "2019-09-05T01:50:06Z",
      "author": {
        "name": "TonicPow"
      }
    }
  ]
}`
}

// TestClient_CreateCampaign will test the method CreateCampaign()
func TestClient_CreateCampaign(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("create a campaign (success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelCampaign)

		err = mockResponseData(http.MethodPost, endpoint, http.StatusCreated, campaign)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.CreateCampaign(campaign)
		assert.NoError(t, err)
		assert.NotNil(t, campaign)
		assert.NotNil(t, response)
		assert.Equal(t, testCampaignID, campaign.ID)
	})

	t.Run("missing advertiser profile id", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()
		campaign.AdvertiserProfileID = 0

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelCampaign)

		err = mockResponseData(http.MethodPost, endpoint, http.StatusCreated, campaign)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.CreateCampaign(campaign)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("missing campaign title", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()
		campaign.Title = ""

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelCampaign)

		err = mockResponseData(http.MethodPost, endpoint, http.StatusCreated, campaign)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.CreateCampaign(campaign)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("missing campaign description", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()
		campaign.Description = ""

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelCampaign)

		err = mockResponseData(http.MethodPost, endpoint, http.StatusCreated, campaign)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.CreateCampaign(campaign)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("missing campaign target url", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()
		campaign.TargetURL = ""

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelCampaign)

		err = mockResponseData(http.MethodPost, endpoint, http.StatusCreated, campaign)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.CreateCampaign(campaign)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("error from api (status code)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelCampaign)

		err = mockResponseData(http.MethodPost, endpoint, http.StatusBadRequest, campaign)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.CreateCampaign(campaign)
		assert.Error(t, err)
		assert.NotNil(t, response)
	})

	t.Run("error from api (api error)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelCampaign)

		apiError := &Error{
			Code:        400,
			Data:        "field_name",
			IPAddress:   "127.0.0.1",
			Message:     "some error message",
			Method:      http.MethodPut,
			RequestGUID: "7f3d97a8fd67ff57861904df6118dcc8",
			StatusCode:  http.StatusBadRequest,
			URL:         endpoint,
		}

		err = mockResponseData(http.MethodPost, endpoint, http.StatusBadRequest, apiError)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.CreateCampaign(campaign)
		assert.Error(t, err)
		assert.Equal(t, apiError.Message, err.Error())
		assert.NotNil(t, response)
	})
}

// ExampleClient_CreateCampaign example using CreateCampaign()
//
// See more examples in /examples/
func ExampleClient_CreateCampaign() {

	// Load the client (using test client for example only)
	client, err := newTestClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// Mock response (for example only)
	responseCampaign := newTestCampaign()
	_ = mockResponseData(
		http.MethodPost,
		fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelCampaign),
		http.StatusCreated,
		responseCampaign,
	)

	// Create campaign (using mocking response)
	if _, err = client.CreateCampaign(responseCampaign); err != nil {
		fmt.Printf("error creating campaign: " + err.Error())
		return
	}
	fmt.Printf("created campaign: %s", responseCampaign.Title)
	// Output:created campaign: TonicPow
}

// BenchmarkClient_CreateCampaign benchmarks the method CreateCampaign()
func BenchmarkClient_CreateCampaign(b *testing.B) {
	client, _ := newTestClient()
	campaign := newTestCampaign()
	_ = mockResponseData(
		http.MethodPost,
		fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelCampaign),
		http.StatusCreated,
		campaign,
	)
	for i := 0; i < b.N; i++ {
		_, _ = client.CreateCampaign(campaign)
	}
}

// TestClient_GetCampaign will test the method GetCampaign()
func TestClient_GetCampaign(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("get a campaign (success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()

		endpoint := fmt.Sprintf(
			"%s/%s/details/?%s=%d", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldID, campaign.ID,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, campaign)
		assert.NoError(t, err)

		var newCampaign *Campaign
		var response *StandardResponse
		newCampaign, response, err = client.GetCampaign(campaign.ID)
		assert.NoError(t, err)
		assert.NotNil(t, newCampaign)
		assert.NotNil(t, response)
		assert.Equal(t, testCampaignID, newCampaign.ID)
	})

	t.Run("missing campaign id", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()
		campaign.ID = 0

		endpoint := fmt.Sprintf(
			"%s/%s/details/?%s=%d", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldID, campaign.ID,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, campaign)
		assert.NoError(t, err)

		var newCampaign *Campaign
		var response *StandardResponse
		newCampaign, response, err = client.GetCampaign(campaign.ID)
		assert.Error(t, err)
		assert.Nil(t, newCampaign)
		assert.Nil(t, response)
	})

	t.Run("error from api (status code)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()

		endpoint := fmt.Sprintf(
			"%s/%s/details/?%s=%d", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldID, campaign.ID,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusBadRequest, campaign)
		assert.NoError(t, err)

		var newCampaign *Campaign
		newCampaign, _, err = client.GetCampaign(campaign.ID)
		assert.Error(t, err)
		assert.Nil(t, newCampaign)
	})

	t.Run("error from api (api error)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()

		endpoint := fmt.Sprintf(
			"%s/%s/details/?%s=%d", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldID, campaign.ID,
		)

		apiError := &Error{
			Code:        400,
			Data:        "field_name",
			IPAddress:   "127.0.0.1",
			Message:     "some error message",
			Method:      http.MethodPut,
			RequestGUID: "7f3d97a8fd67ff57861904df6118dcc8",
			StatusCode:  http.StatusBadRequest,
			URL:         endpoint,
		}

		err = mockResponseData(http.MethodGet, endpoint, http.StatusBadRequest, apiError)
		assert.NoError(t, err)

		var newCampaign *Campaign
		newCampaign, _, err = client.GetCampaign(campaign.ID)
		assert.Error(t, err)
		assert.Nil(t, newCampaign)
		assert.Equal(t, apiError.Message, err.Error())
	})
}

// ExampleClient_GetCampaign example using GetCampaign()
//
// See more examples in /examples/
func ExampleClient_GetCampaign() {

	// Load the client (using test client for example only)
	client, err := newTestClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// Mock response (for example only)
	responseCampaign := newTestCampaign()
	_ = mockResponseData(
		http.MethodGet,
		fmt.Sprintf(
			"%s/%s/details/?%s=%d", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldID, responseCampaign.ID,
		),
		http.StatusOK,
		responseCampaign,
	)

	// Get campaign (using mocking response)
	if responseCampaign, _, err = client.GetCampaign(responseCampaign.ID); err != nil {
		fmt.Printf("error getting campaign: " + err.Error())
		return
	}
	fmt.Printf("campaign: %s", responseCampaign.Title)
	// Output:campaign: TonicPow
}

// BenchmarkClient_GetCampaign benchmarks the method GetCampaign()
func BenchmarkClient_GetCampaign(b *testing.B) {
	client, _ := newTestClient()
	campaign := newTestCampaign()
	_ = mockResponseData(
		http.MethodGet,
		fmt.Sprintf(
			"%s/%s/details/?%s=%d", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldID, campaign.ID,
		),
		http.StatusOK,
		campaign,
	)
	for i := 0; i < b.N; i++ {
		_, _, _ = client.GetCampaign(campaign.ID)
	}
}

// TestClient_GetCampaignBySlug will test the method GetCampaignBySlug()
func TestClient_GetCampaignBySlug(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("get a campaign (success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()

		endpoint := fmt.Sprintf(
			"%s/%s/details/?%s=%s", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldSlug, campaign.Slug,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, campaign)
		assert.NoError(t, err)

		var newCampaign *Campaign
		var response *StandardResponse
		newCampaign, response, err = client.GetCampaignBySlug(campaign.Slug)
		assert.NoError(t, err)
		assert.NotNil(t, newCampaign)
		assert.NotNil(t, response)
		assert.Equal(t, testCampaignID, newCampaign.ID)
	})

	t.Run("missing campaign slug", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()
		campaign.Slug = ""

		endpoint := fmt.Sprintf(
			"%s/%s/details/?%s=%s", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldSlug, campaign.Slug,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, campaign)
		assert.NoError(t, err)

		var newCampaign *Campaign
		var response *StandardResponse
		newCampaign, response, err = client.GetCampaignBySlug(campaign.Slug)
		assert.Error(t, err)
		assert.Nil(t, newCampaign)
		assert.Nil(t, response)
	})

	t.Run("error from api (status code)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()

		endpoint := fmt.Sprintf(
			"%s/%s/details/?%s=%s", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldSlug, campaign.Slug,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusBadRequest, campaign)
		assert.NoError(t, err)

		var newCampaign *Campaign
		newCampaign, _, err = client.GetCampaignBySlug(campaign.Slug)
		assert.Error(t, err)
		assert.Nil(t, newCampaign)
	})

	t.Run("error from api (api error)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()

		endpoint := fmt.Sprintf(
			"%s/%s/details/?%s=%s", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldSlug, campaign.Slug,
		)

		apiError := &Error{
			Code:        400,
			Data:        "field_name",
			IPAddress:   "127.0.0.1",
			Message:     "some error message",
			Method:      http.MethodPut,
			RequestGUID: "7f3d97a8fd67ff57861904df6118dcc8",
			StatusCode:  http.StatusBadRequest,
			URL:         endpoint,
		}

		err = mockResponseData(http.MethodGet, endpoint, http.StatusBadRequest, apiError)
		assert.NoError(t, err)

		var newCampaign *Campaign
		newCampaign, _, err = client.GetCampaignBySlug(campaign.Slug)
		assert.Error(t, err)
		assert.Nil(t, newCampaign)
		assert.Equal(t, apiError.Message, err.Error())
	})
}

// ExampleClient_GetCampaignBySlug example using GetCampaignBySlug()
//
// See more examples in /examples/
func ExampleClient_GetCampaignBySlug() {

	// Load the client (using test client for example only)
	client, err := newTestClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// Mock response (for example only)
	responseCampaign := newTestCampaign()
	_ = mockResponseData(
		http.MethodGet,
		fmt.Sprintf(
			"%s/%s/details/?%s=%s", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldSlug, responseCampaign.Slug,
		),
		http.StatusOK,
		responseCampaign,
	)

	// Get campaign (using mocking response)
	if responseCampaign, _, err = client.GetCampaignBySlug(
		responseCampaign.Slug,
	); err != nil {
		fmt.Printf("error getting campaign: " + err.Error())
		return
	}
	fmt.Printf("campaign: %s", responseCampaign.Title)
	// Output:campaign: TonicPow
}

// BenchmarkClient_GetCampaignBySlug benchmarks the method GetCampaignBySlug()
func BenchmarkClient_GetCampaignBySlug(b *testing.B) {
	client, _ := newTestClient()
	campaign := newTestCampaign()
	_ = mockResponseData(
		http.MethodGet,
		fmt.Sprintf(
			"%s/%s/details/?%s=%s", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldSlug, campaign.Slug,
		),
		http.StatusOK,
		campaign,
	)
	for i := 0; i < b.N; i++ {
		_, _, _ = client.GetCampaignBySlug(campaign.Slug)
	}
}

// TestClient_UpdateCampaign will test the method UpdateCampaign()
func TestClient_UpdateCampaign(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("update a campaign (success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()
		campaign.Title = "TonicPow Title"

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelCampaign)

		err = mockResponseData(http.MethodPut, endpoint, http.StatusOK, campaign)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.UpdateCampaign(campaign)
		assert.NoError(t, err)
		assert.NotNil(t, campaign)
		assert.NotNil(t, response)
		assert.Equal(t, "TonicPow Title", campaign.Title)
	})

	t.Run("missing campaign id", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()
		campaign.ID = 0

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelCampaign)

		err = mockResponseData(http.MethodPut, endpoint, http.StatusOK, campaign)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.UpdateCampaign(campaign)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("error from api (status code)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelCampaign)

		err = mockResponseData(http.MethodPut, endpoint, http.StatusBadRequest, campaign)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.UpdateCampaign(campaign)
		assert.Error(t, err)
		assert.NotNil(t, response)
	})

	t.Run("error from api (api error)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		campaign := newTestCampaign()

		endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelCampaign)

		apiError := &Error{
			Code:        400,
			Data:        "field_name",
			IPAddress:   "127.0.0.1",
			Message:     "some error message",
			Method:      http.MethodPut,
			RequestGUID: "7f3d97a8fd67ff57861904df6118dcc8",
			StatusCode:  http.StatusBadRequest,
			URL:         endpoint,
		}

		err = mockResponseData(http.MethodPut, endpoint, http.StatusBadRequest, apiError)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.UpdateCampaign(campaign)
		assert.Error(t, err)
		assert.Equal(t, apiError.Message, err.Error())
		assert.NotNil(t, response)
	})
}

// ExampleClient_UpdateCampaign example using UpdateCampaign()
//
// See more examples in /examples/
func ExampleClient_UpdateCampaign() {

	// Load the client (using test client for example only)
	client, err := newTestClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// Mock response (for example only)
	responseCampaign := newTestCampaign()
	responseCampaign.Title = "TonicPow Title"
	_ = mockResponseData(
		http.MethodPut,
		fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelCampaign),
		http.StatusOK,
		responseCampaign,
	)

	// Update campaign (using mocking response)
	_, err = client.UpdateCampaign(responseCampaign)
	if err != nil {
		fmt.Printf("error updating campaign: " + err.Error())
		return
	}
	fmt.Printf("campaign: %s", responseCampaign.Title)
	// Output:campaign: TonicPow Title
}

// BenchmarkClient_UpdateCampaign benchmarks the method UpdateCampaign()
func BenchmarkClient_UpdateCampaign(b *testing.B) {
	client, _ := newTestClient()
	campaign := newTestCampaign()
	_ = mockResponseData(
		http.MethodPut,
		fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelCampaign),
		http.StatusOK,
		campaign,
	)
	for i := 0; i < b.N; i++ {
		_, _ = client.UpdateCampaign(campaign)
	}
}

// TestClient_ListCampaigns will test the method ListCampaigns()
func TestClient_ListCampaigns(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("list campaigns (success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignResults(1, 25)

		endpoint := fmt.Sprintf(
			"%s/%s/list?%s=%d&%s=%d&%s=%s&%s=%s&%s=%s&%s=%d&%s=%t",
			EnvironmentDevelopment.apiURL, modelCampaign,
			fieldCurrentPage, 1,
			fieldResultsPerPage, 25,
			fieldSortBy, SortByFieldBalance,
			fieldSortOrder, SortOrderAsc,
			fieldSearchQuery, "",
			fieldMinimumBalance, 0,
			fieldExpired, false,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, results)
		assert.NoError(t, err)

		var newResults *CampaignResults
		var response *StandardResponse
		newResults, response, err = client.ListCampaigns(
			1, 25, SortByFieldBalance, SortOrderAsc, "", 0, false,
		)
		assert.NoError(t, err)
		assert.NotNil(t, newResults)
		assert.NotNil(t, response)
		assert.Equal(t, testCampaignID, newResults.Campaigns[0].ID)
	})

	t.Run("list campaigns (default sorting)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignResults(1, 25)

		endpoint := fmt.Sprintf(
			"%s/%s/list?%s=%d&%s=%d&%s=%s&%s=%s&%s=%s&%s=%d&%s=%t",
			EnvironmentDevelopment.apiURL, modelCampaign,
			fieldCurrentPage, 1,
			fieldResultsPerPage, 25,
			fieldSortBy, SortByFieldCreatedAt,
			fieldSortOrder, SortOrderDesc,
			fieldSearchQuery, "",
			fieldMinimumBalance, 0,
			fieldExpired, false,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, results)
		assert.NoError(t, err)

		var newResults *CampaignResults
		var response *StandardResponse
		newResults, response, err = client.ListCampaigns(
			1, 25, "", "", "", 0, false,
		)
		assert.NoError(t, err)
		assert.NotNil(t, newResults)
		assert.NotNil(t, response)
		assert.Equal(t, testCampaignID, newResults.Campaigns[0].ID)
	})

	t.Run("invalid sort by", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignResults(1, 25)

		endpoint := fmt.Sprintf(
			"%s/%s/list?%s=%d&%s=%d&%s=%s&%s=%s&%s=%s&%s=%d&%s=%t",
			EnvironmentDevelopment.apiURL, modelCampaign,
			fieldCurrentPage, 1,
			fieldResultsPerPage, 25,
			fieldSortBy, SortByFieldBalance,
			fieldSortOrder, SortOrderDesc,
			fieldSearchQuery, "",
			fieldMinimumBalance, 0,
			fieldExpired, false,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, results)
		assert.NoError(t, err)

		var newResults *CampaignResults
		newResults, _, err = client.ListCampaigns(
			1, 25, "bad_field", SortOrderDesc, "", 0, false,
		)
		assert.Error(t, err)
		assert.Nil(t, newResults)
	})

	t.Run("error from api (status code)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignResults(2, 5)

		endpoint := fmt.Sprintf(
			"%s/%s/list?%s=%d&%s=%d&%s=%s&%s=%s&%s=%s&%s=%d&%s=%t",
			EnvironmentDevelopment.apiURL, modelCampaign,
			fieldCurrentPage, 2,
			fieldResultsPerPage, 5,
			fieldSortBy, SortByFieldBalance,
			fieldSortOrder, SortOrderDesc,
			fieldSearchQuery, "",
			fieldMinimumBalance, 0,
			fieldExpired, false,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusBadRequest, results)
		assert.NoError(t, err)

		var newResults *CampaignResults
		newResults, _, err = client.ListCampaigns(
			2, 5, SortByFieldBalance, SortOrderDesc, "", 0, false,
		)
		assert.Error(t, err)
		assert.Nil(t, newResults)
	})

	t.Run("error from api (api error)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		// results := newTestCampaignResults(2, 5)

		endpoint := fmt.Sprintf(
			"%s/%s/list?%s=%d&%s=%d&%s=%s&%s=%s&%s=%s&%s=%d&%s=%t",
			EnvironmentDevelopment.apiURL, modelCampaign,
			fieldCurrentPage, 2,
			fieldResultsPerPage, 5,
			fieldSortBy, SortByFieldBalance,
			fieldSortOrder, SortOrderDesc,
			fieldSearchQuery, "",
			fieldMinimumBalance, 0,
			fieldExpired, false,
		)

		apiError := &Error{
			Code:        400,
			Data:        "field_name",
			IPAddress:   "127.0.0.1",
			Message:     "some error message",
			Method:      http.MethodPut,
			RequestGUID: "7f3d97a8fd67ff57861904df6118dcc8",
			StatusCode:  http.StatusBadRequest,
			URL:         endpoint,
		}

		err = mockResponseData(http.MethodGet, endpoint, http.StatusBadRequest, apiError)
		assert.NoError(t, err)

		var newResults *CampaignResults
		newResults, _, err = client.ListCampaigns(
			2, 5, SortByFieldBalance, SortOrderDesc, "", 0, false,
		)
		assert.Error(t, err)
		assert.Nil(t, newResults)
		assert.Equal(t, apiError.Message, err.Error())
	})
}

// TestClient_ListCampaignsByURL will test the method ListCampaignsByURL()
func TestClient_ListCampaignsByURL(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("list campaigns by url (success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignResults(1, 25)

		endpoint := fmt.Sprintf("%s/%s/list?%s=%s&%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL, modelCampaign,
			fieldTargetURL, testCampaignTargetURL,
			fieldCurrentPage, 1,
			fieldResultsPerPage, 25,
			fieldSortBy, SortByFieldBalance,
			fieldSortOrder, SortOrderDesc,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, results)
		assert.NoError(t, err)

		var newResults *CampaignResults
		newResults, _, err = client.ListCampaignsByURL(
			testCampaignTargetURL, 1, 25, SortByFieldBalance, SortOrderDesc,
		)
		assert.NoError(t, err)
		assert.NotNil(t, newResults)
		assert.Equal(t, testCampaignID, newResults.Campaigns[0].ID)
	})

	t.Run("list campaigns (default sorting)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignResults(2, 5)

		endpoint := fmt.Sprintf("%s/%s/list?%s=%s&%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL, modelCampaign,
			fieldTargetURL, testCampaignTargetURL,
			fieldCurrentPage, 2,
			fieldResultsPerPage, 5,
			fieldSortBy, SortByFieldCreatedAt,
			fieldSortOrder, SortOrderDesc,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, results)
		assert.NoError(t, err)

		var newResults *CampaignResults
		newResults, _, err = client.ListCampaignsByURL(
			testCampaignTargetURL, 2, 5, "", "",
		)
		assert.NoError(t, err)
		assert.NotNil(t, newResults)
		assert.Equal(t, testCampaignID, newResults.Campaigns[0].ID)
	})

	t.Run("invalid sort by", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignResults(1, 25)

		endpoint := fmt.Sprintf("%s/%s/list?%s=%s&%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL, modelCampaign,
			fieldTargetURL, testCampaignTargetURL,
			fieldCurrentPage, 2,
			fieldResultsPerPage, 5,
			fieldSortBy, SortByFieldCreatedAt,
			fieldSortOrder, SortOrderDesc,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, results)
		assert.NoError(t, err)

		var newResults *CampaignResults
		newResults, _, err = client.ListCampaignsByURL(
			testCampaignTargetURL, 2, 5, "bad_field", "",
		)
		assert.Error(t, err)
		assert.Nil(t, newResults)
	})

	t.Run("invalid target url", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignResults(1, 25)

		endpoint := fmt.Sprintf("%s/%s/list?%s=%s&%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL, modelCampaign,
			fieldTargetURL, testCampaignTargetURL,
			fieldCurrentPage, 2,
			fieldResultsPerPage, 5,
			fieldSortBy, SortByFieldCreatedAt,
			fieldSortOrder, SortOrderDesc,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, results)
		assert.NoError(t, err)

		var newResults *CampaignResults
		newResults, _, err = client.ListCampaignsByURL(
			"", 2, 5, SortByFieldCreatedAt, SortOrderDesc,
		)
		assert.Error(t, err)
		assert.Nil(t, newResults)
	})

	t.Run("error from api (status code)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignResults(2, 5)

		endpoint := fmt.Sprintf("%s/%s/list?%s=%s&%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL, modelCampaign,
			fieldTargetURL, testCampaignTargetURL,
			fieldCurrentPage, 2,
			fieldResultsPerPage, 5,
			fieldSortBy, SortByFieldCreatedAt,
			fieldSortOrder, SortOrderDesc,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusBadRequest, results)
		assert.NoError(t, err)

		var newResults *CampaignResults
		newResults, _, err = client.ListCampaignsByURL(
			testCampaignTargetURL, 2, 5, SortByFieldCreatedAt, SortOrderDesc,
		)
		assert.Error(t, err)
		assert.Nil(t, newResults)
	})

	t.Run("error from api (api error)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		endpoint := fmt.Sprintf("%s/%s/list?%s=%s&%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL, modelCampaign,
			fieldTargetURL, testCampaignTargetURL,
			fieldCurrentPage, 2,
			fieldResultsPerPage, 5,
			fieldSortBy, SortByFieldCreatedAt,
			fieldSortOrder, SortOrderDesc,
		)

		apiError := &Error{
			Code:        400,
			Data:        "field_name",
			IPAddress:   "127.0.0.1",
			Message:     "some error message",
			Method:      http.MethodPut,
			RequestGUID: "7f3d97a8fd67ff57861904df6118dcc8",
			StatusCode:  http.StatusBadRequest,
			URL:         endpoint,
		}

		err = mockResponseData(http.MethodGet, endpoint, http.StatusBadRequest, apiError)
		assert.NoError(t, err)

		var newResults *CampaignResults
		newResults, _, err = client.ListCampaignsByURL(
			testCampaignTargetURL, 2, 5, SortByFieldCreatedAt, SortOrderDesc,
		)
		assert.Error(t, err)
		assert.Nil(t, newResults)
		assert.Equal(t, apiError.Message, err.Error())
	})
}

// TestClient_CampaignsFeed will test the method CampaignsFeed()
func TestClient_CampaignsFeed(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("campaigns feeds (rss success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignFeedRSS()

		endpoint := fmt.Sprintf(
			"%s/%s/feed/?%s=%s", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldFeedType, FeedTypeRSS,
		)

		mockResponseFeed(endpoint, http.StatusOK, results)

		var feedResults string
		var response *StandardResponse
		feedResults, response, err = client.CampaignsFeed(FeedTypeRSS)
		assert.NoError(t, err)
		assert.Equal(t, results, feedResults)
		assert.NotNil(t, response)
	})

	t.Run("campaigns feeds (atom success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignFeedAtom()

		endpoint := fmt.Sprintf(
			"%s/%s/feed/?%s=%s", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldFeedType, FeedTypeAtom,
		)

		mockResponseFeed(endpoint, http.StatusOK, results)

		var feedResults string
		feedResults, _, err = client.CampaignsFeed(FeedTypeAtom)
		assert.NoError(t, err)
		assert.Equal(t, results, feedResults)
	})

	t.Run("campaigns feeds (json success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignFeedJSON()

		endpoint := fmt.Sprintf(
			"%s/%s/feed/?%s=%s", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldFeedType, FeedTypeJSON,
		)

		mockResponseFeed(endpoint, http.StatusOK, results)

		var feedResults string
		feedResults, _, err = client.CampaignsFeed(FeedTypeJSON)
		assert.NoError(t, err)
		assert.Equal(t, results, feedResults)
	})

	t.Run("error from api (status code)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignFeedJSON()

		endpoint := fmt.Sprintf(
			"%s/%s/feed/?%s=%s", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldFeedType, FeedTypeJSON,
		)

		mockResponseFeed(endpoint, http.StatusBadRequest, results)

		var feedResults string
		feedResults, _, err = client.CampaignsFeed(FeedTypeJSON)
		assert.Error(t, err)
		assert.Equal(t, "", feedResults)
	})

	t.Run("error from api (api error)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		// results := newTestCampaignFeedJSON()

		endpoint := fmt.Sprintf(
			"%s/%s/feed/?%s=%s", EnvironmentDevelopment.apiURL,
			modelCampaign, fieldFeedType, FeedTypeJSON,
		)

		apiError := &Error{
			Code:        400,
			Data:        "field_name",
			IPAddress:   "127.0.0.1",
			Message:     "some error message",
			Method:      http.MethodPut,
			RequestGUID: "7f3d97a8fd67ff57861904df6118dcc8",
			StatusCode:  http.StatusBadRequest,
			URL:         endpoint,
		}

		err = mockResponseData(http.MethodGet, endpoint, http.StatusBadRequest, apiError)
		assert.NoError(t, err)

		var feedResults string
		feedResults, _, err = client.CampaignsFeed(FeedTypeJSON)
		assert.Error(t, err)
		assert.Equal(t, "", feedResults)
		assert.Equal(t, apiError.Message, err.Error())
	})
}
