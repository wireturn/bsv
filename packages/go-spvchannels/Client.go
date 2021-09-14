package spvchannels

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	ws "github.com/gorilla/websocket"
)

// WSCallBack is a callback to process websocket messages
//    t   : message type
//    msg : message content
//    err : message error
type WSCallBack = func(t int, msg []byte, err error) error

// spvConfig hold configuration for rest api connection
type spvConfig struct {
	insecure    bool // equivalent curl -k
	baseURL     string
	version     string
	user        string
	passwd      string
	token       string
	channelID   string
	procces     WSCallBack
	maxNotified uint64
}

// SPVConfigFunc set the rest api configuration
type SPVConfigFunc func(c *spvConfig)

// WithInsecure skip the TSL check (for dev only)
func WithInsecure() SPVConfigFunc {
	return func(c *spvConfig) {
		c.insecure = true
	}
}

// WithBaseURL provide base url (domain:port) for the rest api
func WithBaseURL(url string) SPVConfigFunc {
	return func(c *spvConfig) {
		c.baseURL = url
	}
}

// WithVersion provide version string for the rest api
func WithVersion(v string) SPVConfigFunc {
	return func(c *spvConfig) {
		c.version = v
	}
}

// WithUser provide username for rest basic authentification
func WithUser(userName string) SPVConfigFunc {
	return func(c *spvConfig) {
		c.user = userName
	}
}

// WithPassword provide password for rest basic authentification
func WithPassword(p string) SPVConfigFunc {
	return func(c *spvConfig) {
		c.passwd = p
	}
}

// WithToken provide token for rest token bearer authentification
func WithToken(t string) SPVConfigFunc {
	return func(c *spvConfig) {
		c.token = t
	}
}

// WithChannelID provide channel id for websocket notification
func WithChannelID(id string) SPVConfigFunc {
	return func(c *spvConfig) {
		c.channelID = id
	}
}

// WithWebsocketCallBack provide the callback function to process notification messages
func WithWebsocketCallBack(f WSCallBack) SPVConfigFunc {
	return func(c *spvConfig) {
		c.procces = f
	}
}

// WithMaxNotified define the max number of notifications that websocket process.
// After receiving enough messages, the websocket will automatically close
func WithMaxNotified(m uint64) SPVConfigFunc {
	return func(c *spvConfig) {
		c.maxNotified = m
	}
}

func defaultSPVConfig() *spvConfig {
	// Set the default options
	cfg := &spvConfig{
		insecure:  false,
		baseURL:   "localhost:5010",
		version:   "v1",
		user:      "dev",
		passwd:    "dev",
		token:     "",
		channelID: "",
		procces:   nil,
	}
	return cfg
}

// Client hold rest api configuration and http connection
type Client struct {
	cfg        *spvConfig
	HTTPClient HTTPClient
}

