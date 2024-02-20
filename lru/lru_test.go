package cxcachelru

import (
	"strconv"
	"testing"
)

func BenchmarkLRUCache_Put(b *testing.B) {
	cache := NewLRUCache(10000) // Assuming a cache size of 10000 for the benchmark
	testData := make([][]byte, b.N)
	for i := range testData {
		testData[i] = []byte(strconv.Itoa(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Put(testData[i], testData[i])
	}
}

func BenchmarkLRUCache_Get(b *testing.B) {
	cache := NewLRUCache(10000) // Pre-fill the cache to avoid misses
	testData := make([][]byte, b.N)
	for i := range testData {
		testData[i] = []byte(strconv.Itoa(i))
		cache.Put(testData[i], testData[i])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(testData[i])
	}
}

func BenchmarkLRUCache_PutAndGet(b *testing.B) {
	cache := NewLRUCache(10000) // A mix of puts and gets
	testData := make([][]byte, b.N)
	for i := range testData {
		testData[i] = []byte(strconv.Itoa(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			cache.Put(testData[i], testData[i])
		} else {
			cache.Get(testData[i])
		}
	}
}
