package cache

import (
	"context"
	"time"
)

// Cache interface defines cache operations
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Clear(ctx context.Context) error
}

// MockCache is a simple in-memory cache implementation for testing
type MockCache struct {
	data map[string]string
}

// NewMockCache creates a new mock cache
func NewMockCache() *MockCache {
	return &MockCache{
		data: make(map[string]string),
	}
}

// Get retrieves a value from cache
func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	return m.data[key], nil
}

// Set stores a value in cache
func (m *MockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.data[key] = value.(string)
	return nil
}

// Delete removes a value from cache
func (m *MockCache) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

// Exists checks if a key exists in cache
func (m *MockCache) Exists(ctx context.Context, key string) (bool, error) {
	_, exists := m.data[key]
	return exists, nil
}

// Clear clears all cache entries
func (m *MockCache) Clear(ctx context.Context) error {
	m.data = make(map[string]string)
	return nil
}
