// entrypoint to start signature aggregator
package sigag

import (
	"context"
	"frost/internal/party/partyclient"
	"frost/internal/sigag/epoch"
	"frost/internal/sigag/rpc"
	"frost/internal/sigag/store"
	"frost/pkg/collections"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type sigag struct {
	logger *logrus.Logger
	port   string
}

func New(opts Options) *sigag {
	return &sigag{
		logger: opts.Logger,
		port:   opts.Port,
	}
}

func (s *sigag) StartSignatureAggregator(ctx context.Context, intialTick time.Duration) error {
	errs, _ := errgroup.WithContext(ctx)

	peerIpList := collections.NewOrderedList[partyclient.PartyClient]()

	store := store.New(peerIpList)

	errs.Go(func() error {
		return rpc.NewServer(store, s.logger).Run(s.port)
	})

	if err := epoch.NewEpochRunner(peerIpList, intialTick, s.logger).Run(); err != nil {
		s.logger.Error("failed while running epoch", zap.Error(err))
		return err
	}

	return errs.Wait()
}
