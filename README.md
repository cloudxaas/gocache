 

<p align="center">

  <img src="https://github.com/cloudxaas/gocache/assets/104323920/5948a699-64c8-47b8-a5d6-5afedb6a3976" width="40%" height="auto" >
    
   <h1 align="center">Super Fast LRU Cache for Golang</h1>
 <h3 align="center">" Let's <u>Acce</u>lerate <u>LRU</u> " - Accelru</h3>
</p>

# One of The Fastest Zero Allocation LRU Cache for Golang 
(for key, value pairs in []byte)

## Accelru (cxlrubytes)

Supposingly having the best cache hit ratio (in zero allocation class) with "optimum" memory usage.

Please contribute to make it better.
Feedback / comments / suggestions on improvement appreciated (stars too).

Check [lru/bytes](https://github.com/cloudxaas/gocache/tree/main/lru/bytes) for details for key, value in []byte (tested).

Check verion X for extreme speed, lower memory footprint at the expense of a bit lower hit ratio
[lru/bytes](https://github.com/cloudxaas/gocache/tree/main/lrux/bytes) for details for key, value in []byte (tested).

Check [lru](https://github.com/cloudxaas/gocache/tree/main/lru) for Any type details (untested).

## Motivation

Most current (year 2024) golang lru implementations are either not as fast as this, or needed capacity count of items as input parameter, this can result in "OOM" or not being able to fully utilize the memory capacity available.

cxlrubytes thus is designed for:
1. High performance
2. Zero allocation (so no garbage collection)
3. Maximizing memory usage (but not being limited by item capacity)

Will do other input parameters in future but currently, converting everything to []byte and using this gives wonderful results.

google-site-verification: google4b244ca4683e045f.html
