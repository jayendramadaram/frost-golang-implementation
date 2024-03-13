package store

import (
	"fmt"
	"frost/internal/party/rpc"
	"sync"
)

type store struct {
	mu           sync.RWMutex
	locked       bool
	currentEpoch uint
}

type Store interface {
	rpc.Store
}

func New() Store {
	return &store{
		mu: sync.RWMutex{},
	}
}

// IsLocked implements Store.
func (s *store) IsLocked() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.locked
}

// Lock implements Store.
func (s *store) Lock() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.locked = true
}

// UnLock implements Store.
func (s *store) UnLock() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.locked = false
}

func (s *store) NewEpoch(epoch uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if epoch <= s.currentEpoch {
		return fmt.Errorf("recevied invalid epoch %s", epoch)
	}
	s.currentEpoch = epoch
	return nil
}
