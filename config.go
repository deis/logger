package main

import (
	"github.com/kelseyhightower/envconfig"
)

const (
	appName = "logger"
)

type config struct {
	StorageType    string `envconfig:"STORAGE_ADAPTER" default:"redis"`
	NumLines       int    `envconfig:"NUMBER_OF_LINES" default:"1000"`
	AggregatorType string `envconfig:"AGGREGATOR_TYPE" default:"nsq"`
}

func parseConfig(appName string) (*config, error) {
	ret := new(config)
	if err := envconfig.Process(appName, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
