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
    "sync/atomic"

    cx "github.com/cloudxaas/gocx"
)

type Cache struct {
    maxMemory      int64
    currentMemory  int64
    evictBatchSize int
    entries        map[uint64]entry
    indexMap       map[string]uint64
    head, tail     uint64
    mu             sync.RWMutex
    indexCounter   uint64
}

type entry struct {
    key, value []byte
    index      uint8
    prev, next uint64
}

const (
    InvalidIndex = ^uint64(0) // Max value for uint64 to represent an invalid index
)

func NewLRUCache(maxMemory int64, evictBatchSize int) *Cache {
    return &Cache{
        maxMemory:      maxMemory,
        evictBatchSize: evictBatchSize,
        entries:        make(map[uint64]entry),
        indexMap:       make(map[string]uint64),
        head:           InvalidIndex,
        tail:           InvalidIndex,
        indexCounter:   0,
    }
}

func (c *Cache) estimateMemory(key, value []byte) int64 {
    return int64(len(key) + len(value) + 10) // Adding constant overhead for index
}

func (c *Cache) adjustMemory(delta int64) {
    atomic.AddInt64(&c.currentMemory, delta)
}

func (c *Cache) Get(key []byte) ([]byte, bool) {
    keyStr := cx.B2s(key)
    
    c.mu.RLock()
    idx, ok := c.indexMap[keyStr]
    if !ok {
        c.mu.RUnlock()
        return nil, false
    }

    if idx != c.head {
        c.mu.RUnlock()
        c.mu.Lock()
        c.moveToFront(idx)
        c.mu.Unlock()
    } else {
        c.mu.RUnlock()
    }
    
    c.mu.RLock()
    entry := c.entries[idx]
    c.mu.RUnlock()
    
    return entry.value, true
}

func (c *Cache) moveToFront(idx uint64) {
    if idx == InvalidIndex || idx == c.head {
        return
    }

    entry := c.entries[idx]
    c.detach(idx)

    entry.prev = InvalidIndex
    entry.next = c.head
    if c.head != InvalidIndex {
        headEntry := c.entries[c.head]
        headEntry.prev = idx
        c.entries[c.head] = headEntry
    }
    c.head = idx
    c.entries[idx] = entry

    if c.tail == InvalidIndex {
        c.tail = idx
    }
}

func (c *Cache) detach(idx uint64) {
    if idx == InvalidIndex {
        return
    }

    entry := c.entries[idx]
    if entry.prev != InvalidIndex {
        prevEntry := c.entries[entry.prev]
        prevEntry.next = entry.next
        c.entries[entry.prev] = prevEntry
    } else {
        c.head = entry.next
    }

    if entry.next != InvalidIndex {
        nextEntry := c.entries[entry.next]
        nextEntry.prev = entry.prev
        c.entries[entry.next] = nextEntry
    } else {
        c.tail = entry.prev
    }

    entry.prev = InvalidIndex
    entry.next = InvalidIndex
    c.entries[idx] = entry
}

func (c *Cache) evict(entrySize int64) {
    evicted := false
    attempts := 0

    for atomic.LoadInt64(&c.currentMemory)+entrySize > c.maxMemory && c.tail != InvalidIndex {
        tailIdx := c.tail

        if tailIdx == InvalidIndex {
            break
        }

        oldKeyStr := cx.B2s(c.entries[tailIdx].key)
        memSize := c.estimateMemory(c.entries[tailIdx].key, c.entries[tailIdx].value)
        c.adjustMemory(-memSize)

        c.detach(tailIdx)

        delete(c.indexMap, oldKeyStr)
        delete(c.entries, tailIdx)
        evicted = true
        attempts++

        if attempts > len(c.entries) {
            break
        }
    }

    if !evicted {
        return
    }
}

func (c *Cache) wrapIndexCounter() {
    if c.indexCounter == InvalidIndex {
        c.indexCounter = 0
    }
}

func (c *Cache) Set(key, value []byte) error {
    keyStr := cx.B2s(key)
    memSize := c.estimateMemory(key, value)

    c.wrapIndexCounter()

    c.mu.Lock()
    defer c.mu.Unlock()

    for atomic.LoadInt64(&c.currentMemory)+memSize > c.maxMemory && c.tail != InvalidIndex {
        c.evict(memSize)
    }

    if atomic.LoadInt64(&c.currentMemory)+memSize > c.maxMemory {
        return nil
    }

    if idx, ok := c.indexMap[keyStr]; ok {
        entry := c.entries[idx]
        entry.value = value
        c.entries[idx] = entry
        c.moveToFront(idx)
    } else {
        entry := entry{key: key, value: value, index: 0, prev: InvalidIndex, next: c.head}
        c.entries[c.indexCounter] = entry
        c.indexMap[keyStr] = c.indexCounter

        if c.head != InvalidIndex {
            headEntry := c.entries[c.head]
            headEntry.prev = c.indexCounter
            c.entries[c.head] = headEntry
        }
        c.head = c.indexCounter

        if c.tail == InvalidIndex {
            c.tail = c.indexCounter
        }

        c.indexCounter++
    }

    c.adjustMemory(memSize)
    return nil
}

func (c *Cache) Del(key []byte) {
    keyStr := cx.B2s(key)
    c.mu.Lock()
    defer c.mu.Unlock()

    if idx, ok := c.indexMap[keyStr]; ok {
        entry := c.entries[idx]

        memSize := c.estimateMemory(entry.key, entry.value)
        c.adjustMemory(-memSize)

        c.detach(idx)

        delete(c.entries, idx)
        delete(c.indexMap, keyStr)
    }
}


