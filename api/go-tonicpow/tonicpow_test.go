package tonicpow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const (
	testAdvertiserID      uint64 = 23
	testAdvertiserName    string = "TonicPow Test"
	testAPIKey            string = "TestAPIKey12345678987654321"
	testAppID             uint64 = 10
	testCampaignID        uint64 = 23
	testCampaignTargetURL string = "https://tonicpow.com"
	testConversionID      uint64 = 99
	testGoalID            uint64 = 13
	testRateCurrency      string = "usd"
	testTncpwSession      string = "TestSessionKey12345678987654321"
	testGoalName          string = "example_goal"
	testUserID            uint64 = 43
)

// TestVersion will test the method Version()
func TestVersion(t *testing.T) {
	t.Parallel()

	t.Run("get version", func(t *testing.T) {
		ver := Version()
		assert.Equal(t, version, ver)
	})
}

// ExampleVersion example using Version()
//
// See more examples in /examples/
func ExampleVersion() {
	fmt.Printf("version: %s", Version())
	// Output:version: v0.6.8
}

// TestUserAgent will test the method UserAgent()
func TestUserAgent(t *testing.T) {
	t.Parallel()

	t.Run("get user agent", func(t *testing.T) {
		agent := UserAgent()
		assert.Equal(t, defaultUserAgent, agent)
	})
}

// ExampleUserAgent example using UserAgent()
//
// See more examples in /examples/
func ExampleUserAgent() {
	fmt.Printf("user agent: %s", UserAgent())
	// Output:user agent: go-tonicpow: v0.6.8
}

// TestGetFeedType will test the method GetFeedType()
func TestGetFeedType(t *testing.T) {
	t.Run("test valid cases", func(t *testing.T) {
		assert.Equal(t, FeedTypeRSS, GetFeedType("rss"))
		assert.Equal(t, FeedTypeJSON, GetFeedType("json"))
		assert.Equal(t, FeedTypeAtom, GetFeedType("atom"))
		assert.Equal(t, FeedTypeRSS, GetFeedType(""))
	})
}

// TestEnvironment_Alias will test the method Alias()
func TestEnvironment_Alias(t *testing.T) {
	assert.Equal(t, environmentStagingAlias, EnvironmentStaging.Alias())
	assert.Equal(t, environmentDevelopmentAlias, EnvironmentDevelopment.Alias())
	assert.Equal(t, environmentLiveAlias, EnvironmentLive.Alias())
	e := Environment{}
	assert.Equal(t, "", e.Alias())
}

// ExampleEnvironment_Alias example using Alias()
//
// See more examples in /examples/
func ExampleEnvironment_Alias() {
	env := EnvironmentLive
	fmt.Printf("name: %s alias: %s", env.Name(), env.Alias())
	// Output:name: live alias: production
}

// TestEnvironment_Name will test the method Name()
func TestEnvironment_Name(t *testing.T) {
	assert.Equal(t, environmentStagingName, EnvironmentStaging.Name())
	assert.Equal(t, environmentDevelopmentName, EnvironmentDevelopment.Name())
	assert.Equal(t, environmentLiveName, EnvironmentLive.Name())
	e := Environment{}
	assert.Equal(t, "", e.Name())
}

// ExampleEnvironment_Name example using Name()
//
// See more examples in /examples/
func ExampleEnvironment_Name() {
	env := EnvironmentLive
	fmt.Printf("name: %s alias: %s", env.Name(), env.Alias())
	// Output:name: live alias: production
}

// TestEnvironment_URL will test the method URL()
func TestEnvironment_URL(t *testing.T) {
	assert.Equal(t, stagingAPIURL, EnvironmentStaging.URL())
	assert.Equal(t, developmentURL, EnvironmentDevelopment.URL())
	assert.Equal(t, liveAPIURL, EnvironmentLive.URL())
	e := Environment{}
	assert.Equal(t, "", e.URL())
}

// ExampleEnvironment_URL example using URL()
//
// See more examples in /examples/
func ExampleEnvironment_URL() {
	env := EnvironmentLive
	fmt.Printf("name: %s url: %s", env.Name(), env.URL())
	// Output:name: live url: https://api.tonicpow.com/v1
}

// mockResponseData is used for mocking the response
func mockResponseData(method, endpoint string, statusCode int, model interface{}) error {
	httpmock.Reset()
	if model != nil && model != "" {
		data, err := json.Marshal(model)
		if err != nil {
			return err
		}
		httpmock.RegisterResponder(method, endpoint, httpmock.NewStringResponder(statusCode, string(data)))
	} else {
		httpmock.RegisterResponder(method, endpoint, httpmock.NewStringResponder(statusCode, ""))
	}

	return nil
}

// mockResponseFeed is used for mocking the response
func mockResponseFeed(endpoint string, statusCode int, feedResults string) {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, endpoint, httpmock.NewStringResponder(statusCode, feedResults))
}
