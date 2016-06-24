package log

import (
	"fmt"

	"github.com/deis/logger/storage"
)

// NewAggregator returns a pointer to an appropriate implementation of the Aggregator interface, as
// determined by the aggregatorType string it is passed.
func NewAggregator(aggregatorType string, storageAdapter storage.Adapter) (Aggregator, error) {
	if aggregatorType == "nsq" {
		return newNSQAggregator(storageAdapter), nil
	}
	return nil, fmt.Errorf("Unrecognized aggregator type: '%s'", aggregatorType)
}
