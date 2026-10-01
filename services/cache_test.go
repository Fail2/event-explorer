package services

import (
	"event-explorer/models"
	"testing"
	"time"
)

func TestGetCacheInstance(t *testing.T) {
	cache := GetCacheInstance()

	if cache == nil {
		t.Fatal("GetCacheInstance() returned nil")
	}

	if cache.items == nil {
		t.Fatal("cache items map is nil")
	}

	// Calling again should return the same singleton instance.
	cache2 := GetCacheInstance()

	if cache != cache2 {
		t.Fatal("GetCacheInstance() did not return the same instance")
	}
}

func TestMemoryCache_SetAndGet(t *testing.T) {
	cache := &MemoryCache{
		items: make(map[string]CacheItem),
	}

	events := []models.UIEvent{
		{
			ID:   "event-1",
			Name: "Test Event",
		},
	}

	cache.Set("London", "GB", "Music", events)

	got, ok := cache.Get("London", "GB", "Music")

	if !ok {
		t.Fatal("expected cache hit, got cache miss")
	}

	if len(got) != 1 {
		t.Fatalf("got %d events, want 1", len(got))
	}

	if got[0].ID != "event-1" {
		t.Errorf("got event ID %q, want event-1", got[0].ID)
	}
}

func TestMemoryCache_GetMissing(t *testing.T) {
	cache := &MemoryCache{
		items: make(map[string]CacheItem),
	}

	got, ok := cache.Get("London", "GB", "Music")

	if ok {
		t.Fatal("expected cache miss, got cache hit")
	}

	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

func TestMemoryCache_GetExpired(t *testing.T) {
	cache := &MemoryCache{
		items: make(map[string]CacheItem),
	}

	key := "London:GB:Music"

	cache.items[key] = CacheItem{
		Events: []models.UIEvent{
			{
				ID:   "expired-event",
				Name: "Expired Event",
			},
		},
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}

	got, ok := cache.Get("London", "GB", "Music")

	if ok {
		t.Fatal("expected expired cache item to be a miss")
	}

	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

func TestMemoryCache_DeleteData(t *testing.T) {
	cache := &MemoryCache{
		items: make(map[string]CacheItem),
	}

	cache.Set(
		"London",
		"GB",
		"Music",
		[]models.UIEvent{
			{
				ID:   "event-1",
				Name: "Test Event",
			},
		},
	)

	count, message := cache.DeleteData("London", "GB", "Music")

	if count != 1 {
		t.Errorf("DeleteData() count = %d, want 1", count)
	}

	if message == "" {
		t.Error("DeleteData() returned empty success message")
	}

	_, ok := cache.Get("London", "GB", "Music")

	if ok {
		t.Error("expected deleted cache item to be unavailable")
	}
}

func TestMemoryCache_DeleteDataMissing(t *testing.T) {
	cache := &MemoryCache{
		items: make(map[string]CacheItem),
	}

	count, message := cache.DeleteData("London", "GB", "Music")

	if count != 0 {
		t.Errorf("DeleteData() count = %d, want 0", count)
	}

	if message == "" {
		t.Error("DeleteData() returned empty message")
	}
}

func TestMemoryCache_ClearData(t *testing.T) {
	cache := &MemoryCache{
		items: make(map[string]CacheItem),
	}

	cache.Set(
		"London",
		"GB",
		"Music",
		[]models.UIEvent{
			{ID: "event-1"},
		},
	)

	cache.Set(
		"London",
		"GB",
		"Sports",
		[]models.UIEvent{
			{ID: "event-2"},
		},
	)

	count, message := cache.ClearData()

	if count != 2 {
		t.Errorf("ClearData() count = %d, want 2", count)
	}

	if message != "Entire system memory cache cleared successfully" {
		t.Errorf("ClearData() message = %q, want expected message", message)
	}

	if len(cache.items) != 0 {
		t.Errorf("cache contains %d items after ClearData(), want 0", len(cache.items))
	}
}

func TestMemoryCache_ClearDataEmpty(t *testing.T) {
	cache := &MemoryCache{
		items: make(map[string]CacheItem),
	}

	count, message := cache.ClearData()

	if count != 0 {
		t.Errorf("ClearData() count = %d, want 0", count)
	}

	if message != "Entire system memory cache cleared successfully" {
		t.Errorf("ClearData() message = %q, want expected message", message)
	}
}

func TestMemoryCache_StartJanitor(t *testing.T) {
	cache := &MemoryCache{
		items: make(map[string]CacheItem),
	}

	cache.items["expired-key"] = CacheItem{
		Events: []models.UIEvent{
			{ID: "expired-event"},
		},
		ExpiresAt: time.Now().Add(-1 * time.Second),
	}

	cache.items["active-key"] = CacheItem{
		Events: []models.UIEvent{
			{ID: "active-event"},
		},
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	go cache.startJanitor(1 * time.Millisecond)

	time.Sleep(20 * time.Millisecond)

	cache.mu.RLock()
	defer cache.mu.RUnlock()

	if _, exists := cache.items["expired-key"]; exists {
		t.Error("janitor did not remove expired cache item")
	}

	if _, exists := cache.items["active-key"]; !exists {
		t.Error("janitor incorrectly removed active cache item")
	}
}
