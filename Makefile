SHELL = /bin/bash

GO = godep go
GOFMT = gofmt -l
GOLINT = golint
GOTEST = $(GO) test --cover --race -v
GOVET = $(GO) vet

SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
DOCKER_HOST = $(shell echo $$DOCKER_HOST)
REGISTRY = $(shell if [ "$$DEV_REGISTRY" == "registry.hub.docker.com" ]; then echo; else echo $$DEV_REGISTRY/; fi)
GIT_SHA = $(shell git rev-parse --short HEAD)

ifndef IMAGE_PREFIX
  IMAGE_PREFIX = deis/
endif

check-registry:
	@if [ -z "$$DEV_REGISTRY" ]; then \
	  echo "DEV_REGISTRY is not exported, try:  make dev-registry"; \
	exit 2; \
	fi

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/logger

GO_FILES = $(wildcard *.go)
GO_PACKAGES = drain storage syslogish weblog
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))
GO_TESTABLE_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,drain drain/simple storage storage/file storage/ringbuffer)

COMPONENT = $(notdir $(repo_path))
IMAGE = $(IMAGE_PREFIX)$(COMPONENT):$(BUILD_TAG)
DEV_IMAGE = $(REGISTRY)$(IMAGE)
BINARY_DEST_DIR = image/bin

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' -o $(BINARY_DEST_DIR)/logger github.com/deis/logger || exit 1
	@$(call check-static-binary,$(BINARY_DEST_DIR)/logger)
	docker build -t $(IMAGE) image

clean: check-docker check-registry
	rm -f image/bin/logger
	docker rmi $(IMAGE)

full-clean: check-docker check-registry
	docker images -q $(IMAGE_PREFIX)$(COMPONENT) | xargs docker rmi -f

dev-release: push set-image

push: check-registry
	docker tag -f $~(IMAGE) $(DEV_IMAGE)
	docker push $(DEV_IMAGE)

set-image: check-deisctl
	deisctl config $(COMPONENT) set image=$(DEV_IMAGE)

release:
	docker push $(IMAGE)

deploy: build dev-release restart

test: test-style test-unit test-functional

test-functional:
	@$(MAKE) -C ../tests/ test-etcd
	GOPATH=`cd ../tests/ && godep path`:$(GOPATH) go test -v ./tests/...

test-style:
# display output, then check
	$(GOFMT) $(GO_PACKAGES) $(GO_FILES)
	@$(GOFMT) $(GO_PACKAGES) $(GO_FILES) | read; if [ $$? == 0 ]; then echo "gofmt check failed."; exit 1; fi
	$(GOVET) $(repo_path) $(GO_PACKAGES_REPO_PATH)
	$(GOLINT) ./...

test-unit: test-style
	$(GOTEST) $(GO_TESTABLE_PACKAGES_REPO_PATH)

coverage:
	go test -coverprofile coverage.out ./syslog
	go tool cover -html=coverage.out
