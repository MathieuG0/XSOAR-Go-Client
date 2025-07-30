package xsoar

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
	"slices"

	"github.com/MathieuG0/XSOAR-Go-Client/cache"
	"github.com/hashicorp/go-cleanhttp"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
)

const (
	userAgent = "go-xsoar"
)

type ClientOption func(*Client) error

type Client struct {
	// HTTP httpClient used to request Cortex XSOAR's API
	httpClient *retryablehttp.Client

	// Base URL of the XSOAR server
	baseURL *url.URL

	// Credentials for basic authentication
	username, password string

	// Credentials for API key authentication
	apiKey string

	// User agent sent in requests
	userAgent string

	// API modules
	Integration *IntegrationModule
	Role        *RoleModule
	User        *UserModule
	Server      *ServerModule
}

func NewClient(options ...ClientOption) (*Client, error) {
	cache := cache.NewCache()
	client := &Client{userAgent: userAgent}

	client.httpClient = retryablehttp.NewClient()
	client.httpClient.RetryMax = 0

	client.setDefaultConfig()

	for _, fn := range options {
		err := fn(client)
		if err != nil {
			return nil, err
		}
	}

	client.Integration = &IntegrationModule{client, cache}
	client.Role = &RoleModule{client, cache}
	client.User = &UserModule{client, cache}
	client.Server = &ServerModule{client, cache}

	return client, nil
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		return c.setBaseURL(baseURL)
	}
}

func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		c.apiKey = apiKey
		return nil
	}
}

func WithBasicAuth(username, password string) ClientOption {
	return func(c *Client) error {
		c.username, c.password = username, password
		return nil
	}
}

func WithoutSSLVerify() ClientOption {
	return func(c *Client) error {
		c.disableSSLVerify()
		return nil
	}
}

func (c *Client) disableSSLVerify() {
	t := cleanhttp.DefaultPooledTransport()
	t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c.httpClient.HTTPClient.Transport = t
}

func (c *Client) setBaseURL(baseURL string) (err error) {
	c.baseURL, err = url.Parse(baseURL)

	return err
}

func (c *Client) setDefaultConfig() error {
	if c.baseURL == nil {
		if err := c.setBaseURL(os.Getenv("DEMISTO_BASE_URL")); err != nil {
			return err
		}
	}

	if c.apiKey == "" && c.username == "" && c.password == "" {
		c.apiKey = os.Getenv("DEMISTO_API_KEY")
	}

	if c.username == "" {
		c.username = os.Getenv("DEMISTO_USERNAME")
	}

	if c.password == "" {
		c.password = os.Getenv("DEMISTO_PASSWORD")
	}

	if os.Getenv("DEMISTO_VERIFY_SSL") == "false" {
		c.disableSSLVerify()
	}

	return nil
}

type RequestOption func(req *retryablehttp.Request) error

func WithHeader(key, value string) RequestOption {
	return func(req *retryablehttp.Request) error {
		req.Header.Add(key, value)
		return nil
	}
}

func WithBody(rawBody any) RequestOption {
	return func(req *retryablehttp.Request) error {
		return req.SetBody(rawBody)
	}
}

func (c *Client) NewRequest(method string, endpoint string, options ...RequestOption) (*retryablehttp.Request, error) {
	req, err := retryablehttp.NewRequest(method, c.baseURL.JoinPath(endpoint).String(), nil)
	if err != nil {
		return nil, err
	}

	for _, fn := range options {
		if err := fn(req); err != nil {
			return nil, err
		}
	}

	req.Header.Add("User-Agent", c.userAgent)

	if c.apiKey != "" {
		req.Header.Add("Authorization", c.apiKey)
	} else {
		req.SetBasicAuth(c.username, c.password)
	}

	return req, nil
}

func (c *Client) Do(req *retryablehttp.Request, okCodes ...int) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if (okCodes != nil && !slices.Contains(okCodes, resp.StatusCode)) || resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("http code %d: %s", resp.StatusCode, GetMessage(resp))
	}

	return resp, nil
}
