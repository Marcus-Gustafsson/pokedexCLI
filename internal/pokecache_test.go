package internal

import (
	"testing"
	"time"
	"fmt"
)

// TestAddGet checks that values added to the cache
// can be reliably retrieved by their key within the cache's lifetime.
func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second // Each cache entry is valid for 5 seconds

	// Test cases: different keys and values to try adding/getting
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://exampleURL.com",
			val: []byte("testdata"),
		},
		{
			key: "https://exampleURL.com/path",
			val: []byte("moretestdata"),
		},
	}

	// Loop over the test cases to verify cache behavior for each
	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)             // Create a new cache with cleanup interval
			cache.Add(c.key, c.val)                 // Add key-value pair to the cache

			val, ok := cache.Get(c.key)             // Try to get the value back
			if !ok {
				t.Errorf("expected to find key")    // Fail if value is missing
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")  // Fail if value is incorrect
				return
			}
		})
	}
}

// TestReapLoop checks if entries expire as expected
// after the configured interval for the cache.
func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond                 // Entry should expire in 5ms
	const waitTime = baseTime + 5*time.Millisecond        // Wait longer than base time to ensure expiry
	cache := NewCache(baseTime)                           // Create a new cache with a short interval

	cache.Add("https://exampleURL.com", []byte("testdata"))  // Add an entry

	// Immediately, the entry should exist
	_, ok := cache.Get("https://exampleURL.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)                                  // Wait for expiration

	// After waiting, the entry should be gone (expired)
	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}