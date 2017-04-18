#!/usr/bin/env bash
#
# Build and push Docker images to Docker Hub and quay.io.
#
cd "$(dirname "$0")" || exit 1
echo "Building docker image and pushing to quay.io!"
DEIS_REGISTRY=quay.io/ make -C .. build docker-push
echo "Building docker image and pushing to docker hub!"
make -C .. build docker-push
