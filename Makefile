.PHONY: all test build db

all: test build

test:
	$(MAKE) -C core test
		
build:
	$(MAKE) -C cmd/web build

db:
	psql w < tables.sql
