package paymail

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// TestClient_GetPublicProfile will test the method GetPublicProfile()
func TestClient_GetPublicProfile(t *testing.T) {
	// t.Parallel() (Cannot run in parallel - issues with overriding the mock client)

	t.Run("successful response", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockGetPublicProfile(http.StatusOK)

		var profile *PublicProfile
		profile, err = client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, testDomain)
		assert.NoError(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, http.StatusOK, profile.StatusCode)
		assert.Equal(t, testName, profile.Name)
		assert.Equal(t, testAvatar, profile.Avatar)
	})

	t.Run("successful response - status not modified", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockGetPublicProfile(http.StatusNotModified)

		var profile *PublicProfile
		profile, err = client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, testDomain)
		assert.NoError(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, http.StatusNotModified, profile.StatusCode)
		assert.Equal(t, testName, profile.Name)
		assert.Equal(t, testAvatar, profile.Avatar)
	})

	t.Run("missing url", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockGetPublicProfile(http.StatusOK)

		var profile *PublicProfile
		profile, err = client.GetPublicProfile("invalid-url", testAlias, testDomain)
		assert.Error(t, err)
		assert.Nil(t, profile)
	})

	t.Run("missing alias", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockGetPublicProfile(http.StatusOK)

		var profile *PublicProfile
		profile, err = client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", "", testDomain)
		assert.Error(t, err)
		assert.Nil(t, profile)
	})

	t.Run("missing domain", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		mockGetPublicProfile(http.StatusOK)

		var profile *PublicProfile
		profile, err = client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, "")
		assert.Error(t, err)
		assert.Nil(t, profile)
	})

	t.Run("bad request", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"public-profile/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": "request failed"}`,
			),
		)

		var profile *PublicProfile
		profile, err = client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, testDomain)
		assert.Error(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, http.StatusBadRequest, profile.StatusCode)
	})

	t.Run("http error", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"public-profile/"+testAlias+"@"+testDomain,
			httpmock.NewErrorResponder(fmt.Errorf("error in request")),
		)

		var profile *PublicProfile
		profile, err = client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, testDomain)
		assert.Error(t, err)
		assert.Nil(t, profile)
	})

	t.Run("error occurred", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"public-profile/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				`{"message": request failed}`,
			),
		)

		var profile *PublicProfile
		profile, err = client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, testDomain)
		assert.Error(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, http.StatusBadRequest, profile.StatusCode)
	})

	t.Run("invalid json", func(t *testing.T) {
		client, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodGet, testServerURL+"public-profile/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusOK,
				`{"name": MrZ,avatar: https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon"}`,
			),
		)

		var profile *PublicProfile
		profile, err = client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, testDomain)
		assert.Error(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, http.StatusOK, profile.StatusCode)
		assert.Equal(t, "", profile.Name)
		assert.Equal(t, "", profile.Avatar)
	})
}

// mockGetPublicProfile is used for mocking the response
func mockGetPublicProfile(statusCode int) {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, testServerURL+"public-profile/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(
			statusCode, `{"name": "`+testName+`","avatar": "`+testAvatar+`"}`,
		),
	)
}

// ExampleClient_GetPublicProfile example using GetPublicProfile()
//
// See more examples in /examples/
func ExampleClient_GetPublicProfile() {
	// Load the client
	client, err := newTestClient()
	if err != nil {
		fmt.Printf("error loading client: %s", err.Error())
		return
	}

	mockGetPublicProfile(http.StatusOK)

	// Get profile
	var profile *PublicProfile
	profile, err = client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, testDomain)
	if err != nil {
		fmt.Printf("error getting profile: " + err.Error())
		return
	}
	fmt.Printf("found profile for: %s", profile.Name)
	// Output:found profile for: MrZ
}

// BenchmarkClient_GetPublicProfile benchmarks the method GetPublicProfile()
func BenchmarkClient_GetPublicProfile(b *testing.B) {
	client, _ := newTestClient()
	mockGetPublicProfile(http.StatusOK)
	for i := 0; i < b.N; i++ {
		_, _ = client.GetPublicProfile(testServerURL+"public-profile/{alias}@{domain.tld}", testAlias, testDomain)
	}
}
