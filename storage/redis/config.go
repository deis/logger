package redis

import (
	"github.com/kelseyhightower/envconfig"
)

const (
	appName = "logger"
)

type config struct {
	RedisHost     string `envconfig:"DEIS_LOGGER_REDIS_SERVICE_HOST" default:""`
	RedisPort     int    `envconfig:"DEIS_LOGGER_REDIS_SERVICE_PORT" default:"6379"`
	RedisPassword string `envconfig:"DEIS_LOGGER_REDIS_PASSWORD" default:""`
	RedisDB       int    `envconfig:"DEIS_LOGGER_REDIS_DB" default:"0"`
}

func parseConfig(appName string) (*config, error) {
	ret := new(config)
	if err := envconfig.Process(appName, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
