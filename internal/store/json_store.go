package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type JSONStore struct {
	mu   sync.RWMutex
	path string
}

func NewJSONStore(path string) (*JSONStore, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, errors.New("data file path cannot be empty")
	}

	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create data directory: %w", err)
		}
	}

	s := &JSONStore{path: path}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		if err := s.writeFile([]byte("[]")); err != nil {
			return nil, fmt.Errorf("initialize data file: %w", err)
		}
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read data file: %w", err)
	}

	if strings.TrimSpace(string(data)) == "" {
		if err := s.writeFile([]byte("[]")); err != nil {
			return nil, fmt.Errorf("initialize empty data file: %w", err)
		}
		return s, nil
	}

	var array []any
	if err := json.Unmarshal(data, &array); err != nil {
		return nil, fmt.Errorf("data file must contain a valid JSON array: %w", err)
	}

	return s, nil
}

func (s *JSONStore) Read(dst any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		return fmt.Errorf("read data file: %w", err)
	}
	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("decode data file: %w", err)
	}
	return nil
}

// Update performs a read-modify-write operation under one lock so concurrent
// requests cannot lose writes.
func (s *JSONStore) Update(dst any, mutate func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		data = []byte("[]")
	} else if err != nil {
		return fmt.Errorf("read data file: %w", err)
	}

	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("decode data file: %w", err)
	}
	if err := mutate(); err != nil {
		return err
	}

	encoded, err := json.MarshalIndent(dst, "", "  ")
	if err != nil {
		return fmt.Errorf("encode data file: %w", err)
	}
	return s.writeFile(encoded)
}

func (s *JSONStore) writeFile(data []byte) error {
	dir := filepath.Dir(s.path)
	tmp, err := os.CreateTemp(dir, ".bills-*.json")
	if err != nil {
		return fmt.Errorf("create temporary data file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if err := tmp.Chmod(0o644); err != nil {
		tmp.Close()
		return fmt.Errorf("set data file permissions: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("write temporary data file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temporary data file: %w", err)
	}
	if err := os.Rename(tmpName, s.path); err != nil {
		return fmt.Errorf("replace data file: %w", err)
	}
	return nil
}
