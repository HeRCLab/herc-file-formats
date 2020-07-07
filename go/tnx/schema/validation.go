package schema

import (
	"fmt"
)

// ParameterValidators - This list is initialized with the official validation
// functions, but users may add additional ones to enforce custom constraints.
// Each function in this list is run on each defined parameter. The TNX object
// given will be the parent of the given parameter, and should be used ONLY in
// a read-only capacity to extract information needed to verify the parameter.
//
// No function in this list should modify either the parameter, or the TNX.  If
// the parameter is valid, it should return nil, and otherwise an error.
var ParameterValidators []func(*TNX, *Parameter) error

// SnapshotValidators - as with ParameterValidators, but instead applies to
// snapshot objects.
var SnapshotValidators []func(*TNX, *Snapshot) error

// Validate checks if the TNX is valid. In order for it to have been loaded by
// the JSON decoder, it must have been well formed. T
func Validate(tnx *TNX) error {
	return nil
}

// ValidateSchema ensures that the TNX schema is supported by this
// implementation.
func ValidateSchema(s []string) error {
	if len(s) != 2 {
		return fmt.Errorf("Schema should contain two components, but has %d", len(s))
	}

	if s[0] != "tnx" {
		return fmt.Errorf("Schema should be 'tnx' but was '%s'", s[0])
	}

	if s[1] != "0" {
		return fmt.Errorf("Schema version '%s' is unsupported", s[1])
	}

	return nil
}

// ValidateTopology ensures that the TNX topology definition is valid.
func ValidateTopology(t Topology) error {
	identifiers := make(map[string]bool)
	outputs := make(map[string]bool)
	inputs := make(map[string]bool)

	for _, n := range t.Nodes {
		if _, ok := identifiers[n.ID]; ok {
			return fmt.Errorf("node ID: '%s' aliases another identifier", n.ID)
		}
		identifiers[n.ID] = true

		for _, id := range n.Inputs {
			if _, ok := identifiers[id]; ok {
				return fmt.Errorf("Input ID '%s' of node '%s' aliases another identifier", id, n.ID)
			}
			identifiers[id] = true
			inputs[id] = true
		}

		for _, id := range n.Outputs {
			if _, ok := identifiers[id]; ok {
				return fmt.Errorf("Output ID '%s' of node '%s' aliases another identifier", id, n.ID)
			}
			identifiers[id] = true
			outputs[id] = true
		}
	}

	for _, l := range t.Links {
		if _, ok := identifiers[l.Source]; !ok {
			return fmt.Errorf("Link '%v' references nonexistent source '%s'", l, l.Source)
		}

		if _, ok := identifiers[l.Target]; !ok {
			return fmt.Errorf("Link '%v' references nonexistent target '%s'", l, l.Target)
		}

		if _, ok := outputs[l.Source]; !ok {
			return fmt.Errorf("Link '%v' source %s' is not an output", l, l.Source)
		}

		if _, ok := inputs[l.Target]; !ok {
			return fmt.Errorf("Link '%v' target %s' is not an input", l, l.Target)
		}
	}

	return nil
}
