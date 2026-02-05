package facade

import (
	"context"
	"time"

	"github.com/rrlian/papertok/backend/internal/core/arxiv"
	"github.com/rrlian/papertok/backend/internal/core/auth"
	"github.com/rrlian/papertok/backend/internal/features/paperfeed"
	"github.com/rrlian/papertok/backend/internal/features/papersearch"
	"github.com/rrlian/papertok/backend/internal/features/userauth"
	"github.com/rrlian/papertok/backend/internal/infra/cache"
	"github.com/rrlian/papertok/backend/internal/infra/database"
	"github.com/rrlian/papertok/backend/internal/infra/httpclient"
	paperRepo "github.com/rrlian/papertok/backend/internal/repository/paper"
	userRepo "github.com/rrlian/papertok/backend/internal/repository/user"
)

// Config holds the configuration for the Facade.
type Config struct {
	// ArXiv configuration
	ArxivBaseURL string
	HTTPTimeout  time.Duration
	CacheTTL     time.Duration
	CacheEnabled bool

	// Auth configuration
	JWTSecret      string
	JWTExpiresIn   time.Duration
	UseInMemoryAuth bool // If true, use in-memory repositories for testing

	// Database configuration
	DB database.DB // MySQL database connection (optional, nil for in-memory)
}

// Paper represents a paper in the API response.
// This is a unified type exposed by the Facade.
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
	ImageURL        string    `json:"imageUrl"`
}

// Facade is the unified entry point for all business operations.
// It manages feature instances and provides a clean API for handlers.
type Facade struct {
	paperFeedSvc   paperfeed.Service
	paperSearchSvc papersearch.Service
	userAuthSvc    *userauth.Impl
	authCoreSvc    auth.Service
}

// New creates a new Facade instance with all dependencies initialized.
func New(cfg Config) *Facade {
	// Initialize infrastructure
	httpClient := httpclient.NewClient(httpclient.Config{
		Timeout: cfg.HTTPTimeout,
	})

	var memCache cache.Cache
	if cfg.CacheEnabled {
		memCache = cache.NewMemoryCache()
	} else {
		memCache = &noopCache{}
	}

	// Initialize repositories
	paperRepository := paperRepo.NewMemoryRepository(memCache)

	// Initialize user repository
	var userRepository userRepo.Repository
	if cfg.UseInMemoryAuth {
		userRepository = userRepo.NewMemoryRepository()
	} else if cfg.DB != nil {
		userRepository = userRepo.NewSQLRepository(cfg.DB)
	} else {
		// Fallback to memory repository if no database is provided
		userRepository = userRepo.NewMemoryRepository()
	}

	// Initialize core services
	arxivSvc := arxiv.NewClient(arxiv.Config{
		BaseURL: cfg.ArxivBaseURL,
		Timeout: cfg.HTTPTimeout,
	}, httpClient)

	authCoreSvc, err := auth.New(auth.Config{
		Secret:             cfg.JWTSecret,
		AccessTokenExpiry:  cfg.JWTExpiresIn,
		RefreshTokenExpiry: cfg.JWTExpiresIn * 7, // 7x longer for refresh token
		Issuer:             "papertok",
		PasswordCost:       10,
	})
	if err != nil {
		panic(err) // In production, handle this gracefully
	}

	// Initialize features
	paperFeedSvc := paperfeed.New(arxivSvc, paperRepository, cfg.CacheTTL)
	paperSearchSvc := papersearch.New(arxivSvc)
	userAuthSvc := userauth.New(authCoreSvc, userRepository)

	return &Facade{
		paperFeedSvc:   paperFeedSvc,
		paperSearchSvc: paperSearchSvc,
		userAuthSvc:    userAuthSvc,
		authCoreSvc:    authCoreSvc,
	}
}

// GetPaperFeed fetches papers for the feed.
func (f *Facade) GetPaperFeed(ctx context.Context, category string, limit, offset int, sortBy string) ([]*Paper, error) {
	papers, err := f.paperFeedSvc.GetFeed(ctx, &paperfeed.FetchRequest{
		Category: category,
		Limit:    limit,
		Offset:   offset,
		SortBy:   sortBy,
	})
	if err != nil {
		return nil, err
	}

	return f.convertFeedPapers(papers), nil
}

// SearchPapers searches papers by keyword.
func (f *Facade) SearchPapers(ctx context.Context, query string, limit int) ([]*Paper, error) {
	papers, err := f.paperSearchSvc.Search(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	return f.convertSearchPapers(papers), nil
}

// GetPaperByID retrieves a single paper by ID.
func (f *Facade) GetPaperByID(ctx context.Context, id string) (*Paper, error) {
	paper, err := f.paperSearchSvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if paper == nil {
		return nil, nil
	}

	return f.convertSearchPaper(paper), nil
}

// UserAuth returns the user authentication service.
func (f *Facade) UserAuth() *userauth.Impl {
	return f.userAuthSvc
}

// AuthCore returns the core authentication service.
func (f *Facade) AuthCore() auth.Service {
	return f.authCoreSvc
}

// convertFeedPapers converts paperfeed.Paper to facade.Paper.
func (f *Facade) convertFeedPapers(papers []*paperfeed.Paper) []*Paper {
	result := make([]*Paper, len(papers))
	for i, p := range papers {
		result[i] = &Paper{
			ID:              p.ID,
			Title:           p.Title,
			Authors:         p.Authors,
			Summary:         p.Summary,
			Published:       p.Published,
			Updated:         p.Updated,
			Categories:      p.Categories,
			PrimaryCategory: p.PrimaryCategory,
			ArxivURL:        p.ArxivURL,
			PDFURL:          p.PDFURL,
			ImageURL:        p.ImageURL,
		}
	}
	return result
}

// convertSearchPapers converts papersearch.Paper to facade.Paper.
func (f *Facade) convertSearchPapers(papers []*papersearch.Paper) []*Paper {
	result := make([]*Paper, len(papers))
	for i, p := range papers {
		result[i] = f.convertSearchPaper(p)
	}
	return result
}

// convertSearchPaper converts a single papersearch.Paper to facade.Paper.
func (f *Facade) convertSearchPaper(p *papersearch.Paper) *Paper {
	return &Paper{
		ID:              p.ID,
		Title:           p.Title,
		Authors:         p.Authors,
		Summary:         p.Summary,
		Published:       p.Published,
		Updated:         p.Updated,
		Categories:      p.Categories,
		PrimaryCategory: p.PrimaryCategory,
		ArxivURL:        p.ArxivURL,
		PDFURL:          p.PDFURL,
		ImageURL:        p.ImageURL,
	}
}

// noopCache is a no-op cache implementation for when caching is disabled.
type noopCache struct{}

func (n *noopCache) Get(key string) (interface{}, bool) { return nil, false }
func (n *noopCache) Set(key string, value interface{}, ttl time.Duration) {}
func (n *noopCache) Delete(key string) {}
func (n *noopCache) Clear() {}
