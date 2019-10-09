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
RUN git config --global http.sslVerify "false"
RUN git clone https://gitlab.com/thorchain/bepswap/chain-service.git /tmp/chainservice

# Install jq to update the chain service config.
RUN apt-get update
RUN apt-get install -y jq

# Generate config.
WORKDIR /tmp/chainservice
RUN mkdir -p /etc/chainservice
RUN cat ./cmd/chainservice/config.json | jq \
  --arg CHAIN_HOST "$CHAIN_HOST" \
  --arg INFLUX_HOST "$INFLUX_HOST" \
  '.influx["host"] = $INFLUX_HOST | \
  .statechain["host"] = $CHAIN_HOST' > /etc/chainservice/config.json
RUN cat /etc/chainservice/config.json

# Compile.
RUN GO111MODULE=on go mod verify

WORKDIR /tmp/chainservice/cmd/chainservice
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chainservice .

#
# Main
#
FROM golang:alpine

ENV PATH="${PATH}:/go/bin"

RUN apk add --update busybox-extras && \
    rm -rf /var/cache/apk/*

# Copy the compiled binaires over.
COPY --from=build /tmp/chainservice/cmd/chainservice/chainservice /go/bin/chainservice

# Copy the chain service config.
COPY --from=build /etc/chainservice /etc/chainservice

# Copy the chain service public folder ie generated docs
COPY --from=build /tmp/chainservice/public /go/bin/public

EXPOSE 8080

CMD ["chainservice", "-c", "/etc/chainservice/config.json"]
