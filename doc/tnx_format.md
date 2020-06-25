tnx(4) -- the Trivial Neural network eXchange Format
====================================================

## Description

A TNX file is a JSON formatted document containing a single object. TNX is
intended to allow encoding a neural network's topology and/or weights. The
objective of NNX is to provide a format for initializing neural networks with
common weights, and to examine their internal values later, for the purpose of
cross-validating different implementations.

At present TNX targets exclusively MLP networks, but in the future could expand
to accommodate other types. Unlike other formats such as ONNX, TNX assumes that
the computation medium used by the network operates analogously to a multilayer
perception, where each layer has a set of weights, a set of biases, and an
activation function.

The JOSN object in an TNX file **must** include the following keys:

* `schema` -- a tuple of two elements identifying a schema and a version level.
  At present, this should be `[tnx, 0]`. Compliant implementations should
  refuse to operate on  an unknown schema or version level.
* `topology` -- a TNX topology definition, see the *Topology Definition*
  section below.

The object **may** include the following keys:

* `snapshot` -- an TNX snapshot definition, see the *Snapshot Definition* below.

Any other keys defined should be ignored.

**NOTE**: key names specified in this document are case-sensitive.

## Topology Definition

An TNX topology definition is a JSON object, which **must** contain the
following keys:

* `layers` -- a list of layer objects in an arbitrary order, containing at a
  minimum three layer objects.
* `neurons` -- a list of neuron objects in an arbitrary order.


A layer object **must** contain the following keys:
* `index` -- integer, a zero-indexed layer number. Layer 0 specifies the input layer,
  and the layer with the highest index is the output layer. Layer numbers must
  be contiguous.
* `activation` -- a string specifying the activation function of the layer,
 see *Activation Function* below.

A neuron object **must** contain the following keys:
* `layer` -- integer, a zero-indexed layer number, which **must** correspond to one of
  the layer sizes in the `layers` list. For example, a neuron in the input
  layer would always have a layer number of `0`.
* `index` -- integer, the zero-indexed neuron number. **Must** be unique within each
  layer number, and be in the range 0 to the number of neurons in the relevant
  layer.


## Snapshot Definition

An TNX snapshot definition is a JSON list which associates various values
with a neuron, identified by a layer number and neuron number. Each element
in the list is be an object, which **must** contain the following keys:

* `layer` -- integer, identifies the layer number, **must** correspond to a layer in the
  topology definition.
* `neuron` -- integer, identifies the neuron number, **must** correspond to a neuron
  in the given layer of the topology definition.

A snapshot object **may** contain any subset of the following keys:

* `weight` -- list float, specifies the weigh neuron's weight vector. The ith
  element of `weight` should correspond to the `ith` neuron in the previous
  layer. The number of elements in the `weight` must be equal to the number
  of neurons in the layer immediately preceding the neuron.
* `bias` -- float, the bias value of the neuron.
* `delta` -- float, the delta value of the neuron, as described by Russel &
  Norvig's back-propagation algorithm. This key is used only for
  cross-validation of MLP implementations.
* `output` -- float, neuron output before the activation function is applied.
  This key is used only for cross-validation of MLP implementations.
* `activation` -- float, neuron output after the activation function is
  applied.
