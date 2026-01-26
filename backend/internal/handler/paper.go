package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rrlian/papertok/backend/internal/model"
	"github.com/rrlian/papertok/backend/internal/service"
)

// PaperHandler handles paper-related requests
type PaperHandler struct {
	arxivService service.ArxivService
}

// NewPaperHandler creates a new paper handler
func NewPaperHandler(arxivService service.ArxivService) *PaperHandler {
	return &PaperHandler{
		arxivService: arxivService,
	}
}

// GetPapers handles GET /api/v1/papers
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

	// Build request
	req := &model.FetchRequest{
		Category:   category,
		MaxResults: limit,
		SortBy:     sortBy,
		Offset:     offset,
	}

	// Fetch papers
	papers, err := h.arxivService.FetchPapers(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error: &model.ErrorInfo{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to fetch papers from arXiv",
				Details: err.Error(),
			},
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data: model.PapersResponse{
			Papers:   papers,
			Total:    len(papers),
			Page:     offset/limit + 1,
			PageSize: limit,
		},
		Timestamp: time.Now().Unix(),
	})
}

// SearchPapers handles GET /api/v1/papers/search
func (h *PaperHandler) SearchPapers(c *gin.Context) {
	// Parse query parameters
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error: &model.ErrorInfo{
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

	// Search papers
	papers, err := h.arxivService.SearchPapers(c.Request.Context(), query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error: &model.ErrorInfo{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to search papers",
				Details: err.Error(),
			},
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data: model.PapersResponse{
			Papers:   papers,
			Total:    len(papers),
			Page:     1,
			PageSize: limit,
		},
		Timestamp: time.Now().Unix(),
	})
}

// GetPaperByID handles GET /api/v1/papers/:id
func (h *PaperHandler) GetPaperByID(c *gin.Context) {
	paperID := c.Param("id")
	if paperID == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error: &model.ErrorInfo{
				Code:    "INVALID_PARAMS",
				Message: "Paper ID is required",
			},
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// Get paper by ID
	paper, err := h.arxivService.GetPaperByID(c.Request.Context(), paperID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error: &model.ErrorInfo{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to fetch paper",
				Details: err.Error(),
			},
			Timestamp: time.Now().Unix(),
		})
		return
	}

	if paper == nil {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Error: &model.ErrorInfo{
				Code:    "NOT_FOUND",
				Message: "Paper not found",
			},
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, model.APIResponse{
		Success:   true,
		Data:      paper,
		Timestamp: time.Now().Unix(),
	})
}
