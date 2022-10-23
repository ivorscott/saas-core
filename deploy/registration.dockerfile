FROM golang:1.18-alpine as base

LABEL org.opencontainers.image.authors="devpie"

ENV CGO_ENABLED=0

WORKDIR /app

RUN mkdir log

COPY go.* ./
COPY cmd/registration ./cmd/registration
COPY internal/registration ./internal/registration
COPY pkg ./pkg

RUN go mod download && go build ./cmd/registration

CMD ["/app/registration"]
