package papersearch

import (
	"context"
	"testing"

	"github.com/rrlian/papertok/backend/internal/core/arxiv"
)

// mockArxivService is a mock implementation for testing.
type mockArxivService struct {
	searchPapers []*arxiv.Paper
	getPaper     *arxiv.Paper
	err          error
}

func (m *mockArxivService) FetchByCategory(ctx context.Context, req *arxiv.FetchRequest) ([]*arxiv.Paper, error) {
	return m.searchPapers, m.err
}

func (m *mockArxivService) Search(ctx context.Context, query string, limit int) ([]*arxiv.Paper, error) {
	return m.searchPapers, m.err
}

func (m *mockArxivService) GetByID(ctx context.Context, id string) (*arxiv.Paper, error) {
	return m.getPaper, m.err
}

func TestImpl_Search(t *testing.T) {
	// Arrange
	mockArxiv := &mockArxivService{
		searchPapers: []*arxiv.Paper{
			{
				ID:              "2301.12345",
				Title:           "Machine Learning Paper",
				Authors:         []string{"Jane Doe"},
				Summary:         "A paper about ML",
				PrimaryCategory: "cs.LG",
			},
			{
				ID:              "2301.67890",
				Title:           "Deep Learning Paper",
				Authors:         []string{"John Smith"},
				Summary:         "A paper about DL",
				PrimaryCategory: "cs.LG",
			},
		},
	}
	svc := New(mockArxiv)

	// Act
	papers, err := svc.Search(context.Background(), "machine learning", 10)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(papers) != 2 {
		t.Errorf("Expected 2 papers, got: %d", len(papers))
	}
}

func TestImpl_GetByID(t *testing.T) {
	// Arrange
	mockArxiv := &mockArxivService{
		getPaper: &arxiv.Paper{
			ID:              "2301.12345",
			Title:           "Specific Paper",
			Authors:         []string{"Author"},
			Summary:         "Paper summary",
			PrimaryCategory: "cs.AI",
		},
	}
	svc := New(mockArxiv)

	// Act
	paper, err := svc.GetByID(context.Background(), "2301.12345")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if paper == nil {
		t.Error("Expected paper, got nil")
	}
	if paper.ID != "2301.12345" {
		t.Errorf("Expected ID '2301.12345', got: %s", paper.ID)
	}
}

func TestImpl_GetByID_NotFound(t *testing.T) {
	// Arrange
	mockArxiv := &mockArxivService{
		getPaper: nil,
	}
	svc := New(mockArxiv)

	// Act
	paper, err := svc.GetByID(context.Background(), "nonexistent")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if paper != nil {
		t.Errorf("Expected nil paper, got: %v", paper)
	}
}
