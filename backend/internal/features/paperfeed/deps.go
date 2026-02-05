package paperfeed

import (
	"context"
	"time"

	"github.com/rrlian/papertok/backend/internal/core/arxiv"
)

// arxivService defines the arXiv service capability required by this feature.
type arxivService interface {
	// FetchByCategory fetches papers from arXiv by category.
	FetchByCategory(ctx context.Context, req *arxiv.FetchRequest) ([]*arxiv.Paper, error)
}

// paperRepository defines the paper repository capability required by this feature.
type paperRepository interface {
	// GetByCategory retrieves cached papers by category.
	GetByCategory(ctx context.Context, category string) ([]*repositoryPaper, bool)

	// SaveByCategory stores papers for a category with TTL.
	SaveByCategory(ctx context.Context, category string, papers []*repositoryPaper, ttl time.Duration)
}

// repositoryPaper represents a paper in the repository layer.
// This is a local alias to avoid import cycles.
type repositoryPaper struct {
	ID              string
	Title           string
	Authors         []string
	Summary         string
	Published       time.Time
	Updated         time.Time
	Categories      []string
	PrimaryCategory string
	ArxivURL        string
	PDFURL          string
	ImageURL        string
}
