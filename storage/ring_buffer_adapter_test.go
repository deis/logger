package storage

import (
	"fmt"
	"testing"
)

func TestRingBufferReadFromNonExistingApp(t *testing.T) {
	// Initialize a new storage adapter
	sa, err := NewRingBufferAdapter(10)
	if err != nil {
		t.Error(err)
	}
	a, ok := sa.(*ringBufferAdapter)
	if !ok {
		t.Fatalf("returned adapter was not a ringBuffer")
	}
	// No logs have been writter; there should be no ringBuffer for app
	messages, err := a.Read(app, 10)
	if messages != nil {
		t.Error("Expected no messages, but got some")
	}
	if err == nil || err.Error() != fmt.Sprintf("Could not find logs for '%s'. No ringbuffer existed for '%s'.", app, app) {
		t.Error("Did not receive expected error message")
	}
}

func TestRingBufferWithBadBufferSizes(t *testing.T) {
	// Initialize with invalid buffer sizes
	for _, size := range []int{-1, 0} {
		a, err := NewRingBufferAdapter(size)
		if a != nil {
			t.Error("Expected no storage adapter, but got one")
		}
		if err == nil || err.Error() != fmt.Sprintf("Invalid ringBuffer size: %d", size) {
			t.Error("Did not receive expected error message")
		}
	}
}

func TestRingBufferLogs(t *testing.T) {
	// Initialize with small buffers
	a, err := NewRingBufferAdapter(10)
	if err != nil {
		t.Error(err)
	}
	// And write a few logs to it, but do NOT fill it up
	for i := 0; i < 5; i++ {
		if err := a.Write(app, fmt.Sprintf("message %d", i)); err != nil {
			t.Error(err)
		}
	}
	// Read more logs than there are
	messages, err := a.Read(app, 8)
	if err != nil {
		t.Error(err)
	}
	// Should only get as many messages as we actually have
	if len(messages) != 5 {
		t.Errorf("only expected 5 log messages, got %d", len(messages))
	}
	// Read fewer logs than there are
	messages, err = a.Read(app, 3)
	if err != nil {
		t.Error(err)
	}
	// Should get the 3 MOST RECENT logs
	if len(messages) != 3 {
		t.Errorf("only expected 5 log messages, got %d", len(messages))
	}
	for i := 0; i < 3; i++ {
		expectedMessage := fmt.Sprintf("message %d", i+2)
		if messages[i] != expectedMessage {
			t.Errorf("expected: \"%s\", got \"%s\"", expectedMessage, messages[i])
		}
	}
	// Overfill the buffer
	for i := 5; i < 11; i++ {
		if err := a.Write(app, fmt.Sprintf("message %d", i)); err != nil {
			t.Error(err)
		}
	}
	// Read more logs than the buffer can hold
	messages, err = a.Read(app, 20)
	if err != nil {
		t.Error(err)
	}
	// Should only get as many messages as the buffer can hold
	if len(messages) != 10 {
		t.Errorf("only expected 10 log messages, got %d", len(messages))
	}
	// And they should only be the 10 MOST RECENT logs
	for i := 0; i < 10; i++ {
		expectedMessage := fmt.Sprintf("message %d", i+1)
		if messages[i] != expectedMessage {
			t.Errorf("expected: \"%s\", got \"%s\"", expectedMessage, messages[i])
		}
	}
}

func TestRingBufferDestroy(t *testing.T) {
	sa, err := NewRingBufferAdapter(10)
	if err != nil {
		t.Error(err)
	}

	a, ok := sa.(*ringBufferAdapter)
	if !ok {
		t.Fatalf("returned adapter was not a ringBuffer")
	}
	// Write a log to create the file
	if err := a.Write(app, "Hello, log!"); err != nil {
		t.Error(err)
	}
	// A ringBuffer should exist for the app
	if _, ok := a.ringBuffers[app]; !ok {
		t.Error("Log ringbuffer was expected to exist, but doesn't.")
	}
	// Now destroy it
	if err := a.Destroy(app); err != nil {
		t.Error(err)
	}
	// Now check that the ringBuffer no longer exists
	if _, ok := a.ringBuffers[app]; ok {
		t.Error("Log ringbuffer still exist, but was expected not to.")
	}
}
