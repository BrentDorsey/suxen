SERVICE=suxen
VERSION=$(shell git describe --tags --always --dirty)
ifeq ($(version),)
	TAG=${VERSION}
else
	TAG=$(version)
endif

PACKAGE=github.com/travelaudience/${SERVICE}
PACKAGE_CMD_DAEMON=${PACKAGE}/cmd/${SERVICE}d
PACKAGES=$(shell go list ./... | grep -v /vendor/)
DOCKER_HUB=quay.io/travelaudience

LDFLAGS = -X 'main.version=$(VERSION)'

.PHONY: build install test get publish yarn yarn-build

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -a -o bin/${SERVICE}d ${PACKAGE_CMD_DAEMON}
	CGO_ENABLED=0 GOOS=darwin go build -ldflags "${LDFLAGS}" -a -o bin/${SERVICE}d_osx ${PACKAGE_CMD_DAEMON}

install:
	go install -ldflags "${LDFLAGS}" ${PACKAGE_CMD_DAEMON}

test:
	go test ./...

get:
	go get -u github.com/golang/dep/cmd/dep
	go get -t golang.org/x/tools/cmd/goimports/...
	go get github.com/vektra/mockery/cmd/mockery
	dep ensure

publish: build yarn-build
ifneq ($(skiplogin),true)
	docker login ${DOCKER_HUB}
endif
	docker build -t ${DOCKER_HUB}/${SERVICE}:${TAG} .
	docker push ${DOCKER_HUB}/${SERVICE}:${TAG}

yarn:
	cd ui; yarn

yarn-build: yarn
	cd ui; yarn run build
