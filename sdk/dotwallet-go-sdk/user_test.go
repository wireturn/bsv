package dotwallet

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// TestClient_UserInfo will test the method UserInfo()
func TestClient_UserInfo(t *testing.T) {

	t.Run("test client, get user access token, get info", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockGetUserAccessToken()
		var userToken *DotAccessToken
		userToken, err = c.GetUserToken(testUserCode)
		assert.NoError(t, err)
		assert.NotNil(t, userToken)

		mockUserInfo(http.StatusOK)
		var user *User
		user, err = c.UserInfo(userToken)
		assert.NoError(t, err)
		assert.NotNil(t, user)

		assert.Equal(t, testUserAvatar, user.Avatar)
		assert.Equal(t, testUserID, user.ID)
		assert.Equal(t, testUserNickname, user.Nickname)
		assert.Equal(t, testUserBsvAddress, user.WebWalletAddress.BSVRegular)
		assert.Equal(t, testUserBtcAddress, user.WebWalletAddress.BTCRegular)
		assert.Equal(t, testUserEthAddress, user.WebWalletAddress.ETHRegular)
	})

	t.Run("failed to unmarshal JSON", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockGetUserAccessToken()
		var userToken *DotAccessToken
		userToken, err = c.GetUserToken(testUserCode)
		assert.NoError(t, err)
		assert.NotNil(t, userToken)

		mockUserInfoBadJSON(http.StatusOK)
		var user *User
		user, err = c.UserInfo(userToken)
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("request failed", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockGetUserAccessToken()
		var userToken *DotAccessToken
		userToken, err = c.GetUserToken(testUserCode)
		assert.NoError(t, err)
		assert.NotNil(t, userToken)

		mockUserInfoFailed(http.StatusOK)
		var user *User
		user, err = c.UserInfo(userToken)
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("token revoked or expired", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockGetUserAccessToken()
		var userToken *DotAccessToken
		userToken, err = c.GetUserToken(testUserCode)
		assert.NoError(t, err)
		assert.NotNil(t, userToken)

		mockAccessTokenRevoked()
		var user *User
		user, err = c.UserInfo(userToken)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

// TestClient_UserReceiveAddress will test the method UserReceiveAddress()
func TestClient_UserReceiveAddress(t *testing.T) {

	t.Run("test client, get receive address - bsv", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockUpdateApplicationAccessToken()
		mockUserReceiveAddress(CoinTypeBSV.String())
		var wallets *Wallets
		wallets, err = c.UserReceiveAddress(testUserID, CoinTypeBSV)
		assert.NoError(t, err)
		assert.NotNil(t, wallets)

		assert.Equal(t, testUserBsvAddress, wallets.PrimaryWallet.Address)
		assert.Equal(t, testUserPaymail, wallets.PrimaryWallet.Paymail)
		assert.Equal(t, testUserID, wallets.PrimaryWallet.UserID)
		assert.Equal(t, CoinTypeBSV.String(), wallets.PrimaryWallet.CoinType)
	})

	t.Run("test client, get receive address - btc", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockUpdateApplicationAccessToken()
		mockUserReceiveAddress(CoinTypeBTC.String())
		var wallets *Wallets
		wallets, err = c.UserReceiveAddress(testUserID, CoinTypeBSV)
		assert.NoError(t, err)
		assert.NotNil(t, wallets)

		assert.Equal(t, CoinTypeBTC.String(), wallets.PrimaryWallet.CoinType)
	})

	t.Run("test client, get receive address - eth", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockUpdateApplicationAccessToken()
		mockUserReceiveAddress(CoinTypeETH.String())
		var wallets *Wallets
		wallets, err = c.UserReceiveAddress(testUserID, CoinTypeBSV)
		assert.NoError(t, err)
		assert.NotNil(t, wallets)

		assert.Equal(t, CoinTypeETH.String(), wallets.PrimaryWallet.CoinType)
	})

	t.Run("failed to unmarshal JSON", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockUpdateApplicationAccessToken()
		mockUserReceiveAddressBadJSON(http.StatusOK)
		var wallets *Wallets
		wallets, err = c.UserReceiveAddress(testUserID, CoinTypeBSV)
		assert.Error(t, err)
		assert.Nil(t, wallets)
	})

	t.Run("request failed", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockUpdateApplicationAccessToken()
		mockUserReceiveAddressFailed(http.StatusOK)
		var wallets *Wallets
		wallets, err = c.UserReceiveAddress(testUserID, CoinTypeBSV)
		assert.Error(t, err)
		assert.Nil(t, wallets)
	})

	t.Run("token revoked or expired", func(t *testing.T) {
		c, err := newTestClient()
		assert.NoError(t, err)
		assert.NotNil(t, c)

		mockUpdateApplicationAccessToken()
		mockAccessTokenRevoked()
		var wallets *Wallets
		wallets, err = c.UserReceiveAddress(testUserID, CoinTypeBSV)
		assert.Error(t, err)
		assert.Nil(t, wallets)
	})
}

