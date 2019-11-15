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

COPY . /tmp/chainservice
WORKDIR /tmp/chainservice

ENV PATH="node_modules/.bin:${PATH}"
RUN env

RUN mkdir -p /etc/chainservice
RUN cat ./cmd/chainservice/config.json | jq \
  --arg CHAIN_HOST "$CHAIN_HOST" \
  --arg PG_HOST "$PG_HOST" \
  '.timescale["host"] = $PG_HOST | \
  .thorchain["host"] = $CHAIN_HOST' > /etc/chainservice/config.json
RUN cat /etc/chainservice/config.json

ENTRYPOINT ["dumb-init"]
CMD ["/bin/sh"]
