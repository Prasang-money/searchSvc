package cache

import (
	"sync"

	"github.com/Prasang-money/searchSvc/models"
)

// in this cache I am not handling the collision for simplicity
// Using LRU for eviction policy
type Cache struct {
	data  map[string]*Node
	dll   *DoublyLinkedList
	size  int
	cap   int
	mutex sync.RWMutex // Mutex for thread-safe operations
}

func NewCache(capacity int) *Cache {
	return &Cache{
		data:  make(map[string]*Node),
		cap:   capacity,
		dll:   &DoublyLinkedList{},
		mutex: sync.RWMutex{},
	}
}

func (cache *Cache) Set(key string, value *models.CountryMetadata) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if node, exists := cache.data[key]; exists {
		node.value = value
		// move existing node to front (no size change)
		if cache.dll.head != node {
			cache.dll.remove(node)
			cache.dll.addToFront(node)
		}
		return
	}
	cache.data[key] = NewNode(key, value)
	cache.dll.addToFront(cache.data[key])
	cache.size++
	if cache.size > cache.cap {
		cache.removeOldest()
	}
}

func (cache *Cache) removeOldest() {
	// No need for lock here as this is only called from Set which already holds the lock
	if cache.dll.tail != nil {
		delete(cache.data, cache.dll.tail.key)
		cache.dll.remove(cache.dll.tail)
		cache.size--
	}
}

func (cache *Cache) Get(key string) (models.CountryMetadata, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if cache.data[key] == nil {
		return models.CountryMetadata{}, false
	}
	//fmt.Println("Cache hit for key:", key)
	cache.dll.remove(cache.data[key])
	cache.dll.addToFront(cache.data[key])
	return *cache.data[key].value, true
}

type DoublyLinkedList struct {
	head *Node
	tail *Node
}

func (dll *DoublyLinkedList) addToFront(node *Node) {
	node.prev = nil
	node.next = dll.head
	if dll.head != nil {
		dll.head.prev = node
	}
	dll.head = node
	if dll.tail == nil {
		dll.tail = node
	}
}

func (dll *DoublyLinkedList) remove(node *Node) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		dll.head = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		dll.tail = node.prev
	}
}

type Node struct {
	key   string
	value *models.CountryMetadata
	prev  *Node
	next  *Node
}

func NewNode(key string, value *models.CountryMetadata) *Node {
	return &Node{
		key:   key,
		value: value,
		prev:  nil,
		next:  nil,
	}
}
