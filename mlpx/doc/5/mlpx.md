mlpx(5) -- the MLP eXchange Format
==================================

## Description

An MLPX file is a JSON formatted document containing a single object. MLPX is
used for the purpose of representing MLP weights, biases, and other snapshot
information. It is intended to be trivial to implement in any language the
supports JSON.

The top-level JSON object **must** include the following keys:

* `schema` --  a tuple of two elements identifying a schema and a version
  level. At present, this should be `["mlpx", 0]`. Compliant implementations
  should refuse to operate on an unknown schema or version level.
* `snapshots` -- a table of Snapshot Definition object, as described below.
  Keys are snapshot IDs, and values are snapshot definition objects.

Note that all snapshot definitions in a given MLPX file **must** be isomorphic.
That is they must have the same number of layers, and all layers must have the
same number of neurons across different snapshots. It is intended that storing
multiple snapshots  is used to represent the progression of a network over time.

Snapshot IDs **must** be numeric, positive integers, although necessarily
stored as strings, except as otherwise noted.

If an MLPX file is used for initialization, it **should** include a snapshot
with the ID `initializer`, which implementations **should** use for
initialization purposes. Implementations **may** allow the user to select an
arbitrary snapshot as an initializer (such as when resuming from a checkpoint),
but **should** use `initializer` as the default.

## Snapshot Definition

A snapshot definition is used to record the state of an MLP at a particular
point in time.

A snapshot defintion is a JSON object which **must** include the following
keys:

* `layers` -- a table of Layer Definition objects, as described below. Keys
  should be layer IDs, and values should be Layer Defintion objects.

The following layer IDs are reserved:

* `input` -- used exclusively for the input layer
* `output` -- used exclusively for the output layer

Other layer IDs may be arbitrary strings.


## Layer Definition

A Layer Definition is a JSON object which **must** include the following keys:

* `predecessor` -- string layer ID, identifies the previous layer. This field
  is ignored for the input layer.
* `successor` -- string layer ID, identifies the next layer. This field is
  ignored for the output layer.
* `neurons` -- integer number of neurons in the given layer.

The following keys **may** be included:

* `weights` -- array of weights for the outputs of the previous layer coming
  into the current layer. Elements should be floating point values. The `(j *
  np + i)`-th element is the weight for the connection TO `j`-th neuron in the
  current layer, FROM the `i`-th neuron in the previous layer, assuming that `np`
  is the number of neurons in the previous layer. The weights layer is not
  significant for the input layer.
* `outputs` -- array of floating point values for the outputs of this layer.
  The `i`-th element corresponds to the output of the `i`-th neuron. Should
  contain as many elements as this layer has neurons.
* `activations` -- floating point values of the same dimensions as the `outputs`
  array, but instead encodes the output values after the activation function
  has been applied.
* `deltas` -- floating point values corresponding to the delta values computed.
   Should have the same number of entries as this layer has neurons.
* `biases` -- floating point values corresponding to the bias values for each
   neuron. Should have the same number of entries as this layer has neurons.
* `activation_function` -- string identifying the activation function used,
   is intended for human readers.

These keys are all optional, to support a variety of different use cases. For
example, an mlpx file used to initialize common layer weights and biases might
omit the `activation`, `deltas`, and `biases` fields.

## Rationale

MLPX is intended to be as straightforward as possible to implement for a
variety of MLP implementations. It needs to support the following use cases:

* Snapshot-ing a neural network while it is being trained, or while it is being
  executed on a given input, especially with the purpose of comparing different
  implementations.
* Storing common initial values for a variety of different implementations.

The format allows storing multiple snapshots so that implementations can save
a detailed log of the evolution of their internal model over time in a single
convenient file.

MLPX may be extended in the future to accommodate further MLP-related
functionality, however to keep the format simple, should not be extended to
support other network topologies. Future work in other types of networks, such
as LSTMs, should instead be handled by separate formats.
