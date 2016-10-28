package storage

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	host := os.Getenv("DEIS_LOGGER_REDIS_SERVICE_HOST")
	password := os.Getenv("DEIS_LOGGER_REDIS_PASSWORD")
	db := os.Getenv("DEIS_LOGGER_REDIS_DB")
	pipelineLength := os.Getenv("DEIS_LOGGER_REDIS_PIPELINE_LENGTH")
	pipelineTimeoutSeconds := os.Getenv("DEIS_LOGGER_REDIS_PIPELINE_TIMEOUT_SECONDS")

	os.Setenv("DEIS_LOGGER_REDIS_PASSWORD", "password")
	os.Setenv("DEIS_LOGGER_REDIS_DB", "2")
	os.Setenv("DEIS_LOGGER_REDIS_PIPELINE_LENGTH", "1")
	os.Setenv("DEIS_LOGGER_REDIS_PIPELINE_TIMEOUT_SECONDS", "2")

	p, err := strconv.Atoi(os.Getenv("DEIS_LOGGER_REDIS_SERVICE_PORT"))
	assert.NoError(t, err)

	c, err := parseConfig("foo")
	assert.NoError(t, err, "error parsing config")
	assert.Equal(t, c.Host, host)
	assert.Equal(t, c.Port, p)
	assert.Equal(t, c.Password, "password")
	assert.Equal(t, c.PipelineLength, 1)
	assert.Equal(t, c.PipelineTimeoutSeconds, 2)

	os.Setenv("DEIS_LOGGER_REDIS_PASSWORD", password)
	os.Setenv("DEIS_LOGGER_REDIS_DB", db)
	os.Setenv("DEIS_LOGGER_REDIS_PIPELINE_LENGTH", pipelineLength)
	os.Setenv("DEIS_LOGGER_REDIS_PIPELINE_TIMEOUT_SECONDS", pipelineTimeoutSeconds)
}
