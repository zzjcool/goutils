# ztype

`ztype` is a Go package that provides generic type utilities, focusing on thread-safe data structure implementations.

## Features

- Supports different types using Go generics
- All data structures are concurrency-safe
- Clean and easy-to-use API

## Data Structures

### SafeMap

`SafeMap` is a thread-safe map implementation that supports generic key-value pairs.

```go
// Create a new SafeMap
m := ztype.NewSafeMap[string, int]()

// Set a key-value pair
m.Set("key", 123)

// Get a value
value, exists := m.Get("key")

// Delete a key
m.Delete("key")

// Iterate over all key-value pairs
m.Range(func(key string, value int) bool {
    fmt.Printf("Key: %s, Value: %d\n", key, value)
    return true // return true to continue iteration, false to stop
})
```

### SafeBiMap

`SafeBiMap` is a thread-safe bidirectional map implementation that supports lookups from both key to value and value to key.

```go
// Create a new SafeBiMap
bm := ztype.NewSafeBiMap[int, string]()

// Set a key-value pair
bm.Set(1, "one")

// Get a value by key
value, exists := bm.Get(1)

// Get a key by value
key, exists := bm.GetKey("one")

// Delete by key
bm.DeleteKey(1)

// Delete by value
bm.DeleteValue("one")

// Iterate by key
bm.RangeByKey(func(key int, value string) bool {
    fmt.Printf("Key: %d, Value: %s\n", key, value)
    return true // return true to continue iteration, false to stop
})

// Iterate by value
bm.RangeByValue(func(value string, key int) bool {
    fmt.Printf("Value: %s, Key: %d\n", value, key)
    return true // return true to continue iteration, false to stop
})

// Get the number of key-value pairs in the map
count := bm.Len()
```

## Thread Safety

All data structures use `sync.RWMutex` to implement read-write locks, ensuring safe access in concurrent environments. This makes these data structures particularly suitable for use in multi-goroutine environments.

## Performance Considerations

- Read operations (such as `Get`) use read locks, allowing multiple concurrent reads
- Write operations (such as `Set`, `Delete`) use write locks, ensuring exclusive access
- For `SafeBiMap`, maintaining bidirectional mappings incurs a slight performance overhead but provides the convenience of bidirectional lookups

## Use Cases

- Concurrent applications requiring thread-safe map data structures
- Applications needing bidirectional lookup functionality (using `SafeBiMap`)
- Applications that need to share data between multiple goroutines
