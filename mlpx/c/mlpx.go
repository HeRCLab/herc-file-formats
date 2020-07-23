package main

// This file implement a C API for for the MLPX library.
//
// The following conventions are used
//
// * No memory is shared, any arrays are copied at FFI boundaries
// * Returns are always an integer 0 on success, or nonzero on failure
// * Error information is retrieved using MLPXGetError()

//#include <stdlib.h>
import "C"

import (
	"fmt"

	"github.com/herclab/herc-file-formats/mlpx/go/mlpx"
)

var nexthandle C.int = 0
var handles map[C.int]*mlpx.MLPX = map[C.int]*mlpx.MLPX{}
var lastError string = ""

//export MLPXGetError
func MLPXGetError() *C.char {
	return C.CString(lastError)
}

//export MLPXOpen
func MLPXOpen(path *C.char, handle *C.int) C.int {
	h := nexthandle
	nexthandle++

	mlp, err := mlpx.ReadJSON(C.GoString(path))
	if err != nil {
		lastError = fmt.Sprintf("%v", err)
		return 1
	}

	handles[h] = mlp
	*handle = h
	return 0
}

//export MLPXClose
func MLPXClose(handle C.int) C.int {
	if _, ok := handles[handle]; ok {
		delete(handles, handle)
		return 0
	} else {
		lastError = fmt.Sprintf("unknown handle %d", int(handle))
		return 1
	}
}

// Returns nil on error and sets lastError
func getMLP(handle C.int) *mlpx.MLPX {
	mlp, ok := handles[handle]
	if !ok {
		lastError = fmt.Sprintf("unknown handle %d", int(handle))
		return nil
	}
	return mlp

}

//export MLPXGetNumSnapshots
func MLPXGetNumSnapshots(handle C.int, snapc *C.int) C.int {
	mlp := getMLP(handle)
	if mlp == nil {
		return 1
	}

	snapIDs := mlp.SortedSnapshotIDs()
	*snapc = C.int(len(snapIDs))

	return 0
}

//export MLPXGetSnapshotIDByIndex
func MLPXGetSnapshotIDByIndex(handle C.int, index C.int, id **C.char) C.int {
	mlp := getMLP(handle)
	if mlp == nil {
		return 1
	}

	snapIDs := mlp.SortedSnapshotIDs()

	if index < 0 || index >= C.int(len(snapIDs)) {
		lastError = fmt.Sprintf("snapshot index out of bounds: %d", index)
		return 1
	}

	*id = C.CString(snapIDs[index])
	return 0
}

//export MLPXGetSnapshotIndexByID
func MLPXGetSnapshotIndexByID(handle C.int, id *C.char, index *C.int) C.int {
	mlp := getMLP(handle)
	if mlp == nil {
		return 1
	}

	snapIDs := mlp.SortedSnapshotIDs()
	for i, v := range snapIDs {
		if v == C.GoString(id) {
			*index = C.int(i)
			return 0
		}
	}

	lastError = fmt.Sprintf("no such snapshot ID '%s'", C.GoString(id))
	return 1
}

// Either get the snapshot, or return nil and set lastError
func getSnapshot(handle, snapshotIndex C.int) *mlpx.Snapshot {
	mlp, ok := handles[handle]
	if !ok {
		lastError = fmt.Sprintf("unknown handle %d", int(handle))
		return nil
	}

	snapIDs := mlp.SortedSnapshotIDs()
	if snapshotIndex < 0 || snapshotIndex >= C.int(len(snapIDs)) {
		lastError = fmt.Sprintf("snapshot index out of bounds: %d", snapshotIndex)
		return nil
	}

	snapshot := mlp.Snapshots[snapIDs[snapshotIndex]]
	return snapshot
}

//export MLPXSnapshotGetNumLayers
func MLPXSnapshotGetNumLayers(handle C.int, snapshotIndex C.int, layerc *C.int) C.int {
	snapshot := getSnapshot(handle, snapshotIndex)
	if snapshot == nil {
		return 1
	}

	*layerc = C.int(len(snapshot.SortedLayerIDs()))

	return 0
}

//export MLPXLayerGetIndexByID
func MLPXLayerGetIndexByID(handle C.int, snapshotIndex C.int, id *C.char, index *C.int) C.int {
	snapshot := getSnapshot(handle, snapshotIndex)
	if snapshot == nil {
		return 1
	}

	for i, v := range snapshot.SortedLayerIDs() {
		if v == C.GoString(id) {
			*index = C.int(i)
			return 0
		}
	}

	lastError = fmt.Sprintf("unknown layer ID '%s'", C.GoString(id))
	return 1
}

// Get the layer, or set lastError and return nil
func getLayer(handle, snapshotIndex, layerIndex C.int) *mlpx.Layer {
	snapshot := getSnapshot(handle, snapshotIndex)
	if snapshot == nil {
		return nil
	}

	layerids := snapshot.SortedLayerIDs()

	if int(layerIndex) < 0 || int(layerIndex) >= len(layerids) {
		lastError = fmt.Sprintf("layer index out of bounds %d", int(layerIndex))
		return nil
	}

	return snapshot.Layers[layerids[int(layerIndex)]]
}

