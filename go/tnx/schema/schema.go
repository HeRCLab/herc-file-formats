// Package schema implements the TNX schema as described in tnx(4).
package schema

import (
	"encoding/json"
)

// TNX implements the top level TNX container object.
type TNX struct {
	Topology   Topology             `json: topology`
	Parameters map[string]Parameter `json: parameters`
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
	Id string `json: id`

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

type Parameter struct {
	Dimensions []int     `json: dimensions`
	Deltas     []float64 `json: deltas`
	Weights    []float64 `json: weights`
	Biases     []float64 `json: biases`
	Activation string    `json: activation`
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
