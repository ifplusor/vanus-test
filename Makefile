GOOS ?=
GOARCH ?=

GO_BUILD_ENV = GO111MODULE=on CGO_ENABLED=0
ifdef GOOS
	GO_BUILD_ENV := $(GO_BUILD_ENV) GOOS=$(GOOS) 
endif
ifdef GOARCH
	GO_BUILD_ENV := $(GO_BUILD_ENV) GOARCH=$(GOARCH)
endif

GO_BUILD_FLAGS = -trimpath

GO_BUILD = $(GO_BUILD_ENV) go build $(GO_BUILD_FLAGS)

.PHONY: build
build: build-vsbench build-integration

.PHONY: build-vsbench
build-vsbench:
	$(GO_BUILD) -o bin/vsbench vsbench/main.go

.PHONY: build-integration
build-integration:
	$(GO_BUILD) -o bin/vanus-integration integration/main.go


CONTAINER_REGISTRY ?= public.ecr.aws
CONTAINER_REPO ?= ${CONTAINER_REGISTRY}/vanus
IMAGE_TAG ?= latest
DOCKER_PLATFORM ?= linux/amd64

.PHONY: docker-build
docker-build: build
	docker build -t $(CONTAINER_REPO)/vanus-test:$(IMAGE_TAG) .

.PHONY: docker-push
docker-push: build
	docker buildx build --platform $(DOCKER_PLATFORM) -t $(CONTAINER_REPO)/vanus-test:$(IMAGE_TAG) . --push

.PHONY: docker-push2
docker-push2:
	docker buildx build --platform $(DOCKER_PLATFORM) -t $(CONTAINER_REPO)/vanus-test:$(IMAGE_TAG) -f Dockerfile.build .
