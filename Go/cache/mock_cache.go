package cache

import "time"

type MockCache struct {
	store map[string]string
}

func NewMockCache() *MockCache {
	client := MockCache{
		store: make(map[string]string),
	}
	return &client
}

func (c *MockCache) Get(key string) string {
	return c.store[key]
}

func (c *MockCache) Set(key string, data string, expiryTime time.Duration) {
	c.store[key] = data
}
