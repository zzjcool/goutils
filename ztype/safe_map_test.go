package ztype_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/zzjcool/goutils/ztype"
)

func TestNewSafeMap(t *testing.T) {
	t.Run("Test creating a new SafeMap with integer keys and string values", func(t *testing.T) {
		intStrMap := ztype.NewSafeMap[int, string]()
		if reflect.TypeOf(intStrMap) != reflect.TypeOf(&ztype.SafeMap[int, string]{}) {
			t.Errorf("Expected SafeMap[int, string], got %T", intStrMap)
		}
	})

	t.Run("Test creating a new SafeMap with string keys and integer values", func(t *testing.T) {
		strIntMap := ztype.NewSafeMap[string, int]()
		if reflect.TypeOf(strIntMap) != reflect.TypeOf(&ztype.SafeMap[string, int]{}) {
			t.Errorf("Expected SafeMap[string, int], got %T", strIntMap)
		}
	})

	t.Run("Test creating a new SafeMap with custom struct keys and float values", func(t *testing.T) {
		type CustomStruct struct {
			Field1 int
			Field2 string
		}
		structFloatMap := ztype.NewSafeMap[CustomStruct, float64]()
		if reflect.TypeOf(structFloatMap) != reflect.TypeOf(&ztype.SafeMap[CustomStruct, float64]{}) {
			t.Errorf("Expected SafeMap[CustomStruct, float64], got %T", structFloatMap)
		}
	})
}

func TestSafeMap_Get(t *testing.T) {
	// Create a new SafeMap instance
	safeMap := ztype.NewSafeMap[int, string]()

	// Add some values to the SafeMap
	safeMap.Set(1, "Value 1")
	safeMap.Set(2, "Value 2")

	// Test the Get method
	value, _ := safeMap.Get(1)
	if value != "Value 1" {
		t.Errorf("Expected value 'Value 1', got '%s'", value)
	}
}

func TestSafeMap_Set(t *testing.T) {
	// Create a new SafeMap instance
	safeMap := ztype.NewSafeMap[int, string]()

	// Test the Set method
	safeMap.Set(1, "Value 1")
	value, _ := safeMap.Get(1)
	if value != "Value 1" {
		t.Errorf("Expected value 'Value 1', got '%s'", value)
	}
}

func TestSafeMap_Delete(t *testing.T) {
	// Create a new SafeMap instance
	safeMap := ztype.NewSafeMap[int, string]()

	// Add some values to the SafeMap
	safeMap.Set(1, "Value 1")
	safeMap.Set(2, "Value 2")

	// Delete a value from the SafeMap
	safeMap.Delete(1)

	// Test the Delete method
	_, ok := safeMap.Get(1)
	if ok {
		t.Error("Expected key 1 to be deleted, but it still exists")
	}
}



func TestSafeMap_Range(t *testing.T) {
	// Create a map with just two entries to simplify testing
	m := ztype.NewSafeMap[string, int]()
	m.Set("foo", 1)
	m.Set("bar", 2)

	// Test that all entries are visited when callback always returns true
	visitedAll := make(map[string]bool)
	m.Range(func(key string, value int) bool {
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
	m2 := ztype.NewSafeMap[string, int]()
	m2.Set("key1", 1)
	m2.Set("key2", 2)
	m2.Set("key3", 3)

	visitCount := 0
	m2.Range(func(key string, value int) bool {
		visitCount++
		return false // stop after first key
	})

	if visitCount != 1 {
		t.Errorf("expected iteration to stop after 1 key, but visited %d keys", visitCount)
	}
}

func TestSafeMap_ConcurrentAccess(t *testing.T) {
	// Create a new SafeMap instance
	safeMap := ztype.NewSafeMap[int, int]()

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Number of goroutines
	numGoroutines := 100

	// Add 10 values to the SafeMap concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			safeMap.Set(key, key*2)
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Use a WaitGroup to wait for all goroutines to finish
	var rwg sync.WaitGroup

	// Number of goroutines
	numReadGoroutines := 50

	// Read values from the SafeMap concurrently
	for i := 0; i < numReadGoroutines; i++ {
		rwg.Add(1)
		go func(key int) {
			defer rwg.Done()
			value, _ := safeMap.Get(key)
			t.Logf("Key: %d, Value: %d", key, value)
		}(i)
	}

	// Wait for all goroutines to finish
	rwg.Wait()
}



func BenchmarkSafeMap_Get(b *testing.B) {
    m := ztype.NewSafeMap[int, string]()
    m.Set(1, "value")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        m.Get(1)
    }
}

func BenchmarkSafeMap_Set(b *testing.B) {
    m := ztype.NewSafeMap[int, string]()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        m.Set(i, "value")
    }
}

func BenchmarkSafeMap_Delete(b *testing.B) {
    m := ztype.NewSafeMap[int, string]()
    m.Set(1, "value")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        m.Delete(1)
    }
}

func BenchmarkSafeMap_Range(b *testing.B) {
    m := ztype.NewSafeMap[int, string]()
    m.Set(1, "value")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        m.Range(func(key int, value string) bool {
            return true
        })
    }
}