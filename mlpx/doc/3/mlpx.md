mlpx(3) -- MLPX C API
=====================

## NAME

`MLPXClose`, `MLPXGetError`, `MLPXGetInitializerSnapshotIndex`,
`MLPXGetNumSnapshots`, `MLPXLayerGetActivation`,
`MLPXLayerGetActivationFunction`, `MLPXLayerGetBias`, `MLPXLayerGetDelta`,
`MLPXLayerGetIDByIndex`, `MLPXLayerGetIndexByID`, `MLPXLayerNeurons`,
`MLPXLayerGetOutput`, `MLPXLayerGetPredecessorIndex`,
`MLPXLayerGetSuccessorIndex`, `MLPXLayerGetWeight`, `MLPXLayerSetActivation`,
`MLPXLayerSetActivationFunction`, `MLPXLayerSetBias`, `MLPXLayerSetDelta`,
`MLPXLayerSetOutput`, `MLPXLyaerSetOutput`, `MLPXLayerSetWeight`,
`MLPXMakeIsomorphicSnapshot`, `MLPXMakeMLPX`, `MLPXMakeSNapshot`, `MLPXOpen`,
`MLPXSnapshotGetIDByIndex`, `MLPXSnapshotGetIndexByID`,
`MLPXSnapshotGetNumLayers`

## CONVENTIONS

The MLPX C API is a wrapper around the Go implementation, and thus bears a few
artifacts of being an FFI-based solution. Because passing non-primitive types
between the C and Go runtimes is complex, the API has been structured to use
only primitive types. MLPX objects are referenced by an integer handle, and
snapshots and layers are referenced by indices.

Snapshot and layer indices are subscripts into the sorted list of snapshot and
layer IDs, respectively. Thus they are guaranteed to be positive contiguous
integers, and simply incrementing or decrementing such an index is sufficient
to move between layers or snapshots in canonical order (for snapshots, based on
the initializer and numeric ID rules,and for layers, topological order).

**NOTE**: creating new layers or snapshots will change which indices are valid.
If you create a layer which is in the middle of the topological sort order for
the MLP, you need to re-index everything. Also keep in mind that you will need
to run the appropriate function to get the layer or snapshot count again.

All functions in the API return an integer indicating success or failure.
Success is indicated by a 0 value, and failure by any other value. Error
information can be retrieved by calling the `MLPXGetError()` function, which
returns the most recent error (and only the most recent error). All return
values are passed by reference as function arguments.

Note that string returns values will be malloc-ed using Cgo, and should be
freed by the caller as appropriate.

## SYNOPSIS

**include <mlpx.h>**

* **char\* MLPXGetError(void);**:
	Returns the most recent error string.

* **int MLPXOpen(char\* path, int\* handle);**:
	Opens an MLPX filef rom disk and loads it into memory, creating a
	handle if appropriate. The handle may be invalid if an error is
	encountered.

* **int MLPXSave(int handle, char\* path);**:
	Save the MLPX object referenced by the handle to the specified path,
	overwriting it if it exists already.

* **int MLPXClose(int handle);**:
	Closes a previously opened MLPX handle, allowing it's memory to be
	garbaged collected on the Go side.

* **int MLPXIsomorphicDuplicate(int sourceHandle, int destHandle, char\* snapid);**:
	Create a topologically identical duplicate of the given MLPX with a new
	handle. The created MLPX will have a single snapshot with the specified
	snapshot ID. Note that the Alpha value and all layer fields will be
	uninitialized in the duplicated MLPX.

* **int MLPXNextSnapshotID(int handle, char\*\* nextid);**:
	Returns the next canonical snapshot ID after the most recent.

* **int MLPXGetNumSnapshots(int handle, int\* snapc);**:
	Retrieves the number of snapshots currently stored in an MLPX handle.

* **int MLPXSnapshotGetIDByIndex(int handle, int index, char\*\* id);**:
	Retrieves the string ID of a snapshot by it's index. This will be a
	newly allocated C string, placed in C memory given the given address.

* **int MLPXSnapshotGetIndexByID(int handle, char\* id, int\* index);**:
	Retrieves the numeirc index of a snapshot given it's string ID.

* **int MLPXSnapshotGetNumLayers(int handle, int snapshotIndex, int\* layerc);**:
	Retrieves the number of layers in a snapshot by it's numeric index.

* **int MLPXSnapshotGetAlpha(int handle, int snapshotIndex, double\* alpha);**:
	Retrieves the defined alpha value for the given snapshot.

* **int MLPXSnapshotSetAlpha(int handle, int snapshotIndex, double alpha);**:
	Modify the alpha value for the given snapshot.

* **int MLPXLayerGetIndexByID(int handle, int snapshotIndex, char\* id, int\* index);**:
	Retrieves the numeric index of a layer within a given snapshot by it's
	string identifier.

