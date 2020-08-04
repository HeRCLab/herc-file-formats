include ./opinionated.mk

PREFIX=/usr/local

build: build-mlpx build-wavegen
.PHONY: build

builddir: clean
> mkdir -p ./build
.pHONY: builddir

build-mlpx: builddir
> $(MAKE) -C ./mlpx build
> cp -R ./mlpx/build/* ./build

build-wavegen: builddir
> $(MAKE) -C ./wavegen build
> cp -R ./wavegen/build/* ./build

lint: lint-mlpx lint-wavegen
.PHONY: lint

lint-mlpx:
> $(MAKE) -C ./mlpx lint
.PHONY: lint-mlpx

lint-wavegen:
> $(MAKE) -C ./wavegen lint
.PHONY: lint-wavegen

test: test-mlpx test-wavegen
.PHONY: test

test-mlpx:
> $(MAKE) -C ./mlpx test
.PHONY: test-mlpx

test-wavegen:
> $(MAKE) -C ./wavegen test
.PHONY: test-wavegen

ci: lint test
.PHONY: ci

install: build
> cp -r ./build/bin/* $(PREFIX)/bin/
> cp -r ./build/man/* $(PREFIX)/man/
> cp -r ./build/lib/* $(PREFIX)/lib/
> cp -r ./build/include/* $(PREFIX)/include/
.PHONY: install

clean:
> rm -rf ./build
.PHONY: clean
