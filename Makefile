BINARY_NAME := resize.exe
BINARY_DIR := .\build
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt

CP=$(shell cp)
DIST=\\192.168.101.66\share\tools

all: build test

build:
	$(GOBUILD) -o $(BINARY_DIR)\$(BINARY_NAME) -v .

deploy: build
	copy $(BINARY_DIR)\$(BINARY_NAME) $(DIST)\$(BINARY_NAME)

test: build
	.\$(BINARY_DIR)\$(BINARY_NAME)

clean:
	$(GOCLEAN)

format:
	$(GOFMT) ./...

.PHONY: build deploy test clean format