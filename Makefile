DOCKER_ORG ?= arschles
DOCKER_VERSION ?= $(shell git rev-parse --short HEAD)
DOCKER_IMAGE_NAME ?= megaboom
DOCKER_IMAGE := quay.io/${DOCKER_ORG}/${DOCKER_IMAGE_NAME}:${DOCKER_VERSION}

TEST_NUM_PODS ?= 1
TEST_NUM_CONCURRENT_PER_POD ?= 500
TEST_NUM_REQS_PER_POD ?= 10000
APP_NAME ?= gotest
DEIS_ROUTER_IP ?= 104.154.76.205

build-binary:
	GOOS=linux GOARCH=amd64 go build -o rootfs/bin/megaboom
docker-build:
	docker build -t ${DOCKER_IMAGE} rootfs
docker-push:
	docker push ${DOCKER_IMAGE}
deis-deploy:
	deis pull ${DOCKER_IMAGE} -a megaboom
test-live:
	curl -XPOST -d '{"num_pods":${TEST_NUM_PODS}, "num_concurrent_per_pod":${TEST_NUM_CONCURRENT_PER_POD}, "num_reqs_per_pod":${TEST_NUM_REQS_PER_POD}, "http_method":"GET", "endpoint":"http://${APP_NAME}.${DEIS_ROUTER_IP}.nip.io","namespace":"default","image":"quay.io/arschles/boom:0.1.0"}' http://megaboom.${DEIS_ROUTER_IP}.nip.io/job
