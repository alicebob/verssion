.PHONY: test

export PGURL=postgresql:///verssion

test:
	rm -rf pgsnap_*.txt
	PGPROXY=1 go test

replay:
	go test
