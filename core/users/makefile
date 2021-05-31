include .env
export $(shell sed 's/=.*//' .env)

.DEFAULT_GOAL := build

test:
	go test ./... -v
.PHONY: test

unit:
	go test ./... -v -short
.PHONY: unit

fmt:
	go fmt ./...
.PHONY: fmt

lint: fmt
	golangci-lint run
.PHONY: lint

vet: fmt
	go vet ./...
 	$(shell shadow ./...)
.PHONY: vet

build: lint vet
	go build -o main .
