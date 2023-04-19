KO_DOCKER_REPO ?= ghcr.io/fluxcd-community
TAG ?= latest

tidy:
	go mod tidy -v

build:
	go build -o bin/github-app-secret ./cmd/github-app-secret

test:
	go test -v ./...

ko-build:
	KO_DOCKER_REPO=ko.local ko build ./cmd/github-app-secret

ko-publish:
	KO_DOCKER_REPO=$(KO_DOCKER_REPO) ko build -B ./cmd/github-app-secret -t $(TAG)
