package log

import (
	"fmt"
	"reflect"
	"testing"
)

type stubStorageAdapter struct {
}

func (a *stubStorageAdapter) Write(app string, message string) error {
	return nil
}

func (a *stubStorageAdapter) Read(app string, lines int) ([]string, error) {
	return []string{}, nil
}

func (a *stubStorageAdapter) Destroy(app string) error {
	return nil
}

func (a *stubStorageAdapter) Reopen() error {
	return nil
}

func TestGetUsingInvalidValues(t *testing.T) {
	_, err := NewAggregator("bogus", &stubStorageAdapter{})
	if err == nil || err.Error() != fmt.Sprintf("Unrecognized aggregator type: '%s'", "bogus") {
		t.Error("Did not receive expected error message")
	}
}

func TestNSQBasedAggregator(t *testing.T) {
	a, err := NewAggregator("nsq", &stubStorageAdapter{})
	if err != nil {
		t.Error(err)
	}
	expected := "*log.nsqAggregator"
	aType := reflect.TypeOf(a).String()
	if aType != expected {
		t.Errorf("Expected a %s, but got a %s", expected, aType)
	}
}
