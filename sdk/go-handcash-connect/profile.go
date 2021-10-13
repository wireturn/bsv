package handcash

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
{
  "publicProfile": {
    "id": "1234567890",
    "handle": "MisterZ",
    "paymail": "MisterZ@beta.handcash.io",
		"bitcoinUnit": "DUR",
    "displayName": "",
    "avatarUrl": "https://beta-cloud.handcash.io/users/profilePicture/MisterZ",
    "localCurrencyCode": "USD"
  },
  "privateProfile": {
    "phoneNumber": "+15554443333",
    "email": "email@domain.com"
  }
}
*/

// GetProfile will get the profile for the associated auth token
//
// Specs: https://github.com/HandCash/handcash-connect-sdk-js/blob/master/src/profile/index.js
func (c *Client) GetProfile(ctx context.Context, token string) (*Profile, error) {

	// Make sure we have an auth token
	if len(token) == 0 {
		return nil, fmt.Errorf("missing auth token")
	}

	// Get the signed request
	signed, err := c.getSignedRequest(
		http.MethodGet,
		endpointProfileCurrent,
		token,
		&requestBody{authToken: token},
		currentISOTimestamp(),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating signed request: %w", err)
	}

	// Make the HTTP request
	response := httpRequest(
		ctx,
		c,
		&httpPayload{
			Data:           []byte(emptyBody),
			ExpectedStatus: http.StatusOK,
			Method:         signed.Method,
			URL:            signed.URI,
		},
		signed,
	)

	// Error in request?
	if response.Error != nil {
		return nil, response.Error
	}

	// Unmarshal into the profile
	profile := new(Profile)
	if err = json.Unmarshal(response.BodyContents, &profile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	} else if profile == nil || profile.PublicProfile.ID == "" {
		return nil, fmt.Errorf("failed to find profile")
	}
	return profile, nil
}
