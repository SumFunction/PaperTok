package paper

import (
	"context"
	"fmt"
	"time"

	"github.com/rrlian/papertok/backend/internal/infra/cache"
)

// MemoryRepository implements Repository using in-memory cache.
type MemoryRepository struct {
	cache cache.Cache
}

// Ensure MemoryRepository implements Repository interface
var _ Repository = (*MemoryRepository)(nil)

// NewMemoryRepository creates a new memory-based paper repository.
func NewMemoryRepository(c cache.Cache) *MemoryRepository {
	return &MemoryRepository{
		cache: c,
	}
}

// GetByCategory retrieves papers by category from cache.
func (r *MemoryRepository) GetByCategory(ctx context.Context, category string) ([]*Paper, bool) {
	key := r.categoryKey(category)
	value, found := r.cache.Get(key)
	if !found {
		return nil, false
	}

	papers, ok := value.([]*Paper)
	if !ok {
		return nil, false
	}

	return papers, true
}

// SaveByCategory stores papers for a category in cache.
func (r *MemoryRepository) SaveByCategory(ctx context.Context, category string, papers []*Paper, ttl time.Duration) {
	key := r.categoryKey(category)
	r.cache.Set(key, papers, ttl)
}

// GetByID retrieves a single paper by ID from cache.
func (r *MemoryRepository) GetByID(ctx context.Context, id string) (*Paper, bool) {
	key := r.paperKey(id)
	value, found := r.cache.Get(key)
	if !found {
		return nil, false
	}

	paper, ok := value.(*Paper)
	if !ok {
		return nil, false
	}

	return paper, true
}

// Save stores a single paper in cache.
func (r *MemoryRepository) Save(ctx context.Context, paper *Paper, ttl time.Duration) {
	key := r.paperKey(paper.ID)
	r.cache.Set(key, paper, ttl)
}

// InvalidateCategory removes cached papers for a category.
func (r *MemoryRepository) InvalidateCategory(ctx context.Context, category string) {
	key := r.categoryKey(category)
	r.cache.Delete(key)
}

// Clear removes all cached papers.
func (r *MemoryRepository) Clear(ctx context.Context) {
	r.cache.Clear()
}

// categoryKey generates a cache key for category-based paper lists.
func (r *MemoryRepository) categoryKey(category string) string {
	return fmt.Sprintf("papers:category:%s", category)
}

// paperKey generates a cache key for a single paper.
func (r *MemoryRepository) paperKey(id string) string {
	return fmt.Sprintf("papers:id:%s", id)
}
