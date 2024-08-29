# Golang app builder
FROM golang:1.23.0-alpine3.20 AS builder

WORKDIR /mnt/src

COPY go.mod go.sum ./
COPY cmd cmd/
COPY web web/
COPY xero xero/

RUN go build -o webserver ./cmd/webserver/main.go


# Golang app runner
FROM alpine:3.20

COPY --from=builder /mnt/src/webserver /srv/webserver

EXPOSE 4000

ENTRYPOINT [ "/srv/webserver" ]