package lruxbytes

import (
	"testing"

	cx "github.com/cloudxaas/gocx"
	"github.com/phuslu/lru"
)

// Helper function to convert byte slices to string for use as keys in the lru.Cache

func BenchmarkPhusluLRUSet(b *testing.B) {
	cache := lru.NewLRUCache[string, []byte](1024 * 100)
	keys := make([][]byte, 100000)
	values := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		values[i] = make([]byte, 1024) // 1 KB values
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(cx.B2s(keys[i%1000]), values[i%1000])
	}
}

func BenchmarkPhusluLRUGet(b *testing.B) {
	cache := lru.NewLRUCache[string, []byte](1024 * 100)
	for i := 0; i < 100000; i++ {
		cache.Set(cx.B2s([]byte{byte(i)}), make([]byte, 1024)) // 1 KB values
	}
	keys := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = []byte{byte(i % 100000)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get(cx.B2s(keys[i]))
	}
}

func BenchmarkPhusluLRUDelete(b *testing.B) {
	cache := lru.NewLRUCache[string, []byte](1024 * 100)
	keys := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		cache.Set(cx.B2s(keys[i]), make([]byte, 1024)) // 1 KB values
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Delete(cx.B2s(keys[i%100000]))
	}
}

func BenchmarkCXLRUBytesSet(b *testing.B) {
	// Updated to include the eviction count of 1
	cache := NewLRUCache(1024*100, 1) // 10 MB max memory, evict 512 items at once
	keys := make([][]byte, 100000)
	values := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		values[i] = make([]byte, 1024) // 1 KB values
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(keys[i%100000], values[i%100000])
	}
}

func BenchmarkCXLRUBytesGet(b *testing.B) {
	// Updated to include the eviction count of 1
	cache := NewLRUCache(1024*100, 1) // 10 MB max memory, evict 512 items at once
	for i := 0; i < 100000; i++ {
		cache.Set([]byte{byte(i)}, make([]byte, 1024)) // 1 KB values
	}
	keys := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = []byte{byte(i % 100000)}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get(keys[i])
	}
}

func BenchmarkCXLRUBytesDel(b *testing.B) {
	// Updated to include the eviction count of 1
	cache := NewLRUCache(1024*100, 1) // 10 MB max memory, evict 512 items at once
	keys := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		cache.Set(keys[i], make([]byte, 1024)) // 1 KB values
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Del(keys[i%100000])
	}
}

// Parallel benchmarks remain unchanged except the constructor
func BenchmarkCXLRUBytesSetParallel(b *testing.B) {
	cache := NewLRUCache(1024*100, 1) // 10 MB max memory, evict 1 items at once
	keys := make([][]byte, 100000)
	values := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		values[i] = make([]byte, 1024) // 1 KB values
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			cache.Set(keys[i%100000], values[i%100000])
		}
	})
}

func BenchmarkCXLRUBytesGetParallel(b *testing.B) {
	cache := NewLRUCache(1024*100, 1) // 10 MB max memory, evict 1 items at once
	for i := 0; i < 100000; i++ {
		cache.Set([]byte{byte(i)}, make([]byte, 1024)) // 1 KB values
	}

	keys := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			_, _ = cache.Get(keys[i%100000])
		}
	})
}

func BenchmarkCXLRUBytesDelParallel(b *testing.B) {
	cache := NewLRUCache(1024*100, 1) // 10 MB max memory, evict 1 items at once
	keys := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		cache.Set(keys[i], make([]byte, 1024)) // 1 KB values
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			cache.Del(keys[i%100000])
		}
	})
}

func BenchmarkCXLRUBytesShardedSet(b *testing.B) {
	// Updated to use ShardedCache with 16 shards, 10 MB total memory, and 1 eviction count
	cache := NewShardedCache(16, 1024*100, 1)
	keys := make([][]byte, 100000)
	values := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		values[i] = make([]byte, 1024) // 1 KB values
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(keys[i%100000], values[i%100000])
	}
}

func BenchmarkCXLRUBytesShardedGet(b *testing.B) {
	// Updated to use ShardedCache
	cache := NewShardedCache(16, 1024*100, 1)
	for i := 0; i < 100000; i++ {
		cache.Set([]byte{byte(i)}, make([]byte, 1024)) // 1 KB values
	}
	keys := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = []byte{byte(i % 100000)}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get(keys[i])
	}
}

func BenchmarkCXLRUBytesShardedDel(b *testing.B) {
	// Updated to use ShardedCache
	cache := NewShardedCache(16, 1024*100, 1)
	keys := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		cache.Set(keys[i], make([]byte, 1024)) // 1 KB values
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Del(keys[i%100000])
	}
}

// Parallel benchmarks for sharded cache
func BenchmarkCXLRUBytesShardedSetParallel(b *testing.B) {
	cache := NewShardedCache(16, 1024*100, 1)
	keys := make([][]byte, 100000)
	values := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		values[i] = make([]byte, 1024) // 1 KB values
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			cache.Set(keys[i%100000], values[i%100000])
		}
	})
}

func BenchmarkCXLRUBytesShardedGetParallel(b *testing.B) {
	cache := NewShardedCache(16, 1024*100, 1)
	for i := 0; i < 100000; i++ {
		cache.Set([]byte{byte(i)}, make([]byte, 1024)) // 1 KB values
	}

	keys := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			_, _ = cache.Get(keys[i%100000])
		}
	})
}

func BenchmarkCXLRUBytesShardedDelParallel(b *testing.B) {
	cache := NewShardedCache(16, 1024*100, 1)
	keys := make([][]byte, 100000)
	for i := 0; i < 100000; i++ {
		keys[i] = []byte{byte(i)}
		cache.Set(keys[i], make([]byte, 1024)) // 1 KB values
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			cache.Del(keys[i%100000])
		}
	})
}
