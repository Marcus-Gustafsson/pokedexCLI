package internal

import (
	"sync"
	"time"
)

// Cache is a threadsafe structure for storing key-value pairs temporarily.
// It expires entries after a given interval.
type Cache struct {
	cache map[string]cacheEntry
	muPtr *sync.RWMutex // RWMutex allows multiple readers, but only one writer.
}

// cacheEntry represents an individual cached value with its creation timestamp.
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// NewCache creates and returns a pointer to a new Cache instance.
// It also starts a background goroutine that periodically removes expired entries.
// The interval determines how often expired entries are reaped.
func NewCache(interval time.Duration) *Cache {
	c := Cache{
		cache: make(map[string]cacheEntry),
		muPtr: &sync.RWMutex{},
	}

	// Start the cleanup goroutine
	go c.reapLoop(interval)

	return &c
}

// Add stores a key-value pair in the cache.
// It records the current time as the creation time.
func (cPtr *Cache) Add(key string, val []byte) {

	cPtr.muPtr.Lock() // Acquire a write lock (exclusive)
	defer cPtr.muPtr.Unlock()

	cPtr.cache[key] = cacheEntry{
		createdAt: time.Now().UTC(),
		val:       val,
	}
}

// Get retrieves the value stored at a given key.
// It returns the value and true if found; otherwise nil and false.
func (cPtr *Cache) Get(key string) ([]byte, bool) {

	cPtr.muPtr.RLock() // Acquire a read lock (concurrent with other readers)
	defer cPtr.muPtr.RUnlock()

	value, ok := cPtr.cache[key]
	if ok {
		return value.val, true
	}
	return nil, false
}

// reapLoop runs forever in a background goroutine,
// purging expired entries from the cache after each interval.
func (cPtr *Cache) reapLoop(interval time.Duration) {
	for {
		time.Sleep(interval) // Wait for the interval to pass

		cPtr.muPtr.Lock() // Acquire a write lock for safe deletion

		now := time.Now()
		for key, entry := range cPtr.cache {
			// If the entry's age is greater than the interval, delete it
			if now.Sub(entry.createdAt) > interval {
				delete(cPtr.cache, key)
			}
		}
		cPtr.muPtr.Unlock()
	}
}
