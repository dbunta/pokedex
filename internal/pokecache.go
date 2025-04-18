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
	cache.entries[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (cache Cache) Get(key string) ([]byte, bool) {
	if val, ok := cache.entries[key]; ok {
		return val.val, true
	}
	return nil, false
}

func (cache Cache) reapLoop(duration time.Duration) {

	for key := range cache.entries {
		delete(cache.entries, key)
	}
}

func NewCache(duration time.Duration) Cache {
	cache := Cache{}
	return cache
}
