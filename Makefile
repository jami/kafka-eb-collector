GO=go
RACE := $(shell test $$(go env GOARCH) != "amd64" || (echo "-race"))
GOFLAGS= 
VERSION := $(shell git rev-parse HEAD)
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
PROJECT=kafka-eb-collector
COLLECTOR_BIN=bin/kafka-collector
COLLECTOR_IMAGE=jami/kafka-collector
COLLECTOR_SOURCE=./src/cli/main.go

all: build/local

help:
	@echo 'Available commands:'
	@echo
	@echo 'Usage:'
	@echo '    make deps     		Install go deps.'
	@echo '    make build/local    	Compile the project.'
	@echo '    make test/local    	Run ginkgo test suites.'
	@echo '    make build/docker    Create docker container'
	@echo '    make clean    		Clean the directory tree.'
	@echo

deps:
	go mod vendor

build/collector:
	CGO_ENABLED=1 GOOS=linux ${GO} build -a -tags musl -o ${COLLECTOR_BIN} ${COLLECTOR_SOURCE}

build/collector/local:
	$(GO) build -ldflags "-X main.Version=$(VERSION)" -o $(COLLECTOR_BIN) $(GOFLAGS) $(RACE) $(COLLECTOR_SOURCE)

build/collector/linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-X main.Version=$(VERSION)" -o "$(COLLECTOR_BIN)_linux" $(GOFLAGS) $(COLLECTOR_SOURCE)

build/collector/docker:
	docker build -t $(COLLECTOR_IMAGE) -f ./docker/collector/Dockerfile .

rebuild/service/collector:
	docker-compose stop kafka-collector
	docker-compose rm -f kafka-collector
	docker-compose build kafka-collector
	docker-compose create kafka-collector
	docker-compose start kafka-collector

compose/up: build/linux
	docker-compose -f docker-compose.yml up -d

compose/down:
	docker-compose -f docker-compose.yml down
