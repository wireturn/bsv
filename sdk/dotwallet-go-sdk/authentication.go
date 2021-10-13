package dotwallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// UpdateApplicationAccessToken will update the application access token
// This will also store the token on the client for easy access (c.options.token)
//
// For more information: https://developers.dotwallet.com/documents/en/#application-authorization
func (c *Client) UpdateApplicationAccessToken() error {

	// Make the request
	response, err := c.Request(
		http.MethodPost,
		getAccessTokenURI,
		&getAccessTokenRequest{
			ClientID:     c.options.clientID,
			GrantType:    grantTypeClientCredentials,
			ClientSecret: c.options.clientSecret,
		},
		http.StatusOK,
		nil,
	)
	if err != nil {
		return err
	}

	// Unmarshal the response
	resp := new(accessTokenResponse)
	if err = json.Unmarshal(
		response.Body, &resp,
	); err != nil {
		return err
	}

	// Error?
	if resp.Code > 0 {
		return fmt.Errorf(resp.Message)
	}

	// Set the token (access token has limited fields vs user token)
	c.options.token = &DotAccessToken{
		AccessToken: resp.Data.AccessToken,
		ExpiresAt:   time.Now().UTC().Unix() + resp.Data.ExpiresIn,
		ExpiresIn:   resp.Data.ExpiresIn,
		TokenType:   resp.Data.TokenType,
	}

	return nil
}

// GetUserToken will get the user's access_token for the first time given the "code" from the oauth2 callback
// This will also store the user's token on the client for easy access (c.options.userToken)
//
// For more information: https://developers.dotwallet.com/documents/en/#user-authorization
func (c *Client) GetUserToken(code string) (*DotAccessToken, error) {

	// Make the request
	response, err := c.Request(
		http.MethodPost,
		getAccessTokenURI,
		&getDotUserTokenRequest{
			getAccessTokenRequest: getAccessTokenRequest{
				ClientID:     c.options.clientID,
				ClientSecret: c.options.clientSecret,
				GrantType:    grantTypeAuthorizationCode,
			},
			Code:        code,
			RedirectURI: c.options.redirectURI,
		},
		http.StatusOK,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response
	tokenResponse := new(accessTokenResponse)
	if err = json.Unmarshal(
		response.Body, &tokenResponse,
	); err != nil {
		return nil, err
	}

	// Error?
	if tokenResponse.Code > 0 {
		return nil, fmt.Errorf(tokenResponse.Message)
	}

	// Set the user token on the client for easy access
	return newUserAccessToken(tokenResponse), nil
}

// RefreshUserToken will refresh the user's auth_token using the refresh_token
// This will also store the user's token on the client for easy access (c.options.userToken)
//
// For more information: https://developers.dotwallet.com/documents/en/#user-authorization
func (c *Client) RefreshUserToken(token *DotAccessToken) (*DotAccessToken, error) {

	// Make the request
	response, err := c.Request(
		http.MethodPost,
		getAccessTokenURI,
		&refreshDotUserTokenRequest{
			getAccessTokenRequest: getAccessTokenRequest{
				ClientID:     c.options.clientID,
				ClientSecret: c.options.clientSecret,
				GrantType:    grantTypeRefreshToken,
			},
			RefreshToken: token.RefreshToken,
		},
		http.StatusOK,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response
	tokenResponse := new(accessTokenResponse)
	if err = json.Unmarshal(
		response.Body, &tokenResponse,
	); err != nil {
		return nil, err
	}

	// Error?
	if tokenResponse.Code > 0 {
		return nil, fmt.Errorf(tokenResponse.Message)
	}

	// Set the user token on the client for easy access
	return newUserAccessToken(tokenResponse), nil
}

// GetAuthorizeURL will return a new url for starting the user oauth2 authorization process
//
// For more information: https://developers.dotwallet.com/documents/en/#user-authorization
func (c *Client) GetAuthorizeURL(state string, scopes []string) string {
	urlValues := &url.Values{}
	urlValues.Add("client_id", c.options.clientID)
	urlValues.Add("redirect_uri", c.options.redirectURI)
	urlValues.Add("response_type", "code")
	urlValues.Add("scope", strings.Join(scopes, " "))
	urlValues.Add("state", state)
	return fmt.Sprintf("%s%s?%s", c.options.host, authorizeURI, urlValues.Encode())
}

// IsTokenExpired will check if the token is expired
func (c *Client) IsTokenExpired(token *DotAccessToken) bool {
	return token.ExpiresAt < time.Now().UTC().Unix()
}

// newUserAccessToken will create the DotAccessToken object which is enriched from the response data
// This is used only for user access_tokens
func newUserAccessToken(tokenResponse *accessTokenResponse) *DotAccessToken {
	return &DotAccessToken{
		AccessToken:           tokenResponse.Data.AccessToken,
		ExpiresAt:             time.Now().UTC().Unix() + tokenResponse.Data.ExpiresIn,
		ExpiresIn:             tokenResponse.Data.ExpiresIn,
		RefreshToken:          tokenResponse.Data.RefreshToken,
		RefreshTokenExpiresAt: time.Now().UTC().Unix() + int64(defaultRefreshTokenExpiresIn.Seconds()),
		RefreshTokenExpiresIn: int64(defaultRefreshTokenExpiresIn.Seconds()),
		Scopes:                strings.Split(tokenResponse.Data.Scope, " "),
		TokenType:             tokenResponse.Data.TokenType,
	}
}
