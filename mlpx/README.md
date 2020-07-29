# mlpx

*MultiLayer Perceptron eXchange*


This project implements a trivial JSON-based format for storing and comparing snapshots of Multilayer Perceptrons. 

## Installation

1. Run `make build`
2. Copy `./build/bin/mlpx` into your `$PATH`, OR proceed to the next step
3. Run `sudo make install`


## Documentation

* [mlpx(5)](./doc/5/mlpx.md)

## Dependencies

To run the linting checks:

* `go get -u github.com/gordonklaus/ineffassign`
* `go get -u github.com/kisielk/errcheck`
* `go get -u golang.org/x/lint/golint`

