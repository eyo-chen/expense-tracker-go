PWD = $(shell pwd)

mock:
	docker run --rm -v "$(PWD)":/src -w /src vektra/mockery --all