// NewClient create a new rest api client by providing fuctional config settings
//
// Example of usage :
//
//
//	client := spv.NewClient(
//		spv.WithBaseURL("localhost:5010"),
//		spv.WithVersion("v1"),
//		spv.WithUser("dev"),
//		spv.WithPassword("dev"),
//		spv.WithInsecure(),
//	)
//
// The full list of functional settings for a rest client are :
//
// To disable the TSL certificate check ( used in dev only )
//
//   WithInsecure()
//
// To set the base url of the server
//
//   WithBaseURL(url string)
//
// To set the version string of the rest api
//
//   WithVersion(v string)
//
// To set the user's name for basic authentification
//
//   WithUser(userName string)
//
// To set the user's password for the basic authentification
//
//   WithPassword(p string)
//
// To set the brearer token authentification (this will ignore the basic authentification if set)
//
// WithToken(t string)
//
func NewClient(opts ...SPVConfigFunc) *Client {

	// Start with the defaults then overwrite config with any set by user
	cfg := defaultSPVConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	httpClient := http.Client{
		Timeout: time.Minute,
	}

	if cfg.insecure {
		httpClient.Transport = &http.Transport{
			// #nosec
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return &Client{
		cfg:        cfg,
		HTTPClient: &httpClient,
	}
}

// errorResponse hold structure of error rest call
type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// successResponse hold structure of success rest call
type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// sendRequest send the http request and receive the response
func (c *Client) sendRequest(req *http.Request, out interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	if c.cfg.token == "" {
		req.SetBasicAuth(c.cfg.user, c.cfg.passwd)
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.cfg.token))
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	fullResponse := successResponse{
		Code: res.StatusCode,
		Data: out,
	}

	if out != nil {
		if err = json.NewDecoder(res.Body).Decode(&fullResponse.Data); err != nil {
			return err
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Websocket client is listening to the stream of notifications, which notifies new messages for a specific channel.
// It does not receive the content of the message itself (even though, the notification itself is a text)
// User has to write an separate engine to pull the new (unread) messages content when it receive a notification.
// This can be easily done through the existing endpoint Messages provided in the rest api

// WSClient is the structure holding the
//    - websocket configuration
//    - websocket connection
//    - number of received notifications
type WSClient struct {
	cfg        *spvConfig
	ws         *ws.Conn
	nbNotified uint64
}

// NewWSClient create a new websocket client by providing fuctional config settings.
//
// Example of usage :
//
//
//	ws := spv.NewWSClient(
//		spv.WithBaseURL("localhost:5010"),
//		spv.WithVersion("v1"),
//		spv.WithChannelID(channelid),
//		spv.WithToken(tok),
//		spv.WithInsecure(),
//		spv.WithWebsocketCallBack(PullUnreadMessages),
//		spv.WithMaxNotified(10),
//	)
//
// The full list of functional settings for a websocket client are :
//
// To disable the TSL certificate check ( used in dev only )
//
//   WithInsecure()
//
// To set the base url of the server
//
//   WithBaseURL(url string)
//
// To set the version string of the server
//
//   WithVersion(v string)
//
// To set channel to be notified
//
//   WithChannelID(channelid string)
//
// To set the token that allow the socket connection
//
//   WithToken(tok string)
//
// To specify a callback function to process the notification
//
//   WithWebsocketCallBack(p PullUnreadMessages)
//
// To set the max number of notifications that user want to receive ( used in test only)
//
// After receiving enough notifications, the socket stop and close
//
//   WithMaxNotified(n uint64)
func NewWSClient(opts ...SPVConfigFunc) *WSClient {
	// Start with the defaults then overwrite config with any set by user
	cfg := defaultSPVConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.procces == nil {
		cfg.procces = func(t int, msg []byte, err error) error {
			return nil
		}
	}

	return &WSClient{
		cfg:        cfg,
		ws:         nil,
		nbNotified: 0,
	}
}

// urlPath return the path part of the connection URL
func (c *WSClient) urlPath() string {
	return fmt.Sprintf("/api/%s/channel/%s/notify", c.cfg.version, c.cfg.channelID)
}

// NbNotified return the number of processed messages
func (c *WSClient) NbNotified() uint64 {
	return c.nbNotified
}

// Run establishes the connection and start listening the notification stream
// process the notification if a callback is provided
func (c *WSClient) Run() error {
	u := url.URL{
		Scheme: "wss",
		Host:   "localhost:5010",
		Path:   c.urlPath(),
	}

	q := u.Query()
	q.Set("token", c.cfg.token)
	u.RawQuery = q.Encode()

	d := ws.DefaultDialer
	if c.cfg.insecure {
		// #nosec
		d = &ws.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
			TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
		}
	}

	conn, httpRESP, err := d.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
		_ = httpRESP.Body.Close()
	}()

	c.ws = conn

	for c.nbNotified < c.cfg.maxNotified {
		t, msg, err := c.ws.ReadMessage()
		if err != nil {
			return fmt.Errorf("%w. Total processed %d messages", err, c.nbNotified)
		}

		c.nbNotified++
		if c.cfg.procces != nil {
			err = c.cfg.procces(t, msg, err)
			if err != nil {
				return fmt.Errorf("%w. Total processed %d messages", err, c.nbNotified)
			}
		}
	}

	return nil
}
