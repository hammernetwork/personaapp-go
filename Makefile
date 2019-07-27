#Go parameters
GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=build/persona

build: 
	$(GOBUILD) -o $(BINARY_NAME) -v
run:
	$(GORUN) main.go