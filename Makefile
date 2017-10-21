.PHONY: all test build db

all: test build

test:
	$(MAKE) -C w test
		
build:
	$(MAKE) -C cmd/wikitree build

db:
	psql w < tables.sql
