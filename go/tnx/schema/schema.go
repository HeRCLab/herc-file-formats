// Package schema implements the TNX schema as described in tnx(4).
//
// The approach used here is a little inflexible, since it relies on all
// possible key types to be defined as appropriate Go structures. This
// is useful because it makes the structures easier to implement, but in the
// future this will probably need to change, since users of this library may
// wish to register their own key types or validation logic.
package schema

// This file defines the relevant data structures to represent a TNX file from
// disk. If you are looking for the"rehydrated" representation of a TNX file,
// you should go to the parent package. If you are looking for validation
// related code, see ./validation.go

import (
	"encoding/json"
)

// NOTE: the TNX parameters and snapshots tables have values as pointers
// because this makes Golang happy when assigning struct members for elements
// of the dictionaries, not for any other technical reason.
//
// However, other items that are points are so that they are nullable, due to
// being optional in the specification. For example, a snapshot which does
// not contain a matrix values can simply set it's matrix pointer to nil.

// TNX implements the top level TNX container object.
type TNX struct {
	// Topology defines to a topology definition as described in tnx(4)
	Topology Topology `json: topology`

	// Parameters defines a parameters table as described in tnx(4)
	Parameters map[string]*Parameter `json: parameters`

	// Snapshots defines a snapshots table as described in tnx(4)
	Snapshots map[string]*Snapshot `json: snapshots`

	// Schema is used to record the schema and version of the TNX.
	Schema []string `json: schema`

	// These fields are used to cache lookup operations on nodes and links.
	//
	// Because they are unexported, the JSON package should ignore them,
	// but we mark them as ignored anyway to be explicit about it.
	linkLookupCache   map[string]*Link `json:"-"`
	nodeIOLookupCache map[string]*Node `json:"-"`
	nodeLookupCache   map[string]*Node `json:"-"`
}

// Topology represents a TNX topology object.
type Topology struct {
	// Nodes is a list of Node objects.
	Nodes []Node `json: nodes`

	// Links is a list of Link objects.
	Links []Link `json: links`
}

// Node represents a TNX node object.
type Node struct {
	// Id should be a unique identification string, not shared by any other
	// TNX node, input, or output.
	ID string `json: id`

	// Operation should be one of the operation strings described in the
	// TNX specification.
	Operation string `json: operation`

	// Inputs should be a list of unique identifier strings.
	Inputs []string `json: inputs`

	// Outputs should be a list of unique identifier strings.
	Outputs []string `json: outputs`
}

// Link represents a TNX link object.
type Link struct {
	// Source must reference a TNX output ID.
	Source string `json: source`

	// Target must reference a TNX output ID.
	Target string `json: target`
}

// Parameter represents the set of all parameters for a specific node. Unused
// parameters should be left as nil.
type Parameter struct {
	// Dimensions represents a dimension list as described in tnx(4)
	Dimensions *[]int `json: dimensions`

	// Deltas represents a deltas list as described in tnx(4)
	Deltas *[]float64 `json: deltas`

	// Weights represents a weights list as described in tnx(4)
	Weights *[]float64 `json: weights`

	// Biases represents a biases list as described in tnx(4)
	Biases *[]float64 `json: biases`

	// Activation represents an activation reference as described in tnx(4)
	Activation *string `json: activation`
}

// Matrix represents a matrix type snapshot value, as described in tnx(4)
type Matrix struct {
	// Name represents the matrix name as described in tnx(4)
	Name string `json: name`

	// Dimensions represents a dimension list as described in tnx(4)
	Dimensions []int `json: dimensions`

	// Data represents a data list as described in tnx(4)
	Data []float64 `json: data`
}

// Snapshot represents a single snapshot object as described in tnx(4)
type Snapshot struct {
	Matrix map[string]*Matrix `json: matrix`
}

// FromJSON de-serializes a TNX object from a JSON file. The TNX returned
// is guaranteed to be well formed, but may not be valid.
func FromJSON(data []byte) (*TNX, error) {
	t := &TNX{}
	err := json.Unmarshal(data, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
