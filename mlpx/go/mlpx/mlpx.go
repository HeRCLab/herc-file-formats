// Package mlpx implements public API for the HeRC MLPX file format.
//
// See MLPX(4) for the complete format specification.
package mlpx

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

// MLPX represents an entire MLPX file
type MLPX struct {

	// Schema is used to represent the schema key
	//
	// Caveat: because of how encoding/json works, the version level
	// is usually encoded as a float64.
	Schema []interface{} `json: "schema"`

	// Snapshots is used to represent the snapshot table.
	Snapshots map[string]*Snapshot `json:"snapshots"`
}

// MakeMLPX creates a new, empty MLPX object
func MakeMLPX() *MLPX {
	return &MLPX{
		Schema:    []interface{}{"mlpx", float64(0)},
		Snapshots: make(map[string]*Snapshot),
	}
}

// IsomorphicDuplicate creates a new MLPX object with a single snapshot with
// the id snapid, which is topologically identical to the source MLPX object.
//
// No snapshot fields are initialized in the duplicate.
//
// The snapshots alpha value will be initialized to 0.0.
func (mlp *MLPX) IsomorphicDuplicate(snapid string) (*MLPX, error) {
	dup := MakeMLPX()

	err := dup.MakeSnapshot(snapid, 0.0)
	if err != nil {
		return nil, err
	}

	if len(mlp.SortedSnapshotIDs()) == 0 {
		return dup, nil
	}

	sourceSnapID := mlp.SortedSnapshotIDs()[0]
	sourceSnap := mlp.Snapshots[sourceSnapID]

	for _, layerid := range sourceSnap.SortedLayerIDs() {
		sourceLayer := sourceSnap.Layers[layerid]
		err := dup.Snapshots[snapid].MakeLayer(
			layerid,
			sourceLayer.Neurons,
			sourceLayer.Predecessor,
			sourceLayer.Successor)

		if err != nil {
			return nil, err
		}
	}

	return dup, nil
}

// Diff returns a list of differences between the given MLPX objects.
//
// The indent parameter will be used to indent any hierarchical data, if
// applicable. The suggested value is "\t".
//
// The epsilon parameter defines the maximum difference of two floating point
// numbers before this algorithm considers them to be different. This should
// usually be a very small number.
func (mlp *MLPX) Diff(other *MLPX, indent string, epsilon float64) []string {
	diffs := []string{}

	// find all common snapshot IDs
	commonSnaps := make(map[string]bool)
	for _, snapid := range mlp.SortedSnapshotIDs() {
		_, ok := other.Snapshots[snapid]
		if ok {
			commonSnaps[snapid] = true
		} else {
			// snapshot in mlp but not other
			diffs = append(diffs,
				fmt.Sprintf("base MLPX has snapshot ID '%s', but other MLPX does not", snapid))
		}
	}
	for _, snapid := range other.SortedSnapshotIDs() {
		_, ok := mlp.Snapshots[snapid]
		if ok {
			commonSnaps[snapid] = true
		} else {
			// snapshot in other but not mlp
			diffs = append(diffs,
				fmt.Sprintf("other MLPX has snapshot ID '%s', but base MLPX does not", snapid))
		}
	}

	// get a nicely sorted list of common snapshot IDs
	commonList := make([]string, 0)
	for k := range commonSnaps {
		commonList = append(commonList, k)
	}
	commonList = SortSnapshotIDs(commonList)

	// now recurse into the common snapshots
	for _, snapid := range commonList {
		// find the differences
		baseSnap := mlp.Snapshots[snapid]
		otherSnap := other.Snapshots[snapid]
		snapDiff := baseSnap.Diff(otherSnap, indent, epsilon)

		// apply indent
		for i, v := range snapDiff {
			snapDiff[i] = fmt.Sprintf("%s%s", indent, v)
		}

		// fold into the main diff list
		if len(snapDiff) > 0 {
			diffs = append(diffs,
				fmt.Sprintf("Snapshot ID '%s' differs", snapid))
			diffs = append(diffs, snapDiff...)
		}

	}

	return diffs
}

// Snapshot represents a single snapshot definition
type Snapshot struct {

	// Parent is the MLPX object which this snapshot belongs to.
	//
	// DANGER: modify this field with care, changing this may corrupt the
	// in-memory representation of the MLP.
	Parent *MLPX `json:"-"`

	// ID is the snapshot ID
	//
	// DANGER: modify this field with care, changing this may corrupt the
	// in-memory representation of the MLP.
	ID string `json:"-"`

	// Alpha is the learning rate within the layer
	Alpha float64 `json:"alpha"`

	// Layers is the list of layers in the snapshot.
	Layers map[string]*Layer `json:"layers"`
}

