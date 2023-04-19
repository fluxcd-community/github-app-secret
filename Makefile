KO_DOCKER_REPO ?= ghcr.io/fluxcd-community
TAG ?= latest

all: build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: build
build: fmt vet
	go build -o bin/github-app-secret ./cmd/github-app-secret

.PHONY: test
test:
	go test -v ./...

.PHONY: ko-build
ko-build:
	KO_DOCKER_REPO=ko.local ko build -B ./cmd/github-app-secret

.PHONY: ko-publish
ko-publish:
	KO_DOCKER_REPO=$(KO_DOCKER_REPO) ko build -B ./cmd/github-app-secret -t $(TAG)
