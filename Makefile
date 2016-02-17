SHELL = /bin/bash
GO = go
GOFMT = gofmt -l
GOLINT = golint
GOTEST = $(GO) test --cover --race -v
GOVET = $(GO) vet
GO_FILES = $(wildcard *.go)
GO_PACKAGES = drain storage syslogish weblog
GO_PACKAGES_REPO_PATH = $(addprefix $(REPO_PATH)/,$(GO_PACKAGES))
GO_TESTABLE_PACKAGES_REPO_PATH = $(addprefix $(REPO_PATH)/,drain drain/simple storage storage/file storage/ringbuffer)

# the filepath to this repository, relative to $GOPATH/src
REPO_PATH = github.com/deis/logger

# The following variables describe the containerized development environment
# and other build options
DEV_ENV_IMAGE := quay.io/deis/go-dev:0.7.0
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}
DEV_ENV_CMD := docker run --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} ${DEV_ENV_IMAGE}
DEV_ENV_CMD_INT := docker run -it --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} ${DEV_ENV_IMAGE}
LDFLAGS := "-s -X main.version=${VERSION}"

BINARY_DEST_DIR = image/bin

DOCKER_HOST = $(shell echo $$DOCKER_HOST)
BUILD_TAG ?= git-$(shell git rev-parse --short HEAD)
SHORT_NAME ?= logger
DEIS_REGISTRY ?= ${DEV_REGISTRY}
IMAGE_PREFIX ?= deis
IMAGE_LATEST := ${DEIS_REGISTRY}${IMAGE_PREFIX}/${SHORT_NAME}:latest
IMAGE := ${DEIS_REGISTRY}${IMAGE_PREFIX}/${SHORT_NAME}:${BUILD_TAG}

info:
	@echo "Build tag:  ${BUILD_TAG}"
	@echo "Registry:   ${DEIS_REGISTRY}"
	@echo "Image:      ${IMAGE}"

check-docker:
	@if [ -z $$(which docker) ]; then \
	  echo "Missing docker client which is required for development"; \
	  exit 2; \
	fi

# Allow developers to step into the containerized development environment
dev: check-docker
	${DEV_ENV_CMD_INT} bash

# Containerized dependency resolution
bootstrap: check-docker
	${DEV_ENV_CMD} glide install

# Containerized build of the binary
build-with-container: check-docker
	mkdir -p ${BINARY_DEST_DIR}
	${DEV_ENV_CMD} make build-binary
	docker build --rm -t ${IMAGE} image

build: build-with-container docker-build

push: docker-push

docker-build: check-docker
	docker build -t $(IMAGE_LATEST) image
	docker tag -f $(IMAGE_LATEST) $(IMAGE)

docker-push: check-docker
	docker push $(IMAGE)

clean: check-docker
	docker rmi $(IMAGE)

update-manifests:
	sed 's#\(image:\) .*#\1 $(IMAGE)#' manifests/deis-logger-rc.yaml > manifests/deis-logger-rc.tmp.yaml

build-binary:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags ${LDFLAGS} -o $(BINARY_DEST_DIR)/logger main.go

test: test-style test-unit

test-style: check-docker
	${DEV_ENV_CMD} make style-check

style-check:
# display output, then check
	$(GOFMT) $(GO_PACKAGES) $(GO_FILES)
	@$(GOFMT) $(GO_PACKAGES) $(GO_FILES) | read; if [ $$? == 0 ]; then echo "gofmt check failed."; exit 1; fi
	$(GOVET) $(REPO_PATH) $(GO_PACKAGES_REPO_PATH)
	$(GOLINT) ./...

test-unit:
	${DEV_ENV_CMD} $(GOTEST) $(GO_TESTABLE_PACKAGES_REPO_PATH)

kube-install:
	kubectl create -f manifests/deis-logger-svc.yaml
	kubectl create -f manifests/deis-logger-rc.yaml
	kubectl create -f manifests/deis-logger-fluentd-daemon.yaml

kube-delete:
	-kubectl delete -f manifests/deis-logger-svc.yaml
	-kubectl delete -f manifests/deis-logger-rc.tmp.yaml
	-kubectl delete -f manifests/deis-logger-fluentd-daemon.yaml

kube-create: update-manifests
	kubectl create -f manifests/deis-logger-svc.yaml
	kubectl create -f manifests/deis-logger-rc.tmp.yaml
	kubectl create -f manifests/deis-logger-fluentd-daemon.yaml

kube-replace: build push update-manifests
	kubectl replace --force -f manifests/deis-logger-rc.tmp.yaml

kube-update: update-manifests
	kubectl delete -f manifests/deis-logger-rc.tmp.yaml
	kubectl create -f manifests/deis-logger-rc.tmp.yaml
