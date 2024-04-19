package cxcachelru

import (
    "sync"
    "testing"
)

// Example type implementing Sizer interface
type IntSizer int

func (i IntSizer) Size() int64 {
    return int64(8) // Simulate an integer size of 8 bytes
}

// TestLRUCache_Parallel ensures the LRU cache operates correctly under concurrent access.
func TestLRUCache_Parallel(t *testing.T) {
    cache := NewLRUCache[IntSizer, IntSizer](1000, 100) // Maximum memory units and batch size

    // Use multiple keys and values for the operations
    keys := []IntSizer{1, 2, 3, 4, 5}
    values := []IntSizer{10, 20, 30, 40, 50}

    t.Run("Parallel Put and Get", func(t *testing.T) {
        t.Parallel()
        var wg sync.WaitGroup
        for i := 0; i < 100; i++ { // Launch 100 goroutines
            wg.Add(1)
            go func() {
                defer wg.Done()
                for j, key := range keys {
                    allocsPut := testing.AllocsPerRun(10, func() {
                        cache.Put(key, values[j])
                    })

                    // Ensure there are no allocations during the Put operations
                    if allocsPut > 0 {
                        t.Errorf("Put operation allocates %f times, want 0 allocations", allocsPut)
                    }

                    value, found := cache.Get(key)
                    if !found || value != values[j] {
                        t.Errorf("Get(%v) = %v, want %v", key, value, values[j])
                    }

                    allocsGet := testing.AllocsPerRun(10, func() {
                        cache.Get(key)
                    })

                    // Ensure minimal allocations during the Get operations
                    if allocsGet > 0 {
                        t.Errorf("Get operation allocates %f times, want 0 allocations", allocsGet)
                    }
                }
            }()
        }
        wg.Wait()
    })

    t.Run("Parallel Delete", func(t *testing.T) {
        t.Parallel()
        var wg sync.WaitGroup
        for i := 0; i < 100; i++ { // Launch 100 goroutines
            wg.Add(1)
            go func() {
                defer wg.Done()
                for _, key := range keys {
                    allocsDelete := testing.AllocsPerRun(10, func() {
                        cache.Delete(key)
                    })

                    // Check memory allocations for Delete operation
                    if allocsDelete > 0 {
                        t.Errorf("Delete operation for key %v allocates %f times, want 0 allocations", key, allocsDelete)
                    }

                    _, found := cache.Get(key)
                    if found {
                        t.Errorf("Get(%v) after Delete should not find the item, but did", key)
                    }
                }
            }()
        }
        wg.Wait()
    })
}
