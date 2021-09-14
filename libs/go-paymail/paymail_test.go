package paymail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Default testing variables
const (
	testAddress   = "1Cat862cjhp8SgLLMvin5gyk5UScasg1P9"
	testAlias     = "mrz"
	testAvatar    = "https://www.gravatar.com/avatar/372bc0ab9b8a8930d4a86b2c5b11f11e?d=identicon"
	testDomain    = "test.com"
	testMessage   = "This is a test message"
	testName      = "MrZ"
	testOutput    = "76a9147f11c8f67a2781df0400ebfb1f31b4c72a780b9d88ac"
	testPubKey    = "02ead23149a1e33df17325ec7a7ba9e0b20c674c57c630f527d69b866aa9b65b10"
	testServerURL = "https://" + testDomain + "/api/v1/" + DefaultServiceName + "/"
)

// TestVersion will test the method Version()
func TestVersion(t *testing.T) {
	t.Parallel()

	t.Run("get version", func(t *testing.T) {
		ver := Version()
		assert.Equal(t, version, ver)
	})
}

// TestUserAgent will test the method UserAgent()
func TestUserAgent(t *testing.T) {
	t.Parallel()

	t.Run("get user agent", func(t *testing.T) {
		agent := UserAgent()
		assert.Equal(t, defaultUserAgent, agent)
	})
}
