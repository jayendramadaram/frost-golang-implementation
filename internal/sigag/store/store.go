package store

import (
	"fmt"
	"frost/internal/party/partyclient"
	"frost/internal/sigag/epoch"
	"frost/internal/sigag/rpc"
	"frost/pkg/collections"
	"sync"

	"github.com/rosedblabs/rosedb/v2"
)

type store struct {
	peerIpList *collections.OrderedList[partyclient.PartyClient]
	locked     bool
	mu         sync.RWMutex
	db         *rosedb.DB
}

// GetEpochParties implements Store.
func (s *store) GetEpochParties() rpc.Parties {
	s.mu.RLock()
	defer s.mu.RUnlock()

	Parties := make(rpc.Parties)
	s.db.Ascend(func(k []byte, v []byte) (bool, error) {
		Parties[string(k)] = string(v)
		return true, nil
	})

	return Parties
}

// PutParties implements Store.
func (s *store) PutParties(parties rpc.Parties) error {
	batch := s.db.NewBatch(rosedb.DefaultBatchOptions)

	for id, url := range parties {
		if err := batch.Put([]byte(id), []byte(url)); err != nil {
			return err
		}
	}

	return batch.Commit()
}

// GetPartyCLients implements Store.
func (s *store) GetPartyCLients() collections.OrderedList[partyclient.PartyClient] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return *s.peerIpList
}

// IsLocked implements Store.
func (s *store) IsLocked() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.locked
}

// Lock implements epoch.Store.
func (s *store) Lock() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.locked = true
}

// UnLock implements epoch.Store.
func (s *store) UnLock() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.locked = false
}

// AddParticipant implements rpc.Store.
func (s *store) AddParticipant(party rpc.RegisterParty) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	participant := partyclient.New(party.Address, party.ReportedIp, party.Port, party.Path)

	containsID := func(item, element partyclient.PartyClient) bool {
		return item.ID() == element.ID()
	}

	if s.peerIpList.Contains(participant, containsID) {
		return fmt.Errorf("address already registered")
	}

	if err := participant.Ping(); err != nil {
		return err
	}

	s.peerIpList.Add(participant)
	return nil
}

// GetParties implements rpc.Store.
func (s *store) GetParties() rpc.Parties {
	s.mu.RLock()
	defer s.mu.RUnlock()

	Parties := make(rpc.Parties)
	for _, v := range s.peerIpList.Items {
		id, url := v.Locate()
		Parties[id] = url
	}

	return Parties
}

type Store interface {
	rpc.Store
	epoch.Store
}

func New(peerIpList *collections.OrderedList[partyclient.PartyClient], db *rosedb.DB) Store {
	return &store{
		peerIpList: peerIpList,
		locked:     false,
		mu:         sync.RWMutex{},
		db:         db,
	}
}

// type Server interface {
// 	Run(port string) error
// 	Register(_ context.Context, params *json.RawMessage) (json.RawMessage, error)
// 	Health(_ context.Context, params *json.RawMessage) (json.RawMessage, error)
// 	GetParties(_ context.Context, params *json.RawMessage) (json.RawMessage, error)
// }
