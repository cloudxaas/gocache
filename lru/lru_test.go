package cxcachelru

import (
    "testing"
)

// IntSizer is a simple integer type that implements the Sizer interface
type IntSizer int

func (i IntSizer) Size() int64 {
    return int64(8) // Assume size of int is 8 bytes
}

// setupCache helps in setting up an LRU cache with initial values for benchmarking
func setupCache(size int64, batchSize int) *Cache[IntSizer, IntSizer] {
    cache := NewLRUCache[IntSizer, IntSizer](size, batchSize)
    for i := 0; i < 100; i++ {
        cache.Put(IntSizer(i), IntSizer(i*10))
    }
    return cache
}

// BenchmarkPut measures the performance of the Put method
func BenchmarkPut(b *testing.B) {
    cache := setupCache(1000, 100)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Put(IntSizer(i), IntSizer(i*10))
    }
}

// BenchmarkGet measures the performance of the Get method
func BenchmarkGet(b *testing.B) {
    cache := setupCache(1000, 100)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Get(IntSizer(i))
    }
}

// BenchmarkDelete measures the performance of the Delete method
func BenchmarkDelete(b *testing.B) {
    cache := setupCache(1000, 100)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Delete(IntSizer(i))
    }
}
