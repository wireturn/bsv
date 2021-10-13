package dotwallet

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// TestClient_GetAuthorizeURL will test the method GetAuthorizeURL()
func TestClient_GetAuthorizeURL(t *testing.T) {
	t.Run("generate url from test client", func(t *testing.T) {

		mockUpdateApplicationAccessToken()
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		var state string
		state, err = c.NewState()
		assert.NoError(t, err)
		assert.Equal(t, 64, len(state))

		authURL := c.GetAuthorizeURL(state, []string{
			ScopeUserInfo,
			ScopeAutoPayBSV,
			ScopeAutoPayBTC,
			ScopeAutoPayETH,
		})

		assert.Equal(t, testHost+authorizeURI+"?"+
			"client_id="+testClientID+"&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Fv1%2Fauth%2Fdotwallet&"+
			"response_type=code&scope=user.info+autopay.bsv+autopay.btc+autopay.eth&"+
			"state="+state, authURL)

	})
}

// TestClient_UpdateApplicationAccessToken will test the method UpdateApplicationAccessToken()
func TestClient_UpdateApplicationAccessToken(t *testing.T) {

	t.Run("test client, get application access token", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockUpdateApplicationAccessToken()
		err = c.UpdateApplicationAccessToken()
		assert.NoError(t, err)

		assert.Equal(t, testTokenType, c.options.token.TokenType)
		assert.Equal(t, testExpiresIn, c.options.token.ExpiresIn)
		assert.Equal(t, time.Now().UTC().Unix()+c.options.token.ExpiresIn, c.options.token.ExpiresAt)
		assert.Equal(t, testAccessToken, c.options.token.AccessToken)
	})

	t.Run("failed to unmarshal JSON", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockAccessTokenBadJSON(testHost, http.StatusOK)
		err = c.UpdateApplicationAccessToken()
		assert.Error(t, err)
	})

	t.Run("request failed", func(t *testing.T) {

		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockAccessTokenFailed(http.StatusBadRequest)
		err = c.UpdateApplicationAccessToken()
		assert.Error(t, err)
	})

	t.Run("token revoked or expired", func(t *testing.T) {

		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockAccessTokenRevoked()
		err = c.UpdateApplicationAccessToken()
		assert.Error(t, err)
	})

	t.Run("token expired but auto-fetch new token", func(t *testing.T) {

		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockExpiredApplicationAccessToken(testHost, http.StatusOK)
		err = c.UpdateApplicationAccessToken()
		assert.NoError(t, err)

		// Wait 2 seconds
		time.Sleep(2 * time.Second)

		// Token is expired
		assert.Equal(t, true, c.IsTokenExpired(c.Token()))

		// Try another request - it will auto-fetch a new token
		mockUpdateApplicationAccessToken()
		mockUserReceiveAddress(CoinTypeBSV.String())
		var wallets *Wallets
		wallets, err = c.UserReceiveAddress(testUserID, CoinTypeBSV)
		assert.NoError(t, err)
		assert.NotNil(t, wallets)
	})
}

// TestClient_GetUserToken will test the method GetUserToken()
func TestClient_GetUserToken(t *testing.T) {

	t.Run("test client, get user access token", func(t *testing.T) {

		mockUpdateApplicationAccessToken()
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockGetUserAccessToken()
		var userToken *DotAccessToken
		userToken, err = c.GetUserToken(testUserCode)
		assert.NoError(t, err)
		assert.NotNil(t, userToken)

		assert.Equal(t, testUserRefreshToken, userToken.RefreshToken)
		assert.Equal(t, testUserAccessToken, userToken.AccessToken)
		assert.Equal(t, testTokenType, userToken.TokenType)
		assert.Equal(t, testExpiresIn, userToken.ExpiresIn)
		assert.Equal(t, []string{"user.info"}, userToken.Scopes)
	})

	t.Run("failed to unmarshal JSON", func(t *testing.T) {

		mockUpdateApplicationAccessToken()
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockAccessTokenBadJSON(testHost, http.StatusBadRequest)
		var userToken *DotAccessToken
		userToken, err = c.GetUserToken(testUserCode)
		assert.Error(t, err)
		assert.Nil(t, userToken)
	})

	t.Run("request failed", func(t *testing.T) {

		mockUpdateApplicationAccessToken()
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockAccessTokenFailed(http.StatusBadRequest)
		var userToken *DotAccessToken
		userToken, err = c.GetUserToken(testUserCode)
		assert.Error(t, err)
		assert.Nil(t, userToken)
	})

	t.Run("token revoked or expired", func(t *testing.T) {

		mockUpdateApplicationAccessToken()
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockAccessTokenRevoked()
		var userToken *DotAccessToken
		userToken, err = c.GetUserToken(testUserCode)
		assert.Error(t, err)
		assert.Nil(t, userToken)
	})
}

