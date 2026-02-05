package arxiv

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPClient is a mock implementation of httpClient for testing.
type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

func TestClient_FetchByCategory(t *testing.T) {
	// Arrange
	xmlResponse := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <entry>
    <id>http://arxiv.org/abs/2301.12345</id>
    <published>2023-01-15T10:00:00Z</published>
    <updated>2023-01-16T10:00:00Z</updated>
    <title>Test Paper Title</title>
    <summary>This is a test summary.</summary>
    <author><name>John Doe</name></author>
    <author><name>Jane Smith</name></author>
    <category term="cs.AI"/>
    <category term="cs.LG"/>
    <link href="http://arxiv.org/abs/2301.12345" type="text/html"/>
    <link href="http://arxiv.org/pdf/2301.12345.pdf" type="application/pdf"/>
  </entry>
</feed>`

	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(xmlResponse)),
		},
	}

	client := NewClient(Config{BaseURL: "http://test.com"}, mockClient)

	// Act
	papers, err := client.FetchByCategory(context.Background(), &FetchRequest{
		Category:   "cs.AI",
		MaxResults: 10,
		SortBy:     "lastUpdatedDate",
	})

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(papers) != 1 {
		t.Errorf("Expected 1 paper, got: %d", len(papers))
	}

	paper := papers[0]
	if paper.ID != "2301.12345" {
		t.Errorf("Expected ID '2301.12345', got: %s", paper.ID)
	}
	if paper.Title != "Test Paper Title" {
		t.Errorf("Expected title 'Test Paper Title', got: %s", paper.Title)
	}
	if len(paper.Authors) != 2 {
		t.Errorf("Expected 2 authors, got: %d", len(paper.Authors))
	}
	if paper.PrimaryCategory != "cs.AI" {
		t.Errorf("Expected primary category 'cs.AI', got: %s", paper.PrimaryCategory)
	}
}

func TestClient_FetchByCategory_Error(t *testing.T) {
	// Arrange
	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("")),
		},
	}

	client := NewClient(Config{BaseURL: "http://test.com"}, mockClient)

	// Act
	_, err := client.FetchByCategory(context.Background(), &FetchRequest{
		Category:   "cs.AI",
		MaxResults: 10,
		SortBy:     "lastUpdatedDate",
	})

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !IsFetchFailed(err) {
		t.Errorf("Expected ErrFetchFailed, got: %v", err)
	}
}

func TestCleanText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello   world  ", "hello world"},
		{"line1\nline2\nline3", "line1 line2 line3"},
		{"  multiple   spaces   ", "multiple spaces"},
	}

	for _, tt := range tests {
		result := cleanText(tt.input)
		if result != tt.expected {
			t.Errorf("cleanText(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestExtractID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http://arxiv.org/abs/2301.12345v1", "2301.12345"},
		{"http://arxiv.org/abs/2301.12345", "2301.12345"},
		{"2301.12345", "2301.12345"},
	}

	for _, tt := range tests {
		result := extractID(tt.input)
		if result != tt.expected {
			t.Errorf("extractID(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}
