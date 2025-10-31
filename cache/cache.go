package cache

import (
	"context"
	"encoding/json"
	"time"

	"nutrition-platform/config"

	"github.com/go-redis/redis/v8"
)

// Cache interface defines cache operations
type Cache interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) error
	Delete(key string) error
	Exists(key string) (bool, error)
	Flush() error
	Close() error
}

// RedisCache implements Cache interface using Redis
type RedisCache struct {
	client *redis.Client
	config config.RedisConfig
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(cfg config.RedisConfig) (Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{
		client: rdb,
		config: cfg,
	}, nil
}

// Get retrieves a value from cache
func (c *RedisCache) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.client.Get(ctx, key).Result()
}

// Set stores a value in cache with expiration
func (c *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.client.Set(ctx, key, value, expiration).Err()
}

// Delete removes a value from cache
func (c *RedisCache) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in cache
func (c *RedisCache) Exists(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := c.client.Exists(ctx, key).Result()
	return result > 0, err
}

// Flush clears all cache entries
func (c *RedisCache) Flush() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.FlushDB(ctx).Err()
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// MemoryCache implements Cache interface using in-memory storage
type MemoryCache struct {
	data map[string]*cacheItem
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewMemoryCache creates a new in-memory cache instance
func NewMemoryCache() Cache {
	cache := &MemoryCache{
		data: make(map[string]*cacheItem),
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves a value from memory cache
func (c *MemoryCache) Get(key string) (string, error) {
	item, exists := c.data[key]
	if !exists {
		return "", nil
	}

	if time.Now().After(item.expiration) {
		delete(c.data, key)
		return "", nil
	}

	// Convert to JSON string
	if str, ok := item.value.(string); ok {
		return str, nil
	}

	jsonBytes, err := json.Marshal(item.value)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// Set stores a value in memory cache with expiration
func (c *MemoryCache) Set(key string, value interface{}, expiration time.Duration) error {
	c.data[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(expiration),
	}
	return nil
}

// Delete removes a value from memory cache
func (c *MemoryCache) Delete(key string) error {
	delete(c.data, key)
	return nil
}

// Exists checks if a key exists in memory cache
func (c *MemoryCache) Exists(key string) (bool, error) {
	item, exists := c.data[key]
	if !exists {
		return false, nil
	}

	if time.Now().After(item.expiration) {
		delete(c.data, key)
		return false, nil
	}

	return true, nil
}

// Flush clears all memory cache entries
func (c *MemoryCache) Flush() error {
	c.data = make(map[string]*cacheItem)
	return nil
}

// Close is a no-op for memory cache
func (c *MemoryCache) Close() error {
	return nil
}

// cleanup removes expired items from memory cache
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		for key, item := range c.data {
			if now.After(item.expiration) {
				delete(c.data, key)
			}
		}
	}
}

// NoOpCache implements Cache interface with no operations
type NoOpCache struct{}

// NewNoOpCache creates a new no-op cache instance
func NewNoOpCache() Cache {
	return &NoOpCache{}
}

// Get is a no-op operation
func (c *NoOpCache) Get(key string) (string, error) {
	return "", nil
}

// Set is a no-op operation
func (c *NoOpCache) Set(key string, value interface{}, expiration time.Duration) error {
	return nil
}

// Delete is a no-op operation
func (c *NoOpCache) Delete(key string) error {
	return nil
}

// Exists is a no-op operation
func (c *NoOpCache) Exists(key string) (bool, error) {
	return false, nil
}

// Flush is a no-op operation
func (c *NoOpCache) Flush() error {
	return nil
}

// Close is a no-op operation
func (c *NoOpCache) Close() error {
	return nil
}
