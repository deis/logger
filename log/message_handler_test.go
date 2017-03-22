package log

import (
	"encoding/json"
	"testing"

	"github.com/deis/logger/storage"
	"github.com/stretchr/testify/assert"
)

var (
	validControllerMessage = `{"log": "INFO [foo]: admin deployed 2fd9226", "time": "2016-10-18T20:29:38+00:00", "stream": "stderr", "docker": {"container_id": "containerId"}, "kubernetes": {"namespace_name": "foo", "pod_id": "podId", "pod_name": "podName", "container_name": "deis-controller", "labels": {"app": "foo",
"heritage": "deis", "type": "web", "version": "v2"}, "host": "host"}}`

	invalidControllerMessage = `{"log": "INFO: admin deployed 2fd9226", "stream": "stderr", "docker": {"container_id": "containerId"}, "kubernetes": {"namespace_name": "foo", "pod_id": "podId", "pod_name": "podName", "container_name": "deis-controller", "labels": {"app": "foo",
"heritage": "deis", "type": "web", "version": "v2"}, "host": "host"}}`

	validAppMessage = `{"log": "test message", "stream": "stderr", "time": "2016-10-18T20:29:38+00:00", "docker": {"container_id": "containerId"}, "kubernetes": {"namespace_name": "foo", "pod_id": "podId", "pod_name": "foo-web-845861952-nzf60", "container_name": "foo-web", "labels": {"app": "foo",
"heritage": "deis", "type": "web", "version": "v2"}, "host": "host"}}`

	badPodNameMessage = `{"log": "test message", "stream": "stderr", "time": "2016-10-18T20:29:38+00:00", "docker": {"container_id": "containerId"}, "kubernetes": {"namespace_name": "foo", "pod_id": "podId", "pod_name": "foo-web-845861952", "container_name": "foo-web", "labels": {"app": "foo",
"heritage": "deis", "type": "web", "version": "v2"}, "host": "host"}}`

	badjson = `{"log":}`
)

func TestValidControllerMessage(t *testing.T) {
	message := new(Message)
	err := json.Unmarshal([]byte(validControllerMessage), message)
	assert.NoError(t, err, "error occured parsing log message")
	assert.True(t, fromController(message), "json is not from controller")
}

func TestInvalidControllerMessage(t *testing.T) {
	message := new(Message)
	err := json.Unmarshal([]byte(validAppMessage), message)
	assert.NoError(t, err, "error occured parsing log message")
	assert.False(t, fromController(message), "valid controller message")
}

func TestGetApplicationFromValidControllerMessage(t *testing.T) {
	message := new(Message)
	err := json.Unmarshal([]byte(validControllerMessage), message)
	assert.NoError(t, err, "error occured parsing log message")
	expected := getApplicationFromControllerMessage(message)
	assert.Equal(t, expected, "foo", "failed to retrieve app from message")
}

func TestBuildControllerLogMessageFromValidMessage(t *testing.T) {
	message := new(Message)
	err := json.Unmarshal([]byte(validControllerMessage), message)
	assert.NoError(t, err, "error occured parsing log message")
	expected := buildControllerLogMessage(message)
	assert.Equal(t, expected,
		"2016-10-18T20:29:38+00:00 deis[controller]: INFO admin deployed 2fd9226",
		"failed to build controller log")
}

func TestBuildApplicationLogMessageFromValidMessage(t *testing.T) {
	message := new(Message)
	err := json.Unmarshal([]byte(validAppMessage), message)
	assert.NoError(t, err, "error occured parsing log message")
	expected := buildApplicationLogMessage(message)
	assert.Equal(t, expected,
		"2016-10-18T20:29:38+00:00 foo[web.v2.nzf60]: test message",
		"failed to build application log")
}

func TestBuildApplicationLogMessageFromInvalidMessage(t *testing.T) {
	message := new(Message)
	err := json.Unmarshal([]byte(badPodNameMessage), message)
	assert.NoError(t, err, "error occured parsing log message")
	expected := buildApplicationLogMessage(message)
	assert.Equal(t, expected,
		"2016-10-18T20:29:38+00:00 foo[web.v2]: test message",
		"failed to build application log")
}

func TestHandleValidAppMessage(t *testing.T) {
	a, err := storage.NewRingBufferAdapter(1)
	assert.NoError(t, err, "error creating ring buffer")
	err = handle([]byte(validAppMessage), a)
	assert.NoError(t, err, "error occured storing log message")
	expected, _ := a.Read("foo", 1)
	assert.Equal(t, expected[0],
		"2016-10-18T20:29:38+00:00 foo[web.v2.nzf60]: test message",
		"failed to aquire application log message")
}

func TestHandleValidControllerMessage(t *testing.T) {
	a, err := storage.NewRingBufferAdapter(1)
	assert.NoError(t, err, "error creating ring buffer")
	err = handle([]byte(validControllerMessage), a)
	assert.NoError(t, err, "error occured storing log message")
	expected, _ := a.Read("foo", 1)
	assert.Equal(t, expected[0],
		"2016-10-18T20:29:38+00:00 deis[controller]: INFO admin deployed 2fd9226",
		"failed to aquire controller log message")
}

func TestHandleInvalidAppMessage(t *testing.T) {
	a, err := storage.NewRingBufferAdapter(1)
	assert.NoError(t, err, "error creating ring buffer")
	err = handle([]byte(validAppMessage), a)
	assert.NoError(t, err, "error occured storing log message")
	expected, _ := a.Read("foo", 1)
	assert.Equal(t, expected[0],
		"2016-10-18T20:29:38+00:00 foo[web.v2.nzf60]: test message",
		"failed to aquire application log message")
}

func TestHandleInvalidControllerMessage(t *testing.T) {
	a, err := storage.NewRingBufferAdapter(1)
	assert.NoError(t, err, "error creating ring buffer")
	err = handle([]byte(badjson), a)
	assert.Error(t, err, "no error occured parsing json")
}
