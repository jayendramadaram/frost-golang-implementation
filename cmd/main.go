package main

import (
	"context"
	"fmt"
	"frost/internal/party"
	"frost/internal/sigag"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetLevel(logrus.InfoLevel)
	logger.Formatter = &logrus.TextFormatter{
		DisableColors: false,
		ForceColors:   true,
	}

	// start signature aggregator
	sigAg := sigag.New(sigag.Options{
		Logger: logger,
		Port:   "8080",
	})

	go func() {
		if err := sigAg.StartSignatureAggregator(context.Background(), 40*time.Second); err != nil {
			logger.Error("failed to start signature aggregator", zap.Error(err))
		}
	}()

	// start nodes
	totalNodes := 5
	for i := 1; i <= totalNodes; i++ {
		go func(i int) {
			if err := party.SpinNewParty(fmt.Sprintf("880%d", i), "http://localhost:8080/", logger); err != nil {
				logger.Error("failed to spin new party", zap.Error(err))
			}
		}(i)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
