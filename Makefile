SHELL = /bin/bash
GO = go
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
REGISTRY = $(shell if [ -z $$REGISTRY ]; then echo deis/; else echo $$REGISTRY/; fi)
ifndef VERSION
  VERSION = git-$(shell git rev-parse --short HEAD)
endif
COMPONENT = logger
IMAGE = $(REGISTRY)$(COMPONENT):$(VERSION)

check-docker:
	@if [ -z $$(which docker) ]; then \
	  echo "Missing docker client which is required for development"; \
	  exit 2; \
	fi

build: build-binary docker-build

push: docker-push

docker-build: check-docker
	docker build -t $(IMAGE) image

docker-push: check-docker
	docker push $(IMAGE)

clean: check-docker
	docker rmi $(IMAGE)

update-manifests:
	sed 's#\(image:\) .*#\1 $(IMAGE)#' manifests/deis-logger-rc.yaml > manifests/deis-logger-rc.tmp.yaml

build-binary:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' -o $(BINARY_DEST_DIR)/logger github.com/deis/logger || exit 1
	@$(call check-static-binary,$(BINARY_DEST_DIR)/logger)

test: test-style test-unit

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

kube-delete:
	-kubectl delete -f manifests/deis-logger-rc.yaml
	-kubectl delete -f manifests/deis-logger-svc.yaml
	-kubectl delete -f manifests/deis-logger-fluentd-daemon.yaml

kube-create: update-manifests
	kubectl create -f manifests/deis-logger-rc.tmp.yaml
	kubectl create -f manifests/deis-logger-svc.yaml
	kubectl create -f manifests/deis-logger-fluentd-daemon.yaml

kube-replace: build push update-manifests
	kubectl replace --force -f manifests/deis-logger-rc.tmp.yaml
	kubectl replace --force -f manifests/deis-logger-svc.yaml
	kubectl replace --force -f manifests/deis-logger-fluentd-daemon.yaml
