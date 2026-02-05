package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rrlian/papertok/backend/internal/api/handlers"
	"github.com/rrlian/papertok/backend/internal/api/middleware"
	"github.com/rrlian/papertok/backend/internal/config"
	"github.com/rrlian/papertok/backend/internal/facade"
	"github.com/rrlian/papertok/backend/internal/infra/database"
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

	// Initialize database if configured
	var db database.DB
	useInMemoryAuth := true // Default to in-memory

	if cfg.Database.Host != "" && cfg.Database.Host != "localhost" {
		// Only use MySQL if a non-localhost host is configured
		// This prevents accidentally requiring MySQL in development
		connector, err := database.New(database.Config{
			Host:         cfg.Database.Host,
			Port:         cfg.Database.Port,
			Username:     cfg.Database.Username,
			Password:     cfg.Database.Password,
			Database:     cfg.Database.Database,
			MaxOpenConns: cfg.Database.MaxOpenConns,
			MaxIdleConns: cfg.Database.MaxIdleConns,
			MaxLifetime:  cfg.Database.MaxLifetime,
		})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		db = connector.DB()
		defer connector.Close()
		useInMemoryAuth = false
		log.Printf("Connected to MySQL database at %s:%d", cfg.Database.Host, cfg.Database.Port)
	} else {
		log.Println("Using in-memory authentication (no database configured)")
	}

	// Initialize Facade with all dependencies
	f := facade.New(facade.Config{
		ArxivBaseURL:   cfg.Arxiv.BaseURL,
		HTTPTimeout:    cfg.Arxiv.Timeout,
		CacheTTL:       cfg.Cache.TTL,
		CacheEnabled:   cfg.Cache.Enabled,
		JWTSecret:      cfg.JWT.Secret,
		JWTExpiresIn:   cfg.JWT.ExpiresIn,
		UseInMemoryAuth: useInMemoryAuth,
		DB:             db,
	})

	// Create handlers
	paperHandler := handlers.NewPaperHandler(f)
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(f.UserAuth())

	// Create router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	// Register public routes
	router.GET("/health", healthHandler.HealthCheck)

	// Auth routes (public)
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/register", authHandler.RegisterHandler)
		authGroup.POST("/login", authHandler.LoginHandler)
		authGroup.POST("/refresh", authHandler.RefreshTokenHandler)

		// Protected auth routes
		protected := authGroup.Group("")
		protected.Use(middleware.AuthMiddleware(f.AuthCore()))
		{
			protected.GET("/profile", authHandler.GetProfileHandler)
		}
	}

	// Paper routes (public for now, can be protected later)
	api := router.Group("/api/v1")
	{
		api.GET("/papers", paperHandler.GetPapers)
		api.GET("/papers/search", paperHandler.SearchPapers)
		api.GET("/papers/:id", paperHandler.GetPaperByID)
	}

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting PaperTok API server on %s", addr)
	log.Printf("Available routes:")
	log.Printf("  GET  /health")
	log.Printf("  POST /api/v1/auth/register")
	log.Printf("  POST /api/v1/auth/login")
	log.Printf("  POST /api/v1/auth/refresh")
	log.Printf("  GET  /api/v1/auth/profile (requires auth)")
	log.Printf("  GET  /api/v1/papers")
	log.Printf("  GET  /api/v1/papers/search")
	log.Printf("  GET  /api/v1/papers/:id")

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
