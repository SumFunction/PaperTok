package paperfeed

import (
	"context"
	"time"
)

// Paper represents a paper in the feed response.
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

// FetchRequest contains parameters for fetching the paper feed.
type FetchRequest struct {
	Category   string
	Limit      int
	Offset     int
	SortBy     string
}

// Service defines the interface for paper feed operations.
type Service interface {
	// GetFeed fetches papers for the feed.
	// @Params:
	//   - ctx: context for cancellation and tracing
	//   - req: fetch request parameters
	// @Returns:
	//   - []*Paper: list of papers for the feed
	//   - error: if fetch fails
	GetFeed(ctx context.Context, req *FetchRequest) ([]*Paper, error)
}
