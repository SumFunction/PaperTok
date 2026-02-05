package arxiv

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Config holds the configuration for the arXiv client.
type Config struct {
	BaseURL string
	Timeout time.Duration
}

// Client implements the arXiv Service interface.
type Client struct {
	baseURL    string
	httpClient httpClient
}

// Ensure Client implements Service interface
var _ Service = (*Client)(nil)

// NewClient creates a new arXiv client.
func NewClient(cfg Config, client httpClient) *Client {
	return &Client{
		baseURL:    cfg.BaseURL,
		httpClient: client,
	}
}

// FetchByCategory fetches papers from arXiv by category.
func (c *Client) FetchByCategory(ctx context.Context, req *FetchRequest) ([]*Paper, error) {
	// Build query parameters
	params := url.Values{}

	// Build search query
	var searchQuery string
	if req.Category == "" {
		// Use wildcard search to get latest papers
		searchQuery = "cat:cs.* OR cat:stat.* OR cat:math.*"
	} else {
		searchQuery = fmt.Sprintf("cat:%s", req.Category)
	}
	params.Add("search_query", searchQuery)
	params.Add("sortBy", req.SortBy)
	params.Add("sortOrder", "descending")
	params.Add("max_results", fmt.Sprintf("%d", req.MaxResults))
	if req.Offset > 0 {
		params.Add("start", fmt.Sprintf("%d", req.Offset))
	}

	// Make request
	reqURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())
	resp, err := c.httpClient.Get(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrFetchFailed, resp.StatusCode)
	}

	// Parse response
	feed, err := c.parseFeed(resp.Body)
	if err != nil {
		return nil, err
	}

	return c.convertFeedToPapers(feed), nil
}

// Search searches papers by keyword.
func (c *Client) Search(ctx context.Context, query string, limit int) ([]*Paper, error) {
	// Build query parameters
	params := url.Values{}
	params.Add("search_query", fmt.Sprintf("all:%s", query))
	params.Add("sortBy", "submittedDate")
	params.Add("sortOrder", "descending")
	params.Add("max_results", fmt.Sprintf("%d", limit))

	// Make request
	reqURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())
	resp, err := c.httpClient.Get(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSearchFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrSearchFailed, resp.StatusCode)
	}

	// Parse response
	feed, err := c.parseFeed(resp.Body)
	if err != nil {
		return nil, err
	}

	return c.convertFeedToPapers(feed), nil
}

// GetByID fetches a single paper by its arXiv ID.
func (c *Client) GetByID(ctx context.Context, id string) (*Paper, error) {
	// Search by ID
	papers, err := c.Search(ctx, fmt.Sprintf("id:%s", id), 1)
	if err != nil {
		return nil, err
	}

	if len(papers) == 0 {
		return nil, nil
	}

	return papers[0], nil
}

// parseFeed parses the arXiv Atom feed from the response body.
func (c *Client) parseFeed(body io.Reader) (*Feed, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read response: %v", ErrInvalidResponse, err)
	}

	var feed Feed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, fmt.Errorf("%w: failed to parse XML: %v", ErrInvalidResponse, err)
	}

	return &feed, nil
}

// convertFeedToPapers converts arXiv feed entries to Paper structs.
func (c *Client) convertFeedToPapers(feed *Feed) []*Paper {
	papers := make([]*Paper, 0, len(feed.Entries))

	for _, entry := range feed.Entries {
		paperID := extractID(entry.ID)
		paper := &Paper{
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
		if published, err := parseTime(entry.Published); err == nil {
			paper.Published = published
		}
		if updated, err := parseTime(entry.Updated); err == nil {
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

		// Generate image URL
		if paperID != "" {
			paper.ImageURL = fmt.Sprintf("https://arxiv.org/html/%s/x1.png", paperID)
		}

		papers = append(papers, paper)
	}

	return papers
}

// parseTime parses an arXiv time string.
func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

// extractID extracts arXiv ID from URL.
func extractID(urlStr string) string {
	parts := strings.Split(urlStr, "/")
	if len(parts) > 0 {
		id := parts[len(parts)-1]
		// Remove version if present (e.g., v1, v2)
		if idx := strings.Index(id, "v"); idx > 0 {
			id = id[:idx]
		}
		return id
	}
	return urlStr
}

// cleanText removes extra whitespace from text.
func cleanText(text string) string {
	return strings.Join(strings.Fields(text), " ")
}
