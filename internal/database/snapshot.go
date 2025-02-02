package database

import (
	"encoding/gob"
	"os"
)

type Snapshot struct {
	file *os.File
}

func NewSnapshot(path string) (*Snapshot, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &Snapshot{file: file}, nil
}

func (s *Snapshot) Write(data map[interface{}]map[string]struct{}) error {
	encoder := gob.NewEncoder(s.file)
	return encoder.Encode(data)
}

func (s *Snapshot) Close() error {
	return s.file.Close()
}
