package config

import (
	"fmt"
	"strings"

	"github.com/tokenized/pkg/bitcoin"
)

// Config holds all configuration for the running service.
type Config struct {
	Net            bitcoin.Network
	IsTest         bool
	NodeAddress    string         // IP address of trusted external full node
	UserAgent      string         // User agent to send to external node
	StartHash      bitcoin.Hash32 // Hash of first block to start processing on initial run
	UntrustedCount int            // The number of untrusted nodes to run for double spend monitoring
	SafeTxDelay    int            // Number of milliseconds without conflict before a tx is "safe"
	ShotgunCount   int            // The number of nodes to attempt to send to when broadcasting
	RequestMempool bool           // request mempool after syncing to bitcoind node

	// Retry attempts when main connection fails.
	MaxRetries int
	RetryDelay int
}

const (
	DefaultMaxRetries = 25
	DefaultRetryDelay = 5000
)

// NewConfig returns a new Config populated from environment variables.
func NewConfig(net bitcoin.Network, isTest bool, host, useragent, starthash string,
	untrustedNodes, safeDelay, shotgunCount, maxRetries, retryDelay int,
	requestMempool bool) (Config, error) {
	result := Config{
		Net:            net,
		IsTest:         isTest,
		NodeAddress:    host,
		UserAgent:      useragent,
		UntrustedCount: untrustedNodes,
		SafeTxDelay:    safeDelay,
		ShotgunCount:   shotgunCount,
		MaxRetries:     maxRetries,
		RetryDelay:     retryDelay,
		RequestMempool: requestMempool,
	}

	hash, err := bitcoin.NewHash32FromStr(starthash)
	if err != nil {
		return result, err
	}
	result.StartHash = *hash
	return result, nil
}

// String returns a custom string representation.
//
// This is important so we don't log sensitive config values.
func (c Config) String() string {
	pairs := map[string]string{
		"NodeAddress": c.NodeAddress,
		"UserAgent":   c.UserAgent,
		"StartHash":   c.StartHash.String(),
		"SafeTxDelay": fmt.Sprintf("%d ms", c.SafeTxDelay),
	}

	parts := []string{}

	for k, v := range pairs {
		parts = append(parts, fmt.Sprintf("%v:%v", k, v))
	}

	return fmt.Sprintf("{%v}", strings.Join(parts, " "))
}
