package main

import (
	l "log"

	"github.com/deis/logger/log"
	"github.com/deis/logger/storage"
	"github.com/deis/logger/weblog"
)

func main() {
	cfg, err := parseConfig(appName)
	if err != nil {
		l.Fatalf("config error: %s: ", err)
	}

	storageAdapter, err := storage.NewAdapter(cfg.StorageType, cfg.NumLines)
	if err != nil {
		l.Fatal("Error creating storage adapter: ", err)
	}
	storageAdapter.Start()
	defer storageAdapter.Stop()

	aggregator, err := log.NewAggregator(cfg.AggregatorType, storageAdapter)
	if err != nil {
		l.Fatal("Error creating log aggregator: ", err)
	}
	err = aggregator.Listen()
	if err != nil {
		l.Fatal("Error starting log aggregator: ", err)
	}
	l.Println("Log aggregator running")

	weblogServer, err := weblog.NewServer(storageAdapter)
	if err != nil {
		l.Fatal("Error creating weblog server: ", err)
	}
	serverErrCh := weblogServer.Listen()
	l.Println("Weblog server running")

	defer aggregator.Stop()
	stoppedCh := aggregator.Stopped()
	select {
	case stopErr := <-stoppedCh:
		if err != nil {
			l.Fatal("Log aggregator has stopped: ", stopErr)
		} else {
			l.Fatal("Log aggregator has stopped with no error")
		}
	case serverErr := <-serverErrCh:
		l.Fatal("Weblog server failed: ", serverErr)
	}
}
