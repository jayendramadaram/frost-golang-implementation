package main

import (
	"context"
	"frost/internal/sigag"
	"os"
	"time"

	"github.com/rosedblabs/rosedb/v2"
	"github.com/sirupsen/logrus"
)

func main() {

	options := rosedb.DefaultOptions
	options.DirPath = "D:/codebases/Ozone/frost-golang/cmd/local/sa/tmp/root_sigag"

	os.Remove(options.DirPath)

	// open a database
	db, err := rosedb.Open(options)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = db.Close()
	}()

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
	sigAg.StartSignatureAggregator(context.Background(), 10*time.Second, 100*time.Second, db, 2)
}
