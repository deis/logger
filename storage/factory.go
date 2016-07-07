package storage

import (
	"fmt"
)

type errUnrecognizedStorageAdapterType struct {
	adapterType string
}

func (e errUnrecognizedStorageAdapterType) Error() string {
	return fmt.Sprintf("Unrecognized storage adapter type: %s", e.adapterType)
}

// NewAdapter returns a pointer to an appropriate implementation of the Adapter interface, as
// determined by the adapterType string it is passed.
func NewAdapter(adapterType string, numLines int) (Adapter, error) {
	if adapterType == "file" {
		adapter, err := NewFileAdapter()
		if err != nil {
			return nil, err
		}
		return adapter, nil
	}
	if adapterType == "memory" {
		adapter, err := NewRingBufferAdapter(numLines)
		if err != nil {
			return nil, err
		}
		return adapter, nil
	}
	if adapterType == "redis" {
		adapter, err := NewRedisStorageAdapter(numLines)
		if err != nil {
			return nil, err
		}
		return adapter, nil
	}
	return nil, errUnrecognizedStorageAdapterType{adapterType: adapterType}
}
