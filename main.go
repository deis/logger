package main

import (
	"log"
	"time"

	"github.com/deis/logger/consumer"
	"github.com/deis/logger/logs"
	"github.com/deis/logger/storage"
	"github.com/deis/logger/weblog"
)

func main() {
	cfg, err := parseConfig(appName)
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	nsqConsumerStopDur := time.Duration(cfg.NSQConsumerStopDurSec) * time.Second

	storageAdapter, err := storage.NewAdapter(cfg.StorageType, cfg.NumLines)
	if err != nil {
		log.Fatal("Error creating storage adapter:", err)
	}
	logger, err := logs.NewLogger(storageAdapter)
	if err != nil {
		log.Fatal("Error creating logger", err)
	}

	weblogServer, err := weblog.NewServer(logger)
	if err != nil {
		log.Fatal("Error creating weblog server", err)
	}

	serverErrCh := weblogServer.Listen()
	log.Println("deis-logger running")

	log.Printf("Listening to NSQ on %s", cfg.nsqURL())
	consumer, err := consumer.NewNSQConsumer(
		cfg.nsqURL(),
		cfg.NSQTopic,
		cfg.NSQChannel,
		cfg.NSQConsumerNumThreads,
		nsqMsgHandler(logger),
	)
	if err != nil {
		log.Fatalf("Error creating new NSQ consumer (%s)", err)
	}
	defer consumer.Stop(nsqConsumerStopDur)
	// a channel that never receives, so that we wait either forever or for the consumer to stop
	alwaysCh := make(chan struct{})
	stoppedCh := consumer.Stopped()
	select {
	case stopErr := <-stoppedCh:
		if err != nil {
			log.Fatalf("NSQ consumer has stopped (%s)", stopErr)
		} else {
			log.Fatalf("NSQ consumer has stopped with no error")
		}
	case serverErr := <-serverErrCh:
		log.Fatalf("logs HTTP server failed (%s)", serverErr)
	case <-alwaysCh:
	}
}

func nsqMsgHandler(logger *logs.Logger) consumer.MessageHandler {
	return consumer.MessageHandlerFunc(func(msg *consumer.Message) error {
		if err := logger.WriteLog(msg.Bytes); err != nil {
			log.Printf("Unable to store message '%s' (%s)", string(msg.Bytes), err)
			return err
		}
		return nil
	})
}
