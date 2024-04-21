<p align="center">

  <img src="https://github.com/cloudxaas/gocache/assets/104323920/5948a699-64c8-47b8-a5d6-5afedb6a3976" width="40%" height="auto">

    
 <h2 align="center">" Let's <u>Acce</u>lerate <u>LRU</u> " - AccelruX</h2>
</p>

# One of The Fastest Zero Allocation LRU Cache in Golang (for key, value pairs in []byte) - AccelruX (cxlruxbytes)

# This X version has higher performance, much lower memory use, at the expense of maybe a bit lower hit ratio
predefine hashing function for hashing keys up to 4 billion keys (can make this to 2^64 if there are requests for it)

Difference with non-x version:

1. non-x uses map[string]int, x uses map[uint32]uint32 (special thanks and credit to phuslu/lru suggestion).

2. non-x does not need a hashing mechanism (so no collision issues), you need to predefine the hashing mechanism for x version.
 
3. non-x can use up to available memory and will save all keys as much as possible (may have better hit ratio), x version's hit ratio depends on the hash mechanism used, memory can be more efficient but hit ratio may suffer at collision of hash

Welcome to the repository for one of the fastest LRU cache implementations available for Go. This LRU cache is uniquely designed to limit the memory usage directly, rather than by the number of entries. This makes it ideal for applications where the memory footprint is critical, such as in embedded systems or high-performance computing environments where resources are tightly managed.

(adjust the limit of eviction count to match usage scenario)
if you have roughly (you think) 1 mil items (capacity), you can set eviction count to 1024 or more. depending on usage patterns too. set the parameters accordingly.

## Features

- **Memory-Size Limited**: Unlike other LRU caches that limit the number of entries, this cache controls the total memory used, allowing for better resource management in memory-constrained environments.
- **High Performance**: Designed with performance in mind, benchmarks demonstrate extremely low latency and zero allocations during operations, ensuring minimal impact on application throughput.
- **Concurrency Safe**: Implements synchronization to manage concurrent access, making it suitable for high-concurrency scenarios.
- **1 Item or Batch evictions**: When cache is filled, eviction by batch or 1 by 1, you can set this value.

## Motivation
Most lru cache available online for golang are set by capacity count, which means you may OOM your program. With this lru, once you set the memory size limit, you do not need to worry about OOM or garbage collection issues with zero allocation. OOM = Out of memory.

## Benchmarks

The cache has been rigorously benchmarked on a system with the following specifications:
- **OS**: Linux
- **Architecture**: AMD64
- **CPU**: AMD Ryzen 5 7640HS w/ Radeon 760M Graphics

## NOTE : AccelruX capacity is set in terms of MEMORY SIZE LIMIT and not NUMBER OF ITEMS LIMIT.

#### Note : Benchmark results (103,400,000 bytes cache with 10 byte key and 1024 byte value, using full cache without eviction [i think], same goes for the rest):
**These benchmarks are for reference only, do test it to verify claims. This is using full caching (without eviction for all)**

Default mode of running (sharded / non-sharded)
```
go test -bench=. -benchmem -benchtime=5s
goos: linux
goarch: amd64
pkg: github.com/cloudxaas/gocache/lrux/bytes
cpu: AMD Ryzen 5 7640HS w/ Radeon 760M Graphics     
BenchmarkHashicorpLRUSet-12       	123297104	        46.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashicorpLRUGet-12       	131854731	        45.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashicorpLRURemove-12    	705777380	         8.485 ns/op	       0 B/op	       0 allocs/op
BenchmarkGoFreeLRUSet-12          	233543793	        25.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkGoFreeLRUGet-12          	250022937	        24.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkGoFreeLRURemove-12       	655460384	         9.461 ns/op	       0 B/op	       0 allocs/op
BenchmarkOtterSet-12              	36643623	       163.1 ns/op	      64 B/op	       1 allocs/op
BenchmarkOtterGet-12              	86207361	        71.35 ns/op	       0 B/op	       0 allocs/op
BenchmarkOtterDelete-12           	503693424	        11.67 ns/op	       0 B/op	       0 allocs/op
BenchmarkPhusluLRUSet-12          	299558884	        19.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkPhusluLRUGet-12          	332732420	        18.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkPhusluLRUDelete-12       	623987602	         9.441 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesSet-12         	378581541	        15.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesGet-12         	396226093	        14.75 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesDel-12         	1000000000	         4.906 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/cloudxaas/gocache/lrux/bytes	123.575s

```
reference :

https://github.com/maypok86/otter

https://github.com/phuslu/lru

https://github.com/hashicorp/golang-lru

https://github.com/elastic/go-freelru

These benchmarks results vary, please adjust parameters to use for your own use case. 

*otter has 1 byte alloc / op, in these supposingly zero alloc benchmark cache

## NOTE : AccelruX capacity is set in terms of MEMORY SIZE LIMIT and not NUMBER OF ITEMS LIMIT.
## The benchmarks to test all case scenario is too wide and varied. Please do your own benchmark on it. There are sharded modes, non-sharded modes etc with different settings designed for different use case. Do your own research and do post comments / suggestions on issue tracker. We respect the purpose of different caches designed for use in different scenarios.  

## Accelru and AccelruX are designed for high performance with zero allocation, limited only by memory size instead of item capacity limit. 

## Usage

To use this cache, check the examples folder included, you can configure your own hash function, use xxh3 if you want faster hashing for larger key values > 24 bytes.

# Roadmap / Todo
- add more types / generic types, generic version is here, performance is kind of sad but usable. will improve.
https://github.com/cloudxaas/gocache/tree/main/lru
- add ttl support

Contributors welcome.
