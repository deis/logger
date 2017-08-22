# Deis Logger
[![Build Status](https://ci.deis.io/job/logger/badge/icon)](https://ci.deis.io/job/logger)
[![codecov.io](https://codecov.io/github/deis/logger/coverage.svg?branch=master)](https://codecov.io/github/deis/logger?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/deis/logger)](https://goreportcard.com/report/github.com/deis/logger)
[![Docker Repository on Quay](https://quay.io/repository/deis/logger/status "Docker Repository on Quay")](https://quay.io/repository/deis/logger)

Deis (pronounced DAY-iss) Workflow is an open source Platform as a Service (PaaS) that adds a developer-friendly layer to any [Kubernetes](http://kubernetes.io) cluster, making it easy to deploy and manage applications on your own servers.

![Deis Graphic](https://getdeis.blob.core.windows.net/get-deis/deis-graphic-small.png)

For more information about the Deis Workflow, please visit the main project page at https://github.com/deis/workflow.

We welcome your input! If you have feedback, please [submit an issue][issues]. If you'd like to participate in development, please read the "Development" section below and [submit a pull request][prs].

## Description
A system logger for use in the [Deis Workflow](https://deis.com/workflow/) open source PaaS.

This Docker image is based on [quay.io/deis/base](https://github.com/deis/docker-base) image. You can see what version we are currently using in the [Dockerfile](rootfs/Dockerfile)

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

[issues]: https://github.com/deis/logger/issues
[prs]: https://github.com/deis/logger/pulls
