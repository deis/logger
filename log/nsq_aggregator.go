package log

import (
	"fmt"
	l "log"
	"time"

	nsq "github.com/nsqio/go-nsq"

	"github.com/deis/logger/storage"
)

type nsqAggregator struct {
	listening bool
	cfg       *config
	consumer  *nsq.Consumer
	handler   nsq.HandlerFunc
}

func newNSQAggregator(storageAdapter storage.Adapter) Aggregator {
	return &nsqAggregator{
		handler: nsq.HandlerFunc(func(msg *nsq.Message) error {
			if err := handle(msg.Body, storageAdapter); err != nil {
				msg.Requeue(-1)
				return newErrNSQHandleFailed(err)
			}
			return nil
		}),
	}
}

// Listen starts the aggregator. Invocations of this function are not concurrency safe and multiple
// serialized invocations have no effect.
func (a *nsqAggregator) Listen() error {
	// Should only ever be called once
	if !a.listening {
		a.listening = true
		var err error
		a.cfg, err = parseConfig(appName)
		if err != nil {
			l.Fatalf("config error: %s: ", err)
		}
		config := nsq.NewConfig()
		consumer, err := nsq.NewConsumer(a.cfg.NSQTopic, a.cfg.NSQChannel, config)
		if err != nil {
			return err
		}
		consumer.AddConcurrentHandlers(a.handler, a.cfg.NSQHandlerCount)
		if err := consumer.ConnectToNSQD(a.cfg.nsqURL()); err != nil {
			return err
		}
		a.consumer = consumer
	}
	return nil
}

// Stop is the Aggregator interface implementation
func (a *nsqAggregator) Stop() error {
	a.consumer.Stop()
	timeout := a.cfg.stopTimeoutDuration()
	tmr := time.NewTimer(timeout)
	defer tmr.Stop()
	select {
	case <-tmr.C:
		return newErrStopTimedOut(timeout)
	case <-a.consumer.StopChan:
		return nil
	}
}

// Stopped is the Aggregator interface implementation
func (a *nsqAggregator) Stopped() <-chan error {
	retCh := make(chan error)
	go func() {
		<-a.consumer.StopChan
		retCh <- nil
	}()
	return retCh
}

type errNSQHandleFailed struct {
	err error
}

func newErrNSQHandleFailed(err error) errNSQHandleFailed {
	return errNSQHandleFailed{err: err}
}

func (e errNSQHandleFailed) Error() string {
	return fmt.Sprintf("handling the NSQ message failed with %s", e.err)
}
