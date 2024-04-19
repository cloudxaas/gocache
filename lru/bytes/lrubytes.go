//The MIT License (MIT)
//
// Copyright (c) 2024 CloudXaaS
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
    maxMemory       int64
    currentMemory   int64
    evictBatchSize  int    // Store the number of items to evict at once
    entries         []entry
    indexMap        map[string]int
    head, tail      int
    mu              sync.Mutex
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
    defer c.mu.Unlock()
    
    keyStr := cx.B2s(key)
    if idx, ok := c.indexMap[keyStr]; ok {
        if idx != c.head {
            c.moveToFront(idx)
        }
        return c.entries[idx].value, true
    }
    return nil, false
}

func (c *Cache) Put(key, value []byte) {
    c.mu.Lock()
    defer c.mu.Unlock()

    keyStr := cx.B2s(key)
    memSize := c.estimateMemory(key, value)

    if idx, ok := c.indexMap[keyStr]; ok {
        oldMemSize := c.estimateMemory(c.entries[idx].key, c.entries[idx].value)
        c.adjustMemory(memSize - oldMemSize)
        c.entries[idx].value = value
        c.moveToFront(idx)
        return
    }

    if c.currentMemory + memSize > c.maxMemory {
        for c.currentMemory + memSize > c.maxMemory {
            c.evict()
        }
    }

    c.entries = append(c.entries, entry{key: key, value: value})
    idx := len(c.entries) - 1
    c.indexMap[keyStr] = idx
    c.adjustMemory(memSize)
    c.moveToFront(idx)
}

func (c *Cache) Delete(key []byte) {
    c.mu.Lock()
    defer c.mu.Unlock()

    keyStr := cx.B2s(key)
    if idx, ok := c.indexMap[keyStr]; ok {
        memSize := c.estimateMemory(c.entries[idx].key, c.entries[idx].value)
        c.adjustMemory(-memSize)
        c.detach(idx)
        delete(c.indexMap, keyStr)
    }
}

func (c *Cache) evict() {
    for i := 0; i < c.evictBatchSize && c.tail != -1; i++ {
        oldKeyStr := cx.B2s(c.entries[c.tail].key)
        memSize := c.estimateMemory(c.entries[c.tail].key, c.entries[c.tail].value)
        c.adjustMemory(-memSize)
        c.detach(c.tail)
        delete(c.indexMap, oldKeyStr)
        if c.tail == -1 { // Check if the tail is -1 after detaching to safely exit the loop
            break
        }
    }
}


func (c *Cache) detach(idx int) {
    if c.entries[idx].prev != -1 {
        c.entries[c.entries[idx].prev].next = c.entries[idx].next
    }
    if c.entries[idx].next != -1 {
        c.entries[c.entries[idx].next].prev = c.entries[idx].prev
    }
    if idx == c.head {
        c.head = c.entries[idx].next
    }
    if idx == c.tail {
        c.tail = c.entries[idx].prev
    }
}

func (c *Cache) moveToFront(idx int) {
    if idx == c.head {
        return
    }
    c.detach(idx)
    c.entries[idx].prev = -1
    c.entries[idx].next = c.head
    if c.head != -1 {
        c.entries[c.head].prev = idx
    }
    c.head = idx
}

// CurrentMemory returns the current memory usage of the cache.
func (c *Cache) CurrentMemory() int64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.currentMemory
}
