package arxiv

import "context"

// Service defines the interface for arXiv API operations.
// This interface is used by Features that need to fetch papers from arXiv.
type Service interface {
	// FetchByCategory fetches papers from arXiv by category.
	// @Params:
	//   - ctx: context for cancellation and tracing
	//   - req: fetch request parameters
	// @Returns:
	//   - []*Paper: list of papers
	//   - error: ErrFetchFailed if request fails
	FetchByCategory(ctx context.Context, req *FetchRequest) ([]*Paper, error)

	// Search searches papers by keyword.
	// @Params:
	//   - ctx: context for cancellation and tracing
	//   - query: search keyword
	//   - limit: maximum number of results
	// @Returns:
	//   - []*Paper: list of matching papers
	//   - error: ErrSearchFailed if request fails
	Search(ctx context.Context, query string, limit int) ([]*Paper, error)

	// GetByID fetches a single paper by its arXiv ID.
	// @Params:
	//   - ctx: context for cancellation and tracing
	//   - id: arXiv paper ID (e.g., "2301.12345")
	// @Returns:
	//   - *Paper: the paper if found, nil otherwise
	//   - error: ErrFetchFailed if request fails
	GetByID(ctx context.Context, id string) (*Paper, error)
}