// mockUserInfo is used for mocking the response
func mockUserInfo(statusCode int) {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s%s", testHost, getUserInfo),
		httpmock.NewStringResponder(
			statusCode, `{"code": 0,"msg": "",
"data": {
        "id": "`+testUserID+`",
        "nickname": "`+testUserNickname+`",
        "avatar": "`+testUserAvatar+`",
        "web_wallet_address": {
            "bsv_regular": "`+testUserBsvAddress+`",
            "btc_regular": "`+testUserBtcAddress+`",
            "eth_regular": "`+testUserEthAddress+`"
        }}}`,
		),
	)
}

// mockUserInfoFailed is used for mocking the response
func mockUserInfoFailed(statusCode int) {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s%s", testHost, getUserInfo),
		httpmock.NewStringResponder(
			statusCode, `{
"code": 74012,"msg": "User info failed",
"data": null,"req_id": "77dbf77797fc727a3e7afe39209e7238"}`,
		),
	)
}

// mockUserInfoBadJSON is used for mocking the response
func mockUserInfoBadJSON(statusCode int) {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s%s", testHost, getUserInfo),
		httpmock.NewStringResponder(
			statusCode, `{"code": 74012,"msg",}`,
		),
	)
}

// mockUserReceiveAddress is used for mocking the response
func mockUserReceiveAddress(coinType string) {
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s%s", testHost, getUserReceiveAddress),
		httpmock.NewStringResponder(
			http.StatusOK, `{"code": 0,"msg": "",
    "data": {
        "primary_wallet": {
            "user_id": "`+testUserID+`",
            "wallet_type": "coin_regular",
            "wallet_index": 0,
            "coin_type": "`+coinType+`",
            "address": "`+testUserBsvAddress+`",
            "paymail": "`+testUserPaymail+`"
        },
        "autopay_wallet": {
            "user_id": "`+testUserID+`",
            "wallet_type": "coin_autopay",
            "wallet_index": 0,
            "coin_type": "`+coinType+`",
            "address": "`+testUserBsvAddress+`",
            "paymail": "`+testUserPaymail+`"
        }}}`,
		),
	)
}

// mockUserReceiveAddressFailed is used for mocking the response
func mockUserReceiveAddressFailed(statusCode int) {
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s%s", testHost, getUserReceiveAddress),
		httpmock.NewStringResponder(
			statusCode, `{
"code": 74012,"msg": "User receive address failed",
"data": null,"req_id": "77dbf77797fc727a3e7afe39209e7238"}`,
		),
	)
}

// mockUserReceiveAddressBadJSON is used for mocking the response
func mockUserReceiveAddressBadJSON(statusCode int) {
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s%s", testHost, getUserReceiveAddress),
		httpmock.NewStringResponder(
			statusCode, `{"code": 74012,"msg",}`,
		),
	)
}
