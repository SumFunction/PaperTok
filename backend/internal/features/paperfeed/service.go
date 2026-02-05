package paperfeed

import (
	"context"
	"time"

	"github.com/rrlian/papertok/backend/internal/core/arxiv"
	paperRepo "github.com/rrlian/papertok/backend/internal/repository/paper"
)

// Impl implements the paperfeed Service interface.
type Impl struct {
	arxivSvc  arxiv.Service
	paperRepo paperRepo.Repository
	cacheTTL  time.Duration
}

// Ensure Impl implements Service interface
var _ Service = (*Impl)(nil)

// New creates a new paperfeed service instance.
func New(arxivSvc arxiv.Service, repo paperRepo.Repository, cacheTTL time.Duration) *Impl {
	return &Impl{
		arxivSvc:  arxivSvc,
		paperRepo: repo,
		cacheTTL:  cacheTTL,
	}
}

// GetFeed fetches papers for the feed.
func (s *Impl) GetFeed(ctx context.Context, req *FetchRequest) ([]*Paper, error) {
	// Try cache first
	if cachedPapers, found := s.paperRepo.GetByCategory(ctx, req.Category); found {
		return s.convertRepoPapers(cachedPapers), nil
	}

	// Fetch from arXiv
	arxivPapers, err := s.arxivSvc.FetchByCategory(ctx, &arxiv.FetchRequest{
		Category:   req.Category,
		MaxResults: req.Limit,
		SortBy:     req.SortBy,
		Offset:     req.Offset,
	})
	if err != nil {
		return nil, err
	}

	// Convert and cache
	repoPapers := s.convertArxivPapers(arxivPapers)
	s.paperRepo.SaveByCategory(ctx, req.Category, repoPapers, s.cacheTTL)

	return s.convertRepoPapers(repoPapers), nil
}

// convertArxivPapers converts arXiv papers to repository papers.
func (s *Impl) convertArxivPapers(papers []*arxiv.Paper) []*paperRepo.Paper {
	result := make([]*paperRepo.Paper, len(papers))
	for i, p := range papers {
		result[i] = &paperRepo.Paper{
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
	return result
}

// convertRepoPapers converts repository papers to feature papers.
func (s *Impl) convertRepoPapers(papers []*paperRepo.Paper) []*Paper {
	result := make([]*Paper, len(papers))
	for i, p := range papers {
		result[i] = &Paper{
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
	return result
}
