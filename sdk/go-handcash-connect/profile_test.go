package handcash

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockHTTPGetProfile for mocking requests
type mockHTTPGetProfile struct{}

// Do is a mock http request
func (m *mockHTTPGetProfile) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Beta
	if req.URL.String() == environments[EnvironmentBeta].APIURL+endpointProfileCurrent {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"publicProfile":{"id":"1234567","handle":"MisterZ","paymail":"MisterZ@beta.handcash.io","displayName":"","avatarUrl":"https://beta-cloud.handcash.io/users/profilePicture/MisterZ","localCurrencyCode":"USD"},"privateProfile":{"phoneNumber":"+15554443333","email":"email@domain.com"}}`)))
	}

	// IAE
	if req.URL.String() == environments[EnvironmentIAE].APIURL+endpointProfileCurrent {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"publicProfile":{"id":"1234567","handle":"MisterZ","paymail":"MisterZ@beta.handcash.io","displayName":"","avatarUrl":"https://iae-cloud.handcash.io/users/profilePicture/MisterZ","localCurrencyCode":"USD"},"privateProfile":{"phoneNumber":"+15554443333","email":"email@domain.com"}}`)))
	}

	// Production
	if req.URL.String() == environments[EnvironmentProduction].APIURL+endpointProfileCurrent {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"publicProfile":{"id":"1234567","handle":"MisterZ","paymail":"MisterZ@beta.handcash.io","displayName":"","avatarUrl":"https://cloud.handcash.io/users/profilePicture/MisterZ","localCurrencyCode":"USD"},"privateProfile":{"phoneNumber":"+15554443333","email":"email@domain.com"}}`)))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPGetProfile for mocking requests
type mockHTTPBadRequest struct{}

// Do is a mock http request
func (m *mockHTTPBadRequest) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	resp.StatusCode = http.StatusBadRequest
	resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"Message":"bad request"}`)))

	// Default is valid
	return resp, nil
}

// mockHTTPInvalidProfileData for mocking requests
type mockHTTPInvalidProfileData struct{}

// Do is a mock http request
func (m *mockHTTPInvalidProfileData) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	resp.StatusCode = http.StatusOK
	resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"invalid":"profile"}`)))

	// Default is valid
	return resp, nil
}

func TestClient_GetProfile(t *testing.T) {
	t.Parallel()

	t.Run("missing auth token", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetProfile{}, EnvironmentBeta)
		assert.NotNil(t, client)
		profile, err := client.GetProfile(context.Background(), "")
		assert.Error(t, err)
		assert.Nil(t, profile)
	})

	t.Run("invalid auth token", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetProfile{}, EnvironmentBeta)
		assert.NotNil(t, client)
		profile, err := client.GetProfile(context.Background(), "0")
		assert.Error(t, err)
		assert.Nil(t, profile)
	})

	t.Run("valid auth token (hex decodes) (beta)", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetProfile{}, EnvironmentBeta)
		assert.NotNil(t, client)
		profile, err := client.GetProfile(context.Background(), "000000")
		assert.NoError(t, err)
		assert.Equal(t, "1234567", profile.PublicProfile.ID)
		assert.Equal(t, "MisterZ", profile.PublicProfile.Handle)
		assert.Equal(t, "MisterZ@beta.handcash.io", profile.PublicProfile.Paymail)
		assert.Equal(t, "", profile.PublicProfile.DisplayName)
		assert.Equal(t, CurrencyUSD, profile.PublicProfile.LocalCurrencyCode)
		assert.Equal(t, "https://beta-cloud.handcash.io/users/profilePicture/MisterZ", profile.PublicProfile.AvatarURL)
		assert.Equal(t, "+15554443333", profile.PrivateProfile.PhoneNumber)
		assert.Equal(t, "email@domain.com", profile.PrivateProfile.Email)
	})

	t.Run("valid auth token (IAE)", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetProfile{}, EnvironmentIAE)
		assert.NotNil(t, client)
		profile, err := client.GetProfile(context.Background(), "000000")
		assert.NoError(t, err)
		assert.Equal(t, "1234567", profile.PublicProfile.ID)
		assert.Equal(t, "MisterZ", profile.PublicProfile.Handle)
		assert.Equal(t, "MisterZ@beta.handcash.io", profile.PublicProfile.Paymail)
		assert.Equal(t, "", profile.PublicProfile.DisplayName)
		assert.Equal(t, CurrencyUSD, profile.PublicProfile.LocalCurrencyCode)
		assert.Equal(t, "https://iae-cloud.handcash.io/users/profilePicture/MisterZ", profile.PublicProfile.AvatarURL)
		assert.Equal(t, "+15554443333", profile.PrivateProfile.PhoneNumber)
		assert.Equal(t, "email@domain.com", profile.PrivateProfile.Email)
	})

	t.Run("valid auth token (production)", func(t *testing.T) {
		client := newTestClient(&mockHTTPGetProfile{}, EnvironmentProduction)
		assert.NotNil(t, client)
		profile, err := client.GetProfile(context.Background(), "000000")
		assert.NoError(t, err)
		assert.Equal(t, "1234567", profile.PublicProfile.ID)
		assert.Equal(t, "MisterZ", profile.PublicProfile.Handle)
		assert.Equal(t, "MisterZ@beta.handcash.io", profile.PublicProfile.Paymail)
		assert.Equal(t, "", profile.PublicProfile.DisplayName)
		assert.Equal(t, CurrencyUSD, profile.PublicProfile.LocalCurrencyCode)
		assert.Equal(t, "https://cloud.handcash.io/users/profilePicture/MisterZ", profile.PublicProfile.AvatarURL)
		assert.Equal(t, "+15554443333", profile.PrivateProfile.PhoneNumber)
		assert.Equal(t, "email@domain.com", profile.PrivateProfile.Email)
	})

	t.Run("bad request", func(t *testing.T) {
		client := newTestClient(&mockHTTPBadRequest{}, EnvironmentBeta)
		assert.NotNil(t, client)
		profile, err := client.GetProfile(context.Background(), "000000")
		assert.Error(t, err)
		assert.Nil(t, profile)
	})

	t.Run("invalid profile data", func(t *testing.T) {
		client := newTestClient(&mockHTTPInvalidProfileData{}, EnvironmentBeta)
		assert.NotNil(t, client)
		profile, err := client.GetProfile(context.Background(), "000000")
		assert.Error(t, err)
		assert.Nil(t, profile)
	})
}
