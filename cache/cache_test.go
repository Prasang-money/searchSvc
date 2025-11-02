package cache

import (
	"sync"
	"testing"

	"github.com/Prasang-money/searchSvc/models"
)

// Test initialization of cache
func TestNewCache(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
	}{
		{"zero capacity", 0},
		{"positive capacity", 5},
		{"large capacity", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewCache(tt.capacity)
			if cache == nil {
				t.Error("NewCache returned nil")
			}
			if cache.cap != tt.capacity {
				t.Errorf("Expected capacity %d, got %d", tt.capacity, cache.cap)
			}
			if cache.size != 0 {
				t.Errorf("Expected initial size 0, got %d", cache.size)
			}
			if cache.data == nil {
				t.Error("Cache map not initialized")
			}
			if cache.dll == nil {
				t.Error("Doubly linked list not initialized")
			}
		})
	}
}

// Test basic Set and Get operations
func TestSetAndGet(t *testing.T) {
	cache := NewCache(3)
	testData := &models.CountryMetadata{
		Name:       "United States",
		Population: 331002651,
		Capital:    "Washington, D.C.",
		Currency:   "USD",
	}

	// Test setting a value
	cache.Set("US", testData)
	if cache.size != 1 {
		t.Errorf("Expected size 1, got %d", cache.size)
	}

	// Test getting the value
	result, exists := cache.Get("US")
	if !exists {
		t.Error("Get returned false for existing key")
	}
	if result.Name != testData.Name {
		t.Errorf("Expected name %s, got %s", testData.Name, result.Name)
	}
	if result.Population != testData.Population {
		t.Errorf("Expected population %d, got %d", testData.Population, result.Population)
	}

	// Test getting non-existent key
	_, exists = cache.Get("XX")
	if exists {
		t.Error("Get returned true for non-existent key")
	}
}

// Test cache eviction
func TestEviction(t *testing.T) {
	cache := NewCache(2)

	// Add first item
	cache.Set("1", &models.CountryMetadata{Name: "First"})
	if cache.size != 1 {
		t.Errorf("Expected size 1, got %d", cache.size)
	}

	// Add second item
	cache.Set("2", &models.CountryMetadata{Name: "Second"})
	if cache.size != 2 {
		t.Errorf("Expected size 2, got %d", cache.size)
	}

	// Add third item, should evict first
	cache.Set("3", &models.CountryMetadata{Name: "Third"})
	if cache.size != 2 {
		t.Errorf("Expected size 2, got %d", cache.size)
	}

	// First item should be evicted
	_, exists := cache.Get("1")
	if exists {
		t.Error("First item should have been evicted")
	}

	// Second and third items should exist
	_, exists = cache.Get("2")
	if !exists {
		t.Error("Second item should still exist")
	}
	_, exists = cache.Get("3")
	if !exists {
		t.Error("Third item should exist")
	}
}

// Test LRU behavior
func TestLRU(t *testing.T) {
	cache := NewCache(2)

	// Add two items
	cache.Set("1", &models.CountryMetadata{Name: "First"})
	cache.Set("2", &models.CountryMetadata{Name: "Second"})

	// Access first item to make it most recently used
	cache.Get("1")

	// Add third item, should evict second item (least recently used)
	cache.Set("3", &models.CountryMetadata{Name: "Third"})

	// Check second item was evicted
	_, exists := cache.Get("2")
	if exists {
		t.Error("Second item should have been evicted")
	}

	// First and third items should exist
	_, exists = cache.Get("1")
	if !exists {
		t.Error("First item should still exist")
	}
	_, exists = cache.Get("3")
	if !exists {
		t.Error("Third item should exist")
	}
}

// Test updating existing items
func TestUpdate(t *testing.T) {
	cache := NewCache(2)

	original := &models.CountryMetadata{
		Name:     "Original",
		Currency: "USD",
	}
	updated := &models.CountryMetadata{
		Name:     "Updated",
		Currency: "EUR",
	}

	// Set original value
	cache.Set("key", original)

	// Update value
	cache.Set("key", updated)

	// Size should still be 1
	if cache.size != 1 {
		t.Errorf("Expected size 1 after update, got %d", cache.size)
	}

	// Check updated value
	result, exists := cache.Get("key")
	if !exists {
		t.Error("Key should exist after update")
	}
	if result.Name != updated.Name {
		t.Errorf("Expected name %s, got %s", updated.Name, result.Name)
	}
	if result.Currency != updated.Currency {
		t.Errorf("Expected currency %s, got %s", updated.Currency, result.Currency)
	}
}

// Test concurrent access
func TestConcurrent(t *testing.T) {
	cache := NewCache(100)
	var wg sync.WaitGroup

	// Test concurrent writes
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := string(rune('A' + i%26))
			cache.Set(key, &models.CountryMetadata{
				Name:       "Test Country",
				Population: i,
				Capital:    "Test Capital",
				Currency:   "TST",
			})
		}(i)
	}

	// Test concurrent reads and writes
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := string(rune('A' + i%26))
			if i%2 == 0 {
				cache.Set(key, &models.CountryMetadata{
					Name:       "Updated Country",
					Population: i,
					Capital:    "Updated Capital",
					Currency:   "UPD",
				})
			} else {
				cache.Get(key)
			}
		}(i)
	}

	wg.Wait()

	// Verify cache state is consistent
	if cache.size > cache.cap {
		t.Errorf("Cache size %d exceeded capacity %d", cache.size, cache.cap)
	}
}

// Test DoublyLinkedList operations
func TestDoublyLinkedList(t *testing.T) {
	dll := &DoublyLinkedList{}

	// Test empty list
	if dll.head != nil || dll.tail != nil {
		t.Error("New list should be empty")
	}

	// Test adding first node
	node1 := NewNode("1", &models.CountryMetadata{Name: "First"})
	dll.addToFront(node1)
	if dll.head != node1 || dll.tail != node1 {
		t.Error("First node should be both head and tail")
	}

	// Test adding second node
	node2 := NewNode("2", &models.CountryMetadata{Name: "Second"})
	dll.addToFront(node2)
	if dll.head != node2 || dll.tail != node1 {
		t.Error("Second node should be head, first node should be tail")
	}

	// Test removing middle node
	node3 := NewNode("3", &models.CountryMetadata{Name: "Third"})
	dll.addToFront(node3)
	dll.remove(node2) // Remove middle node
	if dll.head != node3 || dll.tail != node1 {
		t.Error("After removing middle node, third should be head and first should be tail")
	}
	if node3.next != node1 || node1.prev != node3 {
		t.Error("Links not properly updated after removing middle node")
	}
}

// Test zero capacity cache
func TestZeroCapacity(t *testing.T) {
	cache := NewCache(0)

	// Try to add an item
	cache.Set("key", &models.CountryMetadata{Name: "Test"})

	// Verify nothing was stored
	if cache.size != 0 {
		t.Error("Zero capacity cache should not store items")
	}

	_, exists := cache.Get("key")
	if exists {
		t.Error("Zero capacity cache should not return items")
	}
}
