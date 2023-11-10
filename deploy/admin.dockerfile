FROM golang:1.21-alpine as base

LABEL org.opencontainers.image.authors="devpie"

ENV CGO_ENABLED=0

WORKDIR /app

RUN mkdir log

COPY go.* ./
COPY cmd/admin ./cmd/admin
COPY internal/admin ./internal/admin

RUN go mod download && go build ./cmd/admin

CMD ["/app/admin"]
