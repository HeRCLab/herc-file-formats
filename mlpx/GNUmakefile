include ./opinionated.mk

test:
> $(MAKE) -C ./go test
> $(MAKE) -C ./c test
.PHONY: test

lint:
> $(MAKE) -C ./go lint
.PHONY: lint

clean:
> $(MAKE) -C ./go clean
> $(MAKE) -C ./c clean
.PHONY: clean


