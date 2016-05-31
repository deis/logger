package consumer

import (
	"fmt"
	"time"
)

// ErrStopTimedOut is the error returned if a (Consumer).Stop call times out before the stop is complete
type ErrStopTimedOut struct {
	Timeout time.Duration
}

func (e ErrStopTimedOut) Error() string {
	return fmt.Sprintf("stopping a consumer timed out after %s", e.Timeout)
}

// Consumer is the high level interface which can be implemented to consume messages from a source
type Consumer interface {
	// Stop stops the consumer and blocks until it stops or the given duration passes. In the latter case, returns an error of type ErrStopTimedOut
	Stop(time.Duration) error
	// Stopped returns a channel that will receive when the consumer has stopped. The error received may be nil, which means it was stopped cleanly. A non-nil error means it stopped because it errored.
	Stopped() <-chan error
}
