KO_DOCKER_REPO ?= ghcr.io/fluxcd-community
TAG ?= latest

tidy:
	go mod tidy -v

ko-build:
	ko build --local ./cmd/github-app-secret

ko-publish:
	KO_DOCKER_REPO=$(KO_DOCKER_REPO) ko build -B ./cmd/github-app-secret -t $(TAG)
