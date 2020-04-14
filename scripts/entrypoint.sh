#!/bin/sh

if [ ! -z "$EXTERNAL_IP" ]; then
    openssl ecparam -genkey -name secp384r1 -out server.key
    openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650 -batch
fi

midgard -c /etc/midgard/config.json

exec "$@"
