.PHONY: all test build db

all: test build

test:
	$(MAKE) -C w test
		
build:
	$(MAKE) -C cmd/wikispider build
	$(MAKE) -C cmd/web build

db:
	psql w < tables.sql
