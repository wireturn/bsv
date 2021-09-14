package tonicpow

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestAdvertiserProfile creates a dummy profile for testing
func newTestAdvertiserProfile() *AdvertiserProfile {
	return &AdvertiserProfile{
		HomepageURL: "https://tonicpow.com",
		IconURL:     "https://i.imgur.com/HvVmeWI.png",
		PublicGUID:  "a4503e16b25c29b9cf58eee3ad353410",
		Name:        "TonicPow",
		ID:          testAdvertiserID,
		UserID:      testUserID,
	}
}

// newTestAppResults creates a dummy profile for testing
func newTestAppResults(currentPage, resultsPerPage int) *AppResults {
	return &AppResults{
		Apps:           []*App{newTestApp()},
		CurrentPage:    currentPage,
		Results:        1,
		ResultsPerPage: resultsPerPage,
	}
}

// newTestApp creates a dummy profile for testing
func newTestApp() *App {
	return &App{
		AdvertiserProfileID: testAdvertiserID,
		ID:                  testAppID,
		Name:                "TonicPow App",
		UserID:              testUserID,
	}
}

// TestClient_GetAdvertiserProfile will test the method GetAdvertiserProfile()
func TestClient_GetAdvertiserProfile(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("get an advertiser (success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		profile := newTestAdvertiserProfile()
		profile.Name = testAdvertiserName

		endpoint := fmt.Sprintf("%s/%s/details/%d", EnvironmentDevelopment.apiURL, modelAdvertiser, profile.ID)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, profile)
		assert.NoError(t, err)

		var response *StandardResponse
		profile, response, err = client.GetAdvertiserProfile(profile.ID)
		assert.NoError(t, err)
		assert.NotNil(t, profile)
		assert.NotNil(t, response)
		assert.Equal(t, testAdvertiserName, profile.Name)
	})

	t.Run("missing advertiser profile id", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		profile := newTestAdvertiserProfile()
		profile.ID = 0

		endpoint := fmt.Sprintf("%s/%s/details/%d", EnvironmentDevelopment.apiURL, modelAdvertiser, profile.ID)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, profile)
		assert.NoError(t, err)

		var realProfile *AdvertiserProfile
		var response *StandardResponse
		realProfile, response, err = client.GetAdvertiserProfile(profile.ID)
		assert.Error(t, err)
		assert.Nil(t, realProfile)
		assert.Nil(t, response)
	})

	t.Run("error from api (status code)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		profile := newTestAdvertiserProfile()
		endpoint := fmt.Sprintf("%s/%s/details/%d", EnvironmentDevelopment.apiURL, modelAdvertiser, profile.ID)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusBadRequest, profile)
		assert.NoError(t, err)

		var realProfile *AdvertiserProfile
		realProfile, _, err = client.GetAdvertiserProfile(profile.ID)
		assert.Error(t, err)
		assert.Nil(t, realProfile)
	})

	t.Run("error from api (api error)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		profile := newTestAdvertiserProfile()

		endpoint := fmt.Sprintf("%s/%s/details/%d", EnvironmentDevelopment.apiURL, modelAdvertiser, profile.ID)

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

		var realProfile *AdvertiserProfile
		realProfile, _, err = client.GetAdvertiserProfile(profile.ID)
		assert.Error(t, err)
		assert.Nil(t, realProfile)
		assert.Equal(t, apiError.Message, err.Error())
	})
}

// ExampleClient_GetAdvertiserProfile example using GetAdvertiserProfile()
//
// See more examples in /examples/
func ExampleClient_GetAdvertiserProfile() {

	// Load the client (using test client for example only)
	client, err := newTestClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// For mocking
	responseProfile := newTestAdvertiserProfile()

	// Mock response (for example only)
	_ = mockResponseData(
		http.MethodGet,
		fmt.Sprintf("%s/%s/details/%d", EnvironmentDevelopment.apiURL, modelAdvertiser, responseProfile.ID),
		http.StatusOK,
		responseProfile,
	)

	// Get profile (using mocking response)
	var profile *AdvertiserProfile
	if profile, _, err = client.GetAdvertiserProfile(testAdvertiserID); err != nil {
		fmt.Printf("error getting profile: " + err.Error())
		return
	}
	fmt.Printf("advertiser profile: %s", profile.Name)
	// Output:advertiser profile: TonicPow
}

// BenchmarkClient_GetAdvertiserProfile benchmarks the method GetAdvertiserProfile()
func BenchmarkClient_GetAdvertiserProfile(b *testing.B) {
	client, _ := newTestClient()
	profile := newTestAdvertiserProfile()
	_ = mockResponseData(
		http.MethodGet,
		fmt.Sprintf("%s/%s/details/%d", EnvironmentDevelopment.apiURL, modelAdvertiser, profile.ID),
		http.StatusOK,
		profile,
	)
	for i := 0; i < b.N; i++ {
		_, _, _ = client.GetAdvertiserProfile(profile.ID)
	}
}

