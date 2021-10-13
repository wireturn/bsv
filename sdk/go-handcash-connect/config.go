package handcash

const (

	// version is the current package version
	version = "v0.1.0"

	// defaultUserAgent is the default user agent for all requests
	defaultUserAgent string = "go-handcash-connect: " + version

	// apiVersion of the Handcash Connect SDK
	apiVersion = "v1"

	// emptyBody is the default body if no body is set
	emptyBody = "{}"
)

// Environments for Handcash
const (
	EnvironmentBeta       = "beta"
	EnvironmentIAE        = "iae"
	EnvironmentProduction = "prod"
)

var (
	environments = map[string]*Environment{
		EnvironmentBeta: {
			APIURL:      "https://beta-cloud.handcash.io",
			ClientURL:   "https://beta-app.handcash.io",
			Environment: EnvironmentBeta,
		},
		EnvironmentIAE: {
			APIURL:      "https://iae.cloud.handcash.io",
			ClientURL:   "https://iae-app.handcash.io",
			Environment: EnvironmentIAE,
		},
		EnvironmentProduction: {
			APIURL:      "https://cloud.handcash.io",
			ClientURL:   "https://app.handcash.io",
			Environment: EnvironmentProduction,
		},
	}
)

// Endpoints for various services
//
// Specs: https://github.com/HandCash/handcash-connect-sdk-js/blob/master/src/api/http_request_factory.js
const (
	// endpointProfile is for accessing profile information
	endpointProfile = "/" + apiVersion + "/connect/profile"

	// endpointProfileCurrent is for getting the current user profile
	endpointProfileCurrent = endpointProfile + "/currentUserProfile"

	// endpointPublicProfilesByHandle will return profiles given list of handles
	// endpointPublicProfilesByHandle = endpointProfile + "/publicUserProfiles"

	// endpointGetFriends will return a list of friends
	// endpointGetFriends = endpointProfile + "/friends"

	// endpointGetPermissions will return a list of permissions for the user
	// endpointGetPermissions = endpointProfile + "/permissions"

	// endpointGetEncryptionKeypair will return the public key
	// endpointGetEncryptionKeypair = endpointProfile + "/encryptionKeypair"

	// endpointSignData will sign given data
	// endpointSignData = endpointProfile + "/signData"

	// endpointWallet is for accessing wallet information
	endpointWallet = "/" + apiVersion + "/connect/wallet"

	// endpointGetSpendableBalanceRequest will return a spendable balance amount
	endpointGetSpendableBalanceRequest = endpointWallet + "/spendableBalance"

	// endpointGetPayRequest will create a new pay request
	endpointGetPayRequest = endpointWallet + "/pay"

	// endpointGetPaymentRequest will create a new payment request
	endpointGetPaymentRequest = endpointWallet + "/payment"
)