* **int MLPXLayerGetIDByIndex(int handle, int snapshotIndex, int layerIndex, char\*\* id);**:
	Retrives the string identifer of a layer within a given snapshot by it's
	numeric index.

* **int MLPXLayerGetNeurons(int handle, int snapshotIndex, int layerIndex, int\* neuronc);**:
	Retrieves the number of neurons in a given layer identified by it's
	numeric index within a given snapshot.

* **int MLPXMakeMLPX(int\* handle);**:
	Create a new, empty MLPX object without loading from disk.

* **int MLPXMakeSnapshot(int handle, char\* id, double alpha);**:
	Create a new snapshot in a given MLP. The created snapshot's index can
	later be retrieved using `MLPXSnapshotGetIndexByID()` if needed.

* **int MLPXMakeIsomorphicSnapshot(int handle, char\* id, int toSnapshotIndex);**:
	Create a new snapshot, isomorphic to an existing one. This is the
	preferred way to generate new snapshots in an existing MLPX, since
	`mlpx(5)` specifies that all snapshots must be isomorphic (here,
	isomorphism means both topologically identical, and having the same
	layer sizes in terms of neurons).

* **int MLPXLayerGetPredecessorIndex(int handle, int snapshotIndex, int layerIndex, int\* index);**:
	Retrieve the index of the topological predecessor of the given layer.

* **int MLPXLayerGetSuccessorIndex(int handle, int snapshotIndex, int layerIndex, int\* index);**:
	Retrieve the index of the topological successor of the given layer.

* **int MLPXLayerSetWeight(int handle, int snapshotIndex, int layerIndex, int subscript, double value);**:
	Sets a particular subscript within the weights matrix of the given
	layer. If no weights matrix previously existed, one will be allocated
	and all other values initialized to 0.

* **int MLPXLayerGetWeight(int handle, int snapshotIndex, int layerIndex, int subscript, double\* value);**:
	Retrives the weight of a particular subscript within the weights matrix
	of the given layer. Will return an error status if the weights matrix
	has not yet been allocated.

* **int MLPXLayerSetOutput(int handle, int snapshotIndex, int layerIndex, int subscript, double value);**:
	Sets a particular subscript within the outputs matrix of the given
	layer. If no outputs matrix previously existed, one will be allocated
	and all other values initialized to 0.

* **int MLPXLayerGetOutput(int handle, int snapshotIndex, int layerIndex, int subscript, double\* value);**:
	Retrives the weight of a particular subscript within the outputs matrix
	of the given layer. Will return an error status if the outputs matrix
	has not yet been allocated.

* **int MLPXLayerSetActivation(int handle, int snapshotIndex, int layerIndex, int subscript, double value);**:
	Sets a particular subscript within the activations matrix of the given
	layer. If no activations matrix previously existed, one will be allocated
	and all other values initialized to 0.

* **int MLPXLayerGetActivation(int handle, int snapshotIndex, int layerIndex, int subscript, double\* value);**:
	Retrives the weight of a particular subscript within the activations matrix
	of the given layer. Will return an error status if the activations matrix
	has not yet been allocated.

* **int MLPXLayerSetDelta(int handle, int snapshotIndex, int layerIndex, int subscript, double value);**:
	Sets a particular subscript within the deltas matrix of the given
	layer. If no deltas matrix previously existed, one will be allocated
	and all other values initialized to 0.

* **int MLPXLayerGetDelta(int handle, int snapshotIndex, int layerIndex, int subscript, double\* value);**:
	Retrives the weight of a particular subscript within the deltas matrix
	of the given layer. Will return an error status if the deltas matrix
	has not yet been allocated.

* **int MLPXLayerSetBias(int handle, int snapshotIndex, int layerIndex, int subscript, double value);**:
	Sets a particular subscript within the biases matrix of the given
	layer. If no biases matrix previously existed, one will be allocated
	and all other values initialized to 0.

* **int MLPXLayerGetBias(int handle, int snapshotIndex, int layerIndex, int subscript, double\* value);**:
	Retrives the weight of a particular subscript within the biases matrix
	of the given layer. Will return an error status if the biases matrix
	has not yet been allocated.

* **int MLPXLayerSetActivationFunction(int handle, int snapshotIndex, int layerIndex, char\* funct);**:
	Modifies the activation function of the specified layer. The new value
	is copied into Go memory, so `funct` can safely be free-ed by the
	caller.

* **int MLPXLayerGetActivationFunction(int handle, int snapshotIndex, int layerIndex, char\*\* funct);**:
	Retrieves the activation function of the specified layer in a newly
	malloc-ed C string.

## TODO

* Implement support for setting predecessor and successor IDs for layers.
* Implement support for validation.
* Implement support for saving out the MLPX file to disk.

## AUTHOR

Charles Daniels.

## COPYRIGHT

Copyright 2020 Jason Bakos, Charles Daniels, Philip Conrad, all rihts reserved.

## SEE ALSO

`mlpx(5)`
