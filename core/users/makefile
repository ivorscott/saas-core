include .env
export $(shell sed 's/=.*//' .env)

test:
	go test ./... -v
.PHONY: test

unit:
	go test ./... -v -short
.PHONY: unit