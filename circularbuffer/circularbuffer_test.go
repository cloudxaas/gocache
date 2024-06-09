package circularbuffer

import (
	"strconv"
	"testing"
)

func TestCircularBuffer(t *testing.T) {
	buffer := NewCircularBuffer(2)

	buffer.Add("a")
	buffer.Add("b")

	if key, _ := buffer.Remove(); key != "a" {
		t.Errorf("Expected key 'a', got '%s'", key)
	}

	buffer.Add("c")

	if key, _ := buffer.Remove(); key != "b" {
		t.Errorf("Expected key 'b', got '%s'", key)
	}
	if key, _ := buffer.Remove(); key != "c" {
		t.Errorf("Expected key 'c', got '%s'", key)
	}
	if _, ok := buffer.Remove(); ok {
		t.Errorf("Expected buffer to be empty")
	}
}

func TestCircularBufferPeek(t *testing.T) {
	buffer := NewCircularBuffer(3)

	buffer.Add("a")
	buffer.Add("b")
	buffer.Add("c")

	if key, ok := buffer.Peek(); !ok || key != "a" {
		t.Errorf("Expected key 'a', got '%s'", key)
	}

	buffer.Remove()
	if key, ok := buffer.Peek(); !ok || key != "b" {
		t.Errorf("Expected key 'b', got '%s'", key)
	}

	buffer.Remove()
	if key, ok := buffer.Peek(); !ok || key != "c" {
		t.Errorf("Expected key 'c', got '%s'", key)
	}

	buffer.Remove()
	if _, ok := buffer.Peek(); ok {
		t.Errorf("Expected buffer to be empty")
	}
}

func BenchmarkCircularBuffer_Add(b *testing.B) {
	buffer := NewCircularBuffer(25000000)
	keys := make([]string, 20000000)
	for i := 0; i < 20000000; i++ {
		keys[i] = strconv.Itoa(i)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buffer.Add(keys[i%20000000])
	}
}

func BenchmarkCircularBuffer_Remove(b *testing.B) {
	buffer := NewCircularBuffer(25000000)
	keys := make([]string, 20000000)
	for i := 0; i < 20000000; i++ {
		keys[i] = strconv.Itoa(i)
	}
	for i := 0; i < 20000000; i++ {
		buffer.Add(keys[i])
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buffer.Remove()
	}
}

func BenchmarkCircularBuffer_Peek(b *testing.B) {
	buffer := NewCircularBuffer(25000000)
	keys := make([]string, 20000000)
	for i := 0; i < 20000000; i++ {
		keys[i] = strconv.Itoa(i)
	}
	for i := 0; i < 20000000; i++ {
		buffer.Add(keys[i])
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buffer.Peek()
	}
}
