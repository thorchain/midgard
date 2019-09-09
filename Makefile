GOBIN?=$$GOPATH/bin

all: lint install

install: go.sum
	GO111MODULE=on go install -v ./cmd/etl
	GO111MODULE=on go install -v ./cmd/chainservice

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

lint-pre:
	@test -z $(gofmt -l .) # checks code is in proper format
	@go mod verify

lint: lint-pre
	@golangci-lint run

lint-verbose: lint-pre
	@golangci-lint run -v

build:
	@go build ./...

clean:
	rm ${GOBIN}/etl

test-coverage:
	@go test -mod=readonly -v -coverprofile .testCoverage.txt ./...

coverage-report: test-coverage
	@tool cover -html=.testCoverage.txt

test-short:
	@go test -short ./...

test-internal:
	@go test ./...

test:
	@./scripts/run.sh chain-service make test-internal

clear:
	clear

test-watch: clear
	@./scripts/watch.bash

sh:
	@docker-compose run --rm chain-service /bin/sh

influxdb:
	@docker-compose run --rm -p 8086:8086 --no-deps influxdb

run:
	@${GOBIN}/chainservice -c cmd/chainservice/config.json

up:
	@docker-compose up
