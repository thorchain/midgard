#!/bin/sh

docker-compose run -d --rm -p 5432:5432 pg
go test -cover ./...
exit_code=$?
docker-compose down
exit $exit_code
