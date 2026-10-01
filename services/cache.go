package services

import (
	"event-explorer/models"
	"fmt"
	"sync"
	"time"
)

type CacheItem struct {
	Events    []models.UIEvent
	ExpiresAt time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]CacheItem
}

var (
	instance *MemoryCache
	once     sync.Once
)

func GetCacheInstance() *MemoryCache {
	once.Do(func() {
		instance = &MemoryCache{
			items: make(map[string]CacheItem),
		}
		go instance.startJanitor(5 * time.Minute)
	})
	return instance
}

func (c *MemoryCache) Set(city, countryCode, category string, events []models.UIEvent) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := fmt.Sprintf("%s:%s:%s", city, countryCode, category)
	c.items[key] = CacheItem{
		Events:    events,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
}

func (c *MemoryCache) Get(city, countryCode, category string) ([]models.UIEvent, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := fmt.Sprintf("%s:%s:%s", city, countryCode, category)
	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	fmt.Printf("INFO cache hit key=%s\n", key)
	return item.Events, true
}

func (c *MemoryCache) startJanitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		deletedCount := 0

		for key, item := range c.items {
			if now.After(item.ExpiresAt) {
				delete(c.items, key)
				deletedCount++
			}
		}

		if deletedCount > 0 {
			fmt.Printf("INFO cache janitor cleared %d expired items from memory\n", deletedCount)
		}
		c.mu.Unlock()
	}
}

func (c *MemoryCache) DeleteData(city, countryCode, category string) (int, string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := fmt.Sprintf("%s:%s:%s", city, countryCode, category)

	if _, exists := c.items[key]; exists {
		delete(c.items, key)

		fmt.Printf(
			"INFO cache manual check evicted specific key=%s\n",
			key,
		)

		return 1, fmt.Sprintf(
			"Specific cache item bound to key [%s] evicted successfully",
			key,
		)
	}

	return 0, fmt.Sprintf(
		"No active cache record found matching key [%s]",
		key,
	)
}

func (c *MemoryCache) ClearData() (int, string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := len(c.items)
	c.items = make(map[string]CacheItem)

	if count > 0 {
		fmt.Printf(
			"INFO cache manual flush purged total %d elements\n",
			count,
		)
	}

	return count, "Entire system memory cache cleared successfully"
}
