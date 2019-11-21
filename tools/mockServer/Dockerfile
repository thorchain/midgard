FROM golang:latest

WORKDIR /var/mockServer

COPY . .

RUN GO111MODULE=on go mod verify
RUN go build -o main .

EXPOSE 8081

CMD ["./main"]
