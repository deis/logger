package storage

import (
	"reflect"
	"testing"
)

const (
	app = "test-app"
)

func TestFactoryGetUsingInvalidValues(t *testing.T) {
	const adapterType = "bogus"
	_, err := NewAdapter(adapterType, 1)
	if err == nil {
		t.Fatalf("Did not receive an error message")
	}
	unrecognizedErr, ok := err.(errUnrecognizedStorageAdapterType)
	if !ok {
		t.Fatalf("Expected an errUnrecognizedStorageAdapterType, received %s", err)
	}
	if unrecognizedErr.adapterType != adapterType {
		t.Fatalf("Got an errUnrecognizedStorageAdapterType, but expected adapter type %s, got %s", adapterType, unrecognizedErr.adapterType)
	}
}

func TestFactoryGetFileBasedAdapter(t *testing.T) {
	a, err := NewAdapter("file", 1)
	if err != nil {
		t.Error(err)
	}
	retType, ok := a.(*fileAdapter)
	if !ok {
		t.Fatalf("Expected a *fileAdapter, got %s", reflect.TypeOf(retType).String())
	}
}

func TestFactoryGetMemoryBasedAdapter(t *testing.T) {
	a, err := NewAdapter("memory", 1)
	if err != nil {
		t.Error(err)
	}
	retType, ok := a.(*ringBufferAdapter)
	if !ok {
		t.Fatalf("Expected a *ringBufferAdapter, got %s", reflect.TypeOf(retType).String())
	}
}

func TestGetRedisBasedAdapter(t *testing.T) {
	a, err := NewAdapter("redis", 1)
	if err != nil {
		t.Error(err)
	}
	retType, ok := a.(*redisAdapter)
	if !ok {
		t.Errorf("Expected a redisAdapter, but got a %s", reflect.TypeOf(retType).String())
	}
}
