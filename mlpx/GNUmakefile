include ./opinionated.mk

GOSRC=$(shell find . -iname "*.go" )
PREFIX=/usr/local

BUILD_MAN_PAGES=build/man/man3/mlpx.3 build/man/man5/mlpx.5 build/man/man1/mlpx.1 build/man/man1/mlpx-new.1 build/man/man1/mlpx-validate.1 build/man/man1/mlpx-diff.1
BUILD_BINARIES=build/bin/mlpx build/bin/mlpx-config
BUILD_INCLUDES=build/include/mlpx/mlpx.h
BUILD_LIBS=build/lib/libmlpx.so build/lib/libmlpx.a
BUILD_EVERYTHING=$(BUILD_MAN_PAGES) $(BUILD_BINARIES) $(BUILD_INCLUDES) $(BUILD_LIBS)

test:
> $(MAKE) -C ./go test
> $(MAKE) -C ./c test
.PHONY: test

lint:
> $(MAKE) -C ./go lint
.PHONY: lint


builddirs:
> mkdir -p ./build/bin
> mkdir -p ./build/man/man1
> mkdir -p ./build/man/man3
> mkdir -p ./build/man/man5
> mkdir -p ./build/lib
> mkdir -p ./build/include/mlpx
.PHONY: builddirs

build/bin/mlpx: $(GOSRC) go.mod builddirs
> go build -o $@ ./cmd/mlpx/main.go

build/bin/mlpx-config: ./cmd/mlpx-config/mlpx-config
> cp $< $@
> chmod +x $@
.PHONY: build/bin/mlpx-config

build/man/man3/%.3: doc/3/%.md builddirs
> ronn < $< > $@

build/man/man5/%.5: doc/5/%.md builddirs
> ronn < $< > $@

build/man/man1/mlpx.1: ./build/bin/mlpx builddirs
> help2man --include=include.txt --no-info --no-discard-stderr $< > "$@"

build/man/man1/mlpx-new.1: ./build/bin/mlpx builddirs
> help2man --include=include.txt --no-info --no-discard-stderr "$< new" > "$@"

build/man/man1/mlpx-diff.1: ./build/bin/mlpx builddirs
> help2man --include=include.txt --no-info --no-discard-stderr "$< diff" > "$@"

build/man/man1/mlpx-validate.1: ./build/bin/mlpx builddirs
> help2man --include=include.txt --no-info --no-discard-stderr "$< validate" > "$@"

build/lib/libmlpx.so: builddirs
> $(MAKE) -C ./c mlpx.so
> cp ./c/mlpx.so $@
.PHONY: build/lib/libmlpx.so

build/lib/libmlpx.a: builddirs
> $(MAKE) -C ./c mlpx.a
> cp ./c/mlpx.a $@
.PHONY: build/lib/libmlpx.a

build/include/mlpx/mlpx.h: builddirs
> $(MAKE) -C ./c mlpx.h
> cp ./c/mlpx.h $@
.PHONY: build/lib.mlpx.h

build: $(BUILD_EVERYTHING)
.PHONY: build

install: build
> cp -r ./build/bin/* $(PREFIX)/bin/
> cp -r ./build/man/* $(PREFIX)/man/
> cp -r ./build/lib/* $(PREFIX)/lib/
> cp -r ./build/include/* $(PREFIX)/include/
.PHONY: install

clean:
> $(MAKE) -C ./go clean
> $(MAKE) -C ./c clean
> rm -rf ./build
.PHONY: clean
