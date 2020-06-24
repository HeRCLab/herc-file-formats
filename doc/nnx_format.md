nnx(4) -- the Neural Network eXchange Format
============================================

## Description

An NNX file is a JSON formatted document containing a single object. NNX is
intended to allow encoding a neural network's topology and/or weights. The
objective of NNX is to provide a format for initializing neural networks with
common weights, and to examine their internal values later, for the purpose of
cross-validating different implementations.

At present NNX targets exclusively MLP networks, but in the future could expand
to accommodate other types.

The JOSN object in an NNX file **must** include the following keys:

* `schema` -- a tuple of two elements identifying a schema and a version level.
  At present, this should be `[nnx, 0]`. Compliant implementations should
  refuse to operate on  an unknown schema or version level.
* `topology` -- an NNX topology definition, see the *Topology Definition*
  section below.

The NNX file **may** include the following keys:

* `snapshot` -- an NNX snapshot definition, see the *Snapshot Definition* below.

Any other keys defined should be ignored.

**NOTE**: key names must be case-insensitive, in other words both `Schema` and
`schema` are valid.

## Topology Definition

An NNX topology definition is a JSON object, which **must** contain the
following keys:

* `layers` -- an array of integers indicating the size of each layer, for
  example `[10, 100, 15]` would indicate an input layer size of 10 neurons, a
  hidden layer of 100 neurons, and an output layer of 15 neurons. This list
  must contain at least three elements. The number of elements can be used to
  infer the number of layers in the network.
* `neurons` -- a list of neuron objects.

Any other keys should be ignored. Key names are case-insensitive.

A neuron object **must** contain the following keys:
* `layer` -- a zero-indexed layer number, which **must** correspond to one of
  the layer sizes in the `layers` list. For example, a neuron in the input
  layer would always have a layer number of `0`.
* `index` -- the zero-indexed neuron number. **Must** be unique within each
  layer number, and be in the range 0 to the number of neurons in the relevant
  layer.


## Snapshot Definition

An NNX snapshot definition is a JSON object which **must** contain the following
keys:

* `
