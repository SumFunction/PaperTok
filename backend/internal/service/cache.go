package service

import (
	"sync"
	"time"

	"github.com/rrlian/papertok/backend/internal/model"
)

// CacheService defines the interface for cache service
type CacheService interface {
	Get(key string) ([]*model.Paper, bool)
	Set(key string, papers []*model.Paper, ttl time.Duration)
	Invalidate(key string)
	Clear()
}

// cacheItem represents a cached item with expiration
type cacheItem struct {
	papers     []*model.Paper
	expiration time.Time
}

// memoryCache implements an in-memory cache
type memoryCache struct {
	items map[string]*cacheItem
	mu    sync.RWMutex
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache() CacheService {
	cache := &memoryCache{
		items: make(map[string]*cacheItem),
	}

	// Start cleanup goroutine
	go cache.cleanupExpired()

	return cache
}

// Get retrieves papers from cache
func (c *memoryCache) Get(key string) ([]*model.Paper, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	// Check if expired
	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.papers, true
}

// Set stores papers in cache
func (c *memoryCache) Set(key string, papers []*model.Paper, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &cacheItem{
		papers:     papers,
		expiration: time.Now().Add(ttl),
	}
}

// Invalidate removes a specific key from cache
func (c *memoryCache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear removes all items from cache
func (c *memoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*cacheItem)
}

// cleanupExpired periodically removes expired items
func (c *memoryCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.expiration) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
