FROM golang:alpine

ARG chain_host
ARG pg_host

ENV CHAIN_HOST=$chain_host
ENV PG_HOST=$pg_host

RUN apk update && \
    apk add python python-dev py-pip build-base && \
    apk add curl make git linux-headers jq yarn && \
    pip install dumb-init && \
    rm -rf /var/cache/apk/*

COPY . /tmp/midgard
WORKDIR /tmp/midgard

ENV PATH="node_modules/.bin:${PATH}"
RUN env

RUN mkdir -p /etc/midgard
RUN cat ./cmd/midgard/config.json | jq \
  --arg CHAIN_HOST "$CHAIN_HOST" \
  --arg PG_HOST "$PG_HOST" \
  '.timescale["host"] = $PG_HOST | \
  .thorchain["host"] = $CHAIN_HOST' > /etc/midgard/config.json
RUN cat /etc/midgard/config.json

ENTRYPOINT ["dumb-init"]
CMD ["/bin/sh"]
