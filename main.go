package main

import (
	l "log"
	"net/http"

	_ "net/http/pprof"

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
	defer aggregator.Stop()
	l.Println("Log aggregator running")

	weblogServer := weblog.NewServer(storageAdapter)
	weblogServer.Start()
	defer weblogServer.Close()
	l.Printf("Weblog server serving at %s\n", weblogServer.URL)

	// start a Go Profiler
	go func() {
		l.Println(http.ListenAndServe("0.0.0.0:8099", nil))
	}()

	stoppedCh := aggregator.Stopped()
	select {
	case stopErr := <-stoppedCh:
		if err != nil {
			l.Fatal("Log aggregator has stopped: ", stopErr)
		} else {
			l.Fatal("Log aggregator has stopped with no error")
		}
	}
}
