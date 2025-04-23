package ztype_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/zzjcool/goutils/ztype"
)

func TestNewSafeBiMap(t *testing.T) {
	t.Run("Test creating a new SafeBiMap with integer keys and string values", func(t *testing.T) {
		intStrMap := ztype.NewSafeBiMap[int, string]()
		if reflect.TypeOf(intStrMap) != reflect.TypeOf(&ztype.SafeBiMap[int, string]{}) {
			t.Errorf("Expected SafeBiMap[int, string], got %T", intStrMap)
		}
	})

	t.Run("Test creating a new SafeBiMap with string keys and integer values", func(t *testing.T) {
		strIntMap := ztype.NewSafeBiMap[string, int]()
		if reflect.TypeOf(strIntMap) != reflect.TypeOf(&ztype.SafeBiMap[string, int]{}) {
			t.Errorf("Expected SafeBiMap[string, int], got %T", strIntMap)
		}
	})
}

func TestSafeBiMap_Get(t *testing.T) {
	// Create a new SafeBiMap instance
	safeBiMap := ztype.NewSafeBiMap[int, string]()

	// Add some values to the SafeBiMap
	safeBiMap.Set(1, "Value 1")
	safeBiMap.Set(2, "Value 2")

	// Test the Get method
	value, ok := safeBiMap.Get(1)
	if !ok || value != "Value 1" {
		t.Errorf("Expected value 'Value 1', got '%s'", value)
	}
}

func TestSafeBiMap_GetKey(t *testing.T) {
	// Create a new SafeBiMap instance
	safeBiMap := ztype.NewSafeBiMap[int, string]()

	// Add some values to the SafeBiMap
	safeBiMap.Set(1, "Value 1")
	safeBiMap.Set(2, "Value 2")

	// Test the GetKey method
	key, ok := safeBiMap.GetKey("Value 1")
	if !ok || key != 1 {
		t.Errorf("Expected key 1, got %d", key)
	}
}

func TestSafeBiMap_Set(t *testing.T) {
	// Create a new SafeBiMap instance
	safeBiMap := ztype.NewSafeBiMap[int, string]()

	// Test the Set method
	safeBiMap.Set(1, "Value 1")
	value, ok := safeBiMap.Get(1)
	if !ok || value != "Value 1" {
		t.Errorf("Expected value 'Value 1', got '%s'", value)
	}

	key, ok := safeBiMap.GetKey("Value 1")
	if !ok || key != 1 {
		t.Errorf("Expected key 1, got %d", key)
	}

	// Test overwriting an existing key
	safeBiMap.Set(1, "New Value 1")
	value, ok = safeBiMap.Get(1)
	if !ok || value != "New Value 1" {
		t.Errorf("Expected value 'New Value 1', got '%s'", value)
	}

	// Verify old value is no longer mapped
	_, ok = safeBiMap.GetKey("Value 1")
	if ok {
		t.Error("Expected old value 'Value 1' to be removed from reverse mapping")
	}

	// Test overwriting an existing value
	safeBiMap.Set(2, "Value 2")
	safeBiMap.Set(3, "Value 2") // This should overwrite the key for "Value 2"
	
	// Verify new key is mapped to the value
	key, ok = safeBiMap.GetKey("Value 2")
	if !ok || key != 3 {
		t.Errorf("Expected key 3, got %d", key)
	}

	// Verify old key is no longer mapped
	_, ok = safeBiMap.Get(2)
	if ok {
		t.Error("Expected old key 2 to be removed from forward mapping")
	}
}

func TestSafeBiMap_DeleteKey(t *testing.T) {
	// Create a new SafeBiMap instance
	safeBiMap := ztype.NewSafeBiMap[int, string]()

	// Add some values to the SafeBiMap
	safeBiMap.Set(1, "Value 1")
	safeBiMap.Set(2, "Value 2")

	// Delete a value from the SafeBiMap by key
	safeBiMap.DeleteKey(1)

	// Test the DeleteKey method
	_, ok := safeBiMap.Get(1)
	if ok {
		t.Error("Expected key 1 to be deleted, but it still exists")
	}

	// Verify value is also removed from reverse mapping
	_, ok = safeBiMap.GetKey("Value 1")
	if ok {
		t.Error("Expected value 'Value 1' to be deleted from reverse mapping")
	}
}

func TestSafeBiMap_DeleteValue(t *testing.T) {
	// Create a new SafeBiMap instance
	safeBiMap := ztype.NewSafeBiMap[int, string]()

	// Add some values to the SafeBiMap
	safeBiMap.Set(1, "Value 1")
	safeBiMap.Set(2, "Value 2")

	// Delete a key-value pair by value
	safeBiMap.DeleteValue("Value 1")

	// Test the DeleteValue method
	_, ok := safeBiMap.GetKey("Value 1")
	if ok {
		t.Error("Expected value 'Value 1' to be deleted, but it still exists")
	}

	// Verify key is also removed from forward mapping
	_, ok = safeBiMap.Get(1)
	if ok {
		t.Error("Expected key 1 to be deleted from forward mapping")
	}
}

