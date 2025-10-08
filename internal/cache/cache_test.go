package cache

import (
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := NewCache(time.Minute)

	cache.Set("key1", "value1")

	value, ok := cache.Get("key1")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if value != "value1" {
		t.Fatalf("expected 'value1', got %v", value)
	}
}

func TestCache_GetNonExistent(t *testing.T) {
	cache := NewCache(time.Minute)

	value, ok := cache.Get("nonexistent")
	if ok {
		t.Fatal("expected key to not exist")
	}
	if value != nil {
		t.Fatal("expected value to be nil")
	}
}

func TestCache_Expiration(t *testing.T) {
	cache := NewCache(100 * time.Millisecond)

	cache.Set("key1", "value1")

	// Should exist immediately
	value, ok := cache.Get("key1")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if value != "value1" {
		t.Fatalf("expected 'value1', got %v", value)
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should not exist after expiration
	value, ok = cache.Get("key1")
	if ok {
		t.Fatal("expected key to be expired")
	}
	if value != nil {
		t.Fatal("expected value to be nil")
	}
}

func TestCache_NoExpiration(t *testing.T) {
	cache := NewCache(0) // No expiration

	cache.Set("key1", "value1")

	// Should exist even after time passes
	time.Sleep(100 * time.Millisecond)

	value, ok := cache.Get("key1")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if value != "value1" {
		t.Fatalf("expected 'value1', got %v", value)
	}
}

func TestCache_Invalidate(t *testing.T) {
	cache := NewCache(time.Minute)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	// Both should exist
	_, ok1 := cache.Get("key1")
	_, ok2 := cache.Get("key2")
	if !ok1 || !ok2 {
		t.Fatal("expected both keys to exist")
	}

	// Invalidate one key
	cache.Invalidate("key1")

	// One should exist, one should not
	_, ok1 = cache.Get("key1")
	_, ok2 = cache.Get("key2")
	if ok1 {
		t.Fatal("expected key1 to be invalidated")
	}
	if !ok2 {
		t.Fatal("expected key2 to still exist")
	}
}

func TestCache_InvalidateAll(t *testing.T) {
	cache := NewCache(time.Minute)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// All should exist
	_, ok1 := cache.Get("key1")
	_, ok2 := cache.Get("key2")
	_, ok3 := cache.Get("key3")
	if !ok1 || !ok2 || !ok3 {
		t.Fatal("expected all keys to exist")
	}

	// Invalidate all
	cache.InvalidateAll()

	// None should exist
	_, ok1 = cache.Get("key1")
	_, ok2 = cache.Get("key2")
	_, ok3 = cache.Get("key3")
	if ok1 || ok2 || ok3 {
		t.Fatal("expected all keys to be invalidated")
	}
}

func TestCache_Overwrite(t *testing.T) {
	cache := NewCache(time.Minute)

	cache.Set("key1", "value1")

	value, ok := cache.Get("key1")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if value != "value1" {
		t.Fatalf("expected 'value1', got %v", value)
	}

	// Overwrite with new value
	cache.Set("key1", "value2")

	value, ok = cache.Get("key1")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if value != "value2" {
		t.Fatalf("expected 'value2', got %v", value)
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := NewCache(time.Minute)

	// Test concurrent access
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			key := "key" + string(rune(i))
			cache.Set(key, "value"+string(rune(i)))
			cache.Get(key)
			cache.Invalidate(key)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestCache_DifferentTypes(t *testing.T) {
	cache := NewCache(time.Minute)

	// Test different value types
	cache.Set("string", "hello")
	cache.Set("int", 42)
	cache.Set("bool", true)
	cache.Set("slice", []string{"a", "b", "c"})
	cache.Set("map", map[string]int{"a": 1, "b": 2})

	// Check string
	value, ok := cache.Get("string")
	if !ok || value != "hello" {
		t.Fatal("string value incorrect")
	}

	// Check int
	value, ok = cache.Get("int")
	if !ok || value != 42 {
		t.Fatal("int value incorrect")
	}

	// Check bool
	value, ok = cache.Get("bool")
	if !ok || value != true {
		t.Fatal("bool value incorrect")
	}

	// Check slice
	value, ok = cache.Get("slice")
	if !ok {
		t.Fatal("slice value not found")
	}
	slice, ok := value.([]string)
	if !ok || len(slice) != 3 {
		t.Fatal("slice value incorrect")
	}

	// Check map
	value, ok = cache.Get("map")
	if !ok {
		t.Fatal("map value not found")
	}
	m, ok := value.(map[string]int)
	if !ok || m["a"] != 1 || m["b"] != 2 {
		t.Fatal("map value incorrect")
	}
}
