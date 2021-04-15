ARG GITHUB_USER
ARG GITHUB_TOKEN

FROM golang:1.14-alpine as base

ENV CGO_ENABLED=0 

WORKDIR /

FROM base as dev

ARG GITHUB_USER
ARG GITHUB_TOKEN

RUN apk add --no-cache git \ 
&& echo "machine github.com login $GITHUB_USER password $GITHUB_TOKEN" > ~/.netrc

COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

RUN go env && go build -o main .

CMD ["./main"]


