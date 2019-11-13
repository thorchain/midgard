#
# Chain Service
#

#
# Build
#
FROM golang:1.13 AS build

ARG chain_host
ARG influx_host

ENV CHAIN_HOST=$chain_host
ENV INFLUX_HOST=$influx_host

RUN env

WORKDIR /tmp
# TODO Move away from the embed clone as this will only work for master branchs
RUN git config --global http.sslVerify "false"
RUN git clone https://gitlab.com/thorchain/bepswap/chain-service.git /tmp/midgard

# Install jq to update the chain service config.
RUN apt-get update
RUN apt-get install -y jq

# Generate config.
WORKDIR /tmp/midgard
RUN mkdir -p /etc/midgard
RUN cat ./cmd/midgard/config.json | jq \
  --arg CHAIN_HOST "$CHAIN_HOST" \
  --arg INFLUX_HOST "$INFLUX_HOST" \
  '.influx["host"] = $INFLUX_HOST | \
  .thorchain["host"] = $CHAIN_HOST' > /etc/midgard/config.json
RUN cat /etc/midgard/config.json

# Compile.
RUN GO111MODULE=on go mod verify

WORKDIR /tmp/midgard/cmd/midgard
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o midgard .

#
# Main
#
FROM golang:alpine

ENV PATH="${PATH}:/go/bin"

RUN apk add --update busybox-extras && \
    rm -rf /var/cache/apk/*

# Copy the compiled binaires over.
COPY --from=build /tmp/midgard/cmd/midgard/midgard /go/bin/midgard

# Copy the chain service config.
COPY --from=build /etc/midgard /etc/midgard

# Copy the chain service public folder ie generated docs
COPY --from=build /tmp/midgard/public/ /go/public/

EXPOSE 8080

CMD ["midgard", "-c", "/etc/midgard/config.json"]
