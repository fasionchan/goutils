package stl

import (
	"context"
	"sync"
)

type SyncMap[Key comparable, Value any] struct {
	mu sync.RWMutex
	m  map[Key]Value
}

func NewSyncMap[Key comparable, Value any]() *SyncMap[Key, Value] {
	return &SyncMap[Key, Value]{
		m: make(map[Key]Value),
	}
}

func (m *SyncMap[Key, Value]) Native() map[Key]Value {
	return m.m
}

func (m *SyncMap[Key, Value]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.m)
}

func (m *SyncMap[Key, Value]) Empty() bool {
	return m.Len() == 0
}

func (m *SyncMap[Key, Value]) Delete(key Key) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.m, key)
}

func (m *SyncMap[Key, Value]) Keys() []Key {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return MapKeys(m.m)
}

func (m *SyncMap[Key, Value]) Load(key Key) (value Value, loaded bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, loaded = m.m[key]
	return
}

func (m *SyncMap[Key, Value]) LoadOrCreate(ctx context.Context, key Key, create func (ctx context.Context) (Value, error)) (value Value, loaded bool, err error) {
	value, loaded = m.Load(key)
	if loaded {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	value, loaded = m.m[key]
	if loaded {
		return
	}

	value, err = create(ctx)
	if err != nil {
		return
	}

	m.m[key] = value

	return
}

func (m *SyncMap[Key, Value]) LoadOrCreateLite(key Key, create func () Value) (value Value, loaded bool) {
	value, loaded = m.Load(key)
	if loaded {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	value, loaded = m.m[key]
	if loaded {
		return
	}

	value = create()
	m.m[key] = value

	return
}

func (m *SyncMap[Key, Value]) LoadOrStore(key Key, newValue Value) (value Value, loaded bool) {
	value, loaded = m.Load(key)
	if loaded {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	value, loaded = m.m[key]
	if loaded {
		return
	}

	value = newValue
	m.m[key] = value

	return
}

func (m *SyncMap[Key, Value]) LoadAndDelete(key Key) (value Value, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, loaded = m.m[key]
	delete(m.m, key)
	return
}

func (m *SyncMap[Key, Value]) Store(key Key, value Value) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[key] = value
}

func (m *SyncMap[Key, Value]) Swap(key Key, value Value) (Value, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	oldValue, loaded := m.m[key]
	m.m[key] = value
	return oldValue, loaded
}

func (m *SyncMap[Key, Value]) Values() []Value {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return MapValues(m.m)
}