//export MLPXLayerGetIDByIndex
func MLPXLayerGetIDByIndex(handle, snapshotIndex, index C.int, id **C.char) C.int {
	layer := getLayer(handle, snapshotIndex, index)
	if layer == nil {
		return 1
	}

	*id = C.CString(layer.ID)
	return 0
}

//export MLPXLayerGetNeurons
func MLPXLayerGetNeurons(handle, snapshotIndex, layerIndex C.int, neuronc *C.int) C.int {
	layer := getLayer(handle, snapshotIndex, layerIndex)
	if layer == nil {
		return 1
	}

	*neuronc = C.int(layer.Neurons)
	return 0
}

//export MLPXMakeMLPX
func MLPXMakeMLPX(handle *C.int) C.int {
	h := nexthandle
	nexthandle++

	mlp := mlpx.MakeMLPX()

	handles[h] = mlp
	*handle = h
	return 0
}

//export MLPXMakeSnapshot
func MLPXMakeSnapshot(handle C.int, id *C.char) C.int {
	mlp := getMLP(handle)
	if mlp == nil {
		return 1
	}

	err := mlp.MakeSnapshot(C.GoString(id))
	if err != nil {
		lastError = fmt.Sprintf("%v", err)
		return 1
	}

	return 0
}

//export MLPXMakeIsomorphicSnapshot
func MLPXMakeIsomorphicSnapshot(handle C.int, id *C.char, toSnapshotIndex C.int) C.int {
	mlp := getMLP(handle)
	if mlp == nil {
		return 1
	}

	snapshot := getSnapshot(handle, toSnapshotIndex)
	if snapshot == nil {
		return 1
	}

	err := mlp.MakeIsomorphicSnapshot(C.GoString(id), snapshot.ID)
	if err != nil {
		lastError = fmt.Sprintf("%v", err)
		return 1
	}

	return 0
}

//export MLPXGetInitializerSnapshotIndex
func MLPXGetInitializerSnapshotIndex(handle C.int, index *C.int) C.int {
	mlp := getMLP(handle)
	if mlp == nil {
		return 1
	}

	initializer, err := mlp.Initializer()
	if err != nil {
		lastError = fmt.Sprintf("%v", err)
		return 1
	}

	snapIDs := mlp.SortedSnapshotIDs()
	for i, v := range snapIDs {
		if v == initializer.ID {
			*index = C.int(i)
			return 0
		}
	}

	lastError = fmt.Sprintf("No initializer found")
	return 1

}

//export MLPXLayerGetPredecessorIndex
func MLPXLayerGetPredecessorIndex(handle, snapshotIndex, layerIndex C.int, index *C.int) C.int {
	layer := getLayer(handle, snapshotIndex, layerIndex)
	if layer == nil {
		return 1
	}

	return MLPXLayerGetIndexByID(handle, snapshotIndex, C.CString(layer.Predecessor), index)
}

//export MLPXLayerGetSuccessorIndex
func MLPXLayerGetSuccessorIndex(handle, snapshotIndex, layerIndex C.int, index *C.int) C.int {
	layer := getLayer(handle, snapshotIndex, layerIndex)
	if layer == nil {
		return 1
	}

	return MLPXLayerGetIndexByID(handle, snapshotIndex, C.CString(layer.Successor), index)
}

//export MLPXLayerSetWeight
func MLPXLayerSetWeight(handle, snapshotIndex, layerIndex, subscript C.int, value C.double) C.int {
	layer := getLayer(handle, snapshotIndex, layerIndex)
	if layer == nil {
		return 1
	}

	layer.EnsureWeights()

	if subscript < 0 || subscript > C.int(len(*layer.Weights)) {
		lastError = fmt.Sprintf("subscript '%s' out of bounds for layer ID '%s' weights", subscript, layer.ID)
		return 1
	}

	(*layer.Weights)[int(subscript)] = float64(value)
	return 0
}

//export MLPXLayerGetWeight
func MLPXLayerGetWeight(handle, snapshotIndex, layerIndex, subscript C.int, value *C.double) C.int {
	layer := getLayer(handle, snapshotIndex, layerIndex)
	if layer == nil {
		return 1
	}

	if layer.Weights == nil {
		lastError = fmt.Sprintf("layer ID '%s' has no weights", layer.ID)
		return 1
	}

	if subscript < 0 || subscript > C.int(len(*layer.Weights)) {
		lastError = fmt.Sprintf("subscript '%s' out of bounds for layer ID weights'%s'", subscript, layer.ID)
		return 1
	}

	*value = C.double((*layer.Weights)[int(subscript)])
	return 0
}

