# mlpx

*MultiLayer Perceptron eXchange*


This project implements a trivial JSON-based format for storing and comparing snapshots of Multilayer Perceptrons. 

## Installation

1. Run `make build`
2. Copy `./build/bin/mlpx` into your `$PATH`, OR proceed to the next step
3. Run `sudo make install`
	* If you would like the C library to be installed, you MUST use
	  `make install` or a pre-build package.

## Documentation

* [mlpx(5)](./doc/5/mlpx.md)
* [mlpx(3)](./doc/3/mlpx.md)

## Using the C library

After appropriate installation of the C library, you should have a
`mlpx-config` script available in `$PATH`. You can use this to query the
appropriate `CFLAGS` and `LIBS` values for the MLPX library using `mlpx-config
--cflags` and `mlpx-config --libs` commands respectively.

## Dependencies

To run the linting checks:

* `go get -u github.com/gordonklaus/ineffassign`
* `go get -u github.com/kisielk/errcheck`
* `go get -u golang.org/x/lint/golint`

