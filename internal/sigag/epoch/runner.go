package epoch

import (
	"frost/internal/party/partyclient"
	"frost/internal/sigag/rpc"
	"frost/pkg/collections"
	"time"

	"github.com/sirupsen/logrus"
)

type Runner interface {
	Run(time.Duration) error
}

type runner struct {
	nextepoch uint
	initTick  time.Duration
	logger    *logrus.Logger

	store           Store
	thresholdFactor float64
}

type Store interface {
	Lock()
	UnLock()

	PutParties(parties rpc.Parties) error
	GetPartyCLients() *collections.OrderedList[partyclient.PartyClient]
	RemoveParty(item partyclient.PartyClient) error

	PutThreshold(threshold uint, epoch uint) error
}

func NewEpochRunner(store Store, intialTick time.Duration, thresholdFactor float64, logger *logrus.Logger) Runner {
	return &runner{
		store:     store,
		nextepoch: 1,
		initTick:  intialTick,
		logger:    logger,

		thresholdFactor: thresholdFactor,
	}
}

func (r *runner) Run(epochDuration time.Duration) error {
	// Initial Tick [wait for `initTick` for all nodes to be ready and send register message]
	r.awaitInitialTick()

	// in a for loop with a time out of `epoch timeout` start a new epoch
	for {
		// lock system sigag ✅
		// announce new epoch ✅
		// ack from all parties and they lock their systems ✅
		// store responded parties to the db for this epoch ✅

		// DKG init
		// choose n,k ✅
		// Round 1 ✅
		// send list of parties[] to all parties[] ✅
		// parties do dkg among themselves
		// - each party generates a polynomial using shamir secret sharing library
		// - compute POK of secret and commitments for polynomial generated
		// - each party broadcasts commitments and POC to all other parties O(n^2) network calls
		// - - if failed try again retry(with backoff)
		// - if everyone has acquired N commitments then we can start key gen else perform Round(1) of DKG again
		// SA receives Round 1 ACK

		// Round 2
		// - each participant sends secret share for all N participants O(n^2) network calls and get verified accordingly
		// - participates calculate long lives secrets from
		// finally SA can calculate and publish Group Pubkey and Individual Party Pubkeys

		// preprocess
		// SA requests for nonces from all parties for next 10 txs
		// stores nonces in local db

		// chose a subset
		// send tx to choosen set
		// aggregate sigs
		r.store.Lock()
		parties := r.store.GetPartyCLients()
		partyMap, err := r.AnnounceNewEpoch(parties, r.nextepoch)
		if err != nil {
			return err
		}

		if err := r.store.PutParties(partyMap); err != nil {
			return err
		}

		TotalParties := len(partyMap)
		Threshold := uint((float64(TotalParties) / r.thresholdFactor) + 1)

		if err := r.store.PutThreshold(Threshold, r.nextepoch); err != nil {
			return err
		}

		if err := r.AnnounceDKGInit(r.store.GetPartyCLients(), partyMap, Threshold); err != nil {
			r.logger.Errorf("failed to announce dkg init: %v", err)
			continue
		}

		time.Sleep(epochDuration)
		r.nextepoch++
	}

	// unreachable
	return nil
}

func (r *runner) awaitInitialTick() {
	// unblocks time after `initaltick`
	<-time.After(r.initTick)
}

func (r *runner) AnnounceNewEpoch(parties *collections.OrderedList[partyclient.PartyClient], epoch uint) (rpc.Parties, error) {
	partyMap := make(rpc.Parties)
	for _, v := range parties.Items {
		if err := v.NewEpoch(epoch); err != nil {
			r.logger.Errorf("failed to announce new epoch: %v", err)
			if err := r.store.RemoveParty(v); err != nil {
				r.logger.Errorf("failed to remove party: %v", err)
				return nil, err
			}
			continue
		}
		id, url := v.Locate()
		partyMap[id] = url
	}
	return partyMap, nil
}

func (r *runner) AnnounceDKGInit(parties *collections.OrderedList[partyclient.PartyClient], partyMap rpc.Parties, threshold uint) error {
	for _, v := range parties.Items {
		if err := v.DKGInit(partyMap, threshold); err != nil {
			r.logger.Errorf("failed to announce dkg init: %v", err)
			return err
		}
	}
	return nil
}