// Diff returns a list of differences between the given Snapshot objects.
func (snapshot *Snapshot) Diff(other *Snapshot, indent string, epsilon float64) []string {
	diffs := []string{}

	// compare IDs
	if snapshot.ID != other.ID {
		diffs = append(diffs,
			fmt.Sprintf("IDs do not match: base '%s', other '%s'", snapshot.ID, other.ID))
	}

	// compare Alpha values
	if math.Abs(snapshot.Alpha-other.Alpha) > epsilon {
		diffs = append(diffs,
			fmt.Sprintf("Base snapshot Alpha %f does not match other Alpha %f",
				snapshot.Alpha, other.Alpha))
	}

	// find all layer IDs in common between the snapshots
	commonLayers := make(map[string]bool)
	for _, layerid := range snapshot.SortedLayerIDs() {
		_, ok := other.Layers[layerid]
		if ok {
			commonLayers[layerid] = true
		} else {
			// layer in snapshot but not other
			diffs = append(diffs,
				fmt.Sprintf("base snapshot has layer ID '%s', but other snapshot does not", layerid))
		}
	}
	for _, layerid := range other.SortedLayerIDs() {
		_, ok := snapshot.Layers[layerid]
		if ok {
			commonLayers[layerid] = true
		} else {
			// layer in other but not snapshot
			diffs = append(diffs,
				fmt.Sprintf("other snapshot has layer ID '%s', but base snapshot does not", layerid))
		}
	}

	// We cannot sort this as nicely, since layers cannot be sorted outside
	// of the context of a specific snapshot. Instead we sort one of them
	// and then remove all the indices that aren't in the common list. This
	// does mean that the layer IDs will only be sorted within the context
	// of the base snapshot.
	commonList := make([]string, 0)
	for _, v := range snapshot.SortedLayerIDs() {
		_, ok := commonLayers[v]
		if ok {
			commonList = append(commonList, v)
		}
	}

	// Now we recurse into the common layers
	for _, layerid := range commonList {
		baseLayer := snapshot.Layers[layerid]
		otherLayer := other.Layers[layerid]
		layerDiff := baseLayer.Diff(otherLayer, indent, epsilon)

		// apply indent
		for i, v := range layerDiff {
			layerDiff[i] = fmt.Sprintf("%s%s", indent, v)
		}

		// fold into main diff list
		if len(layerDiff) > 0 {
			diffs = append(diffs,
				fmt.Sprintf("Layer ID '%s' differs", layerid))
			diffs = append(diffs, layerDiff...)
		}
	}

	return diffs
}

// NextSnapshotID returns the next canonical snapshot ID. If no snapshots have
// been taken, it returns "initializer". It otherwise returns an integer
// starting from 0 and increasing contiguously.
func (mlp *MLPX) NextSnapshotID() string {
	snapids := mlp.SortedSnapshotIDs()

	if len(snapids) == 0 {
		return "initializer"
	}

	return fmt.Sprintf("%d", len(snapids))
}

// MakeSnapshot creates a new, empty snapshot in the given MLPX object
func (mlp *MLPX) MakeSnapshot(id string, alpha float64) error {
	if _, ok := mlp.Snapshots[id]; ok {
		return fmt.Errorf("Snapshot with ID '%s' already exists", id)
	}

	mlp.Snapshots[id] = &Snapshot{
		Parent: mlp,
		ID:     id,
		Alpha:  alpha,
		Layers: make(map[string]*Layer),
	}

	return nil
}

// MustMakeSnapshot is a wrapper around MakeSnapshot which errors if it fails
func (mlp *MLPX) MustMakeSnapshot(id string, alpha float64) {
	err := mlp.MakeSnapshot(id, alpha)
	if err != nil {
		panic(err)
	}
}

// MakeIsomorphicSnapshot will create a new snapshot which is topologically
// identical to the one specified. This is the preferred way of creating
// snapshots, once the first has been defined, to guarantee that all snapshots
// are isomorphic, which the spec requires.
func (mlp *MLPX) MakeIsomorphicSnapshot(id, to string) error {
	if _, ok := mlp.Snapshots[to]; !ok {
		return fmt.Errorf("Specified snapshot '%s' does not exist", to)
	}

	err := mlp.MakeSnapshot(id, mlp.Snapshots[to].Alpha)
	if err != nil {
		return err
	}

	for layerid, layer := range mlp.Snapshots[to].Layers {
		err := mlp.Snapshots[id].MakeLayer(layerid, layer.Neurons, layer.Predecessor, layer.Successor)
		if err != nil {
			return err
		}
	}

	return nil
}

