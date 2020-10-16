.PHONY: build

buildrun: lint build run

lint:
	@golangci-lint run

run:
	@./ringier-test

build:
	@go build -ldflags="-s -w"