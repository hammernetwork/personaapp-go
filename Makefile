BASE_PATH = $(shell pwd)
export PATH := $(BASE_PATH)/bin:$(PATH)

# Commands
GOCMD=go
GORUN=$(GOCMD) run
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
GREEN_COLOR   = \033[0;32m
PURPLE_COLOR  = \033[0;35m
DEFAULT_COLOR = \033[m

.PHONY: all help clean test lint fmt build grpc

all: clean fmt build lint

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
	@echo " [$(GREEN_COLOR)clean$(DEFAULT_COLOR)]"
	@$(GOCLEAN)
	@rm -rf $(BIN_DIR)

test:
	@echo " [$(GREEN_COLOR)test$(DEFAULT_COLOR)]"
	@$(GOTEST) -race -v -count=1 ./...

lint:
	@echo " [$(GREEN_COLOR)lint$(DEFAULT_COLOR)]"
	@$(GORUN) ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint/main.go run --issues-exit-code 1

fmt:
	@echo " [$(GREEN_COLOR)format$(DEFAULT_COLOR)]"
	@$(GOFMT) $(ALLPKGS)

build:
	@echo " [$(GREEN_COLOR)build$(DEFAULT_COLOR)]"
	@$(GOBUILD) -o $(BINARY_PATH)

grpc:
	@mkdir -p ./bin
ifeq ("$(wildcard ./bin/$(PROTOCGOGEN))","")
	@echo " $(PURPLE_COLOR)[build protoc-go-gen]$(DEFAULT_COLOR)"
	@$(GOBUILD) -o ./bin/$(PROTOCGOGEN) ./vendor/github.com/golang/protobuf/protoc-gen-go/
endif

	@echo " [$(GREEN_COLOR)grpc$(DEFAULT_COLOR)]"
		@-rm -rf ./pkg/grpcapi
		@mkdir -p ./pkg/grpcapi/vacancy ./pkg/grpcapi/auth ./pkg/grpcapi/company

	@${PROTOC} \
		-I ./api \
		./api/vacancy/*.proto \
		--go_out=plugins=grpc:..

	@${PROTOC} \
        -I ./api \
        ./api/auth/*.proto \
        --go_out=plugins=grpc:..

	@${PROTOC} \
        -I ./api \
        ./api/company/*.proto \
        --go_out=plugins=grpc:..

generate:
	@mkdir -p ./bin
ifeq ("$(wildcard ./bin/$(MOCKERY))","")
	@echo " $(PURPLE_COLOR)[build mockery]$(DEFAULT_COLOR)"
	@$(GOBUILD) -o ./bin/$(MOCKERY) ./vendor/github.com/vektra/mockery/cmd/mockery/
endif
	@echo " [$(GREEN_COLOR)generate$(DEFAULT_COLOR)]"
	@$(GOGENERATE) $(PKGS)
