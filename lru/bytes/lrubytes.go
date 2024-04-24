// The MIT License (MIT)
//
// # Copyright (c) 2024 CloudXaaS
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package lrubytes

import (
	"sync"

	cx "github.com/cloudxaas/gocx"
)

type Cache struct {
	maxMemory      int64
	currentMemory  int64
	evictBatchSize int // Store the number of items to evict at once
	entries        []entry
	indexMap       map[string]int
	head, tail     int
	mu             sync.Mutex
}

type entry struct {
	key, value []byte
	prev, next int
}

// NewLRUCache now accepts an additional parameter for batch eviction size
func NewLRUCache(maxMemory int64, evictBatchSize int) *Cache {
	return &Cache{
		maxMemory:      maxMemory,
		evictBatchSize: evictBatchSize, // Initialize evictBatchSize
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

func (c *Cache) Get(key []byte) ([]byte, bool) {
	c.mu.Lock()
	keyStr := cx.B2s(key)
	if idx, ok := c.indexMap[keyStr]; ok {
		if idx != c.head {
			c.moveToFront(idx)
		}
		c.mu.Unlock()
		return c.entries[idx].value, true
	}
	c.mu.Unlock()
	return nil, false
}

func (c *Cache) Set(key, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	keyStr := cx.B2s(key)
	memSize := c.estimateMemory(key, value)

	// Evict until there's enough space for the new entry
	for c.currentMemory+memSize > c.maxMemory && c.tail != -1 {
		c.evict()
	}

	// If there's still not enough space after eviction, don't add the new entry
	if c.currentMemory+memSize > c.maxMemory {
		return
	}

	// Add the new entry
	if idx, ok := c.indexMap[keyStr]; ok {
		oldMemSize := c.estimateMemory(c.entries[idx].key, c.entries[idx].value)
		c.adjustMemory(memSize - oldMemSize)
		c.entries[idx].value = value
		c.moveToFront(idx)
	} else {
		c.entries = append(c.entries, entry{key: key, value: value, prev: -1, next: -1})
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
}

func (c *Cache) moveToFront(idx int) {
	if idx == c.head {
		return // It's already the head, nothing to do
	}
	c.detach(idx) // Detach from current position

	// Set the previous head
	if c.head != -1 {
		c.entries[c.head].prev = idx
	}
	c.entries[idx].next = c.head
	c.entries[idx].prev = -1
	c.head = idx

	// If there was no head before, this is also the tail
	if c.tail == -1 {
		c.tail = idx
	}

	// If moving the last node to the front, update tail if needed
	if c.tail == idx {
		c.tail = c.entries[idx].prev // Update the tail if the moved node was the tail
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
	// Handle previous link
	if c.entries[idx].prev != -1 {
		c.entries[c.entries[idx].prev].next = c.entries[idx].next
	} else {
		// When removing the head, move the head pointer forward
		c.head = c.entries[idx].next
	}

	// Handle next link
	if c.entries[idx].next != -1 {
		c.entries[c.entries[idx].next].prev = c.entries[idx].prev
	} else {
		// When removing the tail, move the tail pointer backward
		c.tail = c.entries[idx].prev
	}

	// Reset the node's links
	c.entries[idx].prev = -1
	c.entries[idx].next = -1

	// Additional check: if the cache is now empty, reset head and tail
	if c.head == -1 {
		c.tail = -1 // Ensures that the tail is also reset when all items are evicted
	}
}

func (c *Cache) evict() {
	for i := 0; i < c.evictBatchSize && c.tail != -1; i++ {
		oldKeyStr := cx.B2s(c.entries[c.tail].key)
		memSize := c.estimateMemory(c.entries[c.tail].key, c.entries[c.tail].value)
		c.adjustMemory(-memSize)
		c.detach(c.tail)

		if c.tail != -1 { // Verify tail is valid before proceeding
			delete(c.indexMap, oldKeyStr)
		}

		if c.tail == -1 { // Break if no more items to evict
			break
		}
	}
}