//export MLPXLayerSetOutput
func MLPXLayerSetOutput(handle, snapshotIndex, layerIndex, subscript C.int, value C.double) C.int {
	layer := getLayer(handle, snapshotIndex, layerIndex)
	if layer == nil {
		return 1
	}

	layer.EnsureOutputs()

	if subscript < 0 || subscript >= C.int(len(*layer.Outputs)) {
		lastError = fmt.Sprintf("subscript '%s' out of bounds for layer ID outputs'%s'", subscript, layer.ID)
		return 1
	}

	(*layer.Outputs)[int(subscript)] = float64(value)
	return 0
}

//export MLPXLayerGetOutput
func MLPXLayerGetOutput(handle, snapshotIndex, layerIndex, subscript C.int, value *C.double) C.int {
	layer := getLayer(handle, snapshotIndex, layerIndex)
	if layer == nil {
		return 1
	}

	if layer.Outputs == nil {
		lastError = fmt.Sprintf("layer ID '%s' has no output", layer.ID)
		return 1
	}

	if subscript < 0 || subscript >= C.int(len(*layer.Outputs)) {
		lastError = fmt.Sprintf("subscript '%s' out of bounds for layer ID outputs'%s'", subscript, layer.ID)
		return 1
	}

	*value = C.double((*layer.Outputs)[int(subscript)])
	return 0
}

//export MLPXLayerSetActivation
func MLPXLayerSetActivation(handle, snapshotIndex, layerIndex, subscript C.int, value C.double) C.int {
	layer := getLayer(handle, snapshotIndex, layerIndex)
	if layer == nil {
		return 1
	}

	layer.EnsureActivations()

	if subscript < 0 || subscript >= C.int(len(*layer.Activations)) {
		lastError = fmt.Sprintf("subscript '%s' out of bounds for layer ID activation '%s'", subscript, layer.ID)
		return 1
	}

	(*layer.Activations)[int(subscript)] = float64(value)
	return 0
}

//export MLPXLayerGetActivation
func MLPXLayerGetActivation(handle, snapshotIndex, layerIndex, subscript C.int, value *C.double) C.int {
	layer := getLayer(handle, snapshotIndex, layerIndex)
	if layer == nil {
		return 1
	}

	if layer.Activations == nil {
		lastError = fmt.Sprintf("layer ID '%s' has no activations", layer.ID)
		return 1
	}

	if subscript < 0 || subscript >= C.int(len(*layer.Activations)) {
		lastError = fmt.Sprintf("subscript '%s' out of bounds for layer ID activations '%s'", subscript, layer.ID)
		return 1
	}

	*value = C.double((*layer.Activations)[int(subscript)])
	return 0
}

//export MLPXLayerSetDelta
func MLPXLayerSetDelta(handle, snapshotIndex, layerIndex, subscript C.int, value C.double) C.int {
	layer := getLayer(handle, snapshotIndex, layerIndex)
	if layer == nil {
		return 1
	}

	layer.EnsureDeltas()

	if subscript < 0 || subscript >= C.int(len(*layer.Deltas)) {
		lastError = fmt.Sprintf("subscript '%s' out of bounds for layer ID delta '%s'", subscript, layer.ID)
		return 1
	}

	(*layer.Deltas)[int(subscript)] = float64(value)
	return 0
}

//export MLPXLayerGetDelta
func MLPXLayerGetDelta(handle, snapshotIndex, layerIndex, subscript C.int, value *C.double) C.int {
	layer := getLayer(handle, snapshotIndex, layerIndex)
	if layer == nil {
		return 1
	}

	if layer.Deltas == nil {
		lastError = fmt.Sprintf("layer ID '%s' has no deltas", layer.ID)
		return 1
	}

	if subscript < 0 || subscript >= C.int(len(*layer.Deltas)) {
		lastError = fmt.Sprintf("subscript '%s' out of bounds for layer ID deltas '%s'", subscript, layer.ID)
		return 1
	}

	*value = C.double((*layer.Deltas)[int(subscript)])
	return 0
}

//export MLPXLayerSetBias
func MLPXLayerSetBias(handle, snapshotIndex, layerIndex, subscript C.int, value C.double) C.int {
	layer := getLayer(handle, snapshotIndex, layerIndex)
	if layer == nil {
		return 1
	}

	layer.EnsureBiases()

	if subscript < 0 || subscript >= C.int(len(*layer.Biases)) {
		lastError = fmt.Sprintf("subscript '%s' out of bounds for layer ID bias '%s'", subscript, layer.ID)
		return 1
	}

	(*layer.Biases)[int(subscript)] = float64(value)
	return 0
}

//export MLPXLayerGetBias
func MLPXLayerGetBias(handle, snapshotIndex, layerIndex, subscript C.int, value *C.double) C.int {
	layer := getLayer(handle, snapshotIndex, layerIndex)
	if layer == nil {
		return 1
	}

	if layer.Biases == nil {
		lastError = fmt.Sprintf("layer ID '%s' has no biases", layer.ID)
		return 1
	}

	if subscript < 0 || subscript >= C.int(len(*layer.Biases)) {
		lastError = fmt.Sprintf("subscript '%s' out of bounds for layer ID biases '%s'", subscript, layer.ID)
		return 1
	}

	*value = C.double((*layer.Biases)[int(subscript)])
	return 0
}

func main() {}
