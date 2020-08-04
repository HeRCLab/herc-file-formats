include ./opinionated.mk

GOSRC=$(shell find . -iname "*.go" )
PREFIX=/usr/local

BUILD_MAN_PAGES=build/man/3/wavegen.3 build/man/5/wavegen.5  build/man/1/wavegen.1 build/man/1/wavegen-view.1 build/man/1/wavegen-generate.1 build/man/1/wavegen-summarize.1
BUILD_BINARIES=build/bin/wavegen build/bin/wavegen-config
BUILD_INCLUDES=build/include/wavegen/wavegen.h
BUILD_LIBS=build/lib/libwavegen.so build/lib/libwavegen.a
BUILD_EVERYTHING=$(BUILD_MAN_PAGES) $(BUILD_BINARIES) $(BUILD_INCLUDES) $(BUILD_LIBS)

build: $(BUILD_EVERYTHING)
.PHONY: build

test:
> $(MAKE) -C ./go test
> $(MAKE) -C ./c test
.PHONY: test

lint:
> $(MAKE) -C ./go lint
.PHONY: lint


builddirs:
> mkdir -p ./build/bin
> mkdir -p ./build/man/1
> mkdir -p ./build/man/3
> mkdir -p ./build/man/5
> mkdir -p ./build/lib
> mkdir -p ./build/include/wavegen
.PHONY: builddirs

build/bin/wavegen: $(GOSRC) go.mod builddirs
> go build -o $@ ./cmd/wavegen/main.go

build/bin/wavegen-config: ./cmd/wavegen-config/wavegen-config
> cp $< $@
> chmod +x $@
.PHONY: build/bin/wavegen-config

build/man/3/%.3: doc/3/%.md builddirs
> ronn < $< > $@

build/man/5/%.5: doc/5/%.md builddirs
> ronn < $< > $@

build/man/1/wavegen.1: ./build/bin/wavegen builddirs
> help2man --include=include.txt --no-info --no-discard-stderr $< > "$@"

build/man/1/wavegen-generate.1: ./build/bin/wavegen builddirs
> help2man --include=include.txt --no-info --no-discard-stderr "$< generate" > "$@"

build/man/1/wavegen-view.1: ./build/bin/wavegen builddirs
> help2man --include=include.txt --no-info --no-discard-stderr "$< view" > "$@"

build/man/1/wavegen-summarize.1: ./build/bin/wavegen builddirs
> help2man --include=include.txt --no-info --no-discard-stderr "$< summarize" > "$@"

build/lib/libwavegen.so: builddirs
> $(MAKE) -C ./c wavegen.so
> cp ./c/wavegen.so $@
.PHONY: build/lib/libwavegen.so

build/lib/libwavegen.a: builddirs
> $(MAKE) -C ./c wavegen.a
> cp ./c/wavegen.a $@
.PHONY: build/lib/libwavegen.a

build/include/wavegen/wavegen.h: builddirs
> $(MAKE) -C ./c wavegen.h
> cp ./c/wavegen.h $@
.PHONY: build/lib.wavegen.h


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
