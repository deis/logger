package log

import (
	"fmt"
	"time"
)

// ErrStopTimedOut is the error returned if a (Aggregator).Stop call times out before the stop is
// complete
type ErrStopTimedOut struct {
	Timeout time.Duration
}

func newErrStopTimedOut(timeout time.Duration) ErrStopTimedOut {
	return ErrStopTimedOut{Timeout: timeout}
}

func (e ErrStopTimedOut) Error() string {
	return fmt.Sprintf("stopping a consumer timed out after %s", e.Timeout)
}
