FROM golang:1.18-alpine as base

LABEL org.opencontainers.image.authors="devpie"

ENV CGO_ENABLED=0

RUN mkdir log
COPY go.* ./
COPY cmd/user ./cmd/user
COPY internal/user ./internal/user
COPY pkg ./pkg

RUN echo "testing" && ls -l && pwd
RUN go mod download && go build ./cmd/user

CMD ["./user"]
