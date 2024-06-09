package circularbuffer

import (
	"sync"
)

// CircularBuffer represents a circular buffer
type CircularBuffer struct {
	buffer     []string
	head, tail int
	capacity   int
	size       int
	mu         sync.Mutex
}

// NewCircularBuffer initializes a new circular buffer with the given capacity
func NewCircularBuffer(capacity int) *CircularBuffer {
	return &CircularBuffer{
		buffer:   make([]string, capacity+1), // +1 to differentiate full and empty states
		capacity: capacity + 1,
	}
}

// Add inserts a key into the circular buffer
func (cb *CircularBuffer) Add(key string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.buffer[cb.tail] = key
	cb.tail = (cb.tail + 1) % cb.capacity

	if cb.size == cb.capacity-1 {
		cb.head = (cb.head + 1) % cb.capacity
	} else {
		cb.size++
	}
}

// Remove removes and returns the oldest key from the circular buffer
func (cb *CircularBuffer) Remove() (string, bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.size == 0 {
		return "", false // Buffer is empty
	}

	key := cb.buffer[cb.head]
	cb.head = (cb.head + 1) % cb.capacity
	cb.size--

	return key, true
}

// Peek returns the oldest key without removing it
func (cb *CircularBuffer) Peek() (string, bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.size == 0 {
		return "", false // Buffer is empty
	}

	return cb.buffer[cb.head], true
}

// IsEmpty checks if the buffer is empty
func (cb *CircularBuffer) IsEmpty() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.size == 0
}

// IsFull checks if the buffer is full
func (cb *CircularBuffer) IsFull() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.size == cb.capacity-1
}
