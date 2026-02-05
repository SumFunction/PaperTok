package paper

import (
	"context"
	"time"
)

// Paper represents a paper entity in the repository layer.
type Paper struct {
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

// Repository defines the interface for paper data access.
// This abstraction allows for different storage implementations
// (memory, database, etc.)
type Repository interface {
	// GetByCategory retrieves papers by category.
	// Returns cached papers if available, otherwise returns nil.
	GetByCategory(ctx context.Context, category string) ([]*Paper, bool)

	// SaveByCategory stores papers for a category with TTL.
	SaveByCategory(ctx context.Context, category string, papers []*Paper, ttl time.Duration)

	// GetByID retrieves a single paper by ID.
	GetByID(ctx context.Context, id string) (*Paper, bool)

	// Save stores a single paper.
	Save(ctx context.Context, paper *Paper, ttl time.Duration)

	// InvalidateCategory removes cached papers for a category.
	InvalidateCategory(ctx context.Context, category string)

	// Clear removes all cached papers.
	Clear(ctx context.Context)
}
