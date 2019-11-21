all: lint install

GOBIN?=${GOPATH}/bin

API_REST_SPEC=./api/rest/v1/specification/openapi-v1.0.0.yml
API_REST_CODE_GEN_LOCATION=./api/rest/v1/codegen/openapi-v1.0.0.go
API_REST_DOCO_GEN_LOCATION=./public/rest/v1/api.html

bootstrap: node_modules ${GOPATH}/bin/oapi-codegen

.PHONY: config, tools, test

config:
	@echo GOBIN: ${GOBIN}
	@echo GOPATH: ${GOPATH}

# cli tool for openapi
${GOPATH}/bin/oapi-codegen:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen

# node_modules for API dev tools
node_modules:
	yarn

install: bootstrap go.sum build
	GO111MODULE=on go install -v ./cmd/midgard

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

pg:
	@docker-compose run --rm -p 5432:5432 --no-deps pg

# -------------------------------------------- API Targets ------------------------------------

# Open API Makefile targets
openapi3validate:
	./node_modules/.bin/oas-validate -v ${API_REST_SPEC}

oapi-codegen-server: openapi3validate
	@${GOBIN}/oapi-codegen --package=api --generate types,server,spec ${API_REST_SPEC} > ${API_REST_CODE_GEN_LOCATION}

doco:
	./node_modules/.bin/redoc-cli bundle ${API_REST_SPEC} -o ${API_REST_DOCO_GEN_LOCATION}

# -----------------------------------------------------------------------------------------

dev: migration-down migration-up
	go run cmd/midgard/main.go

run-in-docker:
	@${GOBIN}/midgard -c /etc/midgard/config.json

run:
	@${GOBIN}/midgard -c cmd/midgard/config.json

up:
	@docker-compose up --build

clean:
	@rm ${GOBIN}/midgard

run-mocked-endpoint:
	go run tools/mockServer/mockServer.go

# ------------------------------------------- sql migrations ----------------------------------------------

${GOBIN}/sql-migrate:
	go get -v github.com/rubenv/sql-migrate/...

create-database-test:
	PGPASSWORD=password psql -h localhost -U postgres -c "create database midgard_test;"

create-database:
	psql -h localhost -U postgres -c "create database midgard;"

drop-database:
	psql -h localhost -U postgres -c "drop database midgard;"

drop-database-test:
	PGPASSWORD=password psql -h localhost -U postgres -c "drop database midgard_test;"


migration-up: ${GOBIN}/sql-migrate
	${GOBIN}/sql-migrate up

migration-down: ${GOBIN}/sql-migrate
	${GOBIN}/sql-migrate down

migration-up-test: ${GOBIN}/sql-migrate
	${GOBIN}/sql-migrate up --env="test"

migration-down-test: ${GOBIN}/sql-migrate
	${GOBIN}/sql-migrate down --env="test"

run-test-suite:
	@make drop-database-test
	@make create-database-test
	@make migration-up-test