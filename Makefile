APP_NAME = tictactoe

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
	CGO_ENABLED=0 GOOS=linux go build -tags=jsoniter -ldflags="-s -w" -o $(APP_NAME)

default: build