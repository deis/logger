package log

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNsqURL(t *testing.T) {
	c := config{
		NSQHost: "somehost",
		NSQPort: 3333,
	}
	assert.Equal(t, c.nsqURL(), "somehost:3333")
}

func TestStopTimeoutDuration(t *testing.T) {
	c := config{
		StopTimeoutSeconds: 60,
	}
	assert.Equal(t, c.stopTimeoutDuration(), time.Duration(c.StopTimeoutSeconds)*time.Second)
}

func TestParseConfig(t *testing.T) {
	os.Setenv("NSQ_TOPIC", "topic")
	os.Setenv("NSQ_CHANNEL", "channel")
	os.Setenv("NSQ_HANDLER_COUNT", "3")
	os.Setenv("AGGREGATOR_STOP_TIMEOUT_SEC", "2")

	port, err := strconv.Atoi(os.Getenv("DEIS_NSQD_SERVICE_PORT_TRANSPORT"))
	assert.NoError(t, err)

	c, err := parseConfig("foo")
	assert.NoError(t, err)
	assert.Equal(t, c.NSQHost, os.Getenv("DEIS_NSQD_SERVICE_HOST"))
	assert.Equal(t, c.NSQPort, port)
	assert.Equal(t, c.NSQTopic, "topic")
	assert.Equal(t, c.NSQChannel, "channel")
	assert.Equal(t, c.NSQHandlerCount, 3)
	assert.Equal(t, c.StopTimeoutSeconds, 2)
}
