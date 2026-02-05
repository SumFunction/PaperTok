package httpclient

import (
	"context"
	"io"
	"net/http"
	"time"
)

// HTTPClient defines the interface for making HTTP requests.
// This abstraction allows for easier testing and different implementations.
type HTTPClient interface {
	// Get performs an HTTP GET request.
	Get(ctx context.Context, url string) (*http.Response, error)

	// Do performs an HTTP request.
	Do(req *http.Request) (*http.Response, error)
}

// Client implements HTTPClient using the standard http.Client.
type Client struct {
	httpClient *http.Client
}

// Ensure Client implements HTTPClient interface
var _ HTTPClient = (*Client)(nil)

// Config holds the configuration for creating an HTTP client.
type Config struct {
	Timeout time.Duration
}

// NewClient creates a new HTTP client with the given configuration.
func NewClient(cfg Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// Get performs an HTTP GET request.
func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}

// Do performs an HTTP request.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

// ReadBody reads and returns the body of an HTTP response.
// The caller is responsible for closing the response body after this function returns.
func ReadBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
