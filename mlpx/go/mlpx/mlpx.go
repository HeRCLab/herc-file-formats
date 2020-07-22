// Package mlpx implements public API for the HeRC MLPX file format.
//
// See MLPX(4) for the complete format specification.
package mlpx

import (
	"fmt"
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

	// Layers is the list of layers in the snapshot.
	Layers map[string]*Layer `json: "layers"`
}

// MakeSnapshot creates a new, empty snapshot in the given MLPX object
func (mlp *MLPX) MakeSnapshot(id string) error {
	if _, ok := mlp.Snapshots[id]; ok {
		return fmt.Errorf("Snapshot with ID '%s' already exists", id)
	}

	mlp.Snapshots[id] = &Snapshot{
		Parent: mlp,
		ID:     id,
		Layers: make(map[string]*Layer),
	}

	return nil
}

// MakeIsomorphicSnapshot will create a new snapshot which is topologically
// identical to the one specified. This is the preferred way of creating
// snapshots, once the first has been defined, to guarantee that all snapshots
// are isomorphic, which the spec requires.
func (mlp *MLPX) MakeIsomorphicSnapshot(id, to string) error {
	err := mlp.MakeSnapshot(id)
	if err != nil {
		return err
	}

	if _, ok := mlp.Snapshots[to]; !ok {
		return fmt.Errorf("Specified snapshot '%s' does not exist", to)
	}

	for layerid, layer := range mlp.Snapshots[to].Layers {
		err := mlp.Snapshots[id].MakeLayer(layerid, layer.Neurons, layer.Predecessor, layer.Successor)
		if err != nil {
			return err
		}
	}

	return nil
}

// func (mlp *MLPX) GetSuccessor(id string) (*Snapshot, error) {
// }

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
	Predecessor string `json: "predecessor"`

	// Successor is the following layer ID
	Successor string `json: "successor"`

	// Neurons is the number of neurons in the layer
	Neurons int `json: "neurons"`

	// Weights is the weights list for the layers
	Weights *[]float64 `json: "weights"`

	// Outputs is the outputs list for the layer
	Outputs *[]float64 `json: "outputs"`

	// Activations is the activation value list for the layer
	Activations *[]float64 `json: "activations"`

	// Deltas is the deltas list for the layer
	Deltas *[]float64 `json: "deltas"`

	// Biases is the biases list for the layer
	Biases *[]float64 `json: "biases"`

	// ActivationFunction is the human-readable activation function used by
	// the layer
	ActivationFunction string `json: "activation_function"`
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
