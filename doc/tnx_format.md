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

* `parameters` -- a TNX parameterization, see *Parameterization* below.
* `snapshot` -- an TNX snapshot definition, see the *Snapshot Definition* below.
	* **NOTE**: if a `snapshot` key is included, a `parameters` key
	  **must** also be included.

Any other keys defined should be ignored.

**NOTE**: key names specified in this document are case-sensitive.

## Topology Definition

A topology definition represents a DAG. Nodes in the DAG are computations to
be performed, and edges represent the flow of data between nodes in the DAG.

An TNX topology definition is a JSON object, which **must** contain the
following keys:

* `nodes` -- a list of node objects.
* `links` -- a list of link objects.

A node object **must** contain the following keys:

* `id` -- a node ID.
* `operation` -- a string identifying which operation this node performs, see
  *Defined Operations*.
* `inputs` -- a list of IDs. Each ID  defines a unique identifier for an input
  to the given node.
* `outputs` -- a list of IDs. Each ID defines a unique identifier for an output
  from the given node.

A link object **must** contain the following keys:

* `source` -- the source ID of the link, which **must** reference an output ID
  defined in the `nodes` list.
* `target` -- the target ID of the link, which **must** reference an input ID
  defined in the `nodes` list.

An ID is an arbitrary unicode string. IDs must be globally unique within the
TNX file. An ID identifies a node, or an in input or output to a node. For
human readability, it is suggested that inputs to a node should have
identifiers of the form `nodeid<-inputname`, and that outputs should have
identifiers of the form `nodeid->outputname`. This is purely a convention,
implementations **should not** rely on this naming scheme to derive any
information about the structure of the graph.

Inputs to the graph are defined by the `input` operation type, such nodes
should have an empty `inputs` list.

Outputs of the graph are defined by the `output `operation type. Such nodes
should have an empty `outputs` list.

## Parameterization

A parameterization is used to describe parameters of nodes in the graph.  This
is distinct from the concept of a snapshot insofar as that a parameter
describes something about the graph which always holds for a given instance of
the graph, while a snapshot is intended to represent the execution state of
the graph at a particular time.

If one imagines compiling the graph to an FPGA using HLS, changing "parameters"
in the sense used in this document would require re-compiling the graph,
re-generating a bitstream, and re-programming the FPGA. A snapshot would be
representative of the internal state of the memories on the FPGA at a given
instant in time.

A TNX parameterization is a table where keys are node IDs (which **must** be
valid node IDs occurring within the topology definition), and values are
arbitrary JSON objects. The parameterization of a node is specific to it's
operation.

## Snapshot Definition

A snapshot definition is used to declare that state of some item in the graph
at some particular time. Because snapshots are intended to be arbitrary
state-storage mechanisms, the content of each snapshot item will vary according
to what type of node it pertains to.

A TNX snapshot definition is a list of JSON objects, each of which may be of
one of several types. All such objects **must** have a `type` field identifying
it's type.  Other fields are specific to each type.

The following `type` values are allowed:

* `matrix`

A matrix object **must** have the following keys:

* `id` -- an ID which **must** reference a node, node input, or node output
  defined in the topology section.
* `name` -- a string definition the meaning of this matrix within the context
  of the node's operation operation.
* `dimensions` -- a list of integers describing the size of the matrix. The
  length of the list is equal to the number of dimensions.
* `data` -- a list of floating-point values describing the matrix's contents.
  The length of the list must be equal to the product of the matrix dimensions.
  Data is packed assuming the dimensions with the smallest index are more major.
  In other words, if the first element of the `dimensions` was the number of
  rows in the matrix, the data would be stored row-major. To be clear, the
  `data` list is always one-dimensional.

## Operations

An operation refers to a computation which is performed by a node on it's
inputs, which **may** causes it's outputs to change. Not all client
implementations need to implement all operations described in this document,
but **should** throw an informative error and exit if an unsupported operation
is given as input.

An operation is identified by an arbitrary unicode string. Specific
implementations **may** define their own operations, but such custom operations
**must** be prefixed with the characters `x:`. It is guaranteed that no future
version of TNX specification will define operations with names starting with
these characters. Custom operations **should** be named descriptively, to avoid
collisions with other implementations. Operation names beginning with the
characters `e:` are reserved for the official implementation for
experimentation and development purposes.

All node inputs and outputs are snapshot-ed using the `matrix` type. Such
snapshots **must** refer to the input or output ID (rather than the ID of the
node itself). The matrix `name` field is not relevant in this context. When an
operation describes the dimensions of it's inputs or outputs, the relevant
snapshots, if any, **must** use the same dimensions.

There are two special-case operations, used to pass input into and out of the
graph.


### Input

The `input` operation is used to consume input from the outside world.