// MustMakeIsomorphicSnapshot is a wrapper around MakeIsomorphicSnapshot
// that calls panic() if it errors.
func (mlp *MLPX) MustMakeIsomorphicSnapshot(id, to string) {
	err := mlp.MakeIsomorphicSnapshot(id, to)
	if err != nil {
		panic(err)
	}
}

// SortSnapshotIDs applies the canonical sorting algorithm for snapshots to
// the given list of snapshot IDs
func SortSnapshotIDs(snapids []string) []string {
	sort.Slice(snapids, func(i, j int) bool {
		if snapids[i] == "initializer" {
			return true
		}

		if snapids[j] == "initializer" {
			return false
		}

		ii, ierr := strconv.Atoi(snapids[i])
		ji, jerr := strconv.Atoi(snapids[j])

		if ierr == nil && jerr == nil {
			// both IDs are numeric
			return ii < ji
		} else if ierr == nil && jerr != nil {
			// i is numeric, j is not, so i sorts first
			return true
		} else if ierr != nil && jerr == nil {
			// i is non-numeric, j is numeric, so j sorts first
			return false
		} else { //ierr != nil && jerr != nil
			return snapids[i] < snapids[j]
		}
	})
	return snapids
}

// SortedSnapshotIDs returns the list of snapshot IDs, sorted in the canonical
// order for MLPX. That is, the "initializer" snapshot sorts before everything
// else, and numeric snapshot IDs are sorted by numeric value, rather than
// by string comparison"
func (mlp *MLPX) SortedSnapshotIDs() []string {
	snapids := make([]string, 0)
	for k := range mlp.Snapshots {
		snapids = append(snapids, k)
	}

	snapids = SortSnapshotIDs(snapids)

	return snapids
}

// Initializer returns the initializer snapshot if any
func (mlp *MLPX) Initializer() (*Snapshot, error) {
	ids := mlp.SortedSnapshotIDs()
	if len(ids) < 1 {
		return nil, fmt.Errorf("No snapshots available")
	}
	return mlp.Snapshots[ids[0]], nil
}

// Latest returns the most recent snapshot, if any
func (mlp *MLPX) Latest() (*Snapshot, error) {
	ids := mlp.SortedSnapshotIDs()
	if len(ids) < 1 {
		return nil, fmt.Errorf("No snapshots available")
	}
	return mlp.Snapshots[ids[len(ids)-1]], nil
}

// SortedLayerIDs returns a list of layer IDs in sorted order. IDs are sorted
// by their topology.
//
// If the MLPX is invalid, then the behavior of this function is undefined.
// In particular, cycles are not valid in MLPX, and may cause unusual behavior.
func (snapshot *Snapshot) SortedLayerIDs() []string {
	layerids := make([]string, 0)

	// first we need to find the "first" layer, usually this is the input
	// layer
	firstID := ""
	for layerid, layer := range snapshot.Layers {
		// we have to pick *something*, even if it's nonsense
		firstID = layerid
		if layerid == "input" {
			break
		} else if _, ok := snapshot.Layers[layer.Predecessor]; !ok {
			// user is doing something naughty, but we'll allow it
			// and pretend the first layer we find with no valid
			// predecessor is the "input"
			break
		}
	}

	first, ok := snapshot.Layers[firstID]
	if !ok {
		// this implies there are no layers
		return layerids
	}

	current := first

	for {

		// detect cycles and bail out
		for _, v := range layerids {
			if v == current.ID {
				return layerids
			}
		}

		layerids = append(layerids, current.ID)

		next, ok := snapshot.Layers[current.Successor]
		if !ok {
			// we're done, we got to the output layer
			return layerids
		}

		current = next
	}

	return layerids

}

// Successor returns the successor of a given snapshot, being the
// snapshot which occurs next after the specified one.
func (snapshot *Snapshot) Successor(id string) (*Snapshot, error) {
	mlp := snapshot.Parent

	snapids := mlp.SortedSnapshotIDs()

	for i, v := range snapids {
		if id == v {
			if (i + 1) >= len(snapids) {
				return nil, fmt.Errorf("Snapshot '%s' is the final snapshot available", id)
			}
			return mlp.Snapshots[snapids[i+1]], nil
		}
	}

	return nil, fmt.Errorf("no such snapshot '%s'", id)
}

