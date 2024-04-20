package lruxbytes

import (
	"sync"

	"github.com/zeebo/xxh3"
)

type Cache struct {
	maxMemory      int64
	currentMemory  int64
	evictBatchSize int
	entries        []entry
	indexMap       map[uint32]uint32 // Change to uint32 for index
	head, tail     int
	mu             sync.Mutex
}

type entry struct {
	key, value []byte
	prev, next int
}

func NewLRUCache(maxMemory int64, evictBatchSize int) *Cache {
	return &Cache{
		maxMemory:      maxMemory,
		evictBatchSize: evictBatchSize,
		entries:        make([]entry, 0),
		indexMap:       make(map[uint32]uint32), // Adjusted to uint32
		head:           -1,
		tail:           -1,
	}
}

func (c *Cache) estimateMemory(key, value []byte) int64 {
	return int64(len(key) + len(value))
}

func (c *Cache) adjustMemory(delta int64) {
	c.currentMemory += delta
}

func (c *Cache) hashKey(key []byte) uint32 {
	return uint32(xxh3.Hash(key))
}

func (c *Cache) Get(key []byte) ([]byte, bool) {
	c.mu.Lock()
	keyHash := c.hashKey(key)
	if idx, ok := c.indexMap[keyHash]; ok {
		if int(idx) != c.head { // Cast idx to int for comparison
			c.moveToFront(int(idx))
		}
		c.mu.Unlock()
		return c.entries[idx].value, true
	}
	c.mu.Unlock()
	return nil, false
}

func (c *Cache) Set(key, value []byte) {
	c.mu.Lock()
	keyHash := c.hashKey(key)
	memSize := c.estimateMemory(key, value)

	if idx, ok := c.indexMap[keyHash]; ok {
		oldMemSize := c.estimateMemory(c.entries[int(idx)].key, c.entries[int(idx)].value)
		c.adjustMemory(memSize - oldMemSize)
		c.entries[idx].value = value
		c.moveToFront(int(idx))
		c.mu.Unlock()
		return
	}

	if c.currentMemory+memSize > c.maxMemory {
		for c.currentMemory+memSize > c.maxMemory && c.tail != -1 {
			c.evict()
		}
	}

	if c.currentMemory+memSize <= c.maxMemory {
		c.entries = append(c.entries, entry{key: key, value: value, prev: -1, next: -1})
		idx := uint32(len(c.entries) - 1) // Cast length to uint32
		c.indexMap[keyHash] = idx
		c.adjustMemory(memSize)
		c.moveToFront(int(idx)) // Cast idx to int for handling

		if c.head == int(idx) {
			if c.tail == -1 {
				c.tail = int(idx)
			}
		}
	}
	c.mu.Unlock()
}

func (c *Cache) Del(key []byte) {
	c.mu.Lock()
	keyHash := c.hashKey(key)
	if idx, ok := c.indexMap[keyHash]; ok {
		memSize := c.estimateMemory(c.entries[int(idx)].key, c.entries[int(idx)].value)
		c.adjustMemory(-memSize)
		c.detach(int(idx))
		delete(c.indexMap, keyHash)
	}
	c.mu.Unlock()
}

func (c *Cache) moveToFront(idx int) {
	if idx == c.head {
		return
	}
	c.detach(idx)

	if c.head != -1 {
		c.entries[c.head].prev = idx
	}
	c.entries[idx].next = c.head
	c.entries[idx].prev = -1
	c.head = idx

	if c.tail == -1 {
		c.tail = idx
	}

	if c.tail == idx {
		c.tail = c.entries[idx].prev
	}
}

func (c *Cache) detach(idx int) {
	if c.entries[idx].prev != -1 {
		c.entries[c.entries[idx].prev].next = c.entries[idx].next
	} else {
		c.head = c.entries[idx].next
	}

	if c.entries[idx].next != -1 {
		c.entries[c.entries[idx].next].prev = c.entries[idx].prev
	} else {
		c.tail = c.entries[idx].prev
	}

	c.entries[idx].prev = -1
	c.entries[idx].next = -1

	if c.head == -1 {
		c.tail = -1
	}
}

func (c *Cache) evict() {
	for i := 0; i < c.evictBatchSize && c.tail != -1; i++ {
		oldKeyHash := c.hashKey(c.entries[c.tail].key)
		memSize := c.estimateMemory(c.entries[c.tail].key, c.entries[c.tail].value)
		c.adjustMemory(-memSize)
		c.detach(c.tail)

		if c.tail != -1 {
			delete(c.indexMap, oldKeyHash)
		}

		if c.tail == -1 {
			break
		}
	}
}

