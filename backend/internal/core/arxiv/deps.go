package arxiv

import (
	"context"
	"net/http"
)

// httpClient defines the HTTP client capability required by this service.
type httpClient interface {
	// Get performs an HTTP GET request.
	Get(ctx context.Context, url string) (*http.Response, error)
}
