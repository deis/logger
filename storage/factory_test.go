package storage

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestGetUsingInvalidValues(t *testing.T) {
	_, err := NewAdapter("bogus", 1)
	if err == nil || err.Error() != fmt.Sprintf("Unrecognized storage adapter type: '%s'", "bogus") {
		t.Error("Did not receive expected error message")
	}
}

func TestGetFileBasedAdapter(t *testing.T) {
	a, err := NewAdapter("file", 1)
	if err != nil {
		t.Error(err)
	}
	expected := "*file.adapter"
	aType := reflect.TypeOf(a).String()
	if aType != expected {
		t.Errorf("Expected a %s, but got a %s", expected, aType)
	}
}

func TestGetMemoryBasedAdapter(t *testing.T) {
	a, err := NewAdapter("memory", 1)
	if err != nil {
		t.Error(err)
	}
	expected := "*ringbuffer.adapter"
	aType := reflect.TypeOf(a).String()
	if aType != expected {
		t.Errorf("Expected a %s, but got a %s", expected, aType)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
