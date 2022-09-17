FROM golang:1.19-alpine AS builder

WORKDIR /code

ENV CGO_ENABLED=0

COPY go.mod /code/go.mod
COPY go.sum /code/go.sum
RUN go mod download -x

COPY . /code/
RUN go build .
RUN ./net-conn-proxy -h



FROM debian:bullseye-slim

WORKDIR /

COPY --from=builder /code/net-conn-proxy /usl/loca/bin/net-conn-proxy

ENTRYPOINT [ "/usl/loca/bin/net-conn-proxy" ]
