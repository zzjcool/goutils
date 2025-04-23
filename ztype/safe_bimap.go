package ztype

import (
	"sync"
)

// SafeBiMap is a thread-safe bidirectional map that allows lookups from both key->value and value->key.
// It uses generics to support different types for keys and values.
type SafeBiMap[K, V comparable] struct {
	sync.RWMutex
	forward map[K]V // Key to Value mapping
	reverse map[V]K // Value to Key mapping
}

// NewSafeBiMap creates a new instance of SafeBiMap.
func NewSafeBiMap[K, V comparable]() *SafeBiMap[K, V] {
	return &SafeBiMap[K, V]{
		forward: make(map[K]V),
		reverse: make(map[V]K),
	}
}

// Get retrieves a value by its key.
func (m *SafeBiMap[K, V]) Get(key K) (V, bool) {
	m.RLock()
	defer m.RUnlock()
	val, ok := m.forward[key]
	return val, ok
}

// GetKey retrieves a key by its value.
func (m *SafeBiMap[K, V]) GetKey(value V) (K, bool) {
	m.RLock()
	defer m.RUnlock()
	key, ok := m.reverse[value]
	return key, ok
}

// Set adds or updates a key-value pair.
// If the key already exists, the old value is removed from the reverse mapping.
// If the value already exists, the old key is removed from the forward mapping.
func (m *SafeBiMap[K, V]) Set(key K, value V) {
	m.Lock()
	defer m.Unlock()

	// If the key already exists, remove the old value from reverse map
	if oldValue, exists := m.forward[key]; exists {
		delete(m.reverse, oldValue)
	}

	// If the value already exists, remove the old key from forward map
	if oldKey, exists := m.reverse[value]; exists {
		delete(m.forward, oldKey)
	}

	// Update both maps
	m.forward[key] = value
	m.reverse[value] = key
}

// DeleteKey removes a key-value pair by key.
func (m *SafeBiMap[K, V]) DeleteKey(key K) {
	m.Lock()
	defer m.Unlock()
	
	// If key exists, remove the corresponding value from reverse map
	if value, exists := m.forward[key]; exists {
		delete(m.reverse, value)
		delete(m.forward, key)
	}
}

// DeleteValue removes a key-value pair by value.
func (m *SafeBiMap[K, V]) DeleteValue(value V) {
	m.Lock()
	defer m.Unlock()
	
	// If value exists, remove the corresponding key from forward map
	if key, exists := m.reverse[value]; exists {
		delete(m.forward, key)
		delete(m.reverse, value)
	}
}

// RangeByKey iterates through all key-value pairs.
// The iteration stops if f returns false.
func (m *SafeBiMap[K, V]) RangeByKey(f func(key K, value V) bool) {
	m.RLock()
	defer m.RUnlock()
	for k, v := range m.forward {
		if !f(k, v) {
			break
		}
	}
}

// RangeByValue iterates through all value-key pairs.
// The iteration stops if f returns false.
func (m *SafeBiMap[K, V]) RangeByValue(f func(value V, key K) bool) {
	m.RLock()
	defer m.RUnlock()
	for v, k := range m.reverse {
		if !f(v, k) {
			break
		}
	}
}

// Len returns the number of key-value pairs in the map.
func (m *SafeBiMap[K, V]) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.forward)
}
