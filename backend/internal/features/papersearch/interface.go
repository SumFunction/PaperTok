package papersearch

import (
	"context"
	"time"
)

// Paper represents a paper in the search response.
type Paper struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Authors         []string  `json:"authors"`
	Summary         string    `json:"summary"`
	Published       time.Time `json:"published"`
	Updated         time.Time `json:"updated"`
	Categories      []string  `json:"categories"`
	PrimaryCategory string    `json:"primaryCategory"`
	ArxivURL        string    `json:"arxivUrl"`
	PDFURL          string    `json:"pdfUrl"`
	ImageURL        string    `json:"imageUrl"`
}

// Service defines the interface for paper search operations.
type Service interface {
	// Search searches papers by keyword.
	// @Params:
	//   - ctx: context for cancellation and tracing
	//   - query: search keyword
	//   - limit: maximum number of results
	// @Returns:
	//   - []*Paper: list of matching papers
	//   - error: if search fails
	Search(ctx context.Context, query string, limit int) ([]*Paper, error)

	// GetByID retrieves a single paper by ID.
	// @Params:
	//   - ctx: context for cancellation and tracing
	//   - id: paper ID
	// @Returns:
	//   - *Paper: the paper if found, nil otherwise
	//   - error: if fetch fails
	GetByID(ctx context.Context, id string) (*Paper, error)
}
