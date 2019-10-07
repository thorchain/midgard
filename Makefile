all: lint install

GOBIN?=${GOPATH}/bin

API_REST_SPEC=./api/rest/v1/specification/openapi-v1.0.0.yml
API_REST_CODE_GEN_LOCATION=./api/rest/v1/codegen/openapi-v1.0.0.go
API_REST_DOCO_GEN_LOCATION=./public/rest/v1/api.html

bootstrap: node_modules ${GOPATH}/bin/oapi-codegen build install

.PHONY: config
config:
	@echo GOBIN: ${GOBIN}
	@echo GOPATH: ${GOPATH}

# cli tool for openapi
${GOPATH}/bin/oapi-codegen:
	go get -u github.com/deepmap/oapi-codegen/cmd/oapi-codegen

# node_modules for API dev tools
node_modules:
	yarn

install: go.sum
	GO111MODULE=on go install -v ./cmd/chainservice
	GO111MODULE=on go install -v ./cmd/chainservice-api-v1

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

build: oapi-codegen-server doco
	@go build ./...

test-coverage:
	@go test -mod=readonly -v -coverprofile .testCoverage.txt ./...

coverage-report: test-coverage
	@tool cover -html=.testCoverage.txt

test-short:
	@go test -short ./...

test-internal:
	@go test -cover ./...

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


# -------------------------------------------- API Targets ------------------------------------

# Open API Makefile targets
openapi3validate:
	oas-validate -v ${API_REST_SPEC}

oapi-codegen-server: openapi3validate
	oapi-codegen --package=api --generate types,server,spec ${API_REST_SPEC} > ${API_REST_CODE_GEN_LOCATION}

doco:
	redoc-cli bundle ${API_REST_SPEC} -o ${API_REST_DOCO_GEN_LOCATION}

# -----------------------------------------------------------------------------------------

run-in-docker:
	@${GOBIN}/chainservice -c /etc/chainservice/config.json

run:
	@${GOBIN}/chainservice -c cmd/chainservice/config.json

run-api-v1:
	@${GOBIN}/chainservice-api-v1 # -c cmd/chainservice/config.json

up:
	@docker-compose up --build
