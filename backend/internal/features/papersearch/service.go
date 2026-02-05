package papersearch

import (
	"context"

	"github.com/rrlian/papertok/backend/internal/core/arxiv"
)

// Impl implements the papersearch Service interface.
type Impl struct {
	arxivSvc arxiv.Service
}

// Ensure Impl implements Service interface
var _ Service = (*Impl)(nil)

// New creates a new papersearch service instance.
func New(arxivSvc arxiv.Service) *Impl {
	return &Impl{
		arxivSvc: arxivSvc,
	}
}

// Search searches papers by keyword.
func (s *Impl) Search(ctx context.Context, query string, limit int) ([]*Paper, error) {
	arxivPapers, err := s.arxivSvc.Search(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	return s.convertArxivPapers(arxivPapers), nil
}

// GetByID retrieves a single paper by ID.
func (s *Impl) GetByID(ctx context.Context, id string) (*Paper, error) {
	arxivPaper, err := s.arxivSvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if arxivPaper == nil {
		return nil, nil
	}

	return s.convertArxivPaper(arxivPaper), nil
}

// convertArxivPapers converts arXiv papers to feature papers.
func (s *Impl) convertArxivPapers(papers []*arxiv.Paper) []*Paper {
	result := make([]*Paper, len(papers))
	for i, p := range papers {
		result[i] = s.convertArxivPaper(p)
	}
	return result
}

// convertArxivPaper converts a single arXiv paper to feature paper.
func (s *Impl) convertArxivPaper(p *arxiv.Paper) *Paper {
	return &Paper{
		ID:              p.ID,
		Title:           p.Title,
		Authors:         p.Authors,
		Summary:         p.Summary,
		Published:       p.Published,
		Updated:         p.Updated,
		Categories:      p.Categories,
		PrimaryCategory: p.PrimaryCategory,
		ArxivURL:        p.ArxivURL,
		PDFURL:          p.PDFURL,
		ImageURL:        p.ImageURL,
	}
}
