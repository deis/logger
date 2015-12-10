SHELL = /bin/bash

ifndef BUILD_TAG
  BUILD_TAG = git-$(shell git rev-parse --short HEAD)
endif

GO = godep go
GOFMT = gofmt -l
GOLINT = golint
GOTEST = $(GO) test --cover --race -v
GOVET = $(GO) vet
GO_FILES = $(wildcard *.go)
GO_PACKAGES = drain storage syslogish weblog
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))
GO_TESTABLE_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,drain drain/simple storage storage/file storage/ringbuffer)

BINARY_DEST_DIR = image/bin

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/deis/logger

DOCKER_HOST = $(shell echo $$DOCKER_HOST)
REGISTRY = $(shell if [ "$$DEV_REGISTRY" == "registry.hub.docker.com" ]; then echo; else echo $$DEV_REGISTRY/; fi)

COMPONENT = logger
IMAGE_PREFIX = deis/
IMAGE = $(IMAGE_PREFIX)$(COMPONENT):$(BUILD_TAG)
DEV_IMAGE = $(REGISTRY)$(IMAGE)

check-docker:
	@if [ -z $$(which docker) ]; then \
	  echo "Missing \`docker\` client which is required for development"; \
	  exit 2; \
	fi

dev-registry: check-docker
	@docker inspect registry >/dev/null 2>&1 && docker start registry || docker run --restart="always" -d -p 5000:5000 --name registry registry:0.9.1
	@echo
	@echo "To use a local registry for Deis development:"
	@echo "    export DEV_REGISTRY=`docker-machine ip $$(docker-machine active 2>/dev/null) 2>/dev/null || echo $(HOST_IPADDR) `:5000"

build: build-binary docker-build

push: build docker-push

docker-build: check-docker
	docker build -t $(IMAGE) image

docker-push: update-manifests
		docker tag -f $(IMAGE) $(DEV_IMAGE)

update-manifests:
	sed 's#\(image:\) .*#\1 $(DEV_IMAGE)#' manifests/deis-logger-rc.yaml > manifests/deis-logger-rc.tmp.yaml

build-binary:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' -o $(BINARY_DEST_DIR)/logger github.com/deis/logger || exit 1
	@$(call check-static-binary,$(BINARY_DEST_DIR)/logger)

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
