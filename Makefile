all: test build

.PHONY: test
test:
	$(MAKE) -C w test
		
.PHONY: build
build:
	$(MAKE) -C cmd/wikitree build
