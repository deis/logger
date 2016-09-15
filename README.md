# Deis Logger
[![Build Status](https://travis-ci.org/deis/logger.svg?branch=master)](https://travis-ci.org/deis/logger)
[![codecov.io](https://codecov.io/github/deis/logger/coverage.svg?branch=master)](https://codecov.io/github/deis/logger?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/deis/logger)](https://goreportcard.com/report/github.com/deis/logger)
[![Docker Repository on Quay](https://quay.io/repository/deis/logger/status "Docker Repository on Quay")](https://quay.io/repository/deis/logger)

Deis (pronounced DAY-iss) is an open source PaaS that makes it easy to deploy and manage
applications on your own servers. Deis builds on [Kubernetes](http://kubernetes.io/) to provide
a lightweight, [Heroku-inspired](http://heroku.com) workflow.

![Deis Graphic](https://s3-us-west-2.amazonaws.com/get-deis/deis-graphic-small.png)

A system logger for use in the [Deis](http://deis.io) open source PaaS.

This Docker image is based on [quay.io/deis/base](https://github.com/deis/docker-base) image. You can see what version we are currently using in the [Dockerfile](rootfs/Dockerfile)

## Description
The new v2 logger implementation has seen a simplification from the last rewrite. While it still uses much of that code it no longer depends on `etcd`. Instead, we will use kubernetes service discovery to determine where logger is running.

We have also decided to not use `logspout` as the mechanism to get logs from each container to the `logger` component. Now we will use [fluentd](http://fluentd.org) which is a widely supported logging framework with hundreds of plugins. This will allow the end user to configure multiple destinations such as Elastic Search and other Syslog compatible endpoints like [papertrail](http://papertrailapp.com).

## Configuration
The following environment variables can be used to configure logger:

| Name | Default Value |
|------|---------------|
| STORAGE_ADAPTER | "redis" |
| NUMBER_OF_LINES (per app) | "1000" |
| AGGREGATOR_TYPE | "nsq" |
| DEIS_NSQD_SERVICE_HOST | "" |
| DEIS_NSQD_SERVICE_PORT_TRANSPORT | 4150 |
| NSQ_TOPIC | logs |
| NSQ_CHANNEL | consume |
| NSQ_HANDLER_COUNT | 30 |
| AGGREGATOR_STOP_TIMEOUT_SEC | 1 |
| DEIS_LOGGER_REDIS_SERVICE_HOST | "" |
| DEIS_LOGGER_REDIS_SERVICE_PORT | 6379 |
| DEIS_LOGGER_REDIS_PASSWORD | "" |
| DEIS_LOGGER_REDIS_DB | 0 |
| DEIS_LOGGER_REDIS_PIPELINE_LENGTH | 50 |
| DEIS_LOGGER_REDIS_PIPELINE_TIMEOUT_SECONDS | 1 |

## Development
The only assumption this project makes about your environment is that you have a working docker host to build the image against.

### Building binary and image
To build the binary and image run the following make command:

```console
DEV_REGISTRY=quay.io IMAGE_PREFIX=myaccount make build
IMAGE_PREFIX=myaccount make build
DEV_REGISTRY=myhost:5000 make build
```

### Pushing the image
The makefile assumes that you are pushing the image to a remote repository like quay or dockerhub. So you will need to supply the `REGISTRY` environment variable.

```console
DEV_REGISTRY=quay.io IMAGE_PREFIX=myaccount make push
IMAGE_PREFIX=myaccount make push
DEV_REGISTRY=myhost:5000 make push
```

### Kubernetes interactions
* `make install` - Install the recently built docker image into the kubernetes cluster
* `make upgrade` - Upgrade a currently installed image
* `make uninstall` - Uninstall logger from a kubernetes cluster

### Architecture Diagram
```
                        ┌────────┐                            
                        │ Router │                  ┌────────┐
                        └────────┘                  │ Logger │
                            │                       └────────┘
                        Log file                        │    
                            │                           │    
                            ▼                           ▼    
┌────────┐             ┌─────────┐    logs/metrics   ┌─────┐
│App Logs│──Log File──▶│ fluentd │───────topics─────▶│ NSQ │
└────────┘             └─────────┘                   └─────┘
                                                        │    
                                                        │    
┌─────────────┐                                         │    
│ HOST        │                                         ▼    
│  Telegraf   │───┐                                ┌────────┐
└─────────────┘   │                                │Telegraf│
                  │                                └────────┘
┌─────────────┐   │                                    │    
│ HOST        │   │    ┌───────────┐                   │    
│  Telegraf   │───┼───▶│ InfluxDB  │◀────Wire ─────────┘    
└─────────────┘   │    └───────────┘   Protocol       
                  │          ▲                        
┌─────────────┐   │          │                        
│ HOST        │   │          ▼                        
│  Telegraf   │───┘    ┌──────────┐                   
└─────────────┘        │ Grafana  │                   
                       └──────────┘                                        
```

## License

© 2016 Engine Yard, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
