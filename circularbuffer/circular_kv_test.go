package circularbuffer

import (
	"testing"
)

func BenchmarkCircularBufferKV_Add(b *testing.B) {
	cb := NewCircularBufferKV(1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cb.Add("key", "value")
		}
	})
}

func BenchmarkCircularBufferKV_Remove(b *testing.B) {
	cb := NewCircularBufferKV(1000)
	for i := 0; i < 1000; i++ {
		cb.Add("key", "value")
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cb.Remove()
		}
	})
}

func BenchmarkCircularBufferKV_Peek(b *testing.B) {
	cb := NewCircularBufferKV(1000)
	for i := 0; i < 1000; i++ {
		cb.Add("key", "value")
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cb.Peek()
		}
	})
}

func TestCircularBufferKV_AddRemove(t *testing.T) {
	cb := NewCircularBufferKV(3)
	cb.Add("key1", "value1")
	cb.Add("key2", "value2")
	cb.Add("key3", "value3")

	if kv, ok := cb.Remove(); !ok || kv.Key != "key1" || kv.Value != "value1" {
		t.Errorf("expected (key1, value1), got (%s, %s)", kv.Key, kv.Value)
	}

	if kv, ok := cb.Peek(); !ok || kv.Key != "key2" || kv.Value != "value2" {
		t.Errorf("expected (key2, value2) at peek, got (%s, %s)", kv.Key, kv.Value)
	}

	cb.Add("key4", "value4")
	if kv, ok := cb.Remove(); !ok || kv.Key != "key2" || kv.Value != "value2" {
		t.Errorf("expected (key2, value2), got (%s, %s)", kv.Key, kv.Value)
	}

	if cb.IsFull() {
		t.Error("expected buffer to not be full")
	}

	cb.Add("key5", "value5")
	if !cb.IsFull() {
		t.Error("expected buffer to be full")
	}

	cb.Remove()
	cb.Remove()
	cb.Remove()
	if !cb.IsEmpty() {
		t.Error("expected buffer to be empty")
	}
}

func TestCircularBufferKV_IsFull(t *testing.T) {
	cb := NewCircularBufferKV(2) // effective capacity is 2
	cb.Add("key1", "value1")
	cb.Add("key2", "value2")

	if !cb.IsFull() {
		t.Error("expected buffer to be full")
	}
}

func TestCircularBufferKV_Peek(t *testing.T) {
	cb := NewCircularBufferKV(3)
	cb.Add("key1", "value1")
	cb.Add("key2", "value2")

	if kv, ok := cb.Peek(); !ok || kv.Key != "key1" || kv.Value != "value1" {
		t.Errorf("expected (key1, value1), got (%s, %s)", kv.Key, kv.Value)
	}

	cb.Remove()
	if kv, ok := cb.Peek(); !ok || kv.Key != "key2" || kv.Value != "value2" {
		t.Errorf("expected (key2, value2), got (%s, %s)", kv.Key, kv.Value)
	}

	cb.Add("key3", "value3")
	cb.Add("key4", "value4")
	if kv, ok := cb.Peek(); !ok || kv.Key != "key2" || kv.Value != "value2" {
		t.Errorf("expected (key2, value2), got (%s, %s)", kv.Key, kv.Value)
	}
}
