.PHONY: check clean build image publish publish-latest test

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse --short HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
BUILD_DATE := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
DOCKER_REGISTRY := gcr.io
DOCKER_REPOSITORY := trois-six/geneparse

default: clean build render

check:
	@golangci-lint run

clean:
	@rm -rf $(OUTPUT_DIR)

build: clean
	@echo Version: $(VERSION) $(BUILD_DATE)
	protoc --go_out=pkg/geneanet pkg/geneanet/api.proto
	CGO_ENABLED=0 go build -v -ldflags '-X "main.version=${VERSION}" -X "main.commit=${SHA}" -X "main.date=${BUILD_DATE}"'

image:
	docker build -t $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY):$(VERSION) .

publish:
	docker push $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY):$(VERSION)

publish-latest:
	docker tag $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY):$(VERSION) $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY):latest
	docker push $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY):latest

test: clean
	go test -v -cover ./...
