BASE_PATH = $(shell pwd)
export PATH := $(BASE_PATH)/bin:$(PATH)
export GOBIN := $(BASE_PATH)/bin

SHELL := env PATH=$(PATH) /bin/bash

# Commands
GOCMD=go
GORUN=$(GOCMD) run
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

# GRPC
PROTOC       = protoc
PROTOCGOGEN  = protoc-gen-go

BIN_DIR=./bin
BINARY_NAME=personaapp
BINARY_PATH=$(BIN_DIR)/$(BINARY_NAME)

# all the packages without vendor
ALL_PKGS = $(shell go list ./... | grep -v /vendor | grep -v /pkg/grpcapi)

# Colors
GREEN_COLOR   = "\033[0;32m"
PURPLE_COLOR  = "\033[0;35m"
DEFAULT_COLOR = "\033[m"

.PHONY: all help clean test lint fmt build grpc

all: clean fmt build lint test

help:
	@echo 'Usage: make <TARGETS> ... <OPTIONS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               Show this help screen.'
	@echo '    clean              Clean all the artifacts.'
	@echo '    test               Run unit tests.'
	@echo '    lint               Run all linters.'
	@echo '    fmt                Run gofmt on package sources.'
	@echo '    build              Compile packages and dependencies.'
	@echo '    grpc               Generate gRPC code'
	@echo '    generate           Generate mocks'
	@echo ''

clean:
	@echo -e [$(GREEN_COLOR)clean$(DEFAULT_COLOR)]
	@$(GOCLEAN)
	@rm -rf $(BIN_DIR)

test:
	@echo -e [$(GREEN_COLOR)test$(DEFAULT_COLOR)]
	# TODO: -race flag seems not working in Alpine, so it disabled for now. Ping me to investigate further. Maksym Hilliaka
	@$(GOTEST) -v -count=1 ./...

lint:
	@echo -e [$(GREEN_COLOR)lint$(DEFAULT_COLOR)]
	@$(GOINSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint
	@$(GOBIN)/golangci-lint run

fmt:
	@echo -e [$(GREEN_COLOR)format$(DEFAULT_COLOR)]
	@$(GOFMT) $(ALLPKGS)

build:
	@echo -e [$(GREEN_COLOR)build$(DEFAULT_COLOR)]
	@$(GOBUILD) -o $(BINARY_PATH)

grpc:
	@echo -e $(GREEN_COLOR)[grpc]$(DEFAULT_COLOR)
	@$(GOINSTALL) github.com/golang/protobuf/protoc-gen-go

	@-rm -rf ./pkg/grpcapi
	@mkdir -p ./pkg/grpcapi/vacancy ./pkg/grpcapi/auth ./pkg/grpcapi/company

	@${PROTOC} \
        -I ./api \
        ./api/vacancy/vacancy.proto \
        --go_out=plugins=grpc:./pkg/grpcapi

	@${PROTOC} \
        -I ./api \
        ./api/auth/auth.proto \
        --go_out=plugins=grpc:./pkg/grpcapi

	@${PROTOC} \
        -I ./api \
        ./api/company/company.proto \
        --go_out=plugins=grpc:./pkg/grpcapi

generate:
	@mkdir -p ./bin
	@echo -e $(PURPLE_COLOR)[building mockery]$(DEFAULT_COLOR)
	@$(GOINSTALL) github.com/vektra/mockery/cmd/mockery
	@echo -e $(PURPLE_COLOR)[mockery built]$(DEFAULT_COLOR)
	@echo -e [$(GREEN_COLOR)generate$(DEFAULT_COLOR)]
	@$(GOGENERATE) $(PKGS)
