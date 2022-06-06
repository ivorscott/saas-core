FROM golang:1.18-alpine as base

LABEL org.opencontainers.image.authors="devpie"

ENV CGO_ENABLED=0

WORKDIR /core

COPY ../go.* .
RUN go mod download

COPY ../cmd/registration ./cmd/registration
COPY ../internal/registration ./internal/registration
COPY ../pkg ./pkg

RUN go build ./cmd/registration

CMD ["./registration"]
