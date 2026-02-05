package papersearch

import (
	"context"

	"github.com/rrlian/papertok/backend/internal/core/arxiv"
)

// arxivService defines the arXiv service capability required by this feature.
type arxivService interface {
	// Search searches papers by keyword.
	Search(ctx context.Context, query string, limit int) ([]*arxiv.Paper, error)

	// GetByID fetches a single paper by ID.
	GetByID(ctx context.Context, id string) (*arxiv.Paper, error)
}
