#!/bin/sh

if [ ! -z "$EXTERNAL_IP" ]; then
    go run $GOROOT/src/crypto/tls/generate_cert.go --host $EXTERNAL_IP
fi

exec "$@"