func TestSafeBiMap_RangeByKey(t *testing.T) {
	// Create a map with just two entries to simplify testing
	m := ztype.NewSafeBiMap[string, int]()
	m.Set("foo", 1)
	m.Set("bar", 2)

	// Test that all entries are visited when callback always returns true
	visitedAll := make(map[string]bool)
	m.RangeByKey(func(key string, value int) bool {
		visitedAll[key] = true
		return true // continue iteration
	})

	if len(visitedAll) != 2 {
		t.Errorf("expected 2 keys to be visited, but got %d", len(visitedAll))
	}

	if !visitedAll["foo"] || !visitedAll["bar"] {
		t.Errorf("not all keys were visited: %v", visitedAll)
	}

	// Test early termination with a new map
	m2 := ztype.NewSafeBiMap[string, int]()
	m2.Set("key1", 1)
	m2.Set("key2", 2)
	m2.Set("key3", 3)

	visitCount := 0
	m2.RangeByKey(func(key string, value int) bool {
		visitCount++
		return false // stop after first key
	})

	if visitCount != 1 {
		t.Errorf("expected iteration to stop after 1 key, but visited %d keys", visitCount)
	}
}

func TestSafeBiMap_RangeByValue(t *testing.T) {
	// Create a map with just two entries to simplify testing
	m := ztype.NewSafeBiMap[string, int]()
	m.Set("foo", 1)
	m.Set("bar", 2)

	// Test that all entries are visited when callback always returns true
	visitedAll := make(map[int]bool)
	m.RangeByValue(func(value int, key string) bool {
		visitedAll[value] = true
		return true // continue iteration
	})

	if len(visitedAll) != 2 {
		t.Errorf("expected 2 values to be visited, but got %d", len(visitedAll))
	}

	if !visitedAll[1] || !visitedAll[2] {
		t.Errorf("not all values were visited: %v", visitedAll)
	}

	// Test early termination with a new map
	m2 := ztype.NewSafeBiMap[string, int]()
	m2.Set("key1", 1)
	m2.Set("key2", 2)
	m2.Set("key3", 3)

	visitCount := 0
	m2.RangeByValue(func(value int, key string) bool {
		visitCount++
		return false // stop after first value
	})

	if visitCount != 1 {
		t.Errorf("expected iteration to stop after 1 value, but visited %d values", visitCount)
	}
}

func TestSafeBiMap_Len(t *testing.T) {
	// Create a new SafeBiMap instance
	safeBiMap := ztype.NewSafeBiMap[int, string]()

	// Test empty map
	if safeBiMap.Len() != 0 {
		t.Errorf("Expected length 0, got %d", safeBiMap.Len())
	}

	// Add some values
	safeBiMap.Set(1, "Value 1")
	safeBiMap.Set(2, "Value 2")

	// Test length after adding values
	if safeBiMap.Len() != 2 {
		t.Errorf("Expected length 2, got %d", safeBiMap.Len())
	}

	// Delete a value
	safeBiMap.DeleteKey(1)

	// Test length after deleting a value
	if safeBiMap.Len() != 1 {
		t.Errorf("Expected length 1, got %d", safeBiMap.Len())
	}
}

func TestSafeBiMap_ConcurrentAccess(t *testing.T) {
	// Create a new SafeBiMap instance
	safeBiMap := ztype.NewSafeBiMap[int, string]()

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Number of goroutines
	numGoroutines := 100

	// Add values to the SafeBiMap concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			safeBiMap.Set(key, "value-"+string(rune(key)))
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Verify the map contains the expected number of entries
	if safeBiMap.Len() != numGoroutines {
		t.Errorf("Expected length %d, got %d", numGoroutines, safeBiMap.Len())
	}

	// Use a WaitGroup for concurrent reads
	var rwg sync.WaitGroup

	// Read values from the SafeBiMap concurrently
	for i := 0; i < numGoroutines; i++ {
		rwg.Add(1)
		go func(key int) {
			defer rwg.Done()
			value, _ := safeBiMap.Get(key)
			t.Logf("Key: %d, Value: %s", key, value)
		}(i)
	}

	// Wait for all read goroutines to finish
	rwg.Wait()
}

func BenchmarkSafeBiMap_Get(b *testing.B) {
	m := ztype.NewSafeBiMap[int, string]()
	m.Set(1, "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(1)
	}
}

func BenchmarkSafeBiMap_GetKey(b *testing.B) {
	m := ztype.NewSafeBiMap[int, string]()
	m.Set(1, "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.GetKey("value")
	}
}

func BenchmarkSafeBiMap_Set(b *testing.B) {
	m := ztype.NewSafeBiMap[int, string]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(i, "value")
	}
}

func BenchmarkSafeBiMap_DeleteKey(b *testing.B) {
	// Create a map and add a single key-value pair
	m := ztype.NewSafeBiMap[int, string]()
	m.Set(1, "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Add the key back if it was deleted in the previous iteration
		if i > 0 {
			m.Set(1, "value")
		}
		// Benchmark the DeleteKey operation
		m.DeleteKey(1)
	}
}

func BenchmarkSafeBiMap_DeleteValue(b *testing.B) {
	// Create a map and add a single key-value pair
	m := ztype.NewSafeBiMap[int, string]()
	m.Set(1, "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Add the key back if it was deleted in the previous iteration
		if i > 0 {
			m.Set(1, "value")
		}
		// Benchmark the DeleteValue operation
		m.DeleteValue("value")
	}
}

func BenchmarkSafeBiMap_RangeByKey(b *testing.B) {
	m := ztype.NewSafeBiMap[int, string]()
	for i := 0; i < 100; i++ {
		m.Set(i, "value")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.RangeByKey(func(key int, value string) bool {
			return true
		})
	}
}

func BenchmarkSafeBiMap_RangeByValue(b *testing.B) {
	m := ztype.NewSafeBiMap[int, string]()
	for i := 0; i < 100; i++ {
		m.Set(i, "value"+string(rune(i)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.RangeByValue(func(value string, key int) bool {
			return true
		})
	}
}