// Predecessor returns the predecessor of a given snapshot, being the
// snapshot which occurs next after the specified one.
func (snapshot *Snapshot) Predecessor(id string) (*Snapshot, error) {
	mlp := snapshot.Parent
	snapids := make([]string, 0)
	for k := range mlp.Snapshots {
		snapids = append(snapids, k)
	}

	sort.Strings(snapids)
	for i, v := range snapids {
		if id == v {
			if (i - 1) < 0 {
				return nil, fmt.Errorf("Snapshot '%s' is the earliest snapshot available", id)
			}
			return mlp.Snapshots[snapids[i-1]], nil
		}
	}

	return nil, fmt.Errorf("no such snapshot '%s'", id)
}

// Layer represents a single layer definition
type Layer struct {
	// Parent is the Snapshot object which the layer belongs to.
	//
	// DANGER: modify this field with care, changing this may corrupt
	// the in-memory representation of the MLP.
	Parent *Snapshot `json:"-"`

	// ID is the layer ID
	//
	// DANGER: modify this field with care, changing this may corrupt the
	// in-memory representation of the MLP.
	ID string `json:"-"`

	// Predecessor is the preceding layer ID
	Predecessor string `json:"predecessor"`

	// Successor is the following layer ID
	Successor string `json:"successor"`

	// Neurons is the number of neurons in the layer
	Neurons int `json: "neurons"`

	// Weights is the weights list for the layers
	Weights *[]float64 `json:"weights"`

	// Outputs is the outputs list for the layer
	Outputs *[]float64 `json:"outputs"`

	// Activations is the activation value list for the layer
	Activations *[]float64 `json:"activations"`

	// Deltas is the deltas list for the layer
	Deltas *[]float64 `json:"deltas"`

	// Biases is the biases list for the layer
	Biases *[]float64 `json:"biases"`

	// ActivationFunction is the human-readable activation function used by
	// the layer
	ActivationFunction string `json:"activation_function"`
}

// header will only be shown if there is a difference
func diffList(base, other *[]float64, indent string, epsilon float64, header string) []string {
	if base != nil && other == nil {
		return []string{header,
			fmt.Sprintf("%sOther list is nil, base is non-nil", indent),
		}
	}

	if base == nil && other != nil {
		return []string{header,
			fmt.Sprintf("%sBase list is nil, other is non-nil", indent),
		}
	}

	if base == nil && other == nil {
		return []string{}
	}

	if len(*base) != len(*other) {
		return []string{header,
			fmt.Sprintf("%sLists are of different lengths: base %d, other %d", indent, len(*base), len(*other))}
	}

	ndiff := 0
	averagediff := 0.0
	for i, v := range *base {
		diff := math.Abs(v - (*other)[i])
		if diff > epsilon {
			ndiff++
			averagediff += diff
		}
	}

	if ndiff > 0 {
		averagediff = averagediff / float64(ndiff)
		return []string{header,
			fmt.Sprintf("%s%d values differ, with an average difference of %f", indent, ndiff, averagediff)}
	}

	return []string{}
}

// Diff returns a list of differences between the two given layers.
func (layer *Layer) Diff(other *Layer, indent string, epsilon float64) []string {
	diffs := []string{}

	// compare IDs
	if layer.ID != other.ID {
		diffs = append(diffs,
			fmt.Sprintf("IDs do not match: base '%s', other '%s'", layer.ID, other.ID))
	}

	// compare predecessors
	if layer.Predecessor != other.Predecessor {
		diffs = append(diffs,
			fmt.Sprintf("Predecessors do not match: base '%s', other '%s'", layer.Predecessor, other.Predecessor))
	}

	// compare Successors
	if layer.Successor != other.Successor {
		diffs = append(diffs,
			fmt.Sprintf("Successors do not match: base '%s', other '%s'", layer.Successor, other.Successor))
	}

	// compare Neurons
	if layer.Neurons != other.Neurons {
		diffs = append(diffs,
			fmt.Sprintf("Neurons do not match: base '%d', other '%d'", layer.Neurons, other.Neurons))
	}

	// compare lists
	diffs = append(diffs, diffList(layer.Weights, other.Weights, indent, epsilon, "Weight matrices do not match")...)
	diffs = append(diffs, diffList(layer.Outputs, other.Outputs, indent, epsilon, "Output matrices do not match")...)
	diffs = append(diffs, diffList(layer.Activations, other.Activations, indent, epsilon, "Activation matrices do not match")...)
	diffs = append(diffs, diffList(layer.Deltas, other.Deltas, indent, epsilon, "Delta matrices do not match")...)
	diffs = append(diffs, diffList(layer.Biases, other.Biases, indent, epsilon, "Bias matrices do not match")...)

	if layer.ActivationFunction != other.ActivationFunction {
		diffs = append(diffs, fmt.Sprintf(
			"Base activation function '%s' does not match other activation function '%s'",
			layer.ActivationFunction,
			other.ActivationFunction))
	}

	return diffs
}

