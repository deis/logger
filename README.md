# Deis Logger
[![Build Status](https://travis-ci.org/deis/logger.svg?branch=master)](https://travis-ci.org/deis/logger)
[![Go Report Card](http://goreportcard.com/badge/deis/logger)](http://goreportcard.com/report/deis/logger)
[![Docker Repository on Quay](https://quay.io/repository/deis/logger/status "Docker Repository on Quay")](https://quay.io/repository/deis/logger)

Deis (pronounced DAY-iss) is an open source PaaS that makes it easy to deploy and manage
applications on your own servers. Deis builds on [Kubernetes](http://kubernetes.io/) to provide
a lightweight, [Heroku-inspired](http://heroku.com) workflow.

![Deis Graphic](https://s3-us-west-2.amazonaws.com/get-deis/deis-graphic-small.png)

A system logger for use in the [Deis](http://deis.io) open source PaaS.

This Docker image is based on the official
[alpine:3.2](https://registry.hub.docker.com/_/alpine/) image.

## Description
The new v2 logger implementation has seen a simplification from the last rewrite. While it still uses much of that code it no longer depends on `etcd`. Instead, we will use kubernetes service discovery to determine where logger is running.

We have also decided to not use `logspout` as the mechanism to get logs from each container to the `logger` component. Now we will use [fluentd](http://fluentd.org) which is a widely supported logging framework with hundreds of plugins. This will allow the end user to configure multiple destinations such as Elastic Search and other Syslog compatible endpoints like [papertrail](http://papertrailapp.com).

** This image requires that the `Daemon Sets` api be available on the kubernetes cluster** For more information on running the `Daemon Sets` api see the [following](https://github.com/kubernetes/kubernetes/blob/master/docs/api.md#enabling-resources-in-the-extensions-group).

## Running logger v2
The following environment variables can be used to configure logger:

* `STORAGE_ADAPTER`: How to store logs that are sent to the logger interface. Default is `memory`
* `NUMBER_OF_LINES`: How many lines to store in the ring buffer. Default is `1000`.

### Installation
Because of the requirement of Daemon Sets we have chosen not to include the logging components in the main [Deis chart](https://github.com/deis/charts/tree/master/deis).

To install the logging system please do the following:

```
$ helm repo add deis https://github.com/deis/charts
$ helm install deis/deis-logger
```
Watch for the logging components to come up:

```
$ kubectl get pods --namespace=deis
```

You should see output similar to this:

```
NAME                        READY     STATUS    RESTARTS   AGE
deis-builder-knypb          1/1       Running   0          1d
deis-database-ldbam         1/1       Running   0          2h
deis-logger-fluentd-tq187   1/1       Running   0          1d
deis-logger-iwos5           1/1       Running   0          2d
deis-minio-zlmk8            1/1       Running   0          3d
deis-registry-bys9n         1/1       Running   0          3d
deis-router-h5f0i           1/1       Running   0          3d
deis-workflow-2v84b         1/1       Running   3          2d
```

After the logging components are installed you will need to restart the `deis-workflow` pod so it can pick up the new service endpoint for the logger component.

```
$ kubectl delete pod <workflow pod name>
```

The replication controller should immediately restart a new pod. Now you can use the `deis logs` command for your applications.

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
* `DEV_REGISTRY=quay.io IMAGE_PREFIX=myaccount make kube-create`: Does a sed replacement of the image name and creates a tmp manifest file that we will use to deploy logger component to kubernetes. This will also start 2 `fluentd` daemonsets.
* `make kube-delete`: This will remove all the logger components from the kubernetes cluster.
* `DEV_REGISTRY=quay.io IMAGE_PREFIX=myaccount make kube-replace`: This will rebuild the binary and image, push it to the remote registry, and then replace the running components with the new version.

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
