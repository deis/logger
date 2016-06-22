package consumer

import (
	"fmt"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

type errNSQHandleFailed struct {
	err error
}

func (e errNSQHandleFailed) Error() string {
	return fmt.Sprintf("handling the NSQ message failed with %s", e.err)
}

type nsqConsumer struct {
	consumer *nsq.Consumer
}

// Stop is the consumer interface implementation
func (n *nsqConsumer) Stop(timeout time.Duration) error {
	n.consumer.Stop()
	tmr := time.NewTimer(timeout)
	defer tmr.Stop()
	select {
	case <-tmr.C:
		return ErrStopTimedOut{Timeout: timeout}
	case <-n.consumer.StopChan:
		return nil
	}
}

// Stopped is the consumer interface implementation
func (n *nsqConsumer) Stopped() <-chan error {
	retCh := make(chan error)
	go func() {
		<-n.consumer.StopChan
		retCh <- nil
	}()
	return retCh
}

func nsqHandlerFromMessageHandler(hdl MessageHandler) nsq.Handler {
	hdlFunc := func(msg *nsq.Message) error {
		if err := hdl.Handle(&Message{Bytes: msg.Body}); err != nil {
			msg.Requeue(-1)
			return errNSQHandleFailed{err: err}
		}
		return nil
	}
	return nsq.HandlerFunc(hdlFunc)
}

// NewNSQConsumer opens a new NSQ connection to addr and registers hdl as a consumer on the given topic or channel. Returns a new Consumer representing the ongoing NSQ connection and consumption, or a non-nil error if the consumption couldn't be set up for any reason. It is the caller's responsibility to call Stop() on the returned consumer.
func NewNSQConsumer(nsqHostStr, topic, channel string, numThreads int, hdl MessageHandler) (Consumer, error) {
	cfg := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		return nil, err
	}
	consumer.AddConcurrentHandlers(nsqHandlerFromMessageHandler(hdl), numThreads)

	if err := consumer.ConnectToNSQD(nsqHostStr); err != nil {
		return nil, err
	}
	return &nsqConsumer{consumer: consumer}, nil
}
