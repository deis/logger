package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

const (
	appName = "logger"
)

type config struct {
	StorageType           string `envconfig:"STORAGE_ADAPTER" default:"memory"`
	NumLines              int    `envconfig:"NUMBER_OF_LINES" default:"1000"`
	NSQServiceHost        string `envconfig:"DEIS_NSQD_SERVICE_HOST" default:""`           // k8s service discovery env var
	NSQServicePort        int    `envconfig:"DEIS_NSQD_SERVICE_PORT_TRANSPORT" default:""` // k8s service discovery env var
	NSQMaxInFlight        int    `envconfig:"MAX_IN_FLIGHT" default:"30"`
	NSQTopic              string `envconfig:"NSQ_TOPIC" default:"logs"`
	NSQChannel            string `envconfig:"NSQ_CHANNEL" default:"consume"`
	NSQConsumerStopDurSec int    `envconfig:"NSQ_CONSUMER_STOP_WAIT_DUR_SEC" default:"1"`
	NSQConsumerNumThreads int    `envconfig:"NSQ_CONSUMER_NUM_THREADS" default:"5"`
}

func (c config) nsqURL() string {
	return fmt.Sprintf("%s:%d", c.NSQServiceHost, c.NSQServicePort)
}

func parseConfig(appName string) (*config, error) {
	ret := new(config)
	if err := envconfig.Process(appName, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
