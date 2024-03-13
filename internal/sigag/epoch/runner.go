package epoch

import (
	"frost/internal/party/partyclient"
	"frost/pkg/collections"
	"time"

	"go.uber.org/zap"
)

type Runner interface {
	Run() error
}

type runner struct {
	nextepoch uint
	initTick  time.Duration
	logger    *zap.Logger

	peerIpList *collections.OrderedList[partyclient.PartyClient]
}

func NewEpochRunner(peerIpList *collections.OrderedList[partyclient.PartyClient], intialTick time.Duration, logger *zap.Logger) Runner {
	return &runner{
		peerIpList: peerIpList,
		nextepoch:  1,
		initTick:   intialTick,
		logger:     logger,
	}
}

func (r *runner) Run() error {
	// Initial Tick [wait for `initTick` for all nodes to be ready and send register message]
	r.awaitInitialTick()

	// in a for loop with a time out of `epoch timeout` start a new epoch

	// unreachable
	return nil
}

func (r *runner) awaitInitialTick() {
	// unblocks time after `initaltick`
	<-time.After(r.initTick)
}
