package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// NewViperConfig will setup and return a new viper based configuration handler.
func NewViperConfig(appname string) *Config {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return &Config{}
}

// WithServer will setup the web server configuration if required.
func (c *Config) WithServer() *Config {
	viper.SetDefault(EnvServerPort, ":8443")
	viper.SetDefault(EnvServerHost, "payd:8443")
	c.Server = &Server{
		Port:     viper.GetString(EnvServerPort),
		Hostname: viper.GetString(EnvServerHost),
	}
	return c
}

// WithDeployment sets up the deployment configuration if required.
func (c *Config) WithDeployment(appName string) *Config {
	viper.SetDefault(EnvEnvironment, "dev")
	viper.SetDefault(EnvRegion, "test")
	viper.SetDefault(EnvCommit, "test")
	viper.SetDefault(EnvVersion, "test")
	viper.SetDefault(EnvBuildDate, time.Now().UTC())
	viper.SetDefault(EnvMainNet, false)

	c.Deployment = &Deployment{
		Environment: viper.GetString(EnvEnvironment),
		Region:      viper.GetString(EnvRegion),
		Version:     viper.GetString(EnvVersion),
		Commit:      viper.GetString(EnvCommit),
		BuildDate:   viper.GetTime(EnvBuildDate),
		AppName:     appName,
		MainNet:     viper.GetBool(EnvMainNet),
	}
	return c
}

// WithLog sets up and returns log config.
func (c *Config) WithLog() *Config {
	viper.SetDefault(EnvLogLevel, "info")
	c.Logging = &Logging{Level: viper.GetString(EnvLogLevel)}
	return c
}

// WithDb sets up and returns database configuration.
func (c *Config) WithDb() *Config {
	viper.SetDefault(EnvDb, "sqlite")
	viper.SetDefault(EnvDbDsn, "file:data/wallet.db?_foreign_keys=true&pooled=true")
	viper.SetDefault(EnvDbSchema, "data/sqlite/migrations")
	viper.SetDefault(EnvDbMigrate, true)
	c.Db = &Db{
		Type:       DbType(viper.GetString(EnvDb)),
		Dsn:        viper.GetString(EnvDbDsn),
		SchemaPath: viper.GetString(EnvDbSchema),
		MigrateDb:  viper.GetBool(EnvDbMigrate),
	}
	return c
}

// WithHeadersv sets up and returns Headersv configuration.
func (c *Config) WithHeadersv() *Config {
	viper.SetDefault(EnvHeadersvAddress, "headersv:8001")
	viper.SetDefault(EnvHeadersvTimeout, 30)
	c.Headersv = &Headersv{
		Address: viper.GetString(EnvHeadersvAddress),
		Timeout: viper.GetInt(EnvHeadersvTimeout),
	}
	return c
}

// WithPaymail sets up and returns paymail configuration.
func (c *Config) WithPaymail() *Config {
	viper.SetDefault(EnvPaymailEnabled, false)
	viper.SetDefault(EnvPaymailAddress, "test@test.com")
	viper.SetDefault(EnvPaymailIsBeta, false)
	c.Paymail = &Paymail{
		UsePaymail: viper.GetBool(EnvPaymailEnabled),
		IsBeta:     viper.GetBool(EnvPaymailIsBeta),
		Address:    viper.GetString(EnvPaymailAddress),
	}
	return c
}

// WithWallet sets up and returns merchant wallet configuration.
func (c *Config) WithWallet() *Config {
	viper.SetDefault(EnvNetwork, "regtest")
	viper.SetDefault(EnvMerchantName, "payd")
	viper.SetDefault(EnvAvatarURL, "https://media.bitcoinfiles.org/eec638f2e10a533b344d71a20f102bca2dbf2385d3a18d77c303539a7e6b666b")
	viper.SetDefault(EnvMerchantAddress, "1 the street, town, T1 1TT")
	viper.SetDefault(EnvMerchantEmail, "test@ppctl.nchain.com")
	viper.SetDefault(EnvPaymentExpiry, 24)
	c.Wallet = &Wallet{
		Network:            viper.GetString(EnvNetwork),
		MerchantAvatarURL:  viper.GetString(EnvAvatarURL),
		MerchantName:       viper.GetString(EnvMerchantName),
		MerchantEmail:      viper.GetString(EnvMerchantEmail),
		Address:            viper.GetString(EnvMerchantAddress),
		PaymentExpiryHours: viper.GetInt(EnvPaymentExpiry),
	}
	return c
}

// WithMapi will setup Mapi settings.
func (c *Config) WithMapi() *Config {
	viper.SetDefault(EnvMAPIMinerName, "local-mapi")
	viper.SetDefault(EnvMAPIURL, "http://mapi:9014")
	viper.SetDefault(EnvMAPIToken, "")
	c.Mapi = &MApi{
		MinerName: viper.GetString(EnvMAPIMinerName),
		URL:       viper.GetString(EnvMAPIURL),
		Token:     viper.GetString(EnvMAPIToken),
	}
	return c
}
