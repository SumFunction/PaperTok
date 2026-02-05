package cache

import "time"

// Cache defines the interface for caching operations.
// This is a generic cache interface that can be implemented
// with different backends (memory, redis, etc.)
type Cache interface {
	// Get retrieves a value from the cache.
	// Returns the value and true if found, nil and false otherwise.
	Get(key string) (interface{}, bool)

	// Set stores a value in the cache with the given TTL.
	Set(key string, value interface{}, ttl time.Duration)

	// Delete removes a value from the cache.
	Delete(key string)

	// Clear removes all values from the cache.
	Clear()
}
