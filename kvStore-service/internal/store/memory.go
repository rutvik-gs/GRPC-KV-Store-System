package store

import (
	"sync"
)

type InMemoryStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func CreateStore() Store {
	return &InMemoryStore{
		data: make(map[string]string),
	}
}

func (i *InMemoryStore) Set(key, value string) error {
	if key == "" {
		return ErrEmptyKey
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	i.data[key] = value
	return nil
}

func (i *InMemoryStore) Get(key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}

	i.mu.RLock()
	defer i.mu.RUnlock()

	value, exists := i.data[key]
	if !exists {
		return "", ErrKeyNotFound
	}

	return value, nil
}

func (i *InMemoryStore) Delete(key string) error {
	if key == "" {
		return ErrEmptyKey
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	if _, exists := i.data[key]; !exists {
		return ErrKeyNotFound
	}

	delete(i.data, key)
	return nil
}
