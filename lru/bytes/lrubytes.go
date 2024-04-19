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
    "unsafe"
)

type Cache struct {
    maxMemory       int64
    currentMemory   int64
    evictBatchSize  int
    entries         []entry
    freeEntries     []int // stack of indices of free entries
    indexMap        map[uintptr]int // map of byte slice pointers to entry indices
    head, tail      int
    mu              sync.Mutex
}

type entry struct {
    key, value []byte
    prev, next int
}

// NewLRUCache initializes a new LRU Cache with given maximum memory and eviction batch size.
func NewLRUCache(maxMemory int64, evictBatchSize int) *Cache {
    return &Cache{
        maxMemory:      maxMemory,
        evictBatchSize: evictBatchSize,
        entries:        make([]entry, 0),
        indexMap:       make(map[uintptr]int),
        head:           -1,
        tail:           -1,
        freeEntries:    make([]int, 0),
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

    ptr := uintptr(unsafe.Pointer(&key[0]))
    if idx, ok := c.indexMap[ptr]; ok {
        entry := &c.entries[idx]
        if idx != c.head {
            c.moveToFront(idx)
        }
        c.mu.Unlock()
        return entry.value, true
    }
    c.mu.Unlock()
    return nil, false
}

func (c *Cache) Put(key, value []byte) {
    c.mu.Lock()

    ptr := uintptr(unsafe.Pointer(&key[0]))
    memSize := c.estimateMemory(key, value)

    if idx, ok := c.indexMap[ptr]; ok {
        oldMemSize := c.estimateMemory(c.entries[idx].key, c.entries[idx].value)
        c.adjustMemory(memSize - oldMemSize)
        c.entries[idx].value = value
        c.moveToFront(idx)
        c.mu.Unlock()
        return
    }

    if c.currentMemory+memSize > c.maxMemory {
        c.evict()
    }

    var idx int
    if len(c.freeEntries) > 0 {
        idx = c.freeEntries[len(c.freeEntries)-1]
        c.freeEntries = c.freeEntries[:len(c.freeEntries)-1]
        c.entries[idx] = entry{key: key, value: value}
    } else {
        c.entries = append(c.entries, entry{key: key, value: value})
        idx = len(c.entries) - 1
    }

    c.indexMap[ptr] = idx
    c.adjustMemory(memSize)
    c.moveToFront(idx)
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

func (c *Cache) evict() {
    for i := 0; i < c.evictBatchSize && c.tail != -1; i++ {
        idx := c.tail
        memSize := c.estimateMemory(c.entries[idx].key, c.entries[idx].value)
        c.adjustMemory(-memSize)
        c.detach(idx)
        c.freeEntries = append(c.freeEntries, idx)
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

    if c.head == -1 {
        c.tail = -1
    }
}

func (c *Cache) Delete(key []byte) {
    c.mu.Lock()

    ptr := uintptr(unsafe.Pointer(&key[0]))
    if idx, ok := c.indexMap[ptr]; ok {
        memSize := c.estimateMemory(c.entries[idx].key, c.entries[idx].value)
        c.adjustMemory(-memSize)
        c.detach(idx)
        c.freeEntries = append(c.freeEntries, idx)
        delete(c.indexMap, ptr)
    }
    c.mu.Unlock()
}
