#
# Midgard
#

#
# Build
#
FROM golang:1.15 AS build

ARG pg_host
ARG rpc_host
ARG thornode_host

ENV PG_HOST=$pg_host
ENV RPC_HOST=$rpc_host
ENV THORNODE_HOST=$thornode_host

RUN env

WORKDIR /tmp/midgard

COPY  . .

# Install jq to update the chain service config.
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update
RUN apt-get install -y jq apt-utils make yarn

RUN make node_modules
RUN make ${GOPATH}/bin/oapi-codegen
# Make sure swagger.json is updated
RUN make oapi-codegen-server
# Generate api document
RUN make doco

# Generate config.
RUN mkdir -p /etc/midgard
RUN cat ./cmd/midgard/config.json | jq \
  --arg RPC_HOST "$RPC_HOST" \
  --arg THORNODE_HOST "$THORNODE_HOST" \
  --arg PG_HOST "$PG_HOST" \
  '.timescale["host"] = $PG_HOST | \
  .timescale["migrationsDir"] = "/var/midgard/db/migrations/" | \
  .thorchain["rpc_host"] = $RPC_HOST | \
  .thorchain["host"] = $THORNODE_HOST' > /etc/midgard/config.json
RUN cat /etc/midgard/config.json

# Compile.
RUN GO111MODULE=on go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o midgard /tmp/midgard/cmd/midgard

#
# Main
#
FROM golang:alpine

ENV PATH="${PATH}:/go/bin"

RUN apk update
RUN apk add make openssl bind-tools curl

COPY --from=build /tmp/midgard/ .

# Copy the db migrations
COPY --from=build /tmp/midgard/db/ /var/midgard/db/

# Copy the compiled binaires over.
COPY --from=build /tmp/midgard/midgard /go/bin/midgard

# Copy the chain service config.
COPY --from=build /etc/midgard /etc/midgard

# Copy the chain service public folder ie generated docs
COPY --from=build /tmp/midgard/public/ /go/public/

CMD [ "midgard", "-c", "/etc/midgard/config.json" ]
