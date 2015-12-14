# Deis Logger
[![Build Status](https://travis-ci.org/deis/logger.svg?branch=master)](https://travis-ci.org/deis/logger) [![Go Report Card](http://goreportcard.com/badge/deis/logger)](http://goreportcard.com/report/deis/logger)


A system logger for use in the [Deis](http://deis.io) open source PaaS.

This Docker image is based on the official
[alpine:3.2](https://registry.hub.docker.com/_/alpine/) image.

## Description
The new v2 logger implementation has seen a simplification from the last rewrite. While it still uses much of that code it no longer depends on `etcd`. Instead, we will use kubernetes service discovery to determine where logger is running.

We have also decided to not use `logspout` as the mechanism to get logs from each container to the `logger` component. Now we will use [fluentd](http://fluentd.org) which is a widely supported logging framework with hundreds of plugins. This will allow the end user to configure multiple destinations such as Elastic Search and other Syslog compatible endpoints like [papertrail](http://papertrailapp.com).

** This image requires that the `daemonsets` api be available on the kubernetes cluster** For more information on running the `daemonsets` api see the [following](https://github.com/kubernetes/kubernetes/blob/master/docs/api.md#enabling-resources-in-the-extensions-group).

## Running logger v2
The following environment variables can be used to configure logger:

* `LOGGER_ADDR`: The interface to bind the logging adapter to. Default is `0.0.0.0`.
* `LOGGER_PORT`: The port that the logger adapter is listening on. Default is `514`.
* `WEB_ADDR`: The interface to bind the web server to. Default is `0.0.0.0`.
* `WEB_PORT`: The port that the web server is listening on. Default is `8088`.
* `STORAGE_ADAPTER`: How to store logs that are sent to the logger interface. Default is `memory`
* `NUMBER_OF_LINES`: How many lines to store in the ring buffer. Default is `1000`.
* `LOG_PATH`: The path of where to store files when using the file storage adapter. Default is `/data/logs`.
* `DRAIN_URL`: Syslog server that the logger component can send data to. No default.

## Development
The only assumption this project makes about your environment is that you have a working docker host to build the image against.

### Building binary and image
To build the binary and image run the following make command:

```console
REGISTRY=quay.io/myaccount make build
REGISTRY=myaccount make build
REGISTRY=myhost:5000 make build
```

### Pushing the image
The makefile assumes that you are pushing the image to a remote repository like quay or dockerhub. So you will need to supply the `REGISTRY` environment variable.

```console
REGISTRY=quay.io/myaccount make push
REGISTRY=myaccount make push
REGISTRY=myhost:5000 make push
```

### Kubernetes interactions
* `REGISTRY=quay.io/myaccount make kube-create`: Does a sed replacement of the image name and creates a tmp manifest file that we will use to deploy logger component to kubernetes. This will also start 2 `fluentd` daemonsets.
* `make kube-delete`: This will remove all the logger components from the kubernetes cluster.
* `REGISTRY=quay.io/myaccount make kube-replace`: This will rebuild the binary and image, push it to the remote registry, and then replace the running components with the new version.

## License

© 2015 Engine Yard, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
