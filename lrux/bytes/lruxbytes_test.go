package lruxbytes

import (
    "testing"
    "fmt"

    cx "github.com/cloudxaas/gocx"
    "github.com/phuslu/lru"
    "github.com/maypok86/otter"
    "github.com/elastic/go-freelru"
    hashicorp "github.com/hashicorp/golang-lru/v2"
)

const (
	offset32 uint32 = 2166136261
	prime32  uint32 = 16777619
)

// FNV-1a hash function for byte slices
func FNV1aHash(key []byte) uint32 {
	var hash uint32 = offset32
	for _, c := range key {
		hash ^= uint32(c)
		hash *= prime32
	}
	return hash
}

// FNV-1a hash function for strings
func FNV1aHashStr(key string) uint32 {
    var hash uint32 = offset32
    for _, c := range key {
        hash ^= uint32(c)
        hash *= prime32
    }
    return hash
}


var (
    keys   []string
    bKeys  [][]byte
    values [][]byte
    u64Values []uint64
    strValues []string
)

func init() {
    keys = make([]string, 100000)
    bKeys = make([][]byte, 100000)
    values = make([][]byte, 100000)
    u64Values = make([]uint64, 100000)
    strValues = make([]string, 100000)
    for i := 0; i < 100000; i++ {
        keys[i] = fmt.Sprintf("key%d", i)
        bKeys[i] = []byte{byte(i)}
        values[i] = make([]byte, 1024) // 1 KB values
        u64Values[i] = uint64(i)
        strValues[i] = "value"
    }
}

func BenchmarkHashicorpLRUSet(b *testing.B) {
    cache, _ := hashicorp.New[string, []byte](100 * 1000)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Add(keys[i%100000], values[i%100000])
    }
}

func BenchmarkHashicorpLRUGet(b *testing.B) {
    cache, _ := hashicorp.New[string, []byte](100 * 1000)
    for i := 0; i < 100000; i++ {
        cache.Add(keys[i], values[i])
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = cache.Get(keys[i%100000])
    }
}

func BenchmarkHashicorpLRURemove(b *testing.B) {
    cache, _ := hashicorp.New[string, []byte](100 * 1000)
    for i := 0; i < 100000; i++ {
        cache.Add(keys[i], values[i])
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = cache.Remove(keys[i%100000])
    }
}

func BenchmarkGoFreeLRUSet(b *testing.B) {
    lru, err := freelru.New[string, uint64](100 * 1000, FNV1aHashStr)
    if err != nil {
        b.Fatal(err)
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        lru.Add(keys[i%100000], u64Values[i%100000])
    }
}

func BenchmarkGoFreeLRUGet(b *testing.B) {
    lru, err := freelru.New[string, uint64](100 * 1000, FNV1aHashStr)
    if err != nil {
        b.Fatal(err)
    }
    for i := 0; i < 100000; i++ {
        lru.Add(keys[i], u64Values[i])
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = lru.Get(keys[i%100000])
    }
}

func BenchmarkGoFreeLRURemove(b *testing.B) {
    lru, err := freelru.New[string, uint64](100 * 1000, FNV1aHashStr)
    if err != nil {
        b.Fatal(err)
    }
    for i := 0; i < 100000; i++ {
        lru.Add(keys[i], u64Values[i])
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = lru.Remove(keys[i%100000])
    }
}

func BenchmarkOtterSet(b *testing.B) {
    cache, err := otter.MustBuilder[string, string](100 * 1000).
        CollectStats().
        Cost(func(key string, value string) uint32 {
            return 1
        }).
        Build()
    if err != nil {
        b.Fatal(err)
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Set(keys[i%100000], strValues[i%100000])
    }
}

func BenchmarkOtterGet(b *testing.B) {
    cache, err := otter.MustBuilder[string, string](100 * 1000).
        CollectStats().
        Cost(func(key string, value string) uint32 {
            return 1
        }).
        Build()
    if err != nil {
        b.Fatal(err)
    }
    for i := 0; i < 100000; i++ {
        cache.Set(keys[i], strValues[i])
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = cache.Get(keys[i%100000])
    }
}

func BenchmarkOtterDelete(b *testing.B) {
    cache, err := otter.MustBuilder[string, string](100 * 1000).
        CollectStats().
        Cost(func(key string, value string) uint32 {
            return 1
        }).
        Build()
    if err != nil {
        b.Fatal(err)
    }
    for i := 0; i < 100000; i++ {
        cache.Set(keys[i], strValues[i])
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Delete(keys[i%100000])
    }
}

func BenchmarkPhusluLRUSet(b *testing.B) {
    cache := lru.NewLRUCache[string, []byte](100 * 1000)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Set(cx.B2s(bKeys[i%100000]), values[i%100000])
    }
}

func BenchmarkPhusluLRUGet(b *testing.B) {
    cache := lru.NewLRUCache[string, []byte](100 * 1000)
    for i := 0; i < 100000; i++ {
        cache.Set(cx.B2s(bKeys[i]), values[i])
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = cache.Get(cx.B2s(bKeys[i%100000]))
    }
}

func BenchmarkPhusluLRUDelete(b *testing.B) {
    cache := lru.NewLRUCache[string, []byte](100 * 1000)
    for i := 0; i < 100000; i++ {
        cache.Set(cx.B2s(bKeys[i]), values[i])
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Delete(cx.B2s(bKeys[i%100000]))
    }
}

func BenchmarkCXLRUBytesSet(b *testing.B) {
    cache := NewLRUCache(100000 * (10+1024), 1, FNV1aHash)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Set(bKeys[i%100000], values[i%100000])
    }
}

func BenchmarkCXLRUBytesGet(b *testing.B) {
    cache := NewLRUCache(100000 * (10+1024), 1, FNV1aHash)
    for i := 0; i < 100000; i++ {
        cache.Set(bKeys[i], values[i])
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = cache.Get(bKeys[i%100000])
    }
}

func BenchmarkCXLRUBytesDel(b *testing.B) {
    cache := NewLRUCache(100000 * (10+1024), 1, FNV1aHash)
    for i := 0; i < 100000; i++ {
        cache.Set(bKeys[i], values[i])
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Del(bKeys[i%100000])
    }
}
