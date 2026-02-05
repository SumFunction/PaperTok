package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rrlian/papertok/backend/internal/facade"
)

// APIResponse represents a standard API response.
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// ErrorInfo represents detailed error information.
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// PapersResponse represents the response for papers list.
type PapersResponse struct {
	Papers   interface{} `json:"papers"`
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// PaperHandler handles paper-related requests.
type PaperHandler struct {
	facade *facade.Facade
}

// NewPaperHandler creates a new paper handler.
func NewPaperHandler(f *facade.Facade) *PaperHandler {
	return &PaperHandler{
		facade: f,
	}
}

// GetPapers handles GET /api/v1/papers.
func (h *PaperHandler) GetPapers(c *gin.Context) {
	// Parse query parameters
	category := c.DefaultQuery("category", "cs.AI")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	sortBy := c.DefaultQuery("sort_by", "lastUpdatedDate")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Fetch papers via facade
	papers, err := h.facade.GetPaperFeed(c.Request.Context(), category, limit, offset, sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: &ErrorInfo{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to fetch papers from arXiv",
				Details: err.Error(),
			},
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: PapersResponse{
			Papers:   papers,
			Total:    len(papers),
			Page:     offset/limit + 1,
			PageSize: limit,
		},
		Timestamp: time.Now().Unix(),
	})
}

// SearchPapers handles GET /api/v1/papers/search.
func (h *PaperHandler) SearchPapers(c *gin.Context) {
	// Parse query parameters
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error: &ErrorInfo{
				Code:    "INVALID_PARAMS",
				Message: "Query parameter 'query' is required",
			},
			Timestamp: time.Now().Unix(),
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	// Search papers via facade
	papers, err := h.facade.SearchPapers(c.Request.Context(), query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: &ErrorInfo{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to search papers",
				Details: err.Error(),
			},
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: PapersResponse{
			Papers:   papers,
			Total:    len(papers),
			Page:     1,
			PageSize: limit,
		},
		Timestamp: time.Now().Unix(),
	})
}

// GetPaperByID handles GET /api/v1/papers/:id.
func (h *PaperHandler) GetPaperByID(c *gin.Context) {
	paperID := c.Param("id")
	if paperID == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error: &ErrorInfo{
				Code:    "INVALID_PARAMS",
				Message: "Paper ID is required",
			},
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// Get paper by ID via facade
	paper, err := h.facade.GetPaperByID(c.Request.Context(), paperID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: &ErrorInfo{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to fetch paper",
				Details: err.Error(),
			},
			Timestamp: time.Now().Unix(),
		})
		return
	}

	if paper == nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error: &ErrorInfo{
				Code:    "NOT_FOUND",
				Message: "Paper not found",
			},
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      paper,
		Timestamp: time.Now().Unix(),
	})
}
