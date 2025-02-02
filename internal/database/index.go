package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/btree"
)

type Index struct {
	mu       sync.RWMutex
	items    map[interface{}]map[string]struct{}
	btree    *btree.BTree
	dataDir  string
	wal      *WAL
	snapshot *Snapshot
}

// NewIndex initializes the index with WAL and Snapshot support.
func NewIndex(dataDir string) (*Index, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	wal, err := NewWAL(filepath.Join(dataDir, "wal.log"))
	if err != nil {
		return nil, err
	}

	snapshot, err := NewSnapshot(filepath.Join(dataDir, "data.db"))
	if err != nil {
		return nil, err
	}

	idx := &Index{
		items:    make(map[interface{}]map[string]struct{}),
		btree:    btree.New(32),
		dataDir:  dataDir,
		wal:      wal,
		snapshot: snapshot,
	}

	// Recover from WAL
	if err := idx.recover(); err != nil {
		return nil, err
	}

	return idx, nil
}

func (idx *Index) Add(value interface{}, key string) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if err := idx.wal.Write(OpAdd, value, key); err != nil {
		return fmt.Errorf("failed to write to WAL: %w", err)
	}

	if idx.items[value] == nil {
		idx.items[value] = make(map[string]struct{})
	}
	idx.items[value][key] = struct{}{}
	idx.btree.ReplaceOrInsert(&Item{Value: value})

	return nil
}

func (idx *Index) Remove(value interface{}, key string) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if err := idx.wal.Write(OpRemove, value, key); err != nil {
		return fmt.Errorf("failed to write to WAL: %w", err)
	}

	delete(idx.items[value], key)
	if len(idx.items[value]) == 0 {
		delete(idx.items, value)
		idx.btree.Delete(&Item{Value: value})
	}

	return nil
}

func (idx *Index) recover() error {
	entries, err := idx.wal.ReadAll()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		value := string(entry.Value)
		switch entry.Op {
		case OpAdd:
			if idx.items[value] == nil {
				idx.items[value] = make(map[string]struct{})
			}
			idx.items[value][entry.Key] = struct{}{}
		case OpRemove:
			delete(idx.items[value], entry.Key)
			if len(idx.items[value]) == 0 {
				delete(idx.items, value)
			}
		}
	}

	return nil
}

func (idx *Index) Checkpoint() error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if err := idx.snapshot.Write(idx.items); err != nil {
		return err
	}

	return idx.wal.Reset()
}

func (idx *Index) Close() error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if err := idx.Checkpoint(); err != nil {
		return err
	}

	if err := idx.wal.Close(); err != nil {
		return err
	}

	return idx.snapshot.Close()
}
