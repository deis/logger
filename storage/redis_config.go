package storage

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const (
	appName = "logger"
)

type redisConfig struct {
	Host                   string `envconfig:"DEIS_LOGGER_REDIS_SERVICE_HOST" default:""`
	Port                   int    `envconfig:"DEIS_LOGGER_REDIS_SERVICE_PORT" default:"6379"`
	Password               string `envconfig:"DEIS_LOGGER_REDIS_PASSWORD" default:""`
	DB                     int    `envconfig:"DEIS_LOGGER_REDIS_DB" default:"0"`
	PipelineLength         int    `envconfig:"DEIS_LOGGER_REDIS_PIPELINE_LENGTH" default:"50"`
	PipelineTimeoutSeconds int    `envconfig:"DEIS_LOGGER_REDIS_PIPELINE_TIMEOUT_SECONDS" default:"1"`
	PipelineTimeout        time.Duration
}

func parseConfig(appName string) (*redisConfig, error) {
	ret := new(redisConfig)
	if err := envconfig.Process(appName, ret); err != nil {
		return nil, err
	}
	ret.PipelineTimeout = time.Duration(ret.PipelineTimeoutSeconds) * time.Second
	return ret, nil
}
