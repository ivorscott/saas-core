ARG GITHUB_USER
ARG GITHUB_TOKEN

FROM golang:1.15-alpine as base

LABEL org.opencontainers.image.authors="devpie"

ENV CGO_ENABLED=0

WORKDIR /users

FROM base as reload

ARG GITHUB_USER
ARG GITHUB_TOKEN

EXPOSE 4000

RUN apk add --no-cache git \
&& echo "machine github.com login $GITHUB_USER password $GITHUB_TOKEN" > ~/.netrc

COPY go.* ./
RUN go mod download && go mod verify

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN go get github.com/githubnemo/CompileDaemon

COPY . .

RUN go env

FROM reload as dev

RUN go build -o main .

CMD ["./main"]

FROM dev as build-stage

RUN go build -o main . && go test ./... -v --short

FROM aquasec/trivy:0.4.4 as trivy

RUN trivy alpine:3.13.4 && \
    echo "No image vulnerabilities" > result

FROM alpine:3.13.4 as prod

EXPOSE 4000

COPY --from=trivy result secure
COPY --from=build-stage /users/main main

HEALTHCHECK --interval=3s --timeout=3s CMD wget --spider -q http://localhost:4000/v1/health || exit 1
    
CMD ["./main"]