package orderedset

import (
	"sync"
	"sync/atomic"
)

// OrderedSet represents a zero-allocated FIFO set with quick lookup
type OrderedSet struct {
	buffer    []string
	index     map[string]struct{}
	head, tail uint32
	capacity   int
	size       uint32
	mu         sync.Mutex
}

// NewOrderedSet initializes a new ordered set with the given capacity
func NewOrderedSet(capacity int) *OrderedSet {
	return &OrderedSet{
		buffer:   make([]string, capacity+1), // +1 to differentiate full and empty states
		index:    make(map[string]struct{}, capacity),
		capacity: capacity + 1,
	}
}

// Add inserts a key into the set
func (s *OrderedSet) Add(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.index[key]; exists {
		return
	}

	s.buffer[s.tail] = key
	s.index[key] = struct{}{}
	s.tail = (s.tail + 1) % uint32(s.capacity)

	if s.size == uint32(s.capacity-1) {
		delete(s.index, s.buffer[s.head])
		s.head = (s.head + 1) % uint32(s.capacity)
	} else {
		atomic.AddUint32(&s.size, 1)
	}
}

// Remove removes and returns the oldest key from the set
func (s *OrderedSet) Remove() (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if atomic.LoadUint32(&s.size) == 0 {
		return "", false // Set is empty
	}

	key := s.buffer[s.head]
	delete(s.index, key)
	s.head = (s.head + 1) % uint32(s.capacity)
	atomic.AddUint32(&s.size, ^uint32(0)) // Decrement size

	return key, true
}

// Contains checks if a key is in the set
func (s *OrderedSet) Contains(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.index[key]
	return exists
}

// IsEmpty checks if the set is empty
func (s *OrderedSet) IsEmpty() bool {
	return atomic.LoadUint32(&s.size) == 0
}

// IsFull checks if the set is full
func (s *OrderedSet) IsFull() bool {
	return atomic.LoadUint32(&s.size) == uint32(s.capacity-1)
}
