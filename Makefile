.DEFAULT_GOAL := docker_build

.PHONY: build docker_build

SERVICE=statsWriter
GO_PATH_SERVICE_MAIN=./
GO_SERVICE_IMPORT_PATH_SRC=$(shell go list ./src)
GO_SERVICE_IMPORT_PATH=$(GO_SERVICE_IMPORT_PATH_SRC:/src=)
PATH_DOCKER_SERVICE_SOURCES=/go/src/$(GO_SERVICE_IMPORT_PATH)
PATH_DOCKER_FILE=$(realpath ./build/Dockerfile)
GO111MODULE=on

GO_BUILD_LDFLAGS= -extldflags '-static'
GO_BUILD_FLAGS=-v --ldflags "$(GO_BUILD_LDFLAGS)"

build:
	@echo '>>> Building go binary'
	CGO_ENABLED=0 go build $(GO_BUILD_FLAGS) -o $(SERVICE) $(GO_PATH_SERVICE_MAIN)

docker_build:
	@echo ">>> Building docker image"
	docker build \
		-t $(SERVICE) \
		--build-arg GIT_BRANCH=$(GO_SERVICE_BUILD_BRANCH) \
		--build-arg GO_SERVICE_IMPORT_PATH=$(PATH_DOCKER_SERVICE_SOURCES) \
		--build-arg GO111MODULE=$(GO111MODULE) \
		-f $(PATH_DOCKER_FILE) \
		.
