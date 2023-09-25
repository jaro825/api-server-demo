.PHONY: clean build test run all lint

BINARY_NAME=api-server

clean:
	go clean
	rm -f ./cmd/proxy/${BINARY_NAME}

build:
	go build -o ./cmd/proxy/${BINARY_NAME} ./cmd/proxy/*.go

test:
	go test -v ./...

test-integration:
	go test -v ./... --tags=integration

run: build
	./cmd/proxy/${BINARY_NAME}

lint:
	golangci-lint run

docker-build:
	 docker build --rm -t api-server:latest .

docker-compose-up:
	docker-compose up --build

all: clean build test run