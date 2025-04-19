package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mu      sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (cache Cache) Add(key string, val []byte) {
	cache.mu.Lock()
	cache.entries[key] = cacheEntry{createdAt: time.Now(), val: val}
	cache.mu.Unlock()
}

func (cache Cache) Get(key string) ([]byte, bool) {
	cache.mu.Lock()
	if val, ok := cache.entries[key]; ok {
		return val.val, true
	}
	cache.mu.Unlock()
	return nil, false
}

func (cache Cache) reapLoop(duration time.Duration) {
	ticker := time.NewTicker(duration)
	time := time.Now().Add(-duration)
	for range ticker.C {
		cache.mu.Lock()
		for key := range cache.entries {
			if cache.entries[key].createdAt.Before(time) {
				delete(cache.entries, key)
			}
		}
		cache.mu.Unlock()
	}
}

func NewCache(duration time.Duration) Cache {
	cache := Cache{}
	cache.entries = make(map[string]cacheEntry)
	go cache.reapLoop(duration)
	return cache
}
