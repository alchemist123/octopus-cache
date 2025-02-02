package database

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

type Operation uint8

const (
	OpAdd Operation = iota
	OpRemove
)

type WALEntry struct {
	Timestamp int64
	Op        Operation
	Value     []byte
	Key       string
}

type WAL struct {
	file *os.File
}

func NewWAL(path string) (*WAL, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open WAL file: %w", err)
	}
	return &WAL{file: file}, nil
}

func (wal *WAL) Write(op Operation, value interface{}, key string) error {
	entry := WALEntry{
		Timestamp: time.Now().UnixNano(),
		Op:        op,
		Value:     []byte(fmt.Sprintf("%v", value)),
		Key:       key,
	}

	size := uint32(8 + 1 + len(entry.Value) + len(entry.Key))
	if err := binary.Write(wal.file, binary.LittleEndian, size); err != nil {
		return err
	}
	if err := binary.Write(wal.file, binary.LittleEndian, entry.Timestamp); err != nil {
		return err
	}
	if err := binary.Write(wal.file, binary.LittleEndian, entry.Op); err != nil {
		return err
	}
	if _, err := wal.file.Write(entry.Value); err != nil {
		return err
	}
	if _, err := wal.file.WriteString(entry.Key); err != nil {
		return err
	}

	return wal.file.Sync()
}

func (wal *WAL) ReadAll() ([]WALEntry, error) {
	var entries []WALEntry
	var size uint32
	for {
		if err := binary.Read(wal.file, binary.LittleEndian, &size); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}

		var entry WALEntry
		if err := binary.Read(wal.file, binary.LittleEndian, &entry.Timestamp); err != nil {
			return nil, err
		}
		if err := binary.Read(wal.file, binary.LittleEndian, &entry.Op); err != nil {
			return nil, err
		}

		valueBytes := make([]byte, size-9)
		if _, err := wal.file.Read(valueBytes); err != nil {
			return nil, err
		}
		entry.Value = valueBytes

		entries = append(entries, entry)
	}

	return entries, nil
}

func (wal *WAL) Reset() error {
	if err := wal.file.Truncate(0); err != nil {
		return err
	}
	_, err := wal.file.Seek(0, 0)
	return err
}

func (wal *WAL) Close() error {
	return wal.file.Close()
}
