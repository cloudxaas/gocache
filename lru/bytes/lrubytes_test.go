package lrubytes

import (
    "testing"
)



func BenchmarkPut(b *testing.B) {
    // Updated to include the eviction count of 512
    cache := NewLRUCache(1024 * 1024 * 10, 512) // 10 MB max memory, evict 512 items at once
    keys := make([][]byte, 1000)
    values := make([][]byte, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = []byte{byte(i)}
        values[i] = make([]byte, 1024) // 1 KB values
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Put(keys[i%1000], values[i%1000])
    }
}

func BenchmarkGet(b *testing.B) {
    // Updated to include the eviction count of 512
    cache := NewLRUCache(1024 * 1024 * 10, 512) // 10 MB max memory, evict 512 items at once
    for i := 0; i < 1000; i++ {
        cache.Put([]byte{byte(i)}, make([]byte, 1024)) // 1 KB values
    }
    keys := make([][]byte, b.N)
    for i := 0; i < b.N; i++ {
        keys[i] = []byte{byte(i % 1000)}
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = cache.Get(keys[i])
    }
}

func BenchmarkDelete(b *testing.B) {
    // Updated to include the eviction count of 512
    cache := NewLRUCache(1024 * 1024 * 10, 512) // 10 MB max memory, evict 512 items at once
    keys := make([][]byte, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = []byte{byte(i)}
        cache.Put(keys[i], make([]byte, 1024)) // 1 KB values
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Delete(keys[i%1000])
    }
}

// Parallel benchmarks remain unchanged except the constructor
func BenchmarkPutParallel(b *testing.B) {
    cache := NewLRUCache(1024 * 1024 * 10, 512) // 10 MB max memory, evict 512 items at once
    keys := make([][]byte, 1000)
    values := make([][]byte, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = []byte{byte(i)}
        values[i] = make([]byte, 1024) // 1 KB values
    }

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for i := 0; pb.Next(); i++ {
            cache.Put(keys[i%1000], values[i%1000])
        }
    })
}

func BenchmarkGetParallel(b *testing.B) {
    cache := NewLRUCache(1024 * 1024 * 10, 512) // 10 MB max memory, evict 512 items at once
    for i := 0; i < 1000; i++ {
        cache.Put([]byte{byte(i)}, make([]byte, 1024)) // 1 KB values
    }

    keys := make([][]byte, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = []byte{byte(i)}
    }

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for i := 0; pb.Next(); i++ {
            _, _ = cache.Get(keys[i%1000])
        }
    })
}

func BenchmarkDeleteParallel(b *testing.B) {
    cache := NewLRUCache(1024 * 1024 * 10, 512) // 10 MB max memory, evict 512 items at once
    keys := make([][]byte, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = []byte{byte(i)}
        cache.Put(keys[i], make([]byte, 1024)) // 1 KB values
    }

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for i := 0; pb.Next(); i++ {
            cache.Delete(keys[i%1000])
        }
    })
}

func BenchmarkShardedPut(b *testing.B) {
    // Updated to use ShardedCache with 16 shards, 10 MB total memory, and 512 eviction count
    cache := NewShardedCache(16, 1024 * 1024 * 10, 512)
    keys := make([][]byte, 1000)
    values := make([][]byte, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = []byte{byte(i)}
        values[i] = make([]byte, 1024) // 1 KB values
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Put(keys[i%1000], values[i%1000])
    }
}

func BenchmarkShardedGet(b *testing.B) {
    // Updated to use ShardedCache
    cache := NewShardedCache(16, 1024 * 1024 * 10, 512)
    for i := 0; i < 1000; i++ {
        cache.Put([]byte{byte(i)}, make([]byte, 1024)) // 1 KB values
    }
    keys := make([][]byte, b.N)
    for i := 0; i < b.N; i++ {
        keys[i] = []byte{byte(i % 1000)}
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = cache.Get(keys[i])
    }
}

func BenchmarkShardedDelete(b *testing.B) {
    // Updated to use ShardedCache
    cache := NewShardedCache(16, 1024 * 1024 * 10, 512)
    keys := make([][]byte, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = []byte{byte(i)}
        cache.Put(keys[i], make([]byte, 1024)) // 1 KB values
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Delete(keys[i%1000])
    }
}

// Parallel benchmarks for sharded cache
func BenchmarkShardedPutParallel(b *testing.B) {
    cache := NewShardedCache(16, 1024 * 1024 * 10, 512)
    keys := make([][]byte, 1000)
    values := make([][]byte, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = []byte{byte(i)}
        values[i] = make([]byte, 1024) // 1 KB values
    }

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for i := 0; pb.Next(); i++ {
            cache.Put(keys[i%1000], values[i%1000])
        }
    })
}

func BenchmarkShardedGetParallel(b *testing.B) {
    cache := NewShardedCache(16, 1024 * 1024 * 10, 512)
    for i := 0; i < 1000; i++ {
        cache.Put([]byte{byte(i)}, make([]byte, 1024)) // 1 KB values
    }

    keys := make([][]byte, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = []byte{byte(i)}
    }

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for i := 0; pb.Next(); i++ {
            _, _ = cache.Get(keys[i%1000])
        }
    })
}

func BenchmarkShardedDeleteParallel(b *testing.B) {
    cache := NewShardedCache(16, 1024 * 1024 * 10, 512)
    keys := make([][]byte, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = []byte{byte(i)}
        cache.Put(keys[i], make([]byte, 1024)) // 1 KB values
    }

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for i := 0; pb.Next(); i++ {
            cache.Delete(keys[i%1000])
        }
    })
}
