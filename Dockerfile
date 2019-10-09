FROM golang:alpine

ARG chain_host
ARG influx_host

ENV CHAIN_HOST=$chain_host
ENV INFLUX_HOST=$influx_host

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
  --arg INFLUX_HOST "$INFLUX_HOST" \
  '.influx["host"] = $INFLUX_HOST | \
  .statechain["host"] = $CHAIN_HOST' > /etc/chainservice/config.json
RUN cat /etc/chainservice/config.json

ENTRYPOINT ["dumb-init"]
CMD ["/bin/sh"]
