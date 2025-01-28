package database

import (
	"sync"
)

type Index struct {
	mu    sync.RWMutex
	items map[interface{}]map[string]struct{}
}

func NewIndex() *Index {
	return &Index{
		items: make(map[interface{}]map[string]struct{}),
	}
}

func (idx *Index) Add(value interface{}, key string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if idx.items[value] == nil {
		idx.items[value] = make(map[string]struct{})
	}
	idx.items[value][key] = struct{}{}
}

func (idx *Index) Remove(value interface{}, key string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	delete(idx.items[value], key)
	if len(idx.items[value]) == 0 {
		delete(idx.items, value)
	}
}

func (idx *Index) Get(value interface{}) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	keys := make([]string, 0, len(idx.items[value]))
	for key := range idx.items[value] {
		keys = append(keys, key)
	}
	return keys
}
