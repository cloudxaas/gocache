package circularbuffer

import (
	"sync"
)

// KV represents a key-value pair
type KV struct {
	Key   string
	Value string
}

// CircularBufferKV represents a circular buffer for key-value pairs
type CircularBufferKV struct {
	buffer     []KV
	head, tail int
	capacity   int
	size       int
	mu         sync.Mutex
}

// NewCircularBufferKV initializes a new circular buffer for key-value pairs with the given capacity
func NewCircularBufferKV(capacity int) *CircularBufferKV {
	return &CircularBufferKV{
		buffer:   make([]KV, capacity+1), // +1 to differentiate full and empty states
		capacity: capacity + 1,
	}
}

// Add inserts a key-value pair into the circular buffer
func (cb *CircularBufferKV) Add(key, value string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.buffer[cb.tail] = KV{Key: key, Value: value}
	cb.tail = (cb.tail + 1) % cb.capacity

	if cb.size == cb.capacity-1 {
		cb.head = (cb.head + 1) % cb.capacity
	} else {
		cb.size++
	}
}

// Remove removes and returns the oldest key-value pair from the circular buffer
func (cb *CircularBufferKV) Remove() (KV, bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.size == 0 {
		return KV{}, false // Buffer is empty
	}

	kv := cb.buffer[cb.head]
	cb.head = (cb.head + 1) % cb.capacity
	cb.size--

	return kv, true
}

// Peek returns the oldest key-value pair without removing it
func (cb *CircularBufferKV) Peek() (KV, bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.size == 0 {
		return KV{}, false // Buffer is empty
	}

	return cb.buffer[cb.head], true
}

// IsEmpty checks if the buffer is empty
func (cb *CircularBufferKV) IsEmpty() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.size == 0
}

// IsFull checks if the buffer is full
func (cb *CircularBufferKV) IsFull() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.size == cb.capacity-1
}
