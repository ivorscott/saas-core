include .env
export $(shell sed 's/=.*//' .env)

test:
	go test ./... -v
.PHONY: test