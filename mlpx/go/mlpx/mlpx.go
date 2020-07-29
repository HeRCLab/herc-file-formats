// Package mlpx implements public API for the HeRC MLPX file format.
//
// See MLPX(4) for the complete format specification.
package mlpx

import (
	"fmt"
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

// NextSnapshotID returns the next canonical snapshot ID. If no snapshots have
// been taken, it returns "initializer". It otherwise returns an integer
// starting from 0 and increasing contiguously.
func (mlp *MLPX) NextSnapshotID() string {
	snapids := mlp.SortedSnapshotIDs()

	if len(snapids) == 0 {
		return "intializer"
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

// SortedSnapshotIDs returns the list of snapshot IDs, sorted in the canonical
// order for MLPX. That is, the "initializer" snapshot sorts before everything
// else, and numeric snapshot IDs are sorted by numeric value, rather than
// by string comparison"
func (mlp *MLPX) SortedSnapshotIDs() []string {
	snapids := make([]string, 0)
	for k := range mlp.Snapshots {
		snapids = append(snapids, k)
	}

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

// EnsureWeights guarantees that the weights matrix for the layer is non-nil
func (layer *Layer) EnsureWeights() {
	if layer.Weights == nil {
		w := make([]float64, layer.Neurons*layer.Parent.Layers[layer.Predecessor].Neurons)
		layer.Weights = &w
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
