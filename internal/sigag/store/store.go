package store

import (
	"fmt"
	"frost/internal/party/partyclient"
	"frost/internal/sigag/rpc"
	"frost/pkg/collections"
)

type store struct {
	peerIpList *collections.OrderedList[partyclient.PartyClient]
}

// AddParticipant implements rpc.Store.
func (s *store) AddParticipant(party rpc.RegisterParty) error {
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
func (s *store) GetParties() (rpc.Parties, error) {
	Parties := make(rpc.Parties)
	for _, v := range s.peerIpList.Items {
		id, url := v.Locate()
		Parties[id] = url
	}

	return Parties, nil
}

type Store interface {
	rpc.Store
}

func New(peerIpList *collections.OrderedList[partyclient.PartyClient]) Store {
	return &store{peerIpList: peerIpList}
}

// type Server interface {
// 	Run(port string) error
// 	Register(_ context.Context, params *json.RawMessage) (json.RawMessage, error)
// 	Health(_ context.Context, params *json.RawMessage) (json.RawMessage, error)
// 	GetParties(_ context.Context, params *json.RawMessage) (json.RawMessage, error)
// }
