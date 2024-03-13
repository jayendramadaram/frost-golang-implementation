package epoch

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Runner interface {
	Run() error
}

type runner struct {
	nextepoch uint
	initTick  time.Duration
	logger    *logrus.Logger

	store Store
}

type Store interface {
}

func NewEpochRunner(store Store, intialTick time.Duration, logger *logrus.Logger) Runner {
	return &runner{
		store:     store,
		nextepoch: 1,
		initTick:  intialTick,
		logger:    logger,
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
