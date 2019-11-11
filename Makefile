all: lint install

GOBIN?=${GOPATH}/bin

API_REST_SPEC=./api/rest/v1/specification/openapi-v1.0.0.yml
API_REST_CODE_GEN_LOCATION=./api/rest/v1/codegen/openapi-v1.0.0.go
API_REST_DOCO_GEN_LOCATION=./public/rest/v1/api.html

bootstrap: node_modules ${GOPATH}/bin/oapi-codegen

.PHONY: config
config:
	@echo GOBIN: ${GOBIN}
	@echo GOPATH: ${GOPATH}

.PHONY: tools

# cli tool for openapi
${GOPATH}/bin/oapi-codegen:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen

# node_modules for API dev tools
node_modules:
	yarn

install: bootstrap go.sum build
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

build: oapi-codegen-server doco

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

influx_stack:
	@docker-compose run --rm -p 8888:8888 chronograf

influxdb:
	@docker-compose run --rm -p 8086:8086 --no-deps influxdb

# -------------------------------------------- API Targets ------------------------------------

# Open API Makefile targets
openapi3validate:
	./node_modules/.bin/oas-validate -v ${API_REST_SPEC}

oapi-codegen-server: openapi3validate
	@${GOBIN}/oapi-codegen --package=api --generate types,server,spec ${API_REST_SPEC} > ${API_REST_CODE_GEN_LOCATION}

doco:
	./node_modules/.bin/redoc-cli bundle ${API_REST_SPEC} -o ${API_REST_DOCO_GEN_LOCATION}

# -----------------------------------------------------------------------------------------

run-in-docker:
	@${GOBIN}/chainservice -c /etc/chainservice/config.json

run:
	@${GOBIN}/chainservice -c cmd/chainservice/config.json

up:
	@docker-compose up --build

clean:
	@rm ${GOBIN}/chainservice

run_mocked_endpint:
	go run tools/mockServer/mockServer.go