// TestClient_UpdateAdvertiserProfile will test the method UpdateAdvertiserProfile()
func TestClient_UpdateAdvertiserProfile(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	endpoint := fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelAdvertiser)

	t.Run("update an advertiser (success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		profile := newTestAdvertiserProfile()

		profile.Name = testAdvertiserName
		err = mockResponseData(http.MethodPut, endpoint, http.StatusOK, profile)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.UpdateAdvertiserProfile(profile)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, testAdvertiserName, profile.Name)
	})

	t.Run("missing advertiser profile id", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		profile := newTestAdvertiserProfile()

		profile.ID = 0
		err = mockResponseData(http.MethodPut, endpoint, http.StatusOK, profile)
		assert.NoError(t, err)

		var response *StandardResponse
		response, err = client.UpdateAdvertiserProfile(profile)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("error from api (status code)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		profile := newTestAdvertiserProfile()

		err = mockResponseData(http.MethodPut, endpoint, http.StatusBadRequest, profile)
		assert.NoError(t, err)

		_, err = client.UpdateAdvertiserProfile(profile)
		assert.Error(t, err)
	})

	t.Run("error from api (api error)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		profile := newTestAdvertiserProfile()

		apiError := &Error{
			Code:        400,
			Data:        "field_name",
			IPAddress:   "127.0.0.1",
			Message:     "some error message",
			Method:      http.MethodPut,
			RequestGUID: "7f3d97a8fd67ff57861904df6118dcc8",
			StatusCode:  http.StatusBadRequest,
			URL:         fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelAdvertiser),
		}

		err = mockResponseData(http.MethodPut, endpoint, http.StatusBadRequest, apiError)
		assert.NoError(t, err)

		_, err = client.UpdateAdvertiserProfile(profile)
		assert.Error(t, err)
		assert.Equal(t, apiError.Message, err.Error())
	})
}

// ExampleClient_UpdateAdvertiserProfile example using UpdateAdvertiserProfile()
//
// See more examples in /examples/
func ExampleClient_UpdateAdvertiserProfile() {

	// Load the client (using test client for example only)
	client, err := newTestClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	// Start with an existing profile
	profile := newTestAdvertiserProfile()
	profile.Name = testAdvertiserName

	// Mock response (for example only)
	_ = mockResponseData(
		http.MethodPut,
		fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelAdvertiser),
		http.StatusOK,
		profile,
	)

	// Update profile
	_, err = client.UpdateAdvertiserProfile(profile)
	if err != nil {
		fmt.Printf("error updating profile: " + err.Error())
		return
	}
	fmt.Printf("profile updated: %s", profile.Name)
	// Output:profile updated: TonicPow Test
}

// BenchmarkClient_UpdateAdvertiserProfile benchmarks the method UpdateAdvertiserProfile()
func BenchmarkClient_UpdateAdvertiserProfile(b *testing.B) {
	client, _ := newTestClient()
	profile := newTestAdvertiserProfile()
	_ = mockResponseData(
		http.MethodPut,
		fmt.Sprintf("%s/%s", EnvironmentDevelopment.apiURL, modelAdvertiser),
		http.StatusOK,
		profile,
	)
	for i := 0; i < b.N; i++ {
		_, _ = client.UpdateAdvertiserProfile(profile)
	}
}

