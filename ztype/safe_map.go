package ztype

import (
	"sync"
)

type SafeMap[K comparable, V any] struct {
	sync.RWMutex
	data map[K]V
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		data: make(map[K]V),
	}
}

func (m *SafeMap[K, V]) Get(key K) (V,bool) {
	m.RLock()
	defer m.RUnlock()
	val,ok:= m.data[key]
	return val,ok
}

func (m *SafeMap[K, V]) Set(key K, value V) {
	m.Lock()
	defer m.Unlock()
	m.data[key] = value
}

func (m *SafeMap[K, V]) Delete(key K) {
	m.Lock()
	defer m.Unlock()
	delete(m.data, key)
}

func (m *SafeMap[K, V]) Range(f func(key K, value V) bool) {
	m.RLock()
	defer m.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}
