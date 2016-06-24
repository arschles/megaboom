DOCKER_ORG ?= arschles
DOCKER_VERSION ?= $(shell git rev-parse --short HEAD)
DOCKER_IMAGE_NAME ?= megaboom
DOCKER_IMAGE := quay.io/${DOCKER_ORG}/${DOCKER_IMAGE_NAME}:${DOCKER_VERSION}

build-binary:
	GOOS=linux GOARCH=amd64 go build -o rootfs/bin/megaboom
docker-build:
	docker build -t ${DOCKER_IMAGE} rootfs
docker-push:
	docker push ${DOCKER_IMAGE}
