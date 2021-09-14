package client

import (
	"github.com/tokenized/pkg/bitcoin"

	"github.com/pkg/errors"
)

// EnvConfig has environment safe types for importing from environment values.
// Apparently it doesn't use the TextMarshaler interfaces.
type EnvConfig struct {
	ServerAddress    string `envconfig:"SPYNODE_SERVER_ADDRESS" json:"SPYNODE_SERVER_ADDRESS"`
	ServerKey        string `envconfig:"SPYNODE_SERVER_KEY" json:"SPYNODE_SERVER_KEY"`
	ClientKey        string `envconfig:"SPYNODE_CLIENT_KEY" json:"SPYNODE_CLIENT_KEY" masked:"true"`
	StartBlockHeight uint32 `default:"478559" envconfig:"SPYNODE_START_BLOCK_HEIGHT" json:"SPYNODE_START_BLOCK_HEIGHT"`
	ConnectionType   uint8  `default:"1" envconfig:"SPYNODE_CONNECTION_TYPE" json:"SPYNODE_CONNECTION_TYPE"`

	MaxRetries int `default:"50" envconfig:"SPYNODE_MAX_RETRIES"`
	RetryDelay int `default:"2000" envconfig:"SPYNODE_RETRY_DELAY"`

	RequestTimeout int `default:"10000" envconfig:"SPYNODE_REQUEST_TIMEOUT"` // in milliseconds
}

type Config struct {
	ServerAddress    string            `json:"SPYNODE_SERVER_ADDRESS"`
	ServerKey        bitcoin.PublicKey `json:"SPYNODE_SERVER_KEY"`
	ClientKey        bitcoin.Key       `json:"SPYNODE_CLIENT_KEY"`
	StartBlockHeight uint32            `json:"SPYNODE_START_BLOCK_HEIGHT"`
	ConnectionType   uint8             `json:"SPYNODE_CONNECTION_TYPE"`

	MaxRetries int `json:"SPYNODE_MAX_RETRIES"`
	RetryDelay int `json:"SPYNODE_RETRY_DELAY"`

	RequestTimeout int `json:"SPYNODE_REQUEST_TIMEOUT"`
}

func NewConfig(serverAddress string, serverKey bitcoin.PublicKey, clientKey bitcoin.Key,
	startBlockHeight uint32, connectionType uint8) *Config {
	return &Config{
		ServerAddress:    serverAddress,
		ServerKey:        serverKey,
		ClientKey:        clientKey,
		StartBlockHeight: startBlockHeight,
		ConnectionType:   connectionType,
		MaxRetries:       50,
		RetryDelay:       2000,
		RequestTimeout:   10000,
	}
}

func ConvertEnvConfig(env *EnvConfig) (*Config, error) {
	result := &Config{
		ServerAddress:    env.ServerAddress,
		StartBlockHeight: env.StartBlockHeight,
		ConnectionType:   env.ConnectionType,
		MaxRetries:       env.MaxRetries,
		RetryDelay:       env.RetryDelay,
		RequestTimeout:   env.RequestTimeout,
	}

	serverKey, err := bitcoin.PublicKeyFromStr(env.ServerKey)
	if err != nil {
		return nil, errors.Wrap(err, "server key")
	}
	result.ServerKey = serverKey

	clientKey, err := bitcoin.KeyFromStr(env.ClientKey)
	if err != nil {
		return nil, errors.Wrap(err, "client key")
	}
	result.ClientKey = clientKey

	return result, nil
}
