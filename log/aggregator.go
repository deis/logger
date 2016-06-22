package log

// Aggregator is an interface for pluggable components that collect log messages delivered via
// some transport mechanism
type Aggregator interface {
	// Listen causes an aggregator to begin consuming messages via its underlying transport
	// mechanism. Implementations of this must be non-blocking. Stop() can be called to stop the
	// aggregator.
	Listen() error
	// Stop stops the consumer and blocks until it stops or the configured duration passes. In the
	// latter case, returns an error of type ErrStopTimedOut
	Stop() error
	// Stopped returns a channel that will receive when the consumer has stopped. The error received
	// may be nil, which means it was stopped cleanly. A non-nil error means it stopped because it
	// errored.
	Stopped() <-chan error
}
