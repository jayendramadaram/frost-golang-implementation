package main

import (
	"context"
	"frost/internal/sigag"
	"time"

	"github.com/sirupsen/logrus"
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
	sigAg.StartSignatureAggregator(context.Background(), 40*time.Second)
}
