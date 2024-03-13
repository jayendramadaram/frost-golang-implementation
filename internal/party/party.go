package party

import (
	"context"
	"frost/internal/party/rpc"
	"frost/internal/party/store"
	client "frost/internal/sigag/sigagclient"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func SpinNewParty(port string, ServerUrl string, logger *logrus.Logger) error {

	errs, _ := errgroup.WithContext(context.Background())

	store := store.New()
	SigAgClient := client.New(ServerUrl)

	errs.Go(func() error {
		return rpc.NewServer(store, logger, SigAgClient).Run(port)
	})

	if err := SigAgClient.Register(port, "127.0.0.1", port, "/"); err != nil {
		return err
	}

	return errs.Wait()
}
