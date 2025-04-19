package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {

	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("testdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestAddOverwrite(t *testing.T) {
	const baseTime = 5 * time.Minute
	cache := NewCache(baseTime)

	key := "https://example.com"
	cache.Add(key, []byte("testdata"))
	cache.Add(key, []byte("testdata2"))

	val, ok := cache.Get(key)
	if !ok {
		t.Errorf("expected to find key")
		return
	}
	if string(val) != "testdata2" {
		t.Errorf("expected to different value in cache")
		return
	}

}

func TestReapLoop(t *testing.T) {
	const baseTime = 2 * time.Second
	const waitTime = baseTime + 3*time.Second
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("http://example.com")
	if ok {
		t.Errorf("expected to find key")
		return
	}
}
