package main

import (
	"context"
	"frost/internal/sigag"
	"time"

	"go.uber.org/zap"
)

func main() {
	// start signature aggregator
	sigAg := sigag.New(sigag.Options{
		Logger: zap.NewExample(),
		Port:   "8080",
	})
	sigAg.StartSignatureAggregator(context.Background(), 10*time.Second)

	// start nodes
}
