package lrubyteswcounter

import (
	"sync"

	cx "github.com/cloudxaas/gocx"
)

type Cache struct {
	maxMemory      int64
	currentMemory  int64
	evictBatchSize int
	entries        []entry
	indexMap       map[string]int
	head, tail     int
	mu             sync.Mutex
}

type entry struct {
	key, value []byte
	counter    uint8 // Add a counter to the entry struct
	prev, next int
}

func NewLRUCache(maxMemory int64, evictBatchSize int) *Cache {
	return &Cache{
		maxMemory:      maxMemory,
		evictBatchSize: evictBatchSize,
		entries:        make([]entry, 0),
		indexMap:       make(map[string]int),
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

func (c *Cache) Get(key []byte) ([]byte, uint16, bool) {
	c.mu.Lock()

	keyStr := cx.B2s(key)
	if idx, ok := c.indexMap[keyStr]; ok {
		if idx != c.head {
			c.moveToFront(idx)
		}
		// Increment the counter each time a Get is done
		c.entries[idx].counter++
		c.mu.Unlock()
		return c.entries[idx].value, c.entries[idx].counter, true
	}
	c.mu.Unlock()
	return nil, 0, false
}

func (c *Cache) Set(key, value []byte) {
	c.mu.Lock()
	keyStr := cx.B2s(key)
	memSize := c.estimateMemory(key, value)

	for c.currentMemory+memSize > c.maxMemory && c.tail != -1 {
		c.evict()
	}

	if c.currentMemory+memSize > c.maxMemory {
		c.mu.Unlock()
		return
	}

	if idx, ok := c.indexMap[keyStr]; ok {
		oldMemSize := c.estimateMemory(c.entries[idx].key, c.entries[idx].value)
		c.adjustMemory(memSize - oldMemSize)
		c.entries[idx].value = value
		c.moveToFront(idx)
	} else {
		c.entries = append(c.entries, entry{key: key, value: value, prev: -1, next: -1, counter: 0})
		idx := len(c.entries) - 1
		c.indexMap[keyStr] = idx
		c.adjustMemory(memSize)
		c.moveToFront(idx)

		if c.head == idx {
			if c.tail == -1 {
				c.tail = idx
			}
		}
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

func (c *Cache) Del(key []byte) {
	c.mu.Lock()

	keyStr := cx.B2s(key)
	if idx, ok := c.indexMap[keyStr]; ok {
		memSize := c.estimateMemory(c.entries[idx].key, c.entries[idx].value)
		c.adjustMemory(-memSize)
		c.detach(idx)
		delete(c.indexMap, keyStr)
	}	
	c.mu.Unlock()
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
		oldKeyStr := cx.B2s(c.entries[c.tail].key)
		memSize := c.estimateMemory(c.entries[c.tail].key, c.entries[c.tail].value)
		c.adjustMemory(-memSize)
		c.detach(c.tail)

		if c.tail != -1 {
			delete(c.indexMap, oldKeyStr)
		}

		if c.tail == -1 {
			break
		}
	}
}
