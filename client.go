package superhub

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
)

const DefaultBaseURL = "https://api.superhub.host/v1"

type Client struct {
	Credentials Credentials
	BaseURL     string
	HttpClient  *http.Client
}

func (c *Client) GetCredentials() Credentials {
	if c.Credentials == nil {
		return &EmptyCredentials{}
	}

	return c.Credentials
}

func (c *Client) GetBaseURL() string {
	baseURL := DefaultBaseURL

	if c.BaseURL != "" {
		baseURL = c.BaseURL
	}

	return baseURL
}

func (c *Client) GetEndpointURL(endpoint string) (string, error) {
	parsedURL, err := url.Parse(c.GetBaseURL())
	if err != nil {
		return "", fmt.Errorf("parsing url: %s", err)
	}

	parsedURL.Path = path.Join(parsedURL.Path, endpoint)
	return parsedURL.String(), nil
}

func (c *Client) GetHttpClient() *http.Client {
	if c.HttpClient == nil {
		return &http.Client{}
	}

	return c.HttpClient
}

func NewClientWithCredentials(credentials Credentials) *Client {
	return &Client{
		Credentials: credentials,
	}
}

func NewClient() *Client {
	return &Client{}
}
