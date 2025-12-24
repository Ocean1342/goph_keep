package crypter

import (
	"errors"
	"sync"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrEmptyLogin  = errors.New("empty login")
)

type storage struct {
	storage map[string][]byte
	mu      sync.Mutex
}

func (s *storage) Set(login string, key []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if key != nil {
		copiedKey := make([]byte, len(key))
		copy(copiedKey, key)
		s.storage[login] = copiedKey
	} else {
		s.storage[login] = nil
	}
}

func (s *storage) Get(login string) ([]byte, error) {
	if login == "" {
		return nil, ErrEmptyLogin
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.storage[login]; ok {
		return v, nil
	}
	return nil, ErrKeyNotFound
}
