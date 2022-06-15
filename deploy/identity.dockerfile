FROM golang:1.18-alpine as base

LABEL org.opencontainers.image.authors="devpie"

ENV CGO_ENABLED=0

WORKDIR /core

COPY ../go.* ./

COPY ../cmd/identity ./cmd/user
COPY ../internal/user ./internal/user
COPY ../pkg ./pkg

RUN go mod download && go build ./cmd/user

CMD ["./user"]
