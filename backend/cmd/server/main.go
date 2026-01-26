package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rrlian/papertok/backend/internal/config"
	"github.com/rrlian/papertok/backend/internal/handler"
	"github.com/rrlian/papertok/backend/internal/middleware"
	"github.com/rrlian/papertok/backend/internal/model"
	"github.com/rrlian/papertok/backend/internal/service"
	"github.com/rrlian/papertok/backend/pkg/arxiv"
)

func main() {
	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create services
	arxivClient := arxiv.NewClient(cfg.Arxiv.BaseURL, cfg.Arxiv.Timeout)
	var cacheService service.CacheService
	if cfg.Cache.Enabled {
		cacheService = service.NewMemoryCache()
	} else {
		cacheService = &noopCache{}
	}
	arxivService := service.NewArxivService(arxivClient, cacheService, cfg.Cache.TTL)

	// Create handlers
	paperHandler := handler.NewPaperHandler(arxivService)
	healthHandler := handler.NewHealthHandler()

	// Create router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	// Register routes
	router.GET("/health", healthHandler.HealthCheck)

	api := router.Group("/api/v1")
	{
		api.GET("/papers", paperHandler.GetPapers)
		api.GET("/papers/search", paperHandler.SearchPapers)
		api.GET("/papers/:id", paperHandler.GetPaperByID)
	}

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting PaperTok API server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// noopCache is a no-op cache implementation
type noopCache struct{}

func (n *noopCache) Get(key string) ([]*model.Paper, bool) {
	return nil, false
}

func (n *noopCache) Set(key string, papers []*model.Paper, ttl time.Duration) {
}

func (n *noopCache) Invalidate(key string) {
}

func (n *noopCache) Clear() {
}
