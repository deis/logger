package storage

import (
	"container/ring"
	"fmt"
	"log"
	"sync"
)

type ringBuffer struct {
	ring  *ring.Ring
	mutex sync.RWMutex
}

func newRingBuffer(size int) *ringBuffer {
	return &ringBuffer{ring: ring.New(size)}
}

func (rb *ringBuffer) write(message string) {
	// Get a write lock since writing adjusts the value of the internal ring pointer
	rb.mutex.Lock()
	defer rb.mutex.Unlock()
	rb.ring = rb.ring.Next()
	rb.ring.Value = message
}

func (rb *ringBuffer) read(lines int) []string {
	if lines <= 0 {
		return []string{}
	}
	// Only need a read lock because nothing we're about to do affects the internal state of the
	// ringBuffer.  Mutliple reads can happen in parallel.  Only writing requires an exclusive lock.
	rb.mutex.RLock()
	defer rb.mutex.RUnlock()
	var start *ring.Ring
	if lines < rb.ring.Len() {
		start = rb.ring.Move(-1 * (lines - 1))
	} else {
		start = rb.ring.Next()
	}
	data := make([]string, 0, lines)
	start.Do(func(line interface{}) {
		if line == nil || lines <= 0 {
			return
		}
		lines--
		data = append(data, line.(string))
	})
	return data
}

type ringBufferAdapter struct {
	bufferSize  int
	ringBuffers map[string]*ringBuffer
	mutex       sync.Mutex
}

// NewRingBufferAdapter returns a storage adapter that uses an in-memory ring buffer of the given size.
func NewRingBufferAdapter(bufferSize int) (Adapter, error) {
	if bufferSize <= 0 {
		return nil, fmt.Errorf("Invalid ringBuffer size: %d", bufferSize)
	}
	return &ringBufferAdapter{bufferSize: bufferSize, ringBuffers: make(map[string]*ringBuffer)}, nil
}

// Write adds a log message to to an app-specific ringBuffer
func (a *ringBufferAdapter) Write(app string, message string) error {
	// Check first if we might actually have to add to the map of ringBuffer pointers so we can avoid
	// waiting for / obtaining a lock unnecessarily
	a.mutex.Lock()
	defer a.mutex.Unlock()
	rb, ok := a.ringBuffers[app]
	if !ok {
		// Ensure only one goroutine at a time can be adding a ringBuffer to the map of ringBuffers
		// pointers
		rb, ok = a.ringBuffers[app]
		if !ok {
			log.Printf("Creating buffer for app:%v", app)
			rb = newRingBuffer(a.bufferSize)
			a.ringBuffers[app] = rb
		}
	}
	rb.write(message)
	return nil
}

// Read retrieves a specified number of log lines from an app-specific ringBuffer
func (a *ringBufferAdapter) Read(app string, lines int) ([]string, error) {
	rb, ok := a.ringBuffers[app]
	if ok {
		data := rb.read(lines)
		if len(data) == 0 {
			return nil, fmt.Errorf("Could not find logs for '%s'. Ringbuffer existed for '%s', but returned no logs.", app, app)
		}
		return data, nil
	}
	return nil, fmt.Errorf("Could not find logs for '%s'. No ringbuffer existed for '%s'.", app, app)
}

// Destroy deletes stored logs for the specified application
func (a *ringBufferAdapter) Destroy(app string) error {
	// Check first if the map of ringBuffer pointers even contains the ringBuffer we intend to
	// delete so we can avoid waiting for / obtaining a lock unnecessarily
	_, ok := a.ringBuffers[app]
	if ok {
		a.mutex.Lock()
		defer a.mutex.Unlock()
		delete(a.ringBuffers, app)
	}
	return nil
}

func (a *ringBufferAdapter) Reopen() error {
	// No-op
	return nil
}
