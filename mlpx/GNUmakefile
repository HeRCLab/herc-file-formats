include ./opinionated.mk

GOSRC=$(shell find . -iname "*.go" )
PREFIX=/usr/local

test:
> $(MAKE) -C ./go test
> $(MAKE) -C ./c test
.PHONY: test

lint:
> $(MAKE) -C ./go lint
.PHONY: lint


builddirs:
> mkdir -p ./build
> mkdir -p ./build/bin
> mkdir -p ./build/man
> mkdir -p ./build/man/3
> mkdir -p ./build/man/5
.PHONY: builddirs

build/bin/mlpx: $(GOSRC) go.mod builddirs
> go build -o $@ ./cmd/mlpx/main.go 

build/man/3/%.3: doc/3/%.md
> ronn < $< > $@

build/man/5/%.5: doc/5/%.md
> ronn < $< > $@

build: build/bin/mlpx build/man/3/mlpx.3 build/man/5/mlpx.5
.PHONY: build

install: build
> cp -r ./build/bin/* $(PREFIX)/bin/
> cp -r ./build/man/* $(PREFIX)/man/
.PHONY: install

clean:
> $(MAKE) -C ./go clean
> $(MAKE) -C ./c clean
> rm -rf ./build
.PHONY: clean
