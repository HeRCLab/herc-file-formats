package schema

import (
	"fmt"
	// "github.com/google/go-cmp/cmp"
)

// ParameterValidators - This list is initialized with the official validation
// functions, but users may add additional ones to enforce custom constraints.
// Each function in this list is run on each defined parameter. The TNX object
// given will be the parent of the given parameter, and should be used ONLY in
// a read-only capacity to extract information needed to verify the parameter.
//
// No function in this list should modify either the parameter, or the TNX.  If
// the parameter is valid, it should return nil, and otherwise an error.
var ParameterValidators []func(*TNX, *Parameter, string) error

// SnapshotValidators - as with ParameterValidators, but instead applies to
// snapshot objects.
var SnapshotValidators []func(*TNX, *Snapshot) error

// Validate checks if the TNX is valid. In order for it to have been loaded by
// the JSON decoder, it must have been well formed. T
func Validate(tnx *TNX) error {
	err := ValidateSchema(tnx.Schema)
	if err != nil {
		return err
	}

	err = ValidateTopology(tnx.Topology)
	if err != nil {
		return err
	}

	err = ValidateParameters(tnx)
	if err != nil {
		return err
	}

	err = ValidateSnapshots(tnx)
	if err != nil {
		return err
	}

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

// ValidateParameters ensures that all parameters are valid
func ValidateParameters(tnx *TNX) error {
	for id, param := range tnx.Parameters {
		for _, v := range ParameterValidators {
			err := v(tnx, param, id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetEffectiveDimensions retrieves the dimensions list for a given input
// or output, based on the parameterization of the relevant node.
func (tnx *TNX) GetEffectiveDimensions(ioid string) (*[]int, error) {
	node, err := tnx.LookupNodeByIOID(ioid)
	if err != nil {
		return nil, err
	}

	param, err := tnx.GetParameters(node.ID)
	if err != nil {
		return nil, err
	}

	if (node.Operation == "input") || (node.Operation == "output") {
		// Note: assumptions about the inputs and outputs of
		// input and outputnodes are enforced elsewhere.

		if param.Dimensions == nil {
			return nil, fmt.Errorf("Node ID '%s' omits it's dimension parameter", node.ID)
		}
		return param.Dimensions, nil

	} else if node.Operation == "mlplayer" {
		if param.Neurons == nil {
			return nil, fmt.Errorf("Node ID '%s' omits it's neurons parameter", node.ID)
		}

		if tnx.IsInput(ioid) {
			links, err := tnx.LookupLinkByEndpoint(ioid)
			if err != nil {
				return nil, err
			}

			if len(links) < 1 {
				return nil, fmt.Errorf("Input IO '%s' is unconnected, cannot compute effective dimensions", ioid)
			} else if len(links) > 1 {
				return nil, fmt.Errorf("Input IO '%s' has multiple sources, invalid topology", ioid)
			}

			link := links[0]

			// NOTE: we assume that we are the Target, because
			// this will always be true in a valid TNX file.
			prevnode, err := tnx.LookupNodeByIOID(link.Source)
			if err != nil {
				return nil, err
			}

			return tnx.GetEffectiveDimensions(prevnode.ID)

		} else if tnx.IsOutput(ioid) {
			if param.Neurons == nil {
				return nil, fmt.Errorf("Node ID '%s' omits it's neurons parameter", node.ID)
			}

			return &[]int{*param.Neurons}, nil

		} else {

			return nil, fmt.Errorf("IOID '%s' does not refer to an input or an output", ioid)
		}

	}

	return nil, fmt.Errorf("IOID '%s' implements an unknown operation '%s'", ioid, node.Operation)
}

// GetParameters retrieves the parameters for a given node by it's ID.
func (tnx *TNX) GetParameters(nodeID string) (*Parameter, error) {
	node, err := tnx.LookupNodeByID(nodeID)
	if err != nil {
		return nil, err
	}

	param, ok := tnx.Parameters[node.ID]
	if !ok {
		return nil, fmt.Errorf("No parameters found for node '%s'", node.ID)
	}

	return param, nil
}

// Validator for input and output node parameters
func init() {
	ParameterValidators = append(ParameterValidators, func(tnx *TNX, param *Parameter, id string) error {
		node, err := tnx.LookupNodeByID(id)
		if err != nil {
			return fmt.Errorf("Parameter '%v' applies to invalid node ID '%s', error was: %v",
				param, id, err)
		}

		if (node.Operation != "input") && (node.Operation != "output") {
			return nil
		}

		if param.Dimensions == nil {
			return fmt.Errorf("Parametrization of node '%s' must define a dimension list", id)
		}

		// input node can have one output and no inputs
		if node.Operation == "input" {
			if len(node.Inputs) > 0 {
				return fmt.Errorf("Input node '%s' cannot define any inputs", node.ID)
			}

			if len(node.Outputs) == 0 {
				return fmt.Errorf("Input node '%s' is redundant, it defines no outputs", node.ID)
			}

			if len(node.Outputs) > 1 {
				return fmt.Errorf("Input node '%s' defines multiple outputs", node.ID)
			}
		}

		// output node can have one input and no outputs
		if node.Operation == "output" {
			if len(node.Outputs) > 0 {
				return fmt.Errorf("Output node '%s' cannot define any outputs", node.ID)
			}

			if len(node.Inputs) == 0 {
				return fmt.Errorf("Output node '%s' is redundant, it defines no inputs", node.ID)
			}

			if len(node.Inputs) > 1 {
				return fmt.Errorf("Output node '%s' defines multiple inputs", node.ID)
			}
		}

		// calculate link ID as appropriate
		var linkID string = ""
		if node.Operation == "input" {
			linkID = node.Outputs[0]
		} else /*output*/ {
			linkID = node.Inputs[0]
		}

		// make sure the link exists
		links, err := tnx.LookupLinkByEndpoint(linkID)
		if err != nil {
			return fmt.Errorf("Node '%s' I/O '%s' not referenced by any link", id, linkID)
		}

		// TODO
		for _, link := range links {
			// now we can find the node on the other end of the link
			otherID := link.Source // keep in mind this is the ID of one of the other node's I/Os
			if otherID == node.ID {
				otherID = link.Target
			}
			other, err := tnx.LookupNodeByIOID(otherID)
			if err != nil {
				return fmt.Errorf("Link %v references invalid I/O '%s'", link, otherID)
			}

			otherParam, ok := tnx.Parameters[other.ID]
			if !ok {
				return fmt.Errorf("Node '%s' specifies dimension %v, but connected node '%s' is unparameterized",
					node.ID, param.Dimensions, other.ID)
			}

			fmt.Printf("%v", otherParam)
			return nil
		}

		// TODO
		return nil

	})

}

// Validate Snapshots ensures that all snapshots are valid
func ValidateSnapshots(tnx *TNX) error {
	return nil
}
