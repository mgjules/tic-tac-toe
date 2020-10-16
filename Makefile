APP_NAME = tictactoe

.PHONY: build

buildrun: lint build run

lint:
	golangci-lint run

run:
	./$(APP_NAME)

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w"

default: build