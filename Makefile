APP_NAME = tictactoe

BUILD_PATH = github.com/mgjules/tic-tac-toe/build

VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
DATE := $(shell date)

LDFLAGS := -s -w -X '$(BUILD_PATH).Version=$(VERSION)' -X '$(BUILD_PATH).Commit=$(COMMIT)' -X '$(BUILD_PATH).Branch=$(BRANCH)' -X '$(BUILD_PATH).Date=$(DATE)'

.PHONY: build

buildrun: build run

vet:
	go vet ./...

lint:
	golangci-lint run

test:
	go test ./... -cover

run:
	./$(APP_NAME)

build: vet lint test
	CGO_ENABLED=0 GOOS=linux go build -tags=jsoniter -ldflags="$(LDFLAGS)" -o $(APP_NAME)