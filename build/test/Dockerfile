FROM golang:1.13

RUN env

WORKDIR /var/midgard

RUN apt-get update
RUN apt-get install -y dumb-init make git

ENTRYPOINT ["dumb-init"]
CMD ["/bin/sh"]
