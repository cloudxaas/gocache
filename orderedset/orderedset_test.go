package orderedset

import (
	"strconv"
	"testing"
)

func TestOrderedSet(t *testing.T) {
	set := NewOrderedSet(2)

	set.Add("a")
	set.Add("b")

	if !set.Contains("a") {
		t.Errorf("Expected set to contain 'a'")
	}

	set.Add("c")

	if set.Contains("a") {
		t.Errorf("Expected set to not contain 'a' after overflow")
	}
	if !set.Contains("b") {
		t.Errorf("Expected set to contain 'b'")
	}
	if !set.Contains("c") {
		t.Errorf("Expected set to contain 'c'")
	}
}

func BenchmarkOrderedSet_Add(b *testing.B) {
	set := NewOrderedSet(100)
	keys := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		keys[i] = strconv.Itoa(i)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		set.Add(keys[i%1000])
	}
}

func BenchmarkOrderedSet_Contains(b *testing.B) {
	set := NewOrderedSet(100)
	keys := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		keys[i] = strconv.Itoa(i)
		set.Add(keys[i])
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		set.Contains(keys[i%1000])
	}
}
