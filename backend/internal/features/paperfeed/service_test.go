package paperfeed

import (
	"context"
	"testing"
	"time"

	"github.com/rrlian/papertok/backend/internal/core/arxiv"
	paperRepo "github.com/rrlian/papertok/backend/internal/repository/paper"
)

// mockArxivService is a mock implementation for testing.
type mockArxivService struct {
	papers []*arxiv.Paper
	err    error
}

func (m *mockArxivService) FetchByCategory(ctx context.Context, req *arxiv.FetchRequest) ([]*arxiv.Paper, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.papers, nil
}

func (m *mockArxivService) Search(ctx context.Context, query string, limit int) ([]*arxiv.Paper, error) {
	return m.papers, m.err
}

func (m *mockArxivService) GetByID(ctx context.Context, id string) (*arxiv.Paper, error) {
	if len(m.papers) > 0 {
		return m.papers[0], m.err
	}
	return nil, m.err
}

// mockPaperRepository is a mock implementation for testing.
type mockPaperRepository struct {
	papers map[string][]*paperRepo.Paper
}

func newMockPaperRepository() *mockPaperRepository {
	return &mockPaperRepository{
		papers: make(map[string][]*paperRepo.Paper),
	}
}

func (m *mockPaperRepository) GetByCategory(ctx context.Context, category string) ([]*paperRepo.Paper, bool) {
	papers, found := m.papers[category]
	return papers, found
}

func (m *mockPaperRepository) SaveByCategory(ctx context.Context, category string, papers []*paperRepo.Paper, ttl time.Duration) {
	m.papers[category] = papers
}

func (m *mockPaperRepository) GetByID(ctx context.Context, id string) (*paperRepo.Paper, bool) {
	return nil, false
}

func (m *mockPaperRepository) Save(ctx context.Context, paper *paperRepo.Paper, ttl time.Duration) {
}

func (m *mockPaperRepository) InvalidateCategory(ctx context.Context, category string) {
	delete(m.papers, category)
}

func (m *mockPaperRepository) Clear(ctx context.Context) {
	m.papers = make(map[string][]*paperRepo.Paper)
}

func TestImpl_GetFeed_FromArxiv(t *testing.T) {
	// Arrange
	mockArxiv := &mockArxivService{
		papers: []*arxiv.Paper{
			{
				ID:              "2301.12345",
				Title:           "Test Paper",
				Authors:         []string{"John Doe"},
				Summary:         "Test summary",
				PrimaryCategory: "cs.AI",
			},
		},
	}
	mockRepo := newMockPaperRepository()
	svc := New(mockArxiv, mockRepo, 5*time.Minute)

	// Act
	papers, err := svc.GetFeed(context.Background(), &FetchRequest{
		Category: "cs.AI",
		Limit:    10,
		SortBy:   "lastUpdatedDate",
	})

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(papers) != 1 {
		t.Errorf("Expected 1 paper, got: %d", len(papers))
	}
	if papers[0].ID != "2301.12345" {
		t.Errorf("Expected ID '2301.12345', got: %s", papers[0].ID)
	}
}

func TestImpl_GetFeed_FromCache(t *testing.T) {
	// Arrange
	mockArxiv := &mockArxivService{
		papers: []*arxiv.Paper{}, // Empty - should not be called
	}
	mockRepo := newMockPaperRepository()
	mockRepo.papers["cs.AI"] = []*paperRepo.Paper{
		{
			ID:              "cached-paper",
			Title:           "Cached Paper",
			PrimaryCategory: "cs.AI",
		},
	}
	svc := New(mockArxiv, mockRepo, 5*time.Minute)

	// Act
	papers, err := svc.GetFeed(context.Background(), &FetchRequest{
		Category: "cs.AI",
		Limit:    10,
		SortBy:   "lastUpdatedDate",
	})

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(papers) != 1 {
		t.Errorf("Expected 1 paper, got: %d", len(papers))
	}
	if papers[0].ID != "cached-paper" {
		t.Errorf("Expected ID 'cached-paper', got: %s", papers[0].ID)
	}
}
