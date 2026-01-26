package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rrlian/papertok/backend/internal/model"
	"github.com/rrlian/papertok/backend/pkg/arxiv"
)

// ArxivService defines the interface for arXiv service
type ArxivService interface {
	FetchPapers(ctx context.Context, req *model.FetchRequest) ([]*model.Paper, error)
	SearchPapers(ctx context.Context, query string, limit int) ([]*model.Paper, error)
	GetPaperByID(ctx context.Context, paperID string) (*model.Paper, error)
}

// arxivServiceImpl implements ArxivService
type arxivServiceImpl struct {
	client    *arxiv.Client
	cache     CacheService
	cacheTTL  time.Duration
}

// NewArxivService creates a new arXiv service
func NewArxivService(client *arxiv.Client, cache CacheService, cacheTTL time.Duration) ArxivService {
	return &arxivServiceImpl{
		client:   client,
		cache:    cache,
		cacheTTL: cacheTTL,
	}
}

// FetchPapers fetches papers from arXiv or cache
func (s *arxivServiceImpl) FetchPapers(ctx context.Context, req *model.FetchRequest) ([]*model.Paper, error) {
	// Try cache first
	cacheKey := s.buildCacheKey(req.Category)
	if papers, found := s.cache.Get(cacheKey); found {
		return papers, nil
	}

	// Fetch from arXiv
	feed, err := s.client.FetchPapers(req.Category, req.MaxResults, req.SortBy)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch papers: %w", err)
	}

	// Convert feed to papers
	papers := s.convertFeedToPapers(feed)

	// Store in cache
	s.cache.Set(cacheKey, papers, s.cacheTTL)

	return papers, nil
}

// SearchPapers searches papers by keyword
func (s *arxivServiceImpl) SearchPapers(ctx context.Context, query string, limit int) ([]*model.Paper, error) {
	// For search, we don't cache for now
	feed, err := s.client.SearchPapers(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search papers: %w", err)
	}

	return s.convertFeedToPapers(feed), nil
}

// GetPaperByID fetches a single paper by ID from arXiv
func (s *arxivServiceImpl) GetPaperByID(ctx context.Context, paperID string) (*model.Paper, error) {
	// Note: Current cache implementation doesn't support single paper lookup
	// Future: use s.buildPaperCacheKey(paperID) for caching
	// For now, we'll fetch directly from arXiv

	// Construct the arXiv ID for searching
	// The paperID format is usually "2301.12345" or similar
	arxivID := fmt.Sprintf("id:%s", paperID)

	// Search by ID
	feed, err := s.client.SearchPapers(arxivID, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch paper: %w", err)
	}

	// Convert feed to papers
	papers := s.convertFeedToPapers(feed)

	if len(papers) == 0 {
		return nil, nil
	}

	return papers[0], nil
}

// convertFeedToPapers converts arXiv feed to paper models
func (s *arxivServiceImpl) convertFeedToPapers(feed *arxiv.Feed) []*model.Paper {
	papers := make([]*model.Paper, 0, len(feed.Entries))

	for _, entry := range feed.Entries {
		paperID := arxiv.ExtractID(entry.ID)
		paper := &model.Paper{
			ID:      paperID,
			Title:   cleanText(entry.Title),
			Summary: cleanText(entry.Summary),
		}

		// Parse authors
		paper.Authors = make([]string, len(entry.Authors))
		for i, author := range entry.Authors {
			paper.Authors[i] = author.Name
		}

		// Parse dates
		if published, err := arxiv.ParseTime(entry.Published); err == nil {
			paper.Published = published
		}
		if updated, err := arxiv.ParseTime(entry.Updated); err == nil {
			paper.Updated = updated
		}

		// Parse categories
		paper.Categories = make([]string, len(entry.Categories))
		for i, cat := range entry.Categories {
			paper.Categories[i] = cat.Term
		}
		if len(paper.Categories) > 0 {
			paper.PrimaryCategory = paper.Categories[0]
		}

		// Parse links
		for _, link := range entry.Links {
			switch link.Type {
			case "text/html":
				paper.ArxivURL = link.Href
			case "application/pdf":
				paper.PDFURL = link.Href
			}
		}

		// 生成图片 URL：arXiv 论文通常有 HTML 版本，其中包含第一张图
		// 格式：https://arxiv.org/html/{paper_id}/x1.png
		if paperID != "" {
			paper.ImageUrl = fmt.Sprintf("https://arxiv.org/html/%s/x1.png", paperID)
		}

		papers = append(papers, paper)
	}

	return papers
}

// buildPaperCacheKey builds a cache key for a single paper
func (s *arxivServiceImpl) buildPaperCacheKey(paperID string) string {
	return fmt.Sprintf("paper:%s", paperID)
}

// buildCacheKey builds a cache key for the category
func (s *arxivServiceImpl) buildCacheKey(category string) string {
	return fmt.Sprintf("papers:%s", category)
}

// cleanText cleans up text from arXiv
func cleanText(text string) string {
	// Remove extra whitespace
	return compactSpace(text)
}

// compactSpace compacts multiple spaces into one
func compactSpace(s string) string {
	lines := make([]string, 0)
	for _, line := range splitLines(s) {
		line = trimSpace(line)
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}
	return joinLines(lines, " ")
}

func splitLines(s string) []string {
	return strings.Fields(s)
}

func trimSpace(s string) string {
	return strings.TrimSpace(s)
}

func joinLines(lines []string, sep string) string {
	return strings.Join(lines, sep)
}
