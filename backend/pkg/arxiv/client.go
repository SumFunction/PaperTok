package arxiv

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client represents an arXiv API client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new arXiv API client
func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// FetchPapers fetches papers from arXiv API
func (c *Client) FetchPapers(category string, maxResults int, sortBy string) (*Feed, error) {
	// Build query parameters
	params := url.Values{}
	// 如果 category 为空，使用通用搜索获取最新论文
	// 使用 lastUpdatedDate 或 submittedDate 排序获取最新提交
	var searchQuery string
	if category == "" {
		// 使用通配符搜索，获取所有最新论文
		searchQuery = "cat:cs.* OR cat:stat.* OR cat:math.*"
	} else {
		searchQuery = fmt.Sprintf("cat:%s", category)
	}
	params.Add("search_query", searchQuery)
	params.Add("sortBy", sortBy)
	params.Add("sortOrder", "descending")
	params.Add("max_results", fmt.Sprintf("%d", maxResults))

	// Make request
	reqURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())
	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from arXiv: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("arXiv API returned status %d", resp.StatusCode)
	}

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var feed Feed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("failed to parse XML response: %w", err)
	}

	return &feed, nil
}

// FetchPapersWithOffset fetches papers from arXiv API with offset/pagination support
func (c *Client) FetchPapersWithOffset(category string, maxResults int, sortBy string, offset int) (*Feed, error) {
	// Build query parameters
	params := url.Values{}
	// 如果 category 为空，使用通用搜索获取最新论文
	var searchQuery string
	if category == "" {
		searchQuery = "cat:cs.* OR cat:stat.* OR cat:math.*"
	} else {
		searchQuery = fmt.Sprintf("cat:%s", category)
	}
	params.Add("search_query", searchQuery)
	params.Add("sortBy", sortBy)
	params.Add("sortOrder", "descending")
	params.Add("max_results", fmt.Sprintf("%d", maxResults))
	params.Add("start", fmt.Sprintf("%d", offset)) // 添加分页支持

	// Make request
	reqURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())
	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from arXiv: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("arXiv API returned status %d", resp.StatusCode)
	}

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var feed Feed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("failed to parse XML response: %w", err)
	}

	return &feed, nil
}

// SearchPapers searches papers by keyword
func (c *Client) SearchPapers(query string, maxResults int) (*Feed, error) {
	// Build query parameters
	params := url.Values{}
	params.Add("search_query", fmt.Sprintf("all:%s", query))
	params.Add("sortBy", "submittedDate")
	params.Add("sortOrder", "descending")
	params.Add("max_results", fmt.Sprintf("%d", maxResults))

	// Make request
	reqURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())
	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search arXiv: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("arXiv API returned status %d", resp.StatusCode)
	}

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var feed Feed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("failed to parse XML response: %w", err)
	}

	return &feed, nil
}

// Feed represents the arXiv Atom feed
type Feed struct {
	XMLName  xml.Name `xml:"feed"`
	Entries  []Entry `xml:"entry"`
}

// Entry represents an arXiv paper entry
type Entry struct {
	ID        string   `xml:"id"`
	Published string   `xml:"published"`
	Updated   string   `xml:"updated"`
	Title     string   `xml:"title"`
	Summary   string   `xml:"summary"`
	Authors   []Author `xml:"author"`
	Links     []Link   `xml:"link"`
	Categories []Category `xml:"category"`
}

// Author represents an author
type Author struct {
	Name string `xml:"name"`
}

// Link represents a link
type Link struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

// Category represents a category
type Category struct {
	Term string `xml:"term,attr"`
}

// ParseTime parses an arXiv time string
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

// ExtractID extracts arXiv ID from URL
func ExtractID(urlStr string) string {
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