// TestClient_RefreshUserToken will test the method RefreshUserToken()
func TestClient_RefreshUserToken(t *testing.T) {
	t.Run("test client, refresh user access token", func(t *testing.T) {

		mockUpdateApplicationAccessToken()
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockGetUserAccessToken()
		var userToken *DotAccessToken
		userToken, err = c.GetUserToken(testUserCode)
		assert.NoError(t, err)
		assert.NotNil(t, userToken)

		oldToken := userToken.AccessToken
		oldRefreshToken := userToken.RefreshToken

		mockRefreshUserAccessToken(testHost, http.StatusOK)
		userToken, err = c.RefreshUserToken(userToken)
		assert.NoError(t, err)
		assert.NotNil(t, userToken)

		assert.NotEqual(t, oldToken, userToken.AccessToken)
		assert.NotEqual(t, oldRefreshToken, userToken.RefreshToken)
	})

	t.Run("failed to unmarshal JSON", func(t *testing.T) {

		mockUpdateApplicationAccessToken()
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockGetUserAccessToken()
		var userToken *DotAccessToken
		userToken, err = c.GetUserToken(testUserCode)
		assert.NoError(t, err)
		assert.NotNil(t, userToken)

		mockAccessTokenBadJSON(testHost, http.StatusBadRequest)
		userToken, err = c.RefreshUserToken(userToken)
		assert.Error(t, err)
		assert.Nil(t, userToken)
	})

	t.Run("request failed", func(t *testing.T) {

		mockUpdateApplicationAccessToken()
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockGetUserAccessToken()
		var userToken *DotAccessToken
		userToken, err = c.GetUserToken(testUserCode)
		assert.NoError(t, err)
		assert.NotNil(t, userToken)

		mockAccessTokenFailed(http.StatusExpectationFailed)
		userToken, err = c.RefreshUserToken(userToken)
		assert.Error(t, err)
		assert.Nil(t, userToken)
	})

	t.Run("token revoked or expired", func(t *testing.T) {

		mockUpdateApplicationAccessToken()
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockGetUserAccessToken()
		var userToken *DotAccessToken
		userToken, err = c.GetUserToken(testUserCode)
		assert.NoError(t, err)
		assert.NotNil(t, userToken)

		mockAccessTokenRevoked()
		userToken, err = c.RefreshUserToken(userToken)
		assert.Error(t, err)
		assert.Nil(t, userToken)
	})
}

// TestClient_IsTokenExpired will test the method IsTokenExpired()
func TestClient_IsTokenExpired(t *testing.T) {

	t.Run("token is valid and not expired", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockUpdateApplicationAccessToken()
		err = c.UpdateApplicationAccessToken()
		assert.NoError(t, err)

		assert.Equal(t, false, c.IsTokenExpired(c.Token()))
	})

	t.Run("token is expired", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockExpiredApplicationAccessToken(testHost, http.StatusOK)
		err = c.UpdateApplicationAccessToken()
		assert.NoError(t, err)

		// Sleep until the token is expired (1 second)
		time.Sleep(2 * time.Second)

		// Should be expired
		assert.Equal(t, true, c.IsTokenExpired(c.Token()))
	})
}

// mockUpdateApplicationAccessToken is used for mocking the response
func mockUpdateApplicationAccessToken() {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s%s", testHost, getAccessTokenURI),
		httpmock.NewStringResponder(
			http.StatusOK, `{"code": 0,"msg": "","data": {
"access_token": "`+testAccessToken+`",
"expires_in": `+fmt.Sprintf("%d", testExpiresIn)+`,
"token_type": "`+testTokenType+`"}}`,
		),
	)
}

// mockExpiredApplicationAccessToken is used for mocking the response
func mockExpiredApplicationAccessToken(host string, statusCode int) {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s%s", host, getAccessTokenURI),
		httpmock.NewStringResponder(
			statusCode, `{"code": 0,"msg": "","data": {
"access_token": "`+testAccessToken+`",
"expires_in": 1,
"token_type": "`+testTokenType+`"}}`,
		),
	)
}

// mockGetUserAccessToken is used for mocking the response
func mockGetUserAccessToken() {
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s%s", testHost, getAccessTokenURI),
		httpmock.NewStringResponder(
			http.StatusOK, `{"code": 0,"msg": "",
    "data": {
        "access_token": "`+testUserAccessToken+`",
        "expires_in": `+fmt.Sprintf("%d", testExpiresIn)+`,
        "refresh_token": "`+testUserRefreshToken+`",
        "scope": "`+testUserScopes+`",
        "token_type": "`+testTokenType+`"}}`,
		),
	)
}

// mockRefreshUserAccessToken is used for mocking the response
func mockRefreshUserAccessToken(host string, statusCode int) {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s%s", host, getAccessTokenURI),
		httpmock.NewStringResponder(
			statusCode, `{"code": 0,"msg": "",
    "data": {
        "access_token": "`+testUserAccessTokenNew+`",
        "expires_in": `+fmt.Sprintf("%d", testExpiresIn)+`,
        "refresh_token": "`+testUserRefreshTokenNew+`",
        "scope": "`+testUserScopes+`",
        "token_type": "`+testTokenType+`"}}`,
		),
	)
}

// mockGetUserAccessTokenFailed is used for mocking the response
func mockAccessTokenFailed(statusCode int) {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s%s", testHost, getAccessTokenURI),
		httpmock.NewStringResponder(
			statusCode, `{
"code": 74012,"msg": "User authentication failed",
"data": null,"req_id": "77dbf77797fc727a3e7afe39209e7238"}`,
		),
	)
}

// mockAccessTokenRevoked is used for mocking the response
func mockAccessTokenRevoked() {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s%s", testHost, getAccessTokenURI),
		httpmock.NewStringResponder(
			http.StatusOK, `{"code": 74013,
    "msg": "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client","data": null,
    "req_id": "zaf8f65fb72e30caab5406c49f876316"}`,
		),
	)
}

// mockGetUserAccessTokenBadJSON is used for mocking the response
func mockAccessTokenBadJSON(host string, statusCode int) {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s%s", host, getAccessTokenURI),
		httpmock.NewStringResponder(
			statusCode, `{"code": 74012,"msg",}`,
		),
	)
}
