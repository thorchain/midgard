all: lint install

API_SPEC=./api/rest/v1/specification/openapi-v1.0.0.yml
API_CODE_GEN_LOCATION=./api/rest/v1/codegen/openapi-v1.0.0.go
API_DOCO_GEN_LOCATION=./public/rest/v1/api.html
GOBIN=./bin

install: go.sum
	# cli tool for openapi
	npm i
	go get -u github.com/deepmap/oapi-codegen/cmd/oapi-codegen
	GO111MODULE=on go install -v ./cmd/etl

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

dev: build
	./bin/chainservice-api-v1

build: clean oapi-codegen-server doco
	@go build -o ./bin/chainservice-api-v1 ./cmd/chainservice-api-v1/main.go

clean:
	rm -rf ${GOBIN}/*

test-coverage:
	@go test -mod=readonly -v -coverprofile .testCoverage.txt ./...

coverage-report: test-coverage
	@tool cover -html=.testCoverage.txt

test-short:
	@go test -short ./...

test-internal:
	@go test ./...

test:
	@./scripts/run.sh etl make test-internal

clear:
	clear

test-watch: clear
	@./scripts/watch.bash

sh:
	@docker-compose run --rm --no-deps etl /bin/sh

influxdb:
	@docker-compose run --rm -p 8086:8086 --no-deps influxdb


# -------------------------------------------- API Targets ------------------------------------

# Open API Makefile targets
openapi3validate:
	oas-validate -v ${API_SPEC}

# TODO Setup auto versioning outputter
oapi-codegen-server: openapi3validate
	oapi-codegen --package=api --generate types,server,spec ${API_SPEC} > ${API_CODE_GEN_LOCATION}

doco:
	redoc-cli bundle ${API_SPEC} -o ${API_DOCO_GEN_LOCATION}


# -----------------------------------------------------------------------------------------