// TestClient_ListCampaignsByAdvertiserProfile will test the method ListCampaignsByAdvertiserProfile()
func TestClient_ListCampaignsByAdvertiserProfile(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("list campaigns by advertiser (success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignResults(1, 25)

		endpoint := fmt.Sprintf("%s/%s/%s/%d?%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL,
			modelAdvertiser, modelCampaign, results.Campaigns[0].AdvertiserProfileID,
			fieldCurrentPage, 1,
			fieldResultsPerPage, 25,
			fieldSortBy, SortByFieldBalance,
			fieldSortOrder, SortOrderDesc,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, results)
		assert.NoError(t, err)

		var response *StandardResponse
		results, response, err = client.ListCampaignsByAdvertiserProfile(
			results.Campaigns[0].AdvertiserProfileID, 1, 25, SortByFieldBalance, SortOrderDesc,
		)
		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.NotNil(t, response)
		assert.Equal(t, "TonicPow", results.Campaigns[0].AdvertiserProfile.Name)
		assert.Equal(t, testCampaignID, results.Campaigns[0].ID)
		assert.Equal(t, testGoalID, results.Campaigns[0].Goals[0].ID)
		assert.Equal(t, testAdvertiserID, results.Campaigns[0].AdvertiserProfileID)
		assert.Equal(t, testUserID, results.Campaigns[0].AdvertiserProfile.UserID)
		assert.Equal(t, 1, len(results.Campaigns))
		assert.Equal(t, 1, results.Results)
		assert.Equal(t, 25, results.ResultsPerPage)
		assert.Equal(t, 1, results.CurrentPage)
	})

	t.Run("list campaigns by advertiser (default sorting)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignResults(1, 25)

		endpoint := fmt.Sprintf("%s/%s/%s/%d?%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL,
			modelAdvertiser, modelCampaign, results.Campaigns[0].AdvertiserProfileID,
			fieldCurrentPage, 1,
			fieldResultsPerPage, 25,
			fieldSortBy, SortByFieldCreatedAt,
			fieldSortOrder, SortOrderDesc,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, results)
		assert.NoError(t, err)

		var response *StandardResponse
		results, response, err = client.ListCampaignsByAdvertiserProfile(
			results.Campaigns[0].AdvertiserProfileID, 1, 25, "", "",
		)
		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.NotNil(t, response)
		assert.Equal(t, 1, len(results.Campaigns))
		assert.Equal(t, 1, results.Results)
		assert.Equal(t, 25, results.ResultsPerPage)
		assert.Equal(t, 1, results.CurrentPage)
	})

	t.Run("missing profile id", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignResults(2, 5)

		endpoint := fmt.Sprintf("%s/%s/%s/%d?%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL,
			modelAdvertiser, modelCampaign, results.Campaigns[0].AdvertiserProfileID,
			fieldCurrentPage, 2,
			fieldResultsPerPage, 5,
			fieldSortBy, SortByFieldBalance,
			fieldSortOrder, SortOrderDesc,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, results)
		assert.NoError(t, err)

		results, _, err = client.ListCampaignsByAdvertiserProfile(
			0, 2, 5, SortByFieldBalance, SortOrderDesc,
		)
		assert.Error(t, err)
		assert.Nil(t, results)
	})

	t.Run("sort by field is invalid", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignResults(1, 25)

		endpoint := fmt.Sprintf("%s/%s/%s/%d?%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL,
			modelAdvertiser, modelCampaign, results.Campaigns[0].AdvertiserProfileID,
			fieldCurrentPage, 1,
			fieldResultsPerPage, 25,
			fieldSortBy, SortByFieldBalance,
			fieldSortOrder, SortOrderDesc,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, results)
		assert.NoError(t, err)

		results, _, err = client.ListCampaignsByAdvertiserProfile(
			results.Campaigns[0].AdvertiserProfileID, 1, 25, "bad_field_name", SortOrderDesc,
		)
		assert.Error(t, err)
		assert.Nil(t, results)
	})

	t.Run("unexpected status code", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignResults(1, 25)

		endpoint := fmt.Sprintf("%s/%s/%s/%d?%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL,
			modelAdvertiser, modelCampaign, results.Campaigns[0].AdvertiserProfileID,
			fieldCurrentPage, 1,
			fieldResultsPerPage, 25,
			fieldSortBy, SortByFieldBalance,
			fieldSortOrder, SortOrderDesc,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusBadRequest, results)
		assert.NoError(t, err)

		results, _, err = client.ListCampaignsByAdvertiserProfile(
			results.Campaigns[0].AdvertiserProfileID, 1, 25, "bad_field_name", SortOrderDesc,
		)
		assert.Error(t, err)
		assert.Nil(t, results)
	})

	t.Run("api error", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestCampaignResults(1, 25)

		endpoint := fmt.Sprintf("%s/%s/%s/%d?%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL,
			modelAdvertiser, modelCampaign, results.Campaigns[0].AdvertiserProfileID,
			fieldCurrentPage, 1,
			fieldResultsPerPage, 25,
			fieldSortBy, SortByFieldBalance,
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

		results, _, err = client.ListCampaignsByAdvertiserProfile(
			results.Campaigns[0].AdvertiserProfileID, 1, 25, SortByFieldBalance, SortOrderDesc,
		)
		assert.Error(t, err)
		assert.Nil(t, results)
		assert.Equal(t, apiError.Message, err.Error())
	})
}

