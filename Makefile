.PHONY: all test testrace build integration db

all: test build

test:
	go test ./...

testrace:
	go test -race ./...
		
build:
	$(MAKE) -C cmd/web build

integration:
	go test -tags integration ./...

db:
	psql verssion < tables.sql

tidy:
	go mod tidy -compat=1.18
