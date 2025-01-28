package database

import (
	"container/heap"
	"sync"
	"time"
)

type Database struct {
	data       *sync.Map
	indexes    map[string]*Index
	expiryHeap *ExpirationHeap
	mu         sync.RWMutex
}

type Item struct {
	Value   interface{}
	Expiry  time.Time
	Indexes map[string]interface{}
}

func NewDatabase() *Database {
	db := &Database{
		data:       &sync.Map{},
		indexes:    make(map[string]*Index),
		expiryHeap: NewExpirationHeap(),
	}
	go db.expirationWorker()
	return db
}

func (db *Database) Set(key string, value interface{}, ttl time.Duration, indexes map[string]interface{}) {
	item := &Item{
		Value:   value,
		Expiry:  time.Now().Add(ttl),
		Indexes: indexes,
	}

	db.data.Store(key, item)
	db.expiryHeap.Push(&HeapItem{
		Key:    key,
		Expiry: item.Expiry,
	})
	db.updateIndexes(key, indexes)
}

func (db *Database) Get(key string) (interface{}, bool) {
	val, ok := db.data.Load(key)
	if !ok {
		return nil, false
	}

	item := val.(*Item)
	if time.Now().After(item.Expiry) {
		db.Delete(key)
		return nil, false
	}
	return item.Value, true
}

func (db *Database) Delete(key string) {
	val, ok := db.data.LoadAndDelete(key)
	if !ok {
		return
	}

	item := val.(*Item)
	db.removeFromIndexes(key, item.Indexes)
}

func (db *Database) Query(indexName string, value interface{}) []string {
	db.mu.RLock()
	index, exists := db.indexes[indexName]
	db.mu.RUnlock()

	if !exists {
		return nil
	}

	index.mu.RLock()
	defer index.mu.RUnlock()

	keys := make([]string, 0, len(index.items[value]))
	for key := range index.items[value] {
		keys = append(keys, key)
	}
	return keys
}

func (db *Database) updateIndexes(key string, indexes map[string]interface{}) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for field, value := range indexes {
		if _, exists := db.indexes[field]; !exists {
			db.indexes[field] = NewIndex()
		}

		index := db.indexes[field]
		index.Add(value, key)
	}
}

func (db *Database) removeFromIndexes(key string, indexes map[string]interface{}) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for field, value := range indexes {
		if index, exists := db.indexes[field]; exists {
			index.Remove(value, key)
		}
	}
}

func (db *Database) expirationWorker() {
	for {
		time.Sleep(1 * time.Second)

		now := time.Now()
		db.mu.Lock()
		for db.expiryHeap.Len() > 0 {
			item := (*db.expiryHeap)[0]
			if item.Expiry.After(now) {
				break
			}
			heap.Pop(db.expiryHeap)
			db.Delete(item.Key)
		}
		db.mu.Unlock()
	}
}
