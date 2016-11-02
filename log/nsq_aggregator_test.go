package log

import (
	"testing"

	"github.com/deis/logger/storage"
	"github.com/stretchr/testify/assert"
)

func TestAggregator(t *testing.T) {
	storageAdapter, err := storage.NewAdapter("memory", 100)
	assert.NoError(t, err)
	aggregator, err := NewAggregator("nsq", storageAdapter)
	assert.NoError(t, err)
	err = aggregator.Listen()
	assert.NoError(t, err)
	stoppedCh := aggregator.Stopped()
	err = aggregator.Stop()
	assert.NoError(t, err)
	stopErr := <-stoppedCh
	assert.NoError(t, stopErr, "Aggregator stopped with error")
}
