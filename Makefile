.PHONY: all test testrace build release ci db tidy

all: test build

test:
	go test ./...

testrace:
	go test -race ./...
		
build:
	$(MAKE) -C cmd/web build

release:
	GOARG=x86-64 GOOS=openbsd $(MAKE) -C cmd/web build

ci:
	go test ./...

db:
	psql verssion < tables.sql

tidy:
	go mod tidy -compat=1.24
