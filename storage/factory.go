package storage

import (
	"fmt"
	"github.com/deis/logger/storage/file"
	"github.com/deis/logger/storage/ringbuffer"
)

// NewAdapter returns a pointer to an appropriate implementation of the Adapter interface, as
// determined by the storeageAdapterType string it is passed.
func NewAdapter(storeageAdapterType string, numLines int) (Adapter, error) {
	if storeageAdapterType == "file" {
		adapter, err := file.NewStorageAdapter()
		if err != nil {
			return nil, err
		}
		return adapter, nil
	}
	if storeageAdapterType == "memory" {
		adapter, err := ringbuffer.NewStorageAdapter(numLines)
		if err != nil {
			return nil, err
		}
		return adapter, nil
	}
	return nil, fmt.Errorf("Unrecognized storage adapter type: '%s'", storeageAdapterType)
}
