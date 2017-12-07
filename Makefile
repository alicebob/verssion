.PHONY: all test build integration db

all: test build

test:
	$(MAKE) -C core test
		
build:
	$(MAKE) -C cmd/web build

integration: db
	$(MAKE) -C core integration

db:
	psql verssion < tables.sql