// EnsureWeights guarantees that the weights matrix for the layer is non-nil
func (layer *Layer) EnsureWeights() {
	if layer.Weights == nil {
		pred, ok := layer.Parent.Layers[layer.Predecessor]
		if ok {
			w := make([]float64, layer.Neurons*pred.Neurons)
			layer.Weights = &w
		} else {
			// this is the input layer, so the weights don't really
			// matter
			w := make([]float64, layer.Neurons)
			layer.Weights = &w
		}
	}
}

// EnsureOutputs guarantees that the outputs matrix for the layer is non-nil
func (layer *Layer) EnsureOutputs() {
	if layer.Outputs == nil {
		o := make([]float64, layer.Neurons)
		layer.Outputs = &o
	}
}

// EnsureActivations guarantees that the outputs matrix for the layer is non-nil
func (layer *Layer) EnsureActivations() {
	if layer.Activations == nil {
		a := make([]float64, layer.Neurons)
		layer.Activations = &a
	}
}

// EnsureDeltas guarantees that the outputs matrix for the layer is non-nil
func (layer *Layer) EnsureDeltas() {
	if layer.Deltas == nil {
		d := make([]float64, layer.Neurons)
		layer.Deltas = &d
	}
}

// EnsureBiases guarantees that the outputs matrix for the layer is non-nil
func (layer *Layer) EnsureBiases() {
	if layer.Biases == nil {
		b := make([]float64, layer.Neurons)
		layer.Biases = &b
	}
}

// MakeLayer creates a new layer attached to the given snapshot. Where
// appropriate (input/output layers), pred or succ may be empty strings.
//
// Note that referential integrity of pred and succ IS NOT VERIFIED at this
// stage, since either or both referenced layers may not exist yet.
func (snapshot *Snapshot) MakeLayer(id string, neurons int, pred, succ string) error {
	if _, ok := snapshot.Layers[id]; ok {
		return fmt.Errorf("Referenced layer '%s' already exists", id)
	}

	snapshot.Layers[id] = &Layer{
		Parent:      snapshot,
		ID:          id,
		Predecessor: pred,
		Successor:   succ,
		Neurons:     neurons,
	}

	return nil
}

// MustMakeLayer is a wrapper around MakeLayer which panics if it encounters an
// error.
func (snapshot *Snapshot) MustMakeLayer(id string, neurons int, pred, succ string) {
	err := snapshot.MakeLayer(id, neurons, pred, succ)
	if err != nil {
		panic(err)
	}
}

// Version retrieves the MLPX schema version of the given file. It can error if
// the schema has the wrong number of components, if the schema is not
// "mlpx", or if the schema version is not an integer.
//
// If an error occurs, then the integer version level returned is undefined.
func (mlp *MLPX) Version() (int, error) {
	if len(mlp.Schema) != 2 {
		return -3, fmt.Errorf("Schema has incorrect number of components %d, expected 2", len(mlp.Schema))
	}

	schema, ok := mlp.Schema[0].(string)
	if !ok {
		return -1, fmt.Errorf("Schema component 0 is not a string: %v", mlp.Schema[0])
	}

	// By default, we're going to get a float64 back, even though per the
	// spec this is actually an integer.
	version, ok := mlp.Schema[1].(float64)
	if !ok {
		// make things a little more convenient for people using the
		// API by allowing integers.
		versionf, ok := mlp.Schema[1].(int)
		version = float64(versionf)
		if !ok {
			return -1, fmt.Errorf("Schema component 1 is not a number : %v", mlp.Schema[1])
		}
	}

	if schema != "mlpx" {
		return -2, fmt.Errorf("Schema component 0 is '%s', expected 'mlpx'", schema)
	}

	return int(version), nil
}
