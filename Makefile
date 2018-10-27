# service specific vars
SERVICE     := proxge
VERSION     := 0.0.1

TARGET      := ${SERVICE}
COMMIT      := $(shell git rev-parse --short HEAD)
BUILD_TIME  := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
IMAGE_NAME  := ${ORG}/${SERVICE}
PACKAGE 	:= $(shell grep module go.mod | awk '{ print $$2; }')

.PHONY: proto deps test build cont cont-nc all deploy help clean lint
.DEFAULT_GOAL := help

help: ## halp
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: deps lint test build ## get && test && build

build: clean lint ## build service binary file
	@echo "[build] building go binary"
	@go build \
		-ldflags "-s -w \
		-X ${PACKAGE}/pkg.Version=${VERSION} \
		-X ${PACKAGE}/pkg.Commit=${COMMIT} \
		-X ${PACKAGE}/pkg.BuildTime=${BUILD_TIME}" \
		-o ${GOPATH}/bin/${TARGET} ./cmd/${TARGET}
	@${TARGET} -v

clean: ## remove service bin from $GOPATH/bin
	@echo "[clean] removing service files"
	rm -f ${GOPATH}/bin/${TARGET}

deps: ## get service pkg + test deps
	@echo "[deps] getting go deps"
	go get -v -t ./...

lambda: lint ## build the lambda binary
	@echo "[lambda] building go binary"
	@GOOS=linux go build \
		-ldflags "-s -w \
		-X ${PACKAGE}/pkg.Version=${VERSION} \
		-X ${PACKAGE}/pkg.Commit=${COMMIT} \
		-X ${PACKAGE}/pkg.BuildTime=${BUILD_TIME}" \
		-o deploy/${TARGET} ./cmd/${TARGET}
	@zip -j -r deployment.zip deploy/*

lambda-deploy: lambda
	aws lambda update-function-code \
	--function-name proxgeFunction \
	--zip-file fileb://./deployment.zip

lint: ## apply golint
	@echo "[lint] applying go fmt & vet"
	go fmt ./...
	go vet ./...

test: lint ## test service code
	@echo "[test] running tests w/ cover"
	go test ./... -cover

release: clean test build lambda-deploy
