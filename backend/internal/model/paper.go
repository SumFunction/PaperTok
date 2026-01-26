package model

import "time"

// Paper represents an arXiv paper
type Paper struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Authors         []string  `json:"authors"`
	Summary         string    `json:"summary"`
	Published       time.Time `json:"published"`
	Updated         time.Time `json:"updated"`
	Categories      []string  `json:"categories"`
	PrimaryCategory string    `json:"primaryCategory"`
	ArxivURL        string    `json:"arxivUrl"`
	PDFURL          string    `json:"pdfUrl"`
	ImageUrl        string    `json:"imageUrl"` // 封面图 URL
}

// FetchRequest defines parameters for fetching papers
type FetchRequest struct {
	Category   string
	MaxResults int
	SortBy     string // "lastUpdatedDate" or "submittedDate"
	Offset     int    // 分页偏移量（用于实现无限滚动）
}

// PaperDetailRequest defines parameters for fetching a single paper
type PaperDetailRequest struct {
	PaperID string // 论文 ID（例如：2301.12345）
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp int64       `json:"timestamp"` // Unix timestamp
}

// ErrorInfo represents detailed error information
type ErrorInfo struct {
	Code    string `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
	Details string `json:"details,omitempty"` // 详细错误信息
}

// PapersResponse represents the response for papers list
type PapersResponse struct {
	Papers   []*Paper `json:"papers"`
	Total    int      `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
}
