package arxiv

import (
	"encoding/xml"
	"time"
)

// Paper represents a paper fetched from arXiv.
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

// FetchRequest contains parameters for fetching papers.
type FetchRequest struct {
	Category   string // arXiv category (e.g., "cs.AI")
	MaxResults int    // Maximum number of results
	SortBy     string // Sort by: "lastUpdatedDate" or "submittedDate"
	Offset     int    // Pagination offset
}

// Feed represents the arXiv Atom feed response.
type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Entries []Entry  `xml:"entry"`
}

// Entry represents an arXiv paper entry in the feed.
type Entry struct {
	ID         string     `xml:"id"`
	Published  string     `xml:"published"`
	Updated    string     `xml:"updated"`
	Title      string     `xml:"title"`
	Summary    string     `xml:"summary"`
	Authors    []Author   `xml:"author"`
	Links      []Link     `xml:"link"`
	Categories []Category `xml:"category"`
}

// Author represents a paper author.
type Author struct {
	Name string `xml:"name"`
}

// Link represents a link in the paper entry.
type Link struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

// Category represents a paper category.
type Category struct {
	Term string `xml:"term,attr"`
}