// TestClient_ListAppsByAdvertiserProfile will test the method ListAppsByAdvertiserProfile()
func TestClient_ListAppsByAdvertiserProfile(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("list apps by advertiser (success)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestAppResults(1, 25)

		endpoint := fmt.Sprintf(
			"%s/%s/%s/?%s=%d&%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL,
			modelAdvertiser, modelApp,
			fieldID, testAdvertiserID,
			fieldCurrentPage, 1,
			fieldResultsPerPage, 25,
			fieldSortBy, SortByFieldName,
			fieldSortOrder, SortOrderDesc,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, results)
		assert.NoError(t, err)

		var response *StandardResponse
		results, response, err = client.ListAppsByAdvertiserProfile(
			results.Apps[0].AdvertiserProfileID, 1, 25, SortByFieldName, SortOrderDesc,
		)
		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.NotNil(t, response)
		assert.Equal(t, "TonicPow App", results.Apps[0].Name)
		assert.Equal(t, testAppID, results.Apps[0].ID)
		assert.Equal(t, 1, len(results.Apps))
		assert.Equal(t, 1, results.Results)
		assert.Equal(t, 25, results.ResultsPerPage)
		assert.Equal(t, 1, results.CurrentPage)
	})

	t.Run("list apps by advertiser (default sorting)", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestAppResults(1, 25)

		endpoint := fmt.Sprintf(
			"%s/%s/%s/?%s=%d&%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL,
			modelAdvertiser, modelApp,
			fieldID, testAdvertiserID,
			fieldCurrentPage, 1,
			fieldResultsPerPage, 25,
			fieldSortBy, SortByFieldCreatedAt,
			fieldSortOrder, SortOrderDesc,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, results)
		assert.NoError(t, err)

		var response *StandardResponse
		results, response, err = client.ListAppsByAdvertiserProfile(
			results.Apps[0].AdvertiserProfileID, 1, 25, "", "",
		)
		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.NotNil(t, response)
		assert.Equal(t, 1, len(results.Apps))
		assert.Equal(t, 1, results.Results)
		assert.Equal(t, 25, results.ResultsPerPage)
		assert.Equal(t, 1, results.CurrentPage)
	})

	t.Run("missing profile id", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestAppResults(2, 5)

		endpoint := fmt.Sprintf(
			"%s/%s/%s/?%s=%d&%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL,
			modelAdvertiser, modelApp,
			fieldID, testAdvertiserID,
			fieldCurrentPage, 2,
			fieldResultsPerPage, 5,
			fieldSortBy, SortByFieldName,
			fieldSortOrder, SortOrderDesc,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, results)
		assert.NoError(t, err)

		var response *StandardResponse
		results, response, err = client.ListAppsByAdvertiserProfile(
			0, 2, 5, SortByFieldName, SortOrderDesc,
		)
		assert.Error(t, err)
		assert.Nil(t, results)
		assert.Nil(t, response)
	})

	t.Run("invalid sort by field", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestAppResults(1, 25)

		endpoint := fmt.Sprintf(
			"%s/%s/%s/?%s=%d&%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL,
			modelAdvertiser, modelApp,
			fieldID, testAdvertiserID,
			fieldCurrentPage, 1,
			fieldResultsPerPage, 25,
			fieldSortBy, SortByFieldName,
			fieldSortOrder, SortOrderDesc,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusOK, results)
		assert.NoError(t, err)

		results, _, err = client.ListAppsByAdvertiserProfile(
			results.Apps[0].AdvertiserProfileID, 1, 25, "bad_field", SortOrderDesc,
		)
		assert.Error(t, err)
		assert.Nil(t, results)
	})

	t.Run("unexpected status code", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestAppResults(1, 25)

		endpoint := fmt.Sprintf(
			"%s/%s/%s/?%s=%d&%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL,
			modelAdvertiser, modelApp,
			fieldID, testAdvertiserID,
			fieldCurrentPage, 1,
			fieldResultsPerPage, 25,
			fieldSortBy, SortByFieldName,
			fieldSortOrder, SortOrderDesc,
		)

		err = mockResponseData(http.MethodGet, endpoint, http.StatusBadRequest, results)
		assert.NoError(t, err)

		results, _, err = client.ListAppsByAdvertiserProfile(
			results.Apps[0].AdvertiserProfileID, 1, 25, SortByFieldName, SortOrderDesc,
		)
		assert.Error(t, err)
		assert.Nil(t, results)
	})

	t.Run("api error", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		results := newTestAppResults(1, 25)

		endpoint := fmt.Sprintf(
			"%s/%s/%s/?%s=%d&%s=%d&%s=%d&%s=%s&%s=%s",
			EnvironmentDevelopment.apiURL,
			modelAdvertiser, modelApp,
			fieldID, testAdvertiserID,
			fieldCurrentPage, 1,
			fieldResultsPerPage, 25,
			fieldSortBy, SortByFieldName,
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

		results, _, err = client.ListAppsByAdvertiserProfile(
			results.Apps[0].AdvertiserProfileID, 1, 25, SortByFieldName, SortOrderDesc,
		)
		assert.Error(t, err)
		assert.Nil(t, results)
		assert.Equal(t, apiError.Message, err.Error())
	})
}
