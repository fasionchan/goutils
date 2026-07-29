package stl

import (
	"context"
	"sync"
)

type SyncMap[Key comparable, Value any] struct {
	m  Mapping[Key, Value]
	createFunc func (context.Context, Key) (Value, error)
	mu sync.RWMutex
}

func NewSyncMap[Key comparable, Value any]() *SyncMap[Key, Value] {
	return &SyncMap[Key, Value]{
		m: make(map[Key]Value),
	}
}

func NewSyncMapPro[Key comparable, Value any](cap int, createFunc func (context.Context, Key) (Value, error)) *SyncMap[Key, Value] {
	if createFunc == nil {
		createFunc = func (context.Context, Key) (_ Value, _ error) {
			return
		}
	}

	return &SyncMap[Key, Value]{
		m: NewMappingWithCap[Key, Value](cap),
		createFunc: createFunc,
	}
}

func (m *SyncMap[Key, Value]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.m.Clear()
}

func (m *SyncMap[Key, Value]) Native() map[Key]Value {
	return m.m
}

func (m *SyncMap[Key, Value]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.m.Len()
}

func (m *SyncMap[Key, Value]) Empty() bool {
	return m.Len() == 0
}

func (m *SyncMap[Key, Value]) Delete(key Key) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.m.Delete(key)
}

func (m *SyncMap[Key, Value]) Get(key Key) Value {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.m.Get(key)
}

func (m *SyncMap[Key, Value]) Keys() []Key {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return MapKeys(m.m)
}

func (m *SyncMap[Key, Value]) Load(key Key) (value Value, loaded bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.m.Load(key)
}

func (m *SyncMap[Key, Value]) LoadOrCreate(ctx context.Context, key Key, create func(ctx context.Context, key Key) (Value, error)) (value Value, loaded bool, err error) {
	value, loaded = m.Load(key)
	if loaded {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	value, loaded = m.m.Load(key)
	if loaded {
		return
	}

	if create == nil {
		create = m.createFunc
	}

	value, err = create(ctx, key)
	if err != nil {
		return
	}

	m.m[key] = value

	return
}

func (m *SyncMap[Key, Value]) LoadOrCreateLite(ctx context.Context, key Key, create func () Value) (Value, bool) {
	value, loaded, _ := m.LoadOrCreate(ctx, key, func (context.Context, Key) (Value, error) {
		return create(), nil
	})
	return value, loaded
}

func (m *SyncMap[Key, Value]) LoadOrStore(key Key, newValue Value) (value Value, loaded bool) {
	value, loaded = m.Load(key)
	if loaded {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	return m.m.LoadOrStore(key, value)
}

func (m *SyncMap[Key, Value]) LoadAndDelete(key Key) (value Value, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.m.LoadAndDelete(key)
}

func (m *SyncMap[Key, Value]) Store(key Key, value Value) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.m.Store(key, value)
}

func (m *SyncMap[Key, Value]) StoreOk(key Key, value Value) (ok bool) {
	if m == nil {
		return false
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	return m.m.StoreOk(key, value)
}

func (m *SyncMap[Key, Value]) Swap(key Key, value Value) (Value, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.m.Swap(key, value)
}

func (m *SyncMap[Key, Value]) SwapOk(key Key, value Value) (old Value, hasOld bool, ok bool) {
	if m == nil {
		ok = false
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	return m.m.SwapOk(key, value)
}

func (m *SyncMap[Key, Value]) Values() []Value {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m.Values()
}

func (m *SyncMap[Key, Value]) SwapMapping(new Mapping[Key, Value]) Mapping[Key, Value] {
	if m == nil {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	old := m.m
	m.m = new

	return old
}