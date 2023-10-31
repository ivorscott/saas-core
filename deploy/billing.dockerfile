FROM golang:1.21-alpine as base

LABEL org.opencontainers.image.authors="devpie"

ENV CGO_ENABLED=0

WORKDIR /app

RUN mkdir log

COPY go.* ./
COPY cmd/billing ./cmd/billing
COPY internal/billing ./internal/billing
COPY pkg ./pkg

RUN go mod download && go build ./cmd/billing

CMD ["/app/billing"]
