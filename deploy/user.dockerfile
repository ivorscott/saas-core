FROM golang:1.18-alpine as base

LABEL org.opencontainers.image.authors="devpie"

ENV CGO_ENABLED=0

WORKDIR /core

COPY ../go.* ./

COPY ../cmd/project ./cmd/project
COPY ../internal/project ./internal/project
COPY ../pkg ./pkg

RUN go mod download && go build ./cmd/project

CMD ["./project"]
