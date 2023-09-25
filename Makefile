.PHONY: clean build test run all lint mocks

BINARY_NAME=api-server

mocks:
	mockery --config .mockery.yaml

clean:
	go clean
	rm -f ./cmd/server/${BINARY_NAME}

build:
	go build -o ./cmd/server/${BINARY_NAME} ./cmd/server/*.go

test:
	go test -v ./...

test-integration:
	go test -v ./... --tags=integration

run: build
	./cmd/server/${BINARY_NAME} --debug

lint:
	golangci-lint run

docker-build:
	 docker build --rm -t api-server:latest .

docker-compose-up:
	docker-compose up --build

all: clean build test run