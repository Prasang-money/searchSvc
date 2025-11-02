package cache

import (
	"fmt"

	"github.com/Prasang-money/searchSvc/models"
)

// in this cache I am not handling the collision for simplicity
// Using LRU for eviction policy
type Cache struct {
	data map[string]*models.CountryMetadata
	dll  *Node
	size int
	cap  int
}

func NewCache(capacity int) *Cache {
	return &Cache{
		data: make(map[string]*models.CountryMetadata),
		cap:  capacity,
		dll:  &Node{},
	}
}

func (cache *Cache) Set(key string, value models.CountryMetadata) {
	cache.data[key] = &value

}

func (cache *Cache) Get(key string) (models.CountryMetadata, error) {
	if cache.data[key] == nil {
		return models.CountryMetadata{}, fmt.Errorf("key not found")
	}
	return *cache.data[key], nil
}

func remove() {

}

type Node struct {
	key   string
	value models.CountryMetadata
	prev  *Node
	next  *Node
}
