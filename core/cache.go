package core

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CacheItem represents an item in the cache
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// Cache is a thread-safe in-memory cache
type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

// NewCache creates a new cache
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

// Set adds an item to the cache with an expiration time
func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(expiration),
	}
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, exists := c.items[key]
	c.mu.RUnlock()

	if !exists {
		return nil, false
	}

	if time.Now().After(item.Expiration) {
		c.Delete(key)
		return nil, false
	}

	return item.Value, true
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]CacheItem)
}

// CacheMiddleware creates a middleware that caches responses
func CacheMiddleware(cache *Cache, expiration time.Duration) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Generate cache key from request
		cacheKey := ctx.Path() + ":" + ctx.Method()

		// Try to get from cache
		if cached, exists := cache.Get(cacheKey); exists {
			return ctx.JSON(cached)
		}

		// Continue to handler
		err := ctx.Next()
		if err != nil {
			return err
		}

		// Cache the response
		cache.Set(cacheKey, ctx.Response().Body(), expiration)

		return nil
	}
} 