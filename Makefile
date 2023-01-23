KO_DOCKER_REPO ?= darkowlzz
TAG ?= latest

tidy:
	go mod tidy -v

ko-build:
	ko build --local ./cmd/gh-app-secret

ko-publish:
	KO_DOCKER_REPO=$(KO_DOCKER_REPO) ko build -B ./cmd/gh-app-secret -t $(TAG)
