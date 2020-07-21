// Package mlpx implements public API for the HeRC MLPX file format.
//
// See MLPX(4) for the complete format specification.
package mlpx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// MLPX represents an entire MLPX file
type MLPX struct {

	// Schema is used to represent the schema key
	Schema []interface{} `json: "schema"`

	// Snapshots is used to represent the snapshot table.
	Snapshots map[string]*Snapshot `json:"snapshots"`
}

// Snapshot represents a single snapshot definition
type Snapshot struct {

	// Layers is the list of layers in the snapshot.
	Layers map[string]*Layer `json: "layers"`
}

// Layer represents a single layer definition
type Layer struct {

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

// Validate checks the MLPX file for any errors. If none are found, it returns
// nil.
func (mlp *MLPX) Validate() error {

	version, err := mlp.Version()
	if err != nil {
		return err
	}

	if version != 0 {
		return fmt.Errorf("Unknown version number %d", version)
	}

	for snapid, snapshot := range mlp.Snapshots {
		for layerid, layer := range snapshot.Layers {
			// verify integrity of predecessor references
			if layerid != "input" {
				// input layers don't have predecessors
				if _, ok := snapshot.Layers[layer.Predecessor]; !ok {
					return fmt.Errorf("Snapshot '%s', layer '%s': predecessor '%s' references nonexistent layer",
						snapid, layerid, layer.Predecessor)
				}
			}

			// verify integrity of successor references
			if layerid != "output" {
				// output layers don't have successors
				if _, ok := snapshot.Layers[layer.Successor]; !ok {
					return fmt.Errorf("Snapshot '%s', layer '%s': predecessor '%s' references nonexistent layer",
						snapid, layerid, layer.Predecessor)
				}
			}

			// verify size of weights list
			if layerid != "input" && layer.Weights != nil {
				expect := layer.Neurons * snapshot.Layers[layer.Predecessor].Neurons
				if len(*layer.Weights) != expect {
					return fmt.Errorf("Snapshot '%s', layer '%s': weights array of length %d, should be %d",
						snapid, layerid, len(*layer.Weights), expect)
				}
			}

			// verify size of outputs list
			if layer.Outputs != nil {
				if len(*layer.Outputs) != layer.Neurons {
					return fmt.Errorf("Snapshot '%s', layer '%s': output array of length %d, should be %d",
						snapid, layerid, len(*layer.Outputs), layer.Neurons)
				}
			}

			// verify size of activation list
			if layer.Activations != nil {
				if len(*layer.Activations) != layer.Neurons {
					return fmt.Errorf("Snapshot '%s', layer '%s': activation array of length %d, should be %d",
						snapid, layerid, len(*layer.Activations), layer.Neurons)
				}
			}

			// verify size of deltas list
			if layer.Deltas != nil {
				if len(*layer.Deltas) != layer.Neurons {
					return fmt.Errorf("Snapshot '%s', layer '%s': delta array of length %d, should be %d",
						snapid, layerid, len(*layer.Deltas), layer.Neurons)
				}
			}

			// verify size of bias list
			if layer.Biases != nil {
				if len(*layer.Biases) != layer.Neurons {
					return fmt.Errorf("Snapshot '%s', layer '%s': bias array of length %d, should be %d",
						snapid, layerid, len(*layer.Biases), layer.Neurons)
				}
			}
		}
	}

	// if we have made it this far, then we know that each snapshot is
	// internally valid, so all we have left to do is verify that all
	// layers have the same topology.
	if len(mlp.Snapshots) < 2 {
		// if there are no snapshots, there is nothing to do, likewise
		// if there is only one snapshot, we can safely assume it
		// is isomorphic to itself
		return nil
	}

	// we choose a "key" snapshot which we will compare with all the others
	keyid := ""
	for k := range mlp.Snapshots {
		keyid = k
		break
	}
	key := mlp.Snapshots[keyid]

	// we now compare each other snapshot against our chosen key
	for snapid, snapshot := range mlp.Snapshots {

		// make sure all layers in the key are in the snapshot
		for layerid := range key.Layers {
			if _, ok := snapshot.Layers[layerid]; !ok {
				return fmt.Errorf("Snapshot '%s' and '%s' are not isomorphic, snapshot '%s' has layer ID '%s, but snapshot '%s' does not",
					keyid, snapid, keyid, layerid, snapid)
			}
		}

		// and the converse
		for layerid := range snapshot.Layers {
			if _, ok := key.Layers[layerid]; !ok {
				return fmt.Errorf("Snapshot '%s' and '%s' are not isomorphic, snapshot '%s' has layer ID '%s, but snapshot '%s' does not",
					snapid, keyid, snapid, layerid, keyid)
			}
		}

		// finally, check the neuron counts, which also imply the
		// other member fields, and we have already validated they
		// are sized in a way appropriate for their neuron counts
		for layerid, layer := range key.Layers {
			if layer.Neurons != snapshot.Layers[layerid].Neurons {
				return fmt.Errorf("Snapshot '%s' and '%s' have different numbers of neurons (%d, and %d respectively) for layer '%s'",
					keyid, snapid, layer.Neurons, snapshot.Layers[layerid].Neurons, layerid)
			}
		}
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
	version, ok := mlp.Schema[1].(int)
	if !ok {
		return -1, fmt.Errorf("Schema component 1 is not an integer: %v", mlp.Schema[1])
	}

	if schema != "mlpx" {
		return -2, fmt.Errorf("Schema component 0 is '%s', expected 'mlpx'", schema)
	}

	return version, nil
}

// ToJSON converts an existing MLPX object to a JSON string and returns it.
func (mlp *MLPX) ToJSON() ([]byte, error) {
	b, err := json.MarshalIndent(mlp, "", "\t")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// WriteJSON calls ToJSON() and then overwrites the specified path with it's
// return.
func (mlp *MLPX) WriteJSON(path string) error {
	b, err := mlp.ToJSON()
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

// FromJSON reads an in-memory JSON string and generates an MLPX object. It
// does not validate the data which is read.
func FromJSON(data []byte) (*MLPX, error) {
	mlp := &MLPX{}
	err := json.Unmarshal(data, mlp)
	if err != nil {
		return nil, err
	}
	return mlp, err
}

// ReadJSON is a utility function which reads a file from disk, then calls
// FromJSON() on it. It does not validate the MLPX file.
func ReadJSON(path string) (*MLPX, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	mlp, err := FromJSON(data)
	if err != nil {
		return nil, err
	}

	return mlp, nil
}