It's parameterization **must** define a `dimensions` list (as described
previous for `matrix` type snapshots).

How data is passed into an input node from the outside world is implementation
defined.

### Output

The `output` operation is used to transfer information to the outside world.

It's parameterization **must** define a `dimensions` list (as described
previous for `matrix` type snapshots).

When data is sent to the input of an output node, it should be transferred out
of the graph to the outside world in an implementation-defined way.

### MLPLayer

The `mlplayer` operation is used to describe a single layer in a MLP using
back-propogation as described by Russel and Norvig.

It's parameterization **must** include the following keys:

* `neurons` -- positive integer describing the number of neurons in the layer.
  For the remainder of this section *n* is used to describe the number of
  neurons thus defined.

It's parameterization **may** include the following keys:

* `activation` -- a node ID referring to the node which implements the layer's
  activation function. If defined, this **must** be a defined node in the
  topology which has this node's output as an input. This field is used to convey
  semantic intent only.

An `mlplayer` operation **must** have exact one input, which **must** be a
one-dimensional matrix which **should** contain a number of values equal to the
number of neurons in the preceding layer. For the remainder of this section,
*k* will refer to the unknown dimension of the layer's input.

An `mlplayer` operation **must** have exactly one output, which **must** be a
one-dimensional matrix which **must** contain a number of values exactly equal
to it's declared number of neurons.

The following snapshot names **may** be defined for an mlplayer node:

* `deltas` -- a matrix consisting of a vector of length *n* describing the
  back-propagation delta values for the node.
* `weights` -- a matrix of size *k* x *n* describing the weights for the
  layer, with `weights[0<j<k][0<i<n]` being the weight from node *j* in the
  preceding layer to the node *i* in this layer.
* `biases` -- a matrix consisting of a vector of length *n* describing the bias
  values for each neuron in the layer.

**NOTE**: activation functions are accomplished by the relevant separate
operations. Therefore, the output values from layer pre-activation are
snapshot-ed as the outputs from the relevant node.

### ReLU

The `relu` operation is used to describe a ReLU activation function. It's
inputs and outputs are the same dimensions. The *i*-th element of the output is
the ReLU of the *i*-th input element.

### Identity

The `identity` operation works similarly to ReLU, but implements an identity
function.

### Sigmoid

The `sigmoid` operation works similarly to ReLU, but implements a sigmoid
function.

## Example

The following example describes an MLP with three hidden layers of size 25, 15,
and 10. It's input layer has a size of 25, and it's output layer a size of 5.

```json
{
	"schema": ["tnx", 0],
	"topology": {
		"nodes": [
			{
				"id": "input",
				"operation": "output",
				"outputs": "input->output0",
			},
			{
				"id": "hidden1"
				"operation": "mlplayer",
				"inputs: ["hidden1<-input0"],
				"outputs": ["hidden1->output0"],
			},
			{
				"id": "activaton1"
				"operation": "relu",
				"inputs: ["activation1<-input0"],
				"outputs": ["activation1->output0"],
			},
			{
				"id": "hidden2"
				"operation": "mlplayer",
				"inputs: ["hidden2<-input0"],
				"outputs": ["hidden2->output0"],
			},
			{
				"id": "activaton2"
				"operation": "relu",
				"inputs: ["activation2<-input0"],
				"outputs": ["activation2->output0"],
			},
			{
				"id": "hidden3"
				"operation": "mlplayer",
				"inputs: ["hidden3<-input0"],
				"outputs": ["hidden3->output0"],
			},
			{
				"id": "activaton3"
				"operation": "relu",
				"inputs: ["activation3<-input0"],
				"outputs": ["activation3->output0"],
			},
			{
				"id": "output",
				"operation": "output"
				"inputs": "output<-input0"
			}
		],
		"links": [
			{
				"source": "input->output0",
				"target": "hidden1<-input0",
			},
			{
				"source": "hidden1->output0",
				"target": "activation1->input0"
			},
			{
				"source": "activation1->output0",
				"target": "hidden2->input0"
			},
			{
				"source": "hidden2->output0",
				"target": "activation2->input0"
			},
			{
				"source": "activation2->output0",
				"target": "hidden3->input0"
			},
			{
				"source": "hidden3->output0",
				"target": "activation3->input0"
			},
			{
				"source": "activation3->output0",
				"target": "output<-input0""
			},
		]
	},
	"parameters": {
		"input": {
			"dimensions": [25]
		},
		"hidden1": {
			"neurons": 25,
			"activation": "activation1"
		},
		"hidden2": {
			"neurons": 15,
			"activation": "activation2"
		}
		"hidden3": {
			"neurons": 15,
			"activation": "activation3"
		}
		"output": {
			"dimensions": [5]
		},
	}
}
```
