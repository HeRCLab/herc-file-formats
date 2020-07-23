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

//export MLPXGetNumLayers
func MLPXGetNumLayers(handle C.int, snapshotIndex C.int, layerc *C.int) C.int {
	snapshot := getSnapshot(handle, snapshotIndex)
	if snapshot == nil {
		return 1
	}

	*layerc = C.int(len(snapshot.SortedLayerIDs()))

	return 0
}

//export MLPXGetLayerIndexByID
func MLPXGetLayerIndexByID(handle C.int, snapshotIndex C.int, id *C.char, index *C.int) C.int {
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

//export MLPXGetLayerIDByIndex
func MLPXGetLayerIDByIndex(handle, snapshotIndex, index C.int, id **C.char) C.int {
	snapshot := getSnapshot(handle, snapshotIndex)
	if snapshot == nil {
		return 1
	}

	layerids := snapshot.SortedLayerIDs()

	if int(index) < 0 || int(index) >= len(layerids) {
		lastError = fmt.Sprintf("layer index out of bounds %d", int(index))
		return 1
	}

	*id = C.CString(layerids[int(index)])
	return 0
}

func main() {